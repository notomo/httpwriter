package httpwriter

import (
	"io"
	"net/http"
)

type Transport struct {
	Transport http.RoundTripper
	GetWriter func(*http.Request) (io.WriteCloser, error)
	SetWriter func(io.Writer)
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	writer, err := t.GetWriter(req)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	if t.SetWriter != nil {
		t.SetWriter(writer)
	}

	return t.transport().RoundTrip(req)
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}
