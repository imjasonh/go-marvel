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

const (
	basePath = "http://gateway.marvel.com/v1/public"
)

// Client provides methods to get Marvel Comics data.
type Client struct {
	PublicKey, PrivateKey string
	Client                *http.Client
}

func (c Client) fetch(path string, params interface{}, out interface{}) error {
	u := c.baseURL(path, params)
	if u.RawQuery != "" {
		u.RawQuery += "&"
	}
	ts, hash := c.hash()
	u.RawQuery += url.Values(map[string][]string{
		"ts":     []string{fmt.Sprintf("%d", ts)},
		"apikey": []string{c.PublicKey},
		"hash":   []string{hash},
	}).Encode()
	if c.Client == nil {
		c.Client = &http.Client{}
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		slurp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("error response from API: %d\n%s", resp.StatusCode, slurp)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c Client) baseURL(path string, params interface{}) url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   "gateway.marvel.com",
		Path:   "/v1/public" + path,
	}
	if params != nil {
		q, _ := query.Values(params)
		u.RawQuery += "&" + q.Encode()
	}
	return u
}

// See http://developer.marvel.com/documentation/authorization
func (c Client) hash() (int64, string) {
	ts := time.Now().Unix()
	hash := md5.New()
	io.WriteString(hash, fmt.Sprintf("%d%s%s", ts, c.PrivateKey, c.PublicKey))
	return ts, fmt.Sprintf("%x", hash.Sum(nil))
}

// URL represents a public web site URL for a resource.
type URL struct {
	Type *string `json:"type,omitempty"`
	URL  *string `json:"url,omitempty"`
}

// CommonParams provides fields common to all request parameter entities.
type CommonParams struct {
	OrderBy       string `url:"orderBy,omitempty"`
	Offset        int    `url:"offset,omitempty"`
	Limit         int    `url:"limit,omitempty"`
	ModifiedSince string `url:"modifiedSince,omitempty"`
}

// CommonResponse provides fields common to all response entities.
type CommonResponse struct {
	Code            *int    `json:"code,omitempty"`
	ETag            *string `json:"etag,omitempty"`
	Status          *string `json:"status,omitempty"`
	Copyright       *string `json:"copyright,omitempty"`
	AttributionText *string `json:"attributionText,omitempty"`
	AttributionHTML *string `json:"attributionHtml,omitempty"`
}

// CommonList provides fields common to data that lists entities, with pagination.
type CommonList struct {
	Offset *int `json:"offset,omitempty"`
	Limit  *int `json:"limit,omitempty"`
	Total  *int `json:"total,omitempty"`
	Count  *int `json:"count,omitempty"`
}

// ResourceList provides fields common to minimal list of entities.
// Use an entity list's List method to retrieve more complete information.
type ResourceList struct {
	Available     *int    `json:"available,omitempty"`
	Returned      *int    `json:"returned,omitempty"`
	CollectionURI *string `json:"collectionUri,omitempty"`
}

// Image provides data necessary to construct an image URL.
type Image struct {
	Path      *string `json:"path,omitempty"`
	Extension *string `json:"extension,omitempty"`
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

// URL returns a complete URL string for an Image, with the specified Variant.
func (i Image) URL(v Variant) string {
	return fmt.Sprintf("%s/%s.%s", *i.Path, string(v), *i.Extension)
}

type Date string

const dateLayout = "2006-01-02T15:04:05-0700"

// Parse returns a time.Time equivalent to the Date.
func (d Date) Parse() time.Time {
	t, err := time.Parse(dateLayout, string(d))
	if err != nil {
		panic(err)
	}
	return t
}

/////
// Characters
/////

// Character begins to construct a request for information based on a Character.
func (c Client) Character(id int) CharacterResource {
	return CharacterResource{basePath: fmt.Sprintf("/characters/%d", id), client: c}
}

// CharacterResource provides methods to issue requests for a Character.
type CharacterResource struct {
	basePath string
	client   Client
}

// Characters issues a request to search for Characters.
func (c Client) Characters(params CharactersParams) (resp *CharactersResponse, err error) {
	err = c.fetch("/characters", params, &resp)
	return
}

// Get issues a request to get a Character.
func (s CharacterResource) Get() (resp *CharactersResponse, err error) {
	err = s.client.fetch(s.basePath, nil, &resp)
	return
}

// Comics issues a request to search for Comics associated with a Character.
func (s CharacterResource) Comics(params ComicsParams) (resp *ComicsResponse, err error) {
	err = s.client.fetch(s.basePath+"/comics", params, &resp)
	return
}

// Events issues a request to search for Events associated with a Character.
func (s CharacterResource) Events(params EventsParams) (resp *EventsResponse, err error) {
	err = s.client.fetch(s.basePath+"/events", params, &resp)
	return
}

// Series issues a request to search for Series associated with a Character.
func (s CharacterResource) Series(params SeriesParams) (resp *SeriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/series", params, &resp)
	return
}

// Stories issues a request to search for Stories associated with a Character.
func (s CharacterResource) Stories(params StoriesParams) (resp *StoriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/stories", params, &resp)
	return
}

// CharactersParams represents parameters to search for Characters.
type CharactersParams struct {
	CommonParams
	Name           string `url:"name,omitempty"`
	NameStartsWith string `url:"nameStartsWith,omitempty"`
	Comics         []int  `url:"comics,omitempty,comma"`
	Events         []int  `url:"events,omitempty,comma"`
	Stories        []int  `url:"stories,omitempty,comma"`
}

// CharactersResponse represents responses to methods that return Characters.
type CharactersResponse struct {
	CommonResponse
	Data struct {
		CommonList
		Results []Character `json:"results,omitempty"`
	} `json:"data,omitempty"`
}

// Character represents a single Character.
type Character struct {
	ResourceURI *string      `json:"resourceURI,omitempty"`
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

// Get issues a request to get complete information about a Character.
func (c Character) Get(cl Client) (resp *CharactersResponse, err error) {
	err = cl.fetch((*c.ResourceURI)[len(basePath):], nil, &resp)
	return
}

// CharactersList represents a list of Characters.
type CharactersList struct {
	ResourceList
	Items []Character `json:"items,omitempty"`
}

// List issues a request to get complete information about a list of Characters.
func (l CharactersList) List(cl Client) (resp *CharactersResponse, err error) {
	err = cl.fetch((*l.CollectionURI)[len(basePath):], nil, &resp)
	return
}

/////
// Comics
/////

// Comic begins to construct a request for information based on a Comic.
func (c Client) Comic(id int) ComicResource {
	return ComicResource{basePath: fmt.Sprintf("comics/%d", id), client: c}
}

// ComicResource provides methods to issue requests for a Comic.
type ComicResource struct {
	basePath string
	client   Client
}

// Comics issues a request to search for Comics.
func (c Client) Comics(params ComicsParams) (resp *ComicsResponse, err error) {
	err = c.fetch("/comics", params, &resp)
	return
}

// Get issues a request to get a Comic.
func (s ComicResource) Get() (resp *ComicsResponse, err error) {
	err = s.client.fetch(s.basePath, nil, &resp)
	return
}

// Characters issues a request to search for Characters associated with a Comic.
func (s ComicResource) Characters(params CharactersParams) (resp *CharactersResponse, err error) {
	err = s.client.fetch(s.basePath+"/characters", params, &resp)
	return
}

// Events issues a request to search for Events associated with a Comic.
func (s ComicResource) Events(params EventsParams) (resp *EventsResponse, err error) {
	err = s.client.fetch(s.basePath+"/events", params, &resp)
	return
}

// Series issues a request to search for Series associated with a Comic.
func (s ComicResource) Series(params SeriesParams) (resp *SeriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/series", params, &resp)
	return
}

// Stories issues a request to search for Stories associated with a Comic.
func (s ComicResource) Stories(params StoriesParams) (resp *StoriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/stories", params, &resp)
	return
}

// ComicsParams represents parameters to search for Comics.
type ComicsParams struct {
	CommonParams
	Format            string `url:"format,omitempty"`
	FormatType        string `url:"formatType,omitempty"`
	NoVariants        bool   `url:"noVariants,omitempty"`
	DateDescriptor    string `url:"dateDescriptor,omitempty"`
	DateRange         string `url:"dateRange,omitempty"`
	DiamondCode       string `url:"diamondCode,omitempty"`
	DigitalID         string `url:"digitalId,omitempty"`
	UPC               string `url:"upc,omitempty"`
	ISBN              string `url:"isbn,omitempty"`
	EAN               string `url:"ean,omitempty"`
	ISSN              string `url:"issn,omitempty"`
	HasDigitalIssue   bool   `url:"hasDigitalIssue,omitempty"`
	Creators          []int  `url:"creators,omitempty,comma"`
	Characters        []int  `url:"characters,omitempty,comma"`
	Events            []int  `url:"events,omitempty,comma"`
	Stories           []int  `url:"stories,omitempty,comma"`
	SharedAppearances []int  `url:"sharedAppearances,omitempty,comma"`
	Collaborators     []int  `url:"collaborators,omitempty,comma"`
}

// ComicsResponse represents responses to methods that return Comics.
type ComicsResponse struct {
	CommonResponse
	Data struct {
		CommonList
		Results []Comic `json:"results,omitempty"`
	} `json:"data,omitempty"`
}

// Comic represents a single Comic.
type Comic struct {
	ResourceURI        *string  `json:"resourceURI,omitempty"`
	ID                 *int     `json:"id,omitempty"`
	Name               *string  `json:"id,omitempty"`
	DigitalID          *int     `json:"digitalId,omitempty"`
	Title              *string  `json:"title,omitempty"`
	IssueNumber        *float64 `json:"issueNumber,omitempty"`
	VariantDescription *string  `json:"variantDescription,omitEmpty"`
	Description        *string  `json:"description,omitempty"`
	Modified           *Date    `json:"modified,omitempty"`
	ISBN               *string  `json:"isbn,omitempty"`
	UPC                *string  `json:"upc,omitempty"`
	DiamondCode        *string  `json:"diamondCode,omitempty"`
	EAN                *string  `json:"ean,omitempty"`
	ISSN               *string  `json:"issn,omitempty"`
	Format             *string  `json:"format,omitempty"`
	PageCount          *int     `json:"pageCount,omitEmpty"`
	TextObjects        []struct {
		Type     string `json:"text,omitempty"`
		Language string `json:"language,omitempty"`
		Text     string `json:"text,omitempty"`
	} `json:"textObjects,omitempty"`
	URLs            []URL   `json:"urls,omitempty"`
	Series          *Series `json:"series,omitempty"`
	Variants        []Comic `json:"variants,omitempty"`
	Collections     []Comic `json:"collections,omitempty"`
	CollectedIssues []Comic `json:"collectedIssues,omitempty"`
	Dates           []struct {
		Type string `json:"type,omitempty"`
		Date Date   `json:"date,omitempty"`
	} `json:"dates,omitempty"`
	Prices []struct {
		Type  string  `json:"type,omitempty"`
		Price float64 `json:"price,omitempty"`
	} `json:"prices,omitempty"`
	Thumbnail  *Image          `json:"thumbnail,omitempty"`
	Images     []Image         `json:"images,omitempty"`
	Creators   *CreatorsList   `json:"creators,omitempty"`
	Characters *CharactersList `json:"characters,omitempty"`
	Stories    *StoriesList    `json:"stories,omitempty"`
	Events     *EventsList     `json:"events,omitempty"`
}

// Get issues a request to get complete information about a Comic.
func (c Comic) Get(cl Client) (resp *ComicsResponse, err error) {
	err = cl.fetch((*c.ResourceURI)[len(basePath):], nil, &resp)
	return
}

// ComicsList represents a list of Comics.
type ComicsList struct {
	ResourceList
	Items []Comic `json:"items,omitempty"`
}

// List issues a request to get complete information about a list of Comics.
func (l ComicsList) List(cl Client) (resp *ComicsResponse, err error) {
	err = cl.fetch((*l.CollectionURI)[len(basePath):], nil, &resp)
	return
}

/////
// Creators
/////

// Creator begins to construct a request for information based on a Creator.
func (c Client) Creator(id int) CreatorResource {
	return CreatorResource{basePath: fmt.Sprintf("/creators/%d", id), client: c}
}

// CreatorResource provides methods to issue requests for a Creator.
type CreatorResource struct {
	basePath string
	client   Client
}

// Creators issues a request to search for Creators.
func (c Client) Creators(params CreatorsParams) (resp *CreatorsResponse, err error) {
	err = c.fetch("/creators", params, &resp)
	return
}

// Get issues a request to get a Creator.
func (s CreatorResource) Get() (resp *CreatorsResponse, err error) {
	err = s.client.fetch(s.basePath, nil, &resp)
	return
}

// Comics issues a request to search for Comics associated with a Creator.
func (s CreatorResource) Comics(params ComicsParams) (resp *ComicsResponse, err error) {
	err = s.client.fetch(s.basePath+"/comics", params, &resp)
	return
}

// Events issues a request to search for Events associated with a Creator.
func (s CreatorResource) Events(params EventsParams) (resp *EventsResponse, err error) {
	err = s.client.fetch(s.basePath+"/events", params, &resp)
	return
}

// Series issues a request to search for Series associated with a Creator.
func (s CreatorResource) Series(params SeriesParams) (resp *SeriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/series", params, &resp)
	return
}

// Stories issues a request to search for Stories associated with a Creator.
func (s CreatorResource) Stories(params StoriesParams) (resp *StoriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/stories", params, &resp)
	return
}

// CreatorsParams represents parameters to search for Creators.
type CreatorsParams struct {
	CommonParams
	FirstName            string `url:"firstName,omitempty"`
	MiddleName           string `url:"middleName,omitempty"`
	LastName             string `url:"lastName,omitempty"`
	Suffix               string `url:"suffix,omitempty"`
	NameStartsWith       string `url:"nameStartsWith,omitempty"`
	FirstNameStartsWith  string `url:"firstNameStartsWith,omitempty"`
	MiddleNameStartsWith string `url:"middleNameStartsWith,omitempty"`
	LastNameStartsWith   string `url:"lastNameStartsWith,omitempty"`
	Comics               []int  `url:"comics,omitempty,comma"`
	Events               []int  `url:"events,omitempty,comma"`
	Stories              []int  `url:"stories,omitempty,comma"`
}

// CreatorsResponse represents responses to methods that return Creators.
type CreatorsResponse struct {
	CommonResponse
	Data struct {
		CommonList
		Results []Creator `json:"results,omitempty"`
	} `json:"data,omitempty"`
}

// Creator represents a single Creator.
type Creator struct {
	ResourceURI *string      `json:"resourceURI,omitempty"`
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

// Get issues a request to get complete information about a Creator.
func (c Creator) Get(cl Client) (resp *CreatorsResponse, err error) {
	err = cl.fetch((*c.ResourceURI)[len(basePath):], nil, &resp)
	return
}

// CreatorsList represents a list of Creators.
type CreatorsList struct {
	ResourceList
	Items []Creator
}

// List issues a request to get complete information about a list of Creators.
func (l CreatorsList) List(cl Client) (resp *CreatorsResponse, err error) {
	err = cl.fetch((*l.CollectionURI)[len(basePath):], nil, &resp)
	return
}

/////
// Events
/////

// Event begins to construct a request for information based on an Event.
func (c Client) Event(id int) EventResource {
	return EventResource{basePath: fmt.Sprintf("events/%d", id), client: c}
}

// EventResource provides methods to issue requests for an Event.
type EventResource struct {
	basePath string
	client   Client
}

// Events issues a request to search for Events.
func (c Client) Events(params EventsParams) (resp *EventsResponse, err error) {
	err = c.fetch("/events", params, &resp)
	return
}

// Get issues a request to get a Event.
func (s EventResource) Get() (resp *EventsResponse, err error) {
	err = s.client.fetch(s.basePath, nil, &resp)
	return
}

// Characters issues a request to search for Characters associated with a Character.
func (s EventResource) Characters(params CharactersParams) (resp *CharactersResponse, err error) {
	err = s.client.fetch(s.basePath+"/characters", params, &resp)
	return
}

// Comics issues a request to search for Comics associated with a Character.
func (s EventResource) Comics(params ComicsParams) (resp *ComicsResponse, err error) {
	err = s.client.fetch(s.basePath+"/comics", params, &resp)
	return
}

// Creators issues a request to search for Creators associated with a Character.
func (s EventResource) Creators(params CreatorsParams) (resp *CreatorsResponse, err error) {
	err = s.client.fetch(s.basePath+"/creators", params, &resp)
	return
}

// Series issues a request to search for Series' associated with a Character.
func (s EventResource) Series(params SeriesParams) (resp *SeriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/series", params, &resp)
	return
}

// Stories issues a request to search for Stories associated with a Character.
func (s EventResource) Stories(params StoriesParams) (resp *StoriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/stories", params, &resp)
	return
}

// EventsParams represents parameters to search for Events.
type EventsParams struct {
	CommonParams
	Name           string `url:"name,omitempty"`
	NameStartsWith string `url:"nameStartsWith,omitempty"`
	Creators       []int  `url:"creators,omitempty,comma"`
	Characters     []int  `url:"characters,omitempty,comma"`
	Comics         []int  `url:"comics,omitempty,comma"`
	Stories        []int  `url:"stories,omitempty,comma"`
}

// EventsResponse represents responses to methods that return Events.
type EventsResponse struct {
	CommonResponse
	Data struct {
		CommonList
		Results []Event `json:"results,omitempty"`
	} `json:"data,omitempty"`
}

// Event represents a single Event.
type Event struct {
	ResourceURI *string         `json:"resourceURI,omitempty"`
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

// Get issues a request to get complete information about an Event.
func (e Event) Get(cl Client) (resp *EventsResponse, err error) {
	err = cl.fetch((*e.ResourceURI)[len(basePath):], nil, &resp)
	return
}

// EventsList represents a list of Events.
type EventsList struct {
	ResourceList
	Items []Event `json:"items,omitempty"`
}

// List issues a request to get complete information about a list of Events.
func (l EventsList) List(cl Client) (resp *EventsResponse, err error) {
	err = cl.fetch((*l.CollectionURI)[len(basePath):], nil, &resp)
	return
}

/////
// Series
/////
// (This is named SingleSeries because Series (plural) will search for series')

// SingleSeries begins to construct a request for information based on a Series.
func (c Client) SingleSeries(id int) SeriesResource {
	return SeriesResource{basePath: fmt.Sprintf("/series/%d", id), client: c}
}

// SeriesResource provides methods to issue requests for a Series.
type SeriesResource struct {
	basePath string
	client   Client
}

// Series issues a request to search for Series'.
func (c Client) Series(params SeriesParams) (resp *SeriesResponse, err error) {
	err = c.fetch("/series", params, &resp)
	return
}

// Get issues a request to get a single Series.
func (s SeriesResource) Get() (resp *SeriesResponse, err error) {
	err = s.client.fetch(s.basePath, nil, &resp)
	return
}

// Characters issues a request to search for Characters associated with a Series.
func (s SeriesResource) Characters(params CharactersParams) (resp *CharactersResponse, err error) {
	err = s.client.fetch(s.basePath+"/characters", params, resp)
	return
}

// Comics issues a request to search for Comics associated with a Series.
func (s SeriesResource) Comics(params ComicsParams) (resp *ComicsResponse, err error) {
	err = s.client.fetch(s.basePath+"/comics", params, &resp)
	return
}

// Creators issues a request to search for Creators associated with a Series.
func (s SeriesResource) Creators(params CreatorsParams) (resp *CreatorsResponse, err error) {
	err = s.client.fetch(s.basePath+"/creators", params, &resp)
	return
}

// Events issues a request to search for Events associated with a Series.
func (s SeriesResource) Events(params EventsParams) (resp *EventsResponse, err error) {
	err = s.client.fetch(s.basePath+"/events", params, &resp)
	return
}

// Stories issues a request to search for Stories associated with a Series.
func (s SeriesResource) Stories(params StoriesParams) (resp *StoriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/stories", params, &resp)
	return
}

// SeriesParams represents parameters to search for Series'.
type SeriesParams struct {
	CommonParams
	Events          string `url:"events,omitempty"`
	Title           string `url:"title,omitempty"`
	TitleStartsWith string `url:"titleStartsWith,omitempty"`
	StartYear       int    `url:"startYear,omitempty"`
	SeriesType      string `url:"seriesType,omitempty"`
	Contains        string `url:"contains,omitempty"`
	Comics          []int  `url:"comics,omitempty,comma"`
	Creators        []int  `url:"creators,omitempty,comma"`
	Characters      []int  `url:"characters,omitempty,comma"`
}

// SeriesResponse represents responses to methods that return Series'.
type SeriesResponse struct {
	CommonResponse
	Data struct {
		CommonList
		Results []Series `json:"results,omitempty"`
	} `json:"data,omitempty"`
}

// Series represents a single Series.
type Series struct {
	ResourceURI *string         `json:"resourceURI,omitempty"`
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

// Get issues a request to get complete information about a Series.
func (s Series) Get(cl Client) (resp *SeriesResponse, err error) {
	err = cl.fetch((*s.ResourceURI)[len(basePath):], nil, &resp)
	return
}

// SeriesList represents a list of Series'.
type SeriesList struct {
	ResourceList
	Items []Series
}

// List issues a request to get complete information about a list of Series'.
func (l SeriesList) List(cl Client) (resp *SeriesResponse, err error) {
	err = cl.fetch((*l.CollectionURI)[len(basePath):], nil, &resp)
	return
}

/////
// Stories
/////

// Story begins to construct a request for information based on a Story.
func (c Client) Story(id int) StoryResource {
	return StoryResource{basePath: fmt.Sprintf("stories/%d", id), client: c}
}

// StoryResource provides methods to issue requests for a Story.
type StoryResource struct {
	basePath string
	client   Client
}

// Stories issues a request to search for Stories.
func (c Client) Stories(params StoriesParams) (resp *StoriesResponse, err error) {
	err = c.fetch("/stories", params, &resp)
	return
}

// Get issues a request to get a Story.
func (s StoryResource) Get() (resp *StoriesResponse, err error) {
	err = s.client.fetch(s.basePath, nil, &resp)
	return
}

// Characters issues a request to search for Characters associated with a Story.
func (s StoryResource) Characters(params CharactersParams) (resp *CharactersResponse, err error) {
	err = s.client.fetch(s.basePath+"/characters", params, &resp)
	return
}

// Comics issues a request to search for Comics associated with a Story.
func (s StoryResource) Comics(params ComicsParams) (resp *ComicsResponse, err error) {
	err = s.client.fetch(s.basePath+"/comics", params, &resp)
	return
}

// Creators issues a request to search for Creators associated with a Story.
func (s StoryResource) Creators(params CreatorsParams) (resp *CreatorsResponse, err error) {
	err = s.client.fetch(s.basePath+"/creators", params, &resp)
	return
}

// Events issues a request to search for Events associated with a Story.
func (s StoryResource) Events(params EventsParams) (resp *EventsResponse, err error) {
	err = s.client.fetch(s.basePath+"/events", params, &resp)
	return
}

// Series issues a request to search for Series associated with a Story.
func (s StoryResource) Series(params SeriesParams) (resp *SeriesResponse, err error) {
	err = s.client.fetch(s.basePath+"/series", params, &resp)
	return
}

// StoriesParams represents parameters to search for Stories.
type StoriesParams struct {
	CommonParams
	Comics     []int `url:"comics,omitempty,comma"`
	Events     []int `url:"events,omitempty,comma"`
	Creators   []int `url:"creators,omitempty,comma"`
	Characters []int `url:"characters,omitempty,comma"`
}

// StoriesResponse represents responses to methods that return Stories.
type StoriesResponse struct {
	CommonResponse
	Data struct {
		CommonList
		Results []Story `json:"results,omitempty"`
	} `json:"data,omitempty"`
}

// Story represents a single Story.
type Story struct {
	ResourceURI   *string         `json:"resourceURI,omitempty"`
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

// Get issues a request to get complete information about a Story.
func (s Story) Get(cl Client) (resp *StoriesResponse, err error) {
	err = cl.fetch((*s.ResourceURI)[len(basePath):], nil, &resp)
	return
}

// StoriesList represents a list of Stories.
type StoriesList struct {
	ResourceList
	Items []Story `json:"items,omitempty"`
}

// List issues a request to get complete information about a list of Stories.
func (l StoriesList) List(cl Client) (resp *StoriesResponse, err error) {
	err = cl.fetch((*l.CollectionURI)[len(basePath):], nil, &resp)
	return
}
