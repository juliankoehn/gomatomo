package matomo

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Error("Missing error `api url must not be empty or null`")
	}

	matomo, err := New(Options{
		PiwikURL: "localhost",
	})
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(matomo.Options.PiwikURL, "piwik") {
		t.Error("Missing piwik.php in MatomoURL")
	}
}

func HttpHandlerTest(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func TestRequestFunc(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/", nil)
	req.Host = "some-name"

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	matomo, err := New(Options{
		PiwikURL: "localhost",
	})
	if err != nil {
		t.Error(err)
	}

	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	backend := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"alive": true}`)

		matomo.Request(r)
	}))
	backend.Listener.Close()
	backend.Listener = l
	backend.Start()
	defer backend.Close()

	_, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	// dnt
	req.Header.Set("DNT", "1")
	rr, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	_ = rr

	//if status := rr.Status; status != "200" {
	//	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	//}
}
