package request

type User struct {
	Name string `json:"name"`
	QQ   string `json:"qq"`
}

type Grab struct {
	Seat  string `json:"seat"`
	Start string `json:"start"`
	End   string `json:"end"`
	Date  string `json:"date"`
}

type GetRecordReq struct {
	Date string `json:"date"`
}
