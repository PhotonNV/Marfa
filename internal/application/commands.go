package application

type GreetCommand struct {
	ChatID int64
}

type EchoCommand struct {
	ChatID  int64
	Message string
}
