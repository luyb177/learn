package tool

import (
	"fmt"
	"learn/check_status/model"
	"net/url"
	"strings"
	"time"
)

func GetParameters(grab *model.SeatInfo) *model.GrabInfo {

	url_start := url.QueryEscape(grab.Start)
	url_end := url.QueryEscape(grab.End)
	//fmt.Println(start)
	//fmt.Println(end)

	start := fmt.Sprintf("%s+%s", grab.Date, url_start)
	end := fmt.Sprintf("%s+%s", grab.Date, url_end)
	//fmt.Println(partS)
	//fmt.Print(partE)
	//格式 1100
	fr_start := url.QueryEscape(grab.Start)
	fr_end := url.QueryEscape(grab.End)
	start_time := strings.Replace(grab.Start, ":", "", -1)
	end_time := strings.Replace(grab.End, ":", "", -1)

	//fmt.Println(start)
	//fmt.Println(end)
	now := time.Now()
	timeMS := now.UnixNano() / int64(time.Millisecond)
	return &model.GrabInfo{
		Seat:      grab.Seat,
		DevId:     "",
		RoomId:    "",
		Date:      grab.Date,
		Start:     start,
		End:       end,
		FrStart:   fr_start,
		FrEnd:     fr_end,
		StartTime: start_time,
		EndTime:   end_time,
		TimeMs:    fmt.Sprintf("%d", timeMS),
	}
}
