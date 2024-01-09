package main

import (
	"flag"
	"github.com/supermetrolog/proxychecker"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"fmt"
)

func main() {
	filePath := flag.String("filepath", "resources/proxylist.txt", "proxy list source file path")
	workersCount := flag.Int("wc", 100, "workers count")
	timeout := flag.Duration("timeout", 5, "connection timeout")
	flag.Parse()

	config := proxychecker.Config{
		ConnTimeout:  *timeout * time.Second,
		WorkersCount: *workersCount,
	}

	reader, err := FileReader(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	app := proxychecker.NewApp(config, reader)
	app.Run()
}

func FileReader(filePath string) (io.Reader, error) {
	path, err := filepath.Abs(filePath)
	if err != nil {
	 	return nil, fmt.Errorf("create abs filepath error: %v", err)
	}

	return os.OpenFile(path, os.O_RDONLY, 0744)
}
