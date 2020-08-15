package core

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Publisher interface {
	Publish(m MatchEvent) error
}

// WebSource represents a web-based publisher
type WebSource struct {
	client      *http.Client
	endpoint    string
	method      string
	contentType string
}

// Publish writes a MatchEvent to the WebSource.
func (w *WebSource) Publish(m MatchEvent) error {
	var data string

	switch w.contentType {
	case "application/json":
		data = m.Json()
	// Add new cases for additional content types
	default:
		// Nothing
	}

	req, err := http.NewRequest(w.method, w.endpoint, bytes.NewBufferString(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", w.contentType)
	_, err = w.client.Do(req)
	if err != nil {
		return err
	}

	// May need to capture the response and check for status codes like 404 or
	// 503 so that errors can be returned if needed.

	return nil
}

// Returns a new WebSource that will send MatchEvents to the given enpoint
// using the given method and content type.
func NewWebSource(endpoint, method, contentType string) (WebSource, error) {
	var w WebSource

	if !(method == "POST" || method == "GET") {
		return w, fmt.Errorf("method must be POST or GET, received %s", method)
	}

	w.client = &http.Client{Timeout: 10 * time.Second}
	w.endpoint = endpoint
	w.method = method
	w.contentType = contentType

	return w, nil
}

// DelimitedSource represents a Comma or Tab delimited publisher.
type DelimitedSource struct {
	writer *csv.Writer
}

// Publish writes a MatchEvent to the file.
func (w *DelimitedSource) Publish(m MatchEvent) error {
	err := w.writer.Write(m.Line())
	if err != nil {
		return err
	}

	w.writer.Flush()

	return nil
}

// NewDelimitedSource returns a new comma or tab delimited writer. If a header
// is provided, it is written to the file.
func NewDelimitedSource(filename string, delimiter rune, header []string) (DelimitedSource, error) {
	var d DelimitedSource

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return d, err
	}

	d.writer = csv.NewWriter(file)
	d.writer.Comma = delimiter

	if header != nil {
		d.writer.Write(header)
	}

	return d, nil
}
