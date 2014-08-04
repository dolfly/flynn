package discoverd_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/flynn/flynn/discoverd/client"
	"github.com/flynn/flynn/discoverd/client/balancer"
	"github.com/flynn/flynn/discoverd/client/dialer"
)

func TestHTTPClient(t *testing.T) {
	client, cleanup := setup(t)
	defer cleanup()

	hc := dialer.NewHTTPClient(client)
	_, err := hc.Get("http://httpclient/")
	if ue, ok := err.(*url.Error); !ok || ue.Err != balancer.ErrNoServices {
		t.Error("Expected err to be ErrNoServices, got", ue.Err)
	}

	s := httptest.NewServer(nil)
	defer s.Close()
	client.Register("httpclient", s.URL[7:])

	set, _ := discoverd.NewServiceSet("httpclient")
	waitUpdates(t, set, true, 1)()
	set.Close()

	_, err = hc.Get("http://httpclient/")
	if err != nil {
		t.Error("Unexpected error during request:", err)
	}
}