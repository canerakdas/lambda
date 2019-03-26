package static

type emailSender struct {
	Email    string
	Password string
	Host     string
	Port     string
}

const mode = "Development"

var Sender = emailSender{
	Email:    "example@gmail.com",
	Password: "1234",
	Host:     "smtp.gmail.com",
	Port:     ":587",
}
