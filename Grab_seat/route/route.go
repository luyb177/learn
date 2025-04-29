package route

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	r  *gin.Engine
	gr *GrabRoot
	sr *SseRoute
	mr *MonitorRoot
}

func NewApp(gr *GrabRoot, sr *SseRoute, mr *MonitorRoot) *App {
	r := gin.Default()
	gr.GrabGroup(r)
	sr.SseGroup(r)
	mr.NewGroup(r)
	return &App{
		r:  r,
		gr: gr,
		sr: sr,
		mr: mr,
	}
}

func (app *App) Run() {
	app.r.Run()
}
