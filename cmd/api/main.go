package main

const (
	gRpcPort = "50001"
)

type Config struct{}

func main() {
	app := Config{}
	app.gPRCListen()
}
