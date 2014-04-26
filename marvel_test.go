package marvel

import (
	"testing"
)

func TestRequest(t *testing.T) {
	r, err := NewClient("d96b5157cfc7a60cbfaa715dc23c3eb1", "ccbc72b222419e2a4e40b4027f3bcb356142651b").Series(2258, CommonRequest{})
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	for _, iss := range r.Data.Results {
		t.Logf(iss.Modified.Parse().String())
		t.Logf(iss.Thumbnail.URL(PortraitIncredible))
	}
}
