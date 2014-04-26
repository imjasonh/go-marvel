package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "http://gateway.marvel.com/v1/public"

type Client struct {
	public, private string
}

func NewClient(public, private string) Client {
	return Client{public, private}
}

// See http://developer.marvel.com/documentation/authorization
func (c Client) hash() (int64, string) {
	ts := time.Now().Unix()
	hash := md5.New()
	io.WriteString(hash, fmt.Sprintf("%d%s%s", ts, c.private, c.public))
	return ts, fmt.Sprintf("%x", hash.Sum(nil))
}

func (c Client) baseURL() url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   "gateway.marvel.com",
		Path:   "/v1/public/",
	}
	ts, hash := c.hash()
	u.RawQuery = url.Values(map[string][]string{
		"ts":     []string{fmt.Sprintf("%d", ts)},
		"apikey": []string{c.public},
		"hash":   []string{hash},
	}).Encode()
	return u
}

type commonResponse struct {
	Code            int    `json:"code"`
	ETag            string `json:"etag"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
}

type commonList struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
	Count  int `json:"count"`
}

func (c Client) Series(id int64) (r struct {
	commonResponse
	Data struct {
		commonList
		Results []struct {
			ID        int `json:"id"`
			DigitalID int `json:"digitalId"`
		} `json:"results"`
	} `json:"data"`
}, err error) {
	u := c.baseURL()
	u.Path += fmt.Sprintf("series/%d/comics", id)

	resp, herr := http.Get(u.String())
	if herr != nil {
		err = herr
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	return
}
