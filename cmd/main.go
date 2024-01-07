package main

import (
	"flag"
	"github.com/supermetrolog/proxychecker"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	workersCount               = 100
	filePath                   = "resource/proxylist.txt"
	timeout      time.Duration = 5
)

func main() {
	flag.String("filepath", filePath, "proxy list source file path")
	flag.Int("wc", workersCount, "workers count")
	flag.Duration("timeout", timeout, "connection timeout")
	flag.Parse()

	config := proxychecker.Config{
		ConnTimeout:  timeout * time.Second,
		WorkersCount: workersCount,
	}

	reader, err := FileReader(filePath)
	if err != nil {
		log.Fatal(err)
	}

	app := proxychecker.NewApp(config, reader)
	app.Run()
}

func FileReader(filePath string) (io.Reader, error) {
	path, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("Create abs filepath error: %v", err)
	}

	return os.OpenFile(path, os.O_RDONLY, 0744)
}
