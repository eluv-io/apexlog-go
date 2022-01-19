package main

import (
	"errors"
	"os"
	"time"

	log "github.com/eluv-io/apexlog-go"
	"github.com/eluv-io/apexlog-go/handlers/json"
	"github.com/eluv-io/apexlog-go/handlers/multi"
	"github.com/eluv-io/apexlog-go/handlers/text"
)

func main() {
	log.SetHandler(multi.New(
		text.New(os.Stderr),
		json.New(os.Stderr),
	))

	ctx := log.WithFields(log.Fields{
		{Name: "file", Value: "something.png"},
		{Name: "type", Value: "image/png"},
		{Name: "user", Value: "tobi"},
	})

	for range time.Tick(time.Millisecond * 200) {
		ctx.Info("upload")
		ctx.Info("upload complete")
		ctx.Warn("upload retry")
		ctx.WithError(errors.New("unauthorized")).Error("upload failed")
	}
}
