package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"learn/Grab_seat/api/response"
	"sync"
)

var (
	eventChannels = make(map[chan response.GrabSeatEvent]struct{})
	eventMutex    sync.RWMutex
)

type SseService interface {
	SSEServer(c *gin.Context)
}

type SseServiceImpl struct {
}

func NewSseServiceImpl() *SseServiceImpl {
	return &SseServiceImpl{}
}

// BroadcastGrabEvent  全局事件广播
func BroadcastGrabEvent(seatEvent response.GrabSeatEvent) {
	eventMutex.RLock()
	defer eventMutex.RUnlock()
	for ch := range eventChannels {
		select {
		case ch <- seatEvent:
		default:
			// 避免阻塞，丢弃事件或处理满通道
		}
	}
}

func (ssr *SseServiceImpl) SSEServer(c *gin.Context) {
	// 初始化SSE头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	eventChan := make(chan response.GrabSeatEvent)
	eventMutex.Lock()
	eventChannels[eventChan] = struct{}{}
	eventMutex.Unlock()

	defer func() {
		eventMutex.Lock()
		delete(eventChannels, eventChan)
		eventMutex.Unlock()
		close(eventChan)
	}()

	for {
		select {
		case event := <-eventChan:
			switch event.Status {
			case "pending":
				c.SSEvent(event.Status, fmt.Sprintf("用户预定的座位%s已被占用，后台正在监控...", event.Seat))
			case "success":
				c.SSEvent(event.Status, fmt.Sprintf("用户预定的座位%s预约成功，开始时间%s-截止时间%s,请及时前往使用...", event.Seat, event.Start, event.End))
			case "failed":
				c.SSEvent(event.Status, fmt.Sprintf("用户预定的座位%s预约失败，请重新预定...", event.Seat))
			case "remind":
				c.SSEvent(event.Status, fmt.Sprintf("用户预定的座位%s正在预约中，%s", event.Seat, event.Content))
			case "completed":
				c.SSEvent(event.Status, fmt.Sprintf("用户预定的座位%s已经使用完毕", event.Seat))
			default:
				c.SSEvent(event.Status, fmt.Sprintf("未知错误，请重新预定..."))
			}
			c.Writer.Flush()
		case <-c.Done():
			fmt.Println("断开连接")
			return
		}
	}
}
