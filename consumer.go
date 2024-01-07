package proxychecker

// Consumer получает задания из канала, обрабатывает их и кладет результат в результирующий канал
type Consumer struct {
	tasksCh   chan *Task
	resultsCh chan *Result
	checker   *Checker
}

func NewConsumer(tasksCh chan *Task, resultsCh chan *Result, checker *Checker) *Consumer {
	return &Consumer{
		tasksCh:   tasksCh,
		resultsCh: resultsCh,
		checker:   checker,
	}
}

func (c *Consumer) Run() {
	for task := range c.tasksCh {
		c.resultsCh <- c.checker.Check(task.ProxyUrl)
	}
}
