package proxychecker

import "log"

type Consumer struct {
	tasksCh chan *Task
	checker *Checker
}

func NewConsumer(tasksCh chan *Task, checker *Checker) *Consumer {
	return &Consumer{
		tasksCh: tasksCh,
		checker: checker,
	}
}

func (c *Consumer) Run() {
	for task := range c.tasksCh {
		result := c.checker.Check(task.ProxyUrl)
		if result.Err != nil {
			//log.Printf("Fail: \"%s\". Err: %v", task.ProxyUrl.String(), result.Err)
			continue
		}

		log.Printf("OK: \"%s\". IP: %s, CountryCode: %s, Country: %s", task.ProxyUrl.String(), result.ExternalIp, result.CountryCode, result.Country)
	}
}
