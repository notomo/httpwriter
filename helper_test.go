package httpwriter_test

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

type TransportFunc func(*http.Request) (*http.Response, error)

func (t TransportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return t(req)
}

func NewServer() (*httptest.Server, func(*testing.T, http.Client, string)) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	server := httptest.NewServer(handler)
	return server, func(t *testing.T, client http.Client, path string) {
		t.Helper()
		if _, err := client.Get(server.URL + path); err != nil {
			t.Fatal(err)
		}
	}
}

func GetFileContent(t *testing.T, tmpfs fs.FS, fileName string) string {
	t.Helper()

	if err := fstest.TestFS(tmpfs, fileName); err != nil {
		t.Fatal(err)
	}

	f, err := tmpfs.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}

	got, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	return string(got)
}
