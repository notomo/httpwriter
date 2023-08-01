package httpwriter

import (
	"io"
	"net/http"
)

type Transport struct {
	GetWriter        func(*http.Request) (io.WriteCloser, error)
	TransportFactory func(io.Writer) http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	writer, err := t.GetWriter(req)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	transport := t.TransportFactory(writer)
	if transport == nil {
		transport = http.DefaultTransport
	}

	return transport.RoundTrip(req)
}
