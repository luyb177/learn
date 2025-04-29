package main

func main() {
	app, err := InitApp("config/config.yml")
	if err != nil {
		panic(err)
	}
	app.Run()

	//grabInfo := model.GrabInfo{
	//	Seat:      "N2001",
	//	DevId:     "101699933",
	//	RoomId:    "101699189",
	//	Date:      "2025-03-25",
	//	Start:     "2025-03-25 19:30",
	//	End:       "2025-03-25 22:00",
	//	FrStart:   "19%3A30",
	//	FrEnd:     "22%3A00",
	//	StartTime: "1930",
	//	EndTime:   "2200",
	//	TimeMs:    "45648646546",
	//	CheckTime: "today",
	//}
	//viperSetting := config.NewViperSetting("config/config.yml")
	//accountConfig := config.NewAccount(viperSetting)
	//clientClient := client.NewClient(accountConfig)
	//clientClient.Login()
	//monitor := service.NewMonitorServiceImpl(clientClient)
	//monitor.CheckSeatStatus(&grabInfo, true)
	//todo 写swagger文档，开始测试
	// 添加信息发送内容
	// 暴露监视接口
	// 最终信息存储到数据库中
}
