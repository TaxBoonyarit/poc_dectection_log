package notification

import (
	"github.com/juunini/simple-go-line-notify/notify"
)

func sendLine() {
	accessToken := "efpb0Pc6Pyzhcn8745jFQ3zWCNftjfI2u4UsKeLJX3m"
	message := "Hello"
	if err := notify.SendText(accessToken, message); err != nil {
		panic(err)
	}
}
