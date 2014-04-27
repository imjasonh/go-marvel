package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	marvel "github.com/ImJasonH/go-marvel"
)

var (
	seriesID = flag.Int64("series", 2258, "Series ID (default: Uncanny X-Men)")
	apiKey = flag.String("pub", "", "Public API key")
	secret = flag.String("priv", "", "Private API secret")
)

func main() {
	flag.Parse()

	c := marvel.NewClient(*apiKey, *secret)

	offset := 0
	limit := 100
	for {
		r, err := c.Series(*seriesID, marvel.CommonParams{offset, limit})
		if err != nil {
			panic(err)
		}
		for _, iss := range r.Data.Results {
			fetchImage(iss.IssueNumber, iss.Thumbnail.URL(marvel.PortraitIncredible))
			fmt.Printf("%d - %s\n", iss.IssueNumber, iss.Thumbnail.URL(marvel.PortraitIncredible))
		}
		if len(r.Data.Results) < limit {
			return
		}
		offset += limit
	}
}

func fetchImage(num int, url string) {
	f, err := os.Create(fmt.Sprintf("%d.jpg", num))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error: %s -> %d\n", url, resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}
}
