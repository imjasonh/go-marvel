// TODO: For structs that define a ResourceURI, add a method to fetch those
//       contents and parse them into the correct response struct.
//       e.g., Series(123).Data.Results[0].Characters.Items[0].Get()...
// TODO: Add a test to fetch a resource, serialize it into JSON and compare
//       it against the response JSON to catch missing fields
// TODO: Figure out the correct incantation to let Series not have to take an
//       empty struct and instead take nil
// TODO: Find/write Swagger Go client generator?

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

func (c Client) baseURL(path string, params interface{}) url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   "gateway.marvel.com",
		Path:   "/v1/public/" + path,
	}
	if params != nil {
		q, _ := query.Values(params)
		u.RawQuery += "&" + q.Encode()
	}
	return u
}

// Fields common to all response entities
type commonResponse struct {
	Code            int
	ETag            string
	Status          string
	Copyright       string
	AttributionText string
	AttributionHTML string
}

type CommonParams struct {
	Offset int `url:"offset,omitempty"`
	Limit  int `url:"limit,omitempty"`
}

// Fields common to data that lists entities, with pagination
type commonList struct {
	Offset int
	Limit  int
	Total  int
	Count  int
}

type Image struct {
	Path      string
	Extension string
}

type Variant string

var (
	PortraitSmall       = Variant("portrait_small")
	PortraitMedium      = Variant("portrait_medium")
	PortraitXLarge      = Variant("portrait_xlarge")
	PortraitFantastic   = Variant("portrait_fantastic")
	PortraitUncanny     = Variant("portrait_uncanny")
	PortraitIncredible  = Variant("portrait_incredible")
	StandardSmall       = Variant("standard_small")
	StandardMedium      = Variant("standard_medium")
	StandardXLarge      = Variant("standard_xlarge")
	StandardFantastic   = Variant("standard_fantastic")
	StandardUncanny     = Variant("standard_uncanny")
	StandardIncredible  = Variant("standard_incredible")
	LandscapeSmall      = Variant("landscape_small")
	LandscapeMedium     = Variant("landscape_medium")
	LandscapeXLarge     = Variant("landscape_xlarge")
	LandscapeFantastic  = Variant("landscape_fantastic")
	LandscapeUncanny    = Variant("landscape_uncanny")
	LandscapeIncredible = Variant("landscape_incredible")
)

func (i Image) URL(v Variant) string {
	return fmt.Sprintf("%s/%s.%s", i.Path, string(v), i.Extension)
}

type Date string

const dateLayout = "2006-01-02T15:04:05-0700"

func (d Date) Parse() time.Time {
	t, err := time.Parse(dateLayout, string(d))
	if err != nil {
		panic(err)
	}
	return t
}

func (c Client) Series(id int64, params CommonParams) (resp struct {
	commonResponse
	Data struct {
		commonList
		Results []Comic
	}
}, err error) {
	u := c.baseURL(fmt.Sprintf("series/%d/comics", id), params)
	r, err := c.fetch(u)
	if err != nil {
		return
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&resp)
	return
}

func (c Client) fetch(u url.URL) (io.ReadCloser, error) {
	if u.RawQuery != "" {
		u.RawQuery += "&"
	}
	ts, hash := c.hash()
	u.RawQuery += url.Values(map[string][]string{
		"ts":     []string{fmt.Sprintf("%d", ts)},
		"apikey": []string{c.public},
		"hash":   []string{hash},
	}).Encode()

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

type resourceList struct {
	Available     int
	Returned      int
	CollectionURI string
}

type Character struct {
	ResourceURI string
	ID          *int         `json:"id,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Description *string      `json:"description,omitempty"`
	Modified    *Date        `json:"modified,omitempty"`
	URLs        []URL        `json:"urls,omitempty"`
	Thumbnail   *Image       `json:"thumbnail,omitempty"`
	Comics      *ComicsList  `json:"comics,omitempty"`
	Stories     *StoriesList `json:"stories,omitempty"`
	Events      *EventsList  `json:"events,omitempty"`
	Series      *SeriesList  `json:"series,omitempty"`
}

type CharactersList struct {
	resourceList
	Items []Character
}

type Comic struct {
	ResourceURI        string
	ID                 *int            `json:"id,omitempty"`
	Name               *string         `json:"id,omitempty"`
	DigitalID          *int            `json:"digitalId,omitempty"`
	Title              *string         `json:"title,omitempty"`
	IssueNumber        *string         `json:"issueNumber,omitempty"`
	VariantDescription *string         `json:"variantDescription,omitEmpty"`
	Description        *string         `json:"description,omitempty"`
	Modified           *Date           `json:"modified,omitempty"`
	ISBN               *string         `json:"isbn,omitempty"`
	UPC                *string         `json:"upc,omitempty"`
	DiamondCode        *string         `json:"diamondCode,omitempty"`
	EAN                *string         `json:"ean,omitempty"`
	ISSN               *string         `json:"issn,omitempty"`
	Format             *string         `json:"format,omitempty"`
	PageCount          *int            `json:"pageCount,omitEmpty"`
	TextObjects        []TextObject    `json:"textObjects,omitempty"`
	URLs               []URL           `json:"urls,omitempty"`
	Series             *Series         `json:"series,omitempty"`
	Variants           []Comic         `json:"variants,omitempty"`
	Collections        []Comic         `json:"collections,omitempty"`
	CollectedIssues    []Comic         `json:"collectedIssues,omitempty"`
	Dates              []ComicDate     `json:"dates,omitempty"`
	Prices             []ComicPrice    `json:"prices,omitempty"`
	Thumbnail          *Image          `json:"thumbnail,omitempty"`
	Images             []Image         `json:"images,omitempty"`
	Creators           *CreatorsList   `json:"creators,omitempty"`
	Characters         *CharactersList `json:"characters,omitempty"`
	Stories            *StoriesList    `json:"stories,omitempty"`
	Events             *EventsList     `json:"events,omitempty"`
}

type TextObject struct {
	Type     string
	Language string
	Text     string
}

type ComicDate struct {
	Type string
	Date Date
}

type ComicPrice struct {
	Type  string
	Price float64
}

type ComicsList struct {
	resourceList
	Items []Comic
}

type Story struct {
	ResourceURI   string
	ID            *int            `json:"id,omitempty"`
	Name          *string         `json:"name,omitempty"`
	Title         *string         `json:"title,omitempty"`
	Description   *string         `json:"description,omitempty"`
	Type          *string         `json:"type,omitempty"`
	Modified      *Date           `json:"date,omitempty"`
	Thumbnail     *Image          `json:"image,omitempty"`
	Comics        *ComicsList     `json:"comics,omitempty"`
	Series        *SeriesList     `json:"series,omitempty"`
	Events        *EventsList     `json:"events,omitempty"`
	Characters    *CharactersList `json:"characters,omitempty"`
	Creators      *CreatorsList   `json:"creators,omitempty"`
	OriginalIssue Comic
}

type StoriesList struct {
	resourceList
	Items []Story
}

type Event struct {
	ResourceURI string
	ID          *int            `json:"id,omitempty"`
	Title       *string         `json:"title,omitempty"`
	Description *string         `json:"description,omitempty"`
	URLs        []URL           `json:"urls,omitempty"`
	Modified    *Date           `json:"modified,omitempty"`
	Start       *Date           `json:"start,omitempty"`
	End         *Date           `json:"end,omitempty"`
	Thumbnail   *Image          `json:"thumbnail,omitempty"`
	Comics      *ComicsList     `json:"comics,omitempty"`
	Stories     *StoriesList    `json:"stories,omitempty"`
	Series      *SeriesList     `json:"series,omitempty"`
	Characters  *CharactersList `json:"characters,omitempty"`
	Creators    *CreatorsList   `json:"creators,omitempty"`
	Next        *Event          `json:"next,omitempty"`
	Previous    *Event          `json:"next,omitempty"`
}

type EventsList struct {
	resourceList
	Items []Event
}

type Series struct {
	ResourceURI string
	ID          *int            `json:"id,omitempty"`
	Name        *string         `json:"name,omitempty"`
	Title       *string         `json:"title,omitempty"`
	Description *string         `json:"description,omitempty"`
	URLs        []URL           `json:"urls,omitempty"`
	StartYear   *int            `json:"startYear,omitempty"`
	EndYear     *int            `json:"endYear,omitempty"`
	Rating      *string         `json:"rating,omitempty"`
	Modified    *Date           `json:"modified,omitempty"`
	Thumbnail   *Image          `json:"thumbnail,omitempty"`
	Comics      *ComicsList     `json:"comics,omitempty"`
	Stories     *StoriesList    `json:"stories,omitempty"`
	Events      *EventsList     `json:"events,omitempty"`
	Characters  *CharactersList `json:"characters,omitempty"`
	Creators    *CreatorsList   `json:"creators,omitempty"`
	Next        *Series         `json:"next,omitempty"`
	Previous    *Series         `json:"next,omitempty"`
}

type SeriesList struct {
	resourceList
	Items []Series
}

type Creator struct {
	ResourceURI string
	ID          *int         `json:"id,omitempty"`
	Name        *string      `json:"name,omitempty"`
	FirstName   *string      `json:"firstName,omitempty"`
	MiddleName  *string      `json:"middleName,omitempty"`
	LastName    *string      `json:"lastName,omitempty"`
	Suffix      *string      `json:"suffix,omitempty"`
	FullName    *string      `json:"fullName,omitempty"`
	Modified    *Date        `json:"modified,omitempty"`
	URLs        []URL        `json:"urls,omitempty"`
	Thumbnail   *Image       `json:"thumbnail,omitempty"`
	Series      *SeriesList  `json:"series,omitempty"`
	Stories     *StoriesList `json:"stories,omitempty"`
	Comics      *ComicsList  `json:"comics,omitempty"`
	Events      *EventsList  `json:"events,omitempty"`
}

type CreatorsList struct {
	resourceList
	Items []Creator
}

// TODO: Replace with subtypes that know what their response will be, with a Fetch() method to pull down and deserialize correctly.
type URL struct {
	Type, URL string
}
