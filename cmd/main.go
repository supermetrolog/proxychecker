package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/supermetrolog/proxychecker/httpclient"
	"github.com/supermetrolog/proxychecker/ipapi"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Task struct {
	ProxyUrl *url.URL
}

type Result struct {
	Err         error
	ExternalIp  string
	Country     string
	CountryCode string
}

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

func main() {
	workersCount := 100
	filePath := "resource/proxylist.txt"
	timeout := 5 * time.Second

	reader, err := FileReader(filePath)
	if err != nil {
		log.Fatal(err)
	}

	checker := NewChecker(timeout, &http.Client{})
	tasks := make(chan *Task, workersCount)

	go func() {
		err := RunTasksProducer(tasks, reader)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go func() {
			RunTasksConsumer(tasks, checker)
			wg.Done()
		}()
	}

	wg.Wait()
}

func FileReader(filePath string) (io.Reader, error) {
	path, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("Create abs filepath error: %v", err)
	}

	return os.OpenFile(path, os.O_RDONLY, 0744)
}

func RunTasksConsumer(tasks chan *Task, checker *Checker) {
	for task := range tasks {
		result := checker.Check(task.ProxyUrl)
		if result.Err != nil {
			//log.Printf("Fail: \"%s\". Err: %v", task.ProxyUrl.String(), result.Err)
			continue
		}

		log.Printf("OK: \"%s\". IP: %s, CountryCode: %s, Country: %s", task.ProxyUrl.String(), result.ExternalIp, result.CountryCode, result.Country)
	}
}

func RunTasksProducer(tasks chan *Task, reader io.Reader) error {
	buffer := bufio.NewReader(reader)

	for {
		line, _, err := buffer.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read line error: %v", err)
		}

		u, err := url.Parse(string(line))
		if err != nil {
			log.Printf("WARN. Parse proxy url error: %v", err)
			continue
		}

		tasks <- &Task{
			ProxyUrl: u,
		}
	}

	close(tasks)
	return nil
}
