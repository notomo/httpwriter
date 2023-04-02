package httpwriter_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/notomo/httpwriter"
)

func TestTransportWithMemory(t *testing.T) {
	server, request := NewServer()
	defer server.Close()

	memory := httpwriter.Memory{}
	logger := log.New(os.Stdout, "", 0)
	client := http.Client{
		Transport: &httpwriter.Transport{
			Transport: TransportFunc(func(req *http.Request) (*http.Response, error) {
				logger.Printf(req.URL.Path)
				return http.DefaultTransport.RoundTrip(req)
			}),
			GetWriter: httpwriter.NewMemoryWriter(&memory),
			SetWriter: logger.SetOutput,
		},
	}

	request(t, client, "/1")
	request(t, client, "/2")

	{
		b := memory.Buffers[0]
		if got, want := b.String(), "/1\n"; got != want {
			t.Errorf("want %q, but actual: %q", want, got)
		}
		if got, want := b.Request.URL.Path, "/1"; got != want {
			t.Errorf("want %q, but actual: %q", want, got)
		}
	}
	{
		b := memory.Buffers[1]
		if got, want := b.String(), "/2\n"; got != want {
			t.Errorf("want %q, but actual: %q", want, got)
		}
		if got, want := b.Request.URL.Path, "/2"; got != want {
			t.Errorf("want %q, but actual: %q", want, got)
		}
	}
}

func TestTransportWithDirectory(t *testing.T) {
	server, request := NewServer()
	defer server.Close()

	path := t.TempDir()
	directory := httpwriter.Directory{Path: path}
	logger := log.New(os.Stdout, "", 0)
	client := http.Client{
		Transport: &httpwriter.Transport{
			Transport: TransportFunc(func(req *http.Request) (*http.Response, error) {
				logger.Printf(req.URL.Path)
				return http.DefaultTransport.RoundTrip(req)
			}),
			GetWriter: httpwriter.MustDirectoryWriter(&directory),
			SetWriter: logger.SetOutput,
		},
	}

	request(t, client, "/1")
	request(t, client, "/2")

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	tmpfs := os.DirFS(path)
	{
		fileName := fmt.Sprintf("0001_GET@http:@@%s@1", serverURL.Host)
		got := GetFileContent(t, tmpfs, fileName)
		want := "/1\n"
		if got != want {
			t.Errorf("want %q, but actual: %q", want, string(got))
		}
	}
	{
		fileName := fmt.Sprintf("0002_GET@http:@@%s@2", serverURL.Host)
		got := GetFileContent(t, tmpfs, fileName)
		want := "/2\n"
		if string(got) != want {
			t.Errorf("want %q, but actual: %q", want, string(got))
		}
	}
}
