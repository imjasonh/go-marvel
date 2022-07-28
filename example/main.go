package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strings"

	marvel "github.com/imjasonh/go-marvel"
)

var (
	seriesID = flag.Int("series", 2258, "Series ID (default: Uncanny X-Men)")
	apiKey   = flag.String("pub", "", "Public API key")
	secret   = flag.String("priv", "", "Private API secret")
)

func main() {
	flag.Parse()

	if *apiKey == "" || *secret == "" {
		log.Fatal("need -pub and -priv")
	}

	c := marvel.Client{
		PublicKey:  *apiKey,
		PrivateKey: *secret,
	}

	offset := 0
	limit := 100
	imgs := []image.Image{}
	for {
		r, err := c.SingleSeries(*seriesID).Comics(marvel.ComicsParams{
			CommonParams: marvel.CommonParams{
				Offset:  offset,
				Limit:   limit,
				OrderBy: "onsaleDate",
			},
		})
		if err != nil {
			log.Fatal("Getting comics (offset=%d)", offset)
		}
		fmt.Print(".")
		for _, iss := range r.Data.Results {
			url := iss.Thumbnail.URL(marvel.PortraitIncredible)
			img := fetchImage(url)
			if img != nil {
				imgs = append(imgs, img)
				fmt.Printf("fetched %v - %s\n", *iss.IssueNumber, url)
			} else {
				fmt.Printf("skipped %s\n", url)
			}
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

func fetchImage(url string) image.Image {
	if strings.Contains(url, "image_not_available") {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	defer resp.Body.Close()

	b := bufio.NewReaderSize(resp.Body, 1)
	if _, err := b.Peek(1); err == bufio.ErrBufferFull {
		return nil
	}

	img, err := jpeg.Decode(b)
	if err != nil {
		return nil
	}
	return img
}

func writeGIF(filename string, imgs []image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	g := gif.GIF{
		Image: make([]*image.Paletted, len(imgs)),
		Delay: make([]int, len(imgs)),
	}
	b := make([]byte, 0, 1024)
	for i, img := range imgs {
		buf := bytes.NewBuffer(b)
		if err = gif.Encode(buf, img, &gif.Options{NumColors: 256}); err != nil {
			return err
		}
		gimg, err := gif.DecodeAll(buf)
		if err != nil {
			return err
		}
		g.Delay[i] = 0
		g.Image[i] = gimg.Image[0]
	}

	return gif.EncodeAll(f, &g)
}
