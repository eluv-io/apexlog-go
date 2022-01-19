**Warning - incompatible changes.** 

This fork changes the excellent `apex/log` framework in a way that makes it NOT backward compatible with upstream:

* `Fields` is now a slice rather than a map: fields are no more reordered when logging occurs.
* Add a `Trace` level for super detailed logging: the original `Trace` function has been renamed to `Watch`.

Other changes: 

* use `sync.Pool` for entries and field instances whenever possible.
* logging functions now have an optional `kv ...interface{}` vararg parameter expected to be key/value pairs each added as a log field.  Values of type `error` can be passed alone and are automatically  assigned to a key 'error'. 

![Structured logging for golang](assets/title.png)

Package log implements a simple structured logging API inspired by Logrus, designed with centralization in mind. Read more on [Medium](https://medium.com/@tjholowaychuk/apex-log-e8d9627f4a9a#.rav8yhkud).

## Handlers

- __apexlogs__ – handler for [Apex Logs](https://apex.sh/logs/)
- __cli__ – human-friendly CLI output
- __discard__ – discards all logs
- __es__ – Elasticsearch handler
- __graylog__ – Graylog handler
- __json__ – JSON output handler
- __kinesis__ – AWS Kinesis handler
- __level__ – level filter handler
- __logfmt__ – logfmt plain-text formatter
- __memory__ – in-memory handler for tests
- __multi__ – fan-out to multiple handlers
- __papertrail__ – Papertrail handler
- __text__ – human-friendly colored output
- __delta__ – outputs the delta between log calls and spinner

## Example

Example using the [Apex Logs](https://apex.sh/logs/) handler.

```go
package main

import (
	"errors"
	"time"

	"github.com/eluv-io/apexlog-go"
)

func main() {
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
		ctx.Errorf("failed to upload %s", "img.png")
	}
}
```

---

[![Build Status](https://semaphoreci.com/api/v1/projects/d8a8b1c0-45b0-4b89-b066-99d788d0b94c/642077/badge.svg)](https://semaphoreci.com/tj/log)
[![GoDoc](https://godoc.org/github.com/apex/log?status.svg)](https://godoc.org/github.com/apex/log)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)

<a href="https://apex.sh"><img src="http://tjholowaychuk.com:6000/svg/sponsor"></a>
