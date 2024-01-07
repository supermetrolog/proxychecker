package proxychecker

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
)

type Producer struct {
	tasksCh chan *Task
	reader  io.Reader
}

func NewProducer(tasksCh chan *Task, reader io.Reader) *Producer {
	return &Producer{tasksCh: tasksCh, reader: reader}
}

func (p *Producer) Run() error {
	buffer := bufio.NewReader(p.reader)

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

		p.tasksCh <- &Task{
			ProxyUrl: u,
		}
	}

	close(p.tasksCh)
	return nil
}
