package main

import (
	"errors"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/apexlogs"
)

func main() {
	url := os.Getenv("APEX_LOGS_URL")
	projectID := os.Getenv("APEX_LOGS_PROJECT_ID")
	authToken := os.Getenv("APEX_LOGS_AUTH_TOKEN")

	h := apexlogs.New(url, projectID, authToken)

	defer h.Close()

	log.SetLevel(log.DebugLevel)
	log.SetHandler(h)

	ctx := log.WithFields(log.Fields{
		{Name: "file", Value: "something.png"},
		{Name: "type", Value: "image/png"},
		{Name: "user", Value: "tobi"},
	})

	go func() {
		for range time.Tick(time.Second) {
			ctx.Debug("doing stuff")
		}
	}()

	go func() {
		for range time.Tick(100 * time.Millisecond) {
			ctx.Info("uploading")
			ctx.Info("upload complete")
		}
	}()

	go func() {
		for range time.Tick(time.Second) {
			ctx.Warn("upload slow")
		}
	}()

	go func() {
		for range time.Tick(2 * time.Second) {
			err := errors.New("boom")
			ctx.WithError(err).Error("upload failed")
		}
	}()

	select {}
}
