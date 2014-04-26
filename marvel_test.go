package marvel

import (
	"flag"
	"testing"
)

var (
	apiKey = flag.String("pub", "", "Public API key")
	secret = flag.String("priv", "", "Private API secret")
)

func TestRequest(t *testing.T) {
	flag.Parse()

	r, err := NewClient(*apiKey, *secret).Series(2258, CommonRequest{})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	for _, iss := range r.Data.Results {
		t.Logf(iss.Modified.Parse().String())
		t.Logf(iss.Thumbnail.URL(PortraitIncredible))
	}
}
