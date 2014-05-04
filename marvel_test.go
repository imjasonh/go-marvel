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

	c := Client{
		PublicKey:  *apiKey,
		PrivateKey: *secret,
	}
	r, err := c.SingleSeries(2258).Comics(ComicsParams{})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	t.Logf("%+v", r.Data)
	for _, iss := range r.Data.Results {
		t.Logf("%v %s", *iss.IssueNumber, iss.Modified.Parse().String())
		t.Logf(iss.Thumbnail.URL(PortraitIncredible))
	}
	comic, err := r.Data.Results[0].Get(c)
	if err != nil {
		t.Errorf("error getting: %v", err)
	}
	t.Logf("%+v", comic.Data.Results[0])
}
