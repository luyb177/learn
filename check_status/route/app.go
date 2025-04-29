package route

import "github.com/gin-gonic/gin"

type App struct {
	r  *gin.Engine
	cr *CheckRoute
}

func NewApp(cr *CheckRoute) *App {
	r := gin.Default()
	cr.NewGroup(r)
	return &App{r: r, cr: cr}
}

func (app *App) Run() {
	app.r.Run(":7777")
}
