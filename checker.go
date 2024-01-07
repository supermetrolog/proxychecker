package proxychecker

import (
	"errors"
	"fmt"
	"github.com/supermetrolog/proxychecker/httpclient"
	"github.com/supermetrolog/proxychecker/ipapi"
	"net/http"
	"net/url"
	"time"
)

type Checker struct {
	timeout    time.Duration
	httpClient *http.Client
}

func NewChecker(timeout time.Duration, httpClient *http.Client) *Checker {
	return &Checker{
		timeout:    timeout,
		httpClient: httpClient,
	}
}

func (checker *Checker) Check(proxy *url.URL) *Result {
	res := &Result{}

	c := httpclient.New(checker.httpClient, proxy, checker.timeout)

	ipApiRes, err := ipapi.Do(c)

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
