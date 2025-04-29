package request

type Grab struct {
	Seat  string `json:"seat"`
	Start string `json:"start"`
	End   string `json:"end"`
}
