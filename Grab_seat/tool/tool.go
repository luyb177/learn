package tool

import (
	"fmt"
	"learn/Grab_seat/api/request"
	"learn/Grab_seat/model"
	"net/url"
	"strings"
	"time"
)

func GetParameters(grab *request.Grab) *model.GrabInfo {
	start := url.QueryEscape(grab.Start)
	end := url.QueryEscape(grab.End)
	//fmt.Println(star)
	//fmt.Println(endd)
	//获取后面的一部分
	partS := strings.Split(grab.Start, " ")
	partE := strings.Split(grab.End, " ")
	//fmt.Println(partS)
	//fmt.Print(partE)
	//格式 1100
	fr_start := url.QueryEscape(partS[1])
	fr_end := url.QueryEscape(partE[1])
	start_time := strings.Replace(partS[1], ":", "", -1)
	end_time := strings.Replace(partE[1], ":", "", -1)

	//fmt.Println(start)
	//fmt.Println(end)
	now := time.Now()
	timeMS := now.UnixNano() / int64(time.Millisecond)
	return &model.GrabInfo{
		Seat:      grab.Seat,
		DevId:     "",
		RoomId:    "",
		Date:      partE[0],
		Start:     start,
		End:       end,
		FrStart:   fr_start,
		FrEnd:     fr_end,
		StartTime: start_time,
		EndTime:   end_time,
		TimeMs:    fmt.Sprintf("%d", timeMS),
	}
}
