package main

import (
	"os"

	log "github.com/eluv-io/apexlog-go"
	"github.com/pkg/errors"

	"github.com/eluv-io/apexlog-go/handlers/logfmt"
)

func main() {
	log.SetHandler(logfmt.New(os.Stderr))

	filename := "something.png"
	body := []byte("whatever")

	ctx := log.WithField("filename", filename)

	err := upload(filename, body)
	if err != nil {
		ctx.WithError(err).Error("upload failed")
	}
}

// Faux upload.
func upload(name string, b []byte) error {
	err := put("/images/"+name, b)
	if err != nil {
		return errors.Wrap(err, "uploading to s3")
	}

	return nil
}

// Faux PUT.
func put(key string, b []byte) error {
	return errors.New("unauthorized")
}
