package proxychecker

import (
	"errors"
	"fmt"
	"github.com/supermetrolog/proxychecker/pkg/ipapi"
	"net/http"
	"net/url"
	"time"
)

// Checker получает прокси, проверяет его и возвращает результат проверки
type Checker struct {
	timeout     time.Duration
	ipApiClient *ipapi.Client
}

func NewChecker(timeout time.Duration, ipApiClient *ipapi.Client) *Checker {
	return &Checker{timeout: timeout, ipApiClient: ipApiClient}
}

func (c *Checker) Check(proxy *url.URL) *Result {
	res := &Result{
		Proxy: proxy,
	}

	c.ipApiClient.HttpClient = &http.Client{
		Timeout: c.timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	ipApiRes, err := c.ipApiClient.Do()

	if err != nil {
		res.Err = fmt.Errorf("ipapi request error: %v", err)
		return res
	}

	if !ipApiRes.IsOK() {
		res.Err = errors.New("ip api response is not success")
		return res
	}

	res.ExternalIp = ipApiRes.Query
	res.Country = ipApiRes.Country
	res.CountryCode = ipApiRes.CountryCode

	return res
}
