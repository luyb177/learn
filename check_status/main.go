package main

func main() {
	app, err := InitApp("config/config.yml")
	if err != nil {
		panic(err)
	}
	app.Run()
}
