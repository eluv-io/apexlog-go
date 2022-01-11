package main

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/kinesis"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
)

func main() {
	log.SetHandler(multi.New(
		text.New(os.Stderr),
		kinesis.New("logs"),
	))

	ctx := log.WithFields(log.Fields{
		{Name: "file", Value: "something.png"},
		{Name: "type", Value: "image/png"},
		{Name: "user", Value: "tobi"},
	})

	for range time.Tick(time.Millisecond * 100) {
		ctx.Info("upload")
		ctx.Info("upload complete")
	}
}
