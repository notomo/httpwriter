package httpwriter_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/notomo/httpwriter"
)

func ExampleDirectory() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	server := httptest.NewServer(handler)
	defer server.Close()

	path := os.TempDir()
	directory := httpwriter.Directory{Path: path}
	client := http.Client{
		Transport: &httpwriter.Transport{
			TransportFactory: func(w io.Writer) http.RoundTripper {
				logger := log.New(os.Stdout, "", 0)
				logger.SetOutput(w)
				return TransportFunc(func(req *http.Request) (*http.Response, error) {
					logger.Printf("-> %s", req.URL.Path)
					return http.DefaultTransport.RoundTrip(req)
				})
			},
			GetWriter: httpwriter.MustDirectoryWriter(&directory),
		},
	}

	_, err := client.Get(server.URL + "/hello")
	if err != nil {
		panic(err)
	}

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		panic(err)
	}

	tmpfs := os.DirFS(path)
	fileName := fmt.Sprintf("0001_GET@http:@@%s@hello", serverURL.Host)
	f, err := tmpfs.Open(fileName)
	if err != nil {
		panic(err)
	}
	fileContent, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(fileContent))
	// Output: -> /hello
}

func ExampleMemory() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	server := httptest.NewServer(handler)
	defer server.Close()

	memory := httpwriter.Memory{}
	client := http.Client{
		Transport: &httpwriter.Transport{
			TransportFactory: func(w io.Writer) http.RoundTripper {
				logger := log.New(os.Stdout, "", 0)
				logger.SetOutput(w)
				return TransportFunc(func(req *http.Request) (*http.Response, error) {
					logger.Printf("-> %s", req.URL.Path)
					return http.DefaultTransport.RoundTrip(req)
				})
			},
			GetWriter: httpwriter.NewMemoryWriter(&memory),
		},
	}

	{
		_, err := client.Get(server.URL + "/hello1")
		if err != nil {
			panic(err)
		}
	}
	{
		_, err := client.Get(server.URL + "/hello2")
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(memory.Buffers[1].String())
	// Output: -> /hello2
}
