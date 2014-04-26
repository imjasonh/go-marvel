package marvel

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

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

func (c Client) baseURL(req interface{}) url.URL {
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
	if req != nil {
		q, _ := query.Values(req)
		u.RawQuery += "&" + q.Encode()
	}
	return u
}

// Fields common to all response entities
type commonResponse struct {
	Code            int    `json:"code"`
	ETag            string `json:"etag"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
}

type CommonRequest struct {
	Offset int `url:"offset,omitempty"`
	Limit  int `url:"limit,omitempty"`
}

// Fields common to data that lists entities, with pagination
type commonList struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
	Count  int `json:"count"`
}

func (c Client) Series(id int64, req CommonRequest) (resp struct {
	commonResponse
	Data struct {
		commonList
		Results []struct {
			ID          int
			Title       string
			Description string
			ResourceURI string
			URLs        []struct {
				Type string
				URL  string
			}
			StartYear int
			EndYear   int
			Rating    string
			//Modified  Date
			Thumbnail struct {
				Path      string
				Extension string
			}
			Comics struct {
				Available     int
				Returned      int
				CollectionURI string
				Items         []struct {
					ResourceURI string
					Name        string
				}
			}
			Stories struct {
				Available     int
				Returned      int
				CollectionURI string
				Items         []struct {
					ResourceURI string
					Name        string
					Type        string
				}
			}
			Events struct {
				Available     int
				Returned      int
				CollectionURI string
				Items         []struct {
					ResourceURI string
					Name        string
					Type        string
				}
			}
			Characters struct {
				Available     int
				Returned      int
				CollectionURI string
				Items         []struct {
					ResourceURI string
					Name        string
					Type        string
				}
			}
			Creators struct {
				Available     int
				Returned      int
				CollectionURI string
				Items         []struct {
					ResourceURI string
					Name        string
					Type        string
				}
			}
			Next struct {
				ResourceURI string
				Name        string
			}
			Previous struct {
				ResourceURI string
				Name        string
			}
		}
	}
}, err error) {
	u := c.baseURL(req)
	u.Path += fmt.Sprintf("series/%d/comics", id)
	r, err := c.fetch(u)
	if err != nil {
		return
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&resp)
	return
}

func (c Client) fetch(u url.URL) (io.ReadCloser, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		slurp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error response from API: %d\n%s", resp.StatusCode, slurp)
	}
	return resp.Body, nil
}
