package ipapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	HttpUrl       = "http://ip-api.com/json"
	StatusSuccess = "success"
)

type Response struct {
	Status        string  `json:"status"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Offset        int     `json:"offset"`
	Currency      string  `json:"currency"`
	Isp           string  `json:"isp"`
	Org           string  `json:"org"`
	As            string  `json:"as"`
	Asname        string  `json:"asname"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
	Query         string  `json:"query"`
}

func (r *Response) IsOK() bool {
	return r.Status == StatusSuccess
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Url        string
	HttpClient HttpClient
}

func NewDefaultClient() *Client {
	return &Client{Url: HttpUrl, HttpClient: &http.Client{}}
}

func (c *Client) Do() (ipApiRes *Response, err error) {
	req, err := http.NewRequest("GET", c.Url, nil)

	if err != nil {
		return nil, fmt.Errorf("create new request error: %v", err)
	}

	res, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http do error: %v", err)
	}

	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if err == nil && errClose != nil {
			err = fmt.Errorf("close response reader error: %v", errClose)
		}
	}(res.Body)

	err = json.NewDecoder(res.Body).Decode(&ipApiRes)
	if err != nil {
		return nil, fmt.Errorf("json decode error: %v", err)
	}

	return ipApiRes, nil
}
