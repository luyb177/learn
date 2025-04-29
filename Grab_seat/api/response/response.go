package response

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginRes struct {
	Ret int    `json:"ret"`
	Act string `json:"act"`
	Msg string `json:"msg"`
}

type GrabSeatEvent struct {
	Seat    string `json:"seat"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Status  string `json:"status"` // pending 监控中 // success 成功 // failed 失败
	Content string `json:"content"`
}

// Res 爬取的结构体
type Res struct {
	Ret  int    `json:"ret"`
	Act  string `json:"act"`
	Msg  string `json:"msg"`
	Data []Info `json:"data"`
}

type Info struct {
	RoomName string `json:"roomName"`
	Title    string `json:"title"`
	Ts       []T    `json:"ts"`
}

type T struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	State  string `json:"state"`
	Title  string `json:"title"`
	Owner  string `json:"owner"`
	Occupy bool   `json:"occupy"`
}
