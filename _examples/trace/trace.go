package main

import (
	"os"
	"time"

	log "github.com/eluv-io/apexlog-go"
	"github.com/eluv-io/apexlog-go/handlers/text"
)

func work(ctx log.Interface) (err error) {
	path := "Readme.md"
	defer ctx.WithField("path", path).Watch("opening").Stop(&err)
	_, err = os.Open(path)
	return
}

func main() {
	log.SetHandler(text.New(os.Stderr))

	ctx := log.WithFields(log.Fields{
		{Name: "app", Value: "myapp"},
		{Name: "env", Value: "prod"},
	})

	for range time.Tick(time.Second) {
		_ = work(ctx)
	}
}
