package proxychecker

import (
	"github.com/supermetrolog/proxychecker/pkg/ipapi"
	"io"
	"log"
	"net/url"
	"sync"
	"time"
)

type Task struct {
	ProxyUrl *url.URL
}

type Result struct {
	Proxy       *url.URL
	Err         error
	ExternalIp  string
	Country     string
	CountryCode string
}

type Config struct {
	ConnTimeout  time.Duration
	WorkersCount int
}

type App struct {
	config      Config
	proxyReader io.Reader
}

func NewApp(config Config, proxyReader io.Reader) *App {
	return &App{
		config:      config,
		proxyReader: proxyReader,
	}
}

func (a *App) Run() {
	checker := NewChecker(a.config.ConnTimeout, ipapi.NewDefaultClient())
	tasks := make(chan *Task, a.config.WorkersCount)
	results := make(chan *Result, a.config.WorkersCount)

	producer := NewProducer(tasks, a.proxyReader)
	consumer := NewConsumer(tasks, results, checker)
	processor := NewProcessor(results)

	go func() {
		defer close(tasks)
		err := producer.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer close(results)
		wg := sync.WaitGroup{}
		wg.Add(a.config.WorkersCount)

		for i := 0; i < a.config.WorkersCount; i++ {
			go func() {
				defer wg.Done()
				consumer.Run()
			}()
		}

		wg.Wait()
	}()

	processor.Run()
}
