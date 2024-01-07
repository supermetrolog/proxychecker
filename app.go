package proxychecker

import (
	"io"
	"log"
	"net/http"
	"net/url"
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
	checker := NewChecker(a.config.ConnTimeout, &http.Client{})
	tasks := make(chan *Task, a.config.WorkersCount)

	producer := NewProducer(tasks, a.proxyReader)
	consumer := NewConsumer(tasks, checker)

	go func() {
		err := producer.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(a.config.WorkersCount)

	for i := 0; i < a.config.WorkersCount; i++ {
		go func() {
			consumer.Run()
			wg.Done()
		}()
	}

	wg.Wait()
}
