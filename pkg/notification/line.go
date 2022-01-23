package notification

import (
	"github.com/juunini/simple-go-line-notify/notify"
)

func SendLine(msg string) {
	accessToken := "efpb0Pc6Pyzhcn8745jFQ3zWCNftjfI2u4UsKeLJX3m"
	if err := notify.SendText(accessToken, msg); err != nil {
		panic(err)
	}
}
