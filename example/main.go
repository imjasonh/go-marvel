package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"net/http"
	"os"

	marvel "github.com/ImJasonH/go-marvel"
)

var (
	seriesID = flag.Int64("series", 2258, "Series ID (default: Uncanny X-Men)")
	apiKey   = flag.String("pub", "", "Public API key")
	secret   = flag.String("priv", "", "Private API secret")
)

func main() {
	flag.Parse()

	c := marvel.NewClient(*apiKey, *secret)

	offset := 0
	limit := 100
	imgs := []image.Image{}
	for {
		r, err := c.Series(*seriesID, marvel.CommonParams{offset, limit})
		if err != nil {
			panic(err)
		}
		for _, iss := range r.Data.Results {
			img, err := fetchImage(iss.Thumbnail.URL(marvel.PortraitIncredible))
			if err != nil {
				fmt.Printf("error: %v", err)
				return
			}
			imgs = append(imgs, img)
			fmt.Printf("fetched %v - %s\n", *iss.IssueNumber, iss.Thumbnail.URL(marvel.PortraitIncredible))
		}
		if len(r.Data.Results) < limit {
			break
		}
		offset += limit
	}

	if err := writeGIF(fmt.Sprintf("%d.gif", *seriesID), imgs); err != nil {
		fmt.Printf("error: %v", err)
	}
}

func fetchImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s -> %d\n", url, resp.StatusCode)
	}
	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	return img, err
}

func writeGIF(filename string, imgs []image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	pimgs := []*image.Paletted{}
	// TODO: convert imgs into pimgs

	return gif.EncodeAll(f, &gif.GIF{Image:pimgs})
}
