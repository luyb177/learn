package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"golang.org/x/sync/semaphore"
	"io"
	"learn/Grab_seat/api/request"
	"learn/Grab_seat/api/response"
	"learn/Grab_seat/client"
	"learn/Grab_seat/config"
	"learn/Grab_seat/dao"
	"learn/Grab_seat/model"
	"learn/Grab_seat/tool"
	"log"
	"sync"
	"time"
)

type GrabService interface {
	SendMsg(grab *request.Grab) error
}

type GrabServiceImpl struct {
	Producer     sarama.AsyncProducer
	Consumer     sarama.Consumer
	client       *client.Client
	cfg          *config.KafkaConfig
	GrabDao      dao.GrabDAO
	ContentDao   dao.ContentDAO
	ms           MonitorService
	errorChannel chan error
	// 新增上下文相关字段
	ctx      context.Context    // 主上下文
	cancel   context.CancelFunc // 取消函数
	ctxMutex sync.Mutex         // 上下文锁
}

func NewGrabServiceImpl(cfg *config.KafkaConfig, GrabDao dao.GrabDAO, Content dao.ContentDAO, client *client.Client, ms MonitorService) *GrabServiceImpl {
	// 配置生产者
	producer, err := newKafkaProducer([]string{cfg.Addr})
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	// 配置消费者
	consumer, err := newKafkaConsumer([]string{cfg.Addr})
	if err != nil {
		// 确保生产者资源被正确释放
		producer.Close()
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}

	var ctx = context.Background()
	// 创建可取消的上下文
	childCtx, cancel := context.WithCancel(ctx)
	// 让client获取登录信息
	client.Login()
	// 获取其他goroutine的错误消息
	errorChan := make(chan error)
	var grabServiceImpl = GrabServiceImpl{
		Producer:     producer,
		Consumer:     consumer,
		cfg:          cfg,
		GrabDao:      GrabDao,
		ContentDao:   Content,
		client:       client,
		ms:           ms,
		errorChannel: errorChan,
		cancel:       cancel,
		ctx:          childCtx,
	}

	// 启动消费协程
	go func() {
		grabServiceImpl.StartConsume()
	}()

	// 开启错误处理协程
	go func() {
		grabServiceImpl.CheckError()
	}()
	return &grabServiceImpl
}

// 新建 Kafka 生产者
func newKafkaProducer(brokers []string) (sarama.AsyncProducer, error) {
	ProducerConfig := sarama.NewConfig()
	ProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	ProducerConfig.Producer.Compression = sarama.CompressionGZIP
	ProducerConfig.Producer.Return.Successes = true
	ProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewAsyncProducer(brokers, ProducerConfig)
	if err != nil {
		log.Printf("Error creating Kafka producer: %v", err)
		return nil, err
	}

	log.Println("Kafka producer initialized successfully")
	return producer, nil
}

// 新建 Kafka 消费者
func newKafkaConsumer(brokers []string) (sarama.Consumer, error) {
	ConsumerConfig := sarama.NewConfig()
	ConsumerConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, ConsumerConfig)
	if err != nil {
		log.Printf("Error creating Kafka consumer: %v", err)
		return nil, err
	}

	log.Println("Kafka consumer initialized successfully")
	return consumer, nil
}

// SendMsg 发送消息到 Kafka
func (gsr *GrabServiceImpl) SendMsg(grab *request.Grab) error {
	// 序列化抓取请求
	grabJson, err := json.Marshal(grab)
	if err != nil {
		return fmt.Errorf("failed to marshal grab: %w", err)
	}

	// 创建 Kafka 消息
	msg := &sarama.ProducerMessage{
		Topic: "grab_seat",
		Key:   sarama.StringEncoder("book"),
		Value: sarama.ByteEncoder(grabJson),
	}

	// 发送消息到 Kafka 生产者
	select {
	case gsr.Producer.Input() <- msg:
	default:
		return fmt.Errorf("kafka producer input channel is full")
	}

	go func() {
		for {
			select {
			case successMsg := <-gsr.Producer.Successes():
				log.Printf("Message sent successfully: topic=%s, partition=%d, offset=%d, grab=%+v",
					successMsg.Topic, successMsg.Partition, successMsg.Offset, grab)
			case errMsg := <-gsr.Producer.Errors():
				log.Printf("Failed to send message: %v, grab=%+v", errMsg, grab)
			}
		}
	}()

	return nil
}

// StartConsume 开始消费消息
func (gsr *GrabServiceImpl) StartConsume() {

	partitions, err := gsr.Consumer.Partitions("grab_seat")

	if err != nil {
		log.Println("Error getting list of partitions: ", err)
	}

	// 使用协程池限制并发数
	sem := semaphore.NewWeighted(20)

	for _, partition := range partitions {
		go func(partition int32) {
			partitionConsumer, err := gsr.Consumer.ConsumePartition("grab_seat", partition, sarama.OffsetNewest)
			if err != nil {
				log.Println(err)
			}
			defer partitionConsumer.Close()

			for msg := range partitionConsumer.Messages() {
				// 控制并发数
				if err := sem.Acquire(gsr.ctx, 1); err != nil {
					log.Println("Error acquiring semaphore: ", err)
					continue
				}

				go func(m *sarama.ConsumerMessage) {
					defer sem.Release(1)

					// 监听全局上下文取消
					select {
					case <-gsr.ctx.Done():
						log.Println("Canceling message processing due to shutdown")
						return
					default:
					}

					log.Printf("Received message: partition=%d, offset=%d, key=%s, value=%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

					var grab request.Grab
					if err := json.Unmarshal(m.Value, &grab); err != nil {
						log.Println("Failed to unmarshal message:", err)
						return
					}

					if err := gsr.StartBook(&grab); err != nil {
						log.Println("Failed to process booking:", err)
						// 将错误发送到统一通道处理
						gsr.errorChannel <- err
					}
				}(msg)
			}
		}(partition)

	}
}

// StartBook 开始预约 根据时间判断是今天还是明天
// 选择相应的函数处理
func (gsr *GrabServiceImpl) StartBook(grab *request.Grab) error {
	//获取参数
	devId, roomId, err := gsr.GrabDao.FindSeatId(grab.Seat)
	if err != nil {
		return err
	}
	grabInfo := tool.GetParameters(grab)
	grabInfo.DevId = devId
	grabInfo.RoomId = roomId

	// 获取时间
	start, err := time.Parse("2006-01-02 15:04", grab.Start)
	if err != nil {
		return err
	}

	// 获取当前时间
	now := time.Now()
	// 判断预约日期是今天还是明天
	if start.Year() == now.Year() && start.Month() == now.Month() && start.Day() == now.Day() {
		// 今天 - 判断时间
		if start.Before(now) {
			// 自动更改时间
			// 更改为当前时间
			grab.Start = now.Format("2006-01-02 15:04")
			grabInfo = tool.GetParameters(grab)
			return gsr.bookImmediately(grabInfo)
		}
		return gsr.bookWithCheckStatus(grabInfo)

	} else if start.Year() == now.Year() && start.Month() == now.Month() && start.Day() == now.Day()+1 {
		// 预约明天的预约
		return gsr.bookTomorrow(grabInfo)
	} else {
		// 既不是今天也不是明天
		// 发送消息
		BroadcastGrabEvent(response.GrabSeatEvent{
			Seat:    grab.Seat,
			Start:   grab.Start,
			End:     grab.End,
			Status:  "failed",
			Content: "预约时间有误，暂时只能预约今天或明天的座位",
		})
		// 最终消息存到数据库
		err := gsr.ContentDao.AddContent(&model.Content{
			Seat:    grab.Seat,
			Start:   grab.Start,
			End:     grab.End,
			Status:  "failed",
			Content: "预约时间有误，暂时只能预约今天或明天的座位",
		})
		if err != nil {
			log.Println(err.Error())
		}
		return errors.New("can only book seats for today or tomorrow")
	}

}

// bookImmediately 今天的预约先检测座位状态
func (gsr *GrabServiceImpl) bookWithCheckStatus(grabInfo *model.GrabInfo) error {
	// 获取座位状态
	ok := gsr.ms.CheckOneSeatStatus(grabInfo)

	if ok {
		// 座位被占用
		// 每5分钟检测一次
		gsr.MonitorAndSendMessage(grabInfo)
		return nil
	} else {
		// 座位未被占用
		// 开始直接抢
		path := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?dialogid=&dev_id=" + grabInfo.DevId + "&lab_id=&kind_id=&room_id=&type=dev&prop=&test_id=&term=&Vnumber=&classkind=&test_name=&start=" + grabInfo.Start + "&end=" + grabInfo.End + "&start_time=" + grabInfo.StartTime + "&end_time=" + grabInfo.EndTime + "&up_file=&memo=&act=set_resv&_=" + grabInfo.TimeMs

		respon, err := gsr.client.Client.Get(path)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		defer respon.Body.Close()
		body, _ := io.ReadAll(respon.Body)
		var loginRes response.LoginRes
		err = json.Unmarshal(body, &loginRes)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		fmt.Println(loginRes)

		if loginRes.Ret == 1 {
			BroadcastGrabEvent(response.GrabSeatEvent{
				Seat:   grabInfo.Seat,
				Start:  grabInfo.Start,
				End:    grabInfo.End,
				Status: "success",
			})

			err := gsr.ContentDao.AddContent(&model.Content{
				Seat:    grabInfo.Seat,
				Start:   grabInfo.Start,
				End:     grabInfo.End,
				Status:  "success",
				Content: "预约成功",
			})

			if err != nil {
				log.Println(err.Error())
			}
			// 开一个协程
			go func() {
				defer func() {
					// todo
					// 监控使用完毕
				}()
			}()
			return nil
		} else if loginRes.Ret == -1 {
			//cookie失效,重新登录
			gsr.RetryLogin(grabInfo)
			return nil
		} else if loginRes.Ret == 0 {
			log.Println(loginRes.Msg)
		}
		return nil
	}

}

// 立即预约没有检测座位状态
func (gsr *GrabServiceImpl) bookImmediately(grabInfo *model.GrabInfo) error {

	path := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/reserve.aspx?dialogid=&dev_id=" + grabInfo.DevId + "&lab_id=&kind_id=&room_id=&type=dev&prop=&test_id=&term=&Vnumber=&classkind=&test_name=&start=" + grabInfo.Start + "&end=" + grabInfo.End + "&start_time=" + grabInfo.StartTime + "&end_time=" + grabInfo.EndTime + "&up_file=&memo=&act=set_resv&_=" + grabInfo.TimeMs

	respon, err := gsr.client.Client.Get(path)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer respon.Body.Close()
	body, _ := io.ReadAll(respon.Body)
	var loginRes response.LoginRes
	err = json.Unmarshal(body, &loginRes)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	fmt.Println(loginRes)
	if loginRes.Ret == 1 {
		BroadcastGrabEvent(response.GrabSeatEvent{
			Seat:   grabInfo.Seat,
			Start:  grabInfo.Start,
			End:    grabInfo.End,
			Status: "success",
		})
		err := gsr.ContentDao.AddContent(&model.Content{
			Seat:    grabInfo.Seat,
			Start:   grabInfo.Start,
			End:     grabInfo.End,
			Status:  "success",
			Content: "预约成功",
		})
		if err != nil {
			log.Println(err.Error())
		}
		return nil
	} else if loginRes.Ret == -1 {
		// cookie 失效,重新登录
		gsr.RetryLogin(grabInfo)
		return nil
	} else if loginRes.Ret == 0 {
		log.Println(loginRes.Msg)
	}
	return nil
}

func (gsr *GrabServiceImpl) bookTomorrow(grabInfo *model.GrabInfo) error {
	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), 17, 50, 0, 0, now.Location())
	sixPm := time.Date(now.Year(), now.Month(), now.Day(), 18, 00, 0, 0, now.Location())
	// 如果当前时间已经晚于目标时间
	if now.Equal(targetTime) {
		// 直接开始高频监视
		BroadcastGrabEvent(response.GrabSeatEvent{
			Seat:    grabInfo.Seat,
			Start:   grabInfo.Start,
			End:     grabInfo.End,
			Status:  "remind",
			Content: "正在高频监视你预定的座位中...",
		})

		err := gsr.ContentDao.AddContent(&model.Content{
			Seat:    grabInfo.Seat,
			Start:   grabInfo.Start,
			End:     grabInfo.End,
			Status:  "remind",
			Content: "离座位开始预定还有10分钟，正在监视中...",
		})
		if err != nil {
			log.Println(err.Error())
		}
	} else if now.After(targetTime) {
		if now.After(sixPm) {
			return gsr.bookWithCheckStatus(grabInfo)
		}
		return gsr.startHighFrequencyCheck(grabInfo)
	}

	// 计算精确休眠时间（含安全缓冲）
	sleepDuration := targetTime.Sub(now) - 30*time.Second // 提前30秒缓冲
	if sleepDuration > 0 {
		// 使用context控制安全休眠
		ctx, cancel := context.WithTimeout(context.Background(), sleepDuration)
		defer cancel()

		select {
		case <-ctx.Done(): // 正常唤醒
		case <-gsr.ctx.Done(): // 全局上下文取消（优雅退出）
			return fmt.Errorf("process canceled")
		}

		BroadcastGrabEvent(response.GrabSeatEvent{
			Seat:    grabInfo.Seat,
			Start:   grabInfo.Start,
			End:     grabInfo.End,
			Status:  "remind",
			Content: "离座位开始预定还有10分钟，正在监视中...",
		})
		err := gsr.ContentDao.AddContent(&model.Content{
			Seat:    grabInfo.Seat,
			Start:   grabInfo.Start,
			End:     grabInfo.End,
			Status:  "remind",
			Content: "离座位开始预定还有10分钟，正在监视中...",
		})
		if err != nil {
			log.Println(err.Error())
		}

	}
	// 进入高频检查阶段
	return gsr.startHighFrequencyCheck(grabInfo)
}

// 高频检查阶段
func (gsr *GrabServiceImpl) startHighFrequencyCheck(grabInfo *model.GrabInfo) error {
	deadline := time.Now().Add(11 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	now := time.Now()
	sixPM := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	for {
		select {
		case <-ticker.C:
			if now.Equal(sixPM) {
				BroadcastGrabEvent(response.GrabSeatEvent{
					Seat:    grabInfo.Seat,
					Start:   grabInfo.Start,
					End:     grabInfo.End,
					Status:  "remind",
					Content: "6点啦，正在为你预约座位",
				})
				err := gsr.ContentDao.AddContent(&model.Content{
					Seat:    grabInfo.Seat,
					Start:   grabInfo.Start,
					End:     grabInfo.End,
					Status:  "remind",
					Content: "6点啦，正在为你预约座位",
				})
				if err != nil {
					log.Println(err.Error())
				}
				if err := gsr.bookImmediately(grabInfo); err != nil {
					// 加入重试逻辑
					return gsr.retryBooking(grabInfo, 3) // 重试3次
				}
				return nil
			} else if now.After(sixPM) {
				BroadcastGrabEvent(response.GrabSeatEvent{
					Seat:    grabInfo.Seat,
					Start:   grabInfo.Start,
					End:     grabInfo.End,
					Status:  "remind",
					Content: "现在已经是" + now.Format(time.DateTime) + ",正在为你预约座位",
				})
				err := gsr.ContentDao.AddContent(&model.Content{
					Seat:    grabInfo.Seat,
					Start:   grabInfo.Start,
					End:     grabInfo.End,
					Status:  "remind",
					Content: "现在已经是" + now.Format(time.DateTime) + ",正在为你预约座位",
				})
				if err != nil {
					log.Println(err.Error())
				}
				if err := gsr.bookImmediately(grabInfo); err != nil {
					// 加入重试逻辑
					return gsr.retryBooking(grabInfo, 3) // 重试3次
				}
				return nil
			}
		case <-time.After(time.Until(deadline)):
			return fmt.Errorf("high frequency check timeout")
		case <-gsr.ctx.Done():
			return fmt.Errorf("process canceled")
		}
	}
}

// 带锁的重试逻辑（防止并发冲突）
func (gsr *GrabServiceImpl) retryBooking(grabInfo *model.GrabInfo, maxRetries int) error {
	gsr.ctxMutex.Lock() // 使用互斥锁保证原子性
	defer gsr.ctxMutex.Unlock()

	for i := 0; i < maxRetries; i++ {
		if err := gsr.bookImmediately(grabInfo); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond) // 休眠后再次尝试
	}
	return fmt.Errorf("max retries exceeded")
}

// MonitorAndSendMessage 监控被占用的座位
func (gsr *GrabServiceImpl) MonitorAndSendMessage(grabInfo *model.GrabInfo) {
	BroadcastGrabEvent(response.GrabSeatEvent{
		Seat:   grabInfo.Seat,
		Start:  grabInfo.Start,
		End:    grabInfo.End,
		Status: "pending",
	})

	// 每5分钟检测一次
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ok := gsr.ms.CheckOneSeatStatus(grabInfo)
			if !ok {
				err := gsr.bookImmediately(grabInfo)
				if err != nil {
					log.Println(err.Error())
				}
			}
		case <-gsr.ctx.Done():
			log.Println("app stop")
			return
		default:
			//防止堵塞
		}
	}
}

// RetryLogin 重试登录
func (gsr *GrabServiceImpl) RetryLogin(grabInfo *model.GrabInfo) {
	gsr.client.Login()
	err := gsr.bookImmediately(grabInfo)
	if err != nil {
		log.Println(err.Error())
	}
}

// TailAfter 追踪预约使用情况
func (gsr *GrabServiceImpl) TailAfter(grabInfo *model.GrabInfo) {
	// 5分钟检测一次
	ticker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			ok := gsr.ms.CheckOneSeatStatus(grabInfo)
			if !ok {
				BroadcastGrabEvent(response.GrabSeatEvent{
					Seat:   grabInfo.Seat,
					Start:  grabInfo.Start,
					End:    grabInfo.End,
					Status: "completed",
				})
			}
		case <-gsr.ctx.Done():
			log.Println("app stop")
			return
		default:
			//防止堵塞

		}
	}
}

// CheckError 检查错误
func (gsr *GrabServiceImpl) CheckError() {
	for {
		select {
		case err := <-gsr.errorChannel:
			log.Println(err.Error())
		case <-gsr.ctx.Done():
			log.Println("app stop")
			return
		default:
			//防止阻塞
		}
	}
}
