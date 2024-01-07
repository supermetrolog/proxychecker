package proxychecker

import "log"

// Processor получает результаты и обрабатывает их (выводит в консоль)
type Processor struct {
	resultsCh chan *Result
}

func NewProcessor(resultsCh chan *Result) *Processor {
	return &Processor{resultsCh: resultsCh}
}

func (p *Processor) Run() {
	for result := range p.resultsCh {
		go p.process(result)
	}
}

func (p *Processor) process(result *Result) {
	if result.Err != nil {
		log.Printf("Fail: \"%s\". Err: %v", result.Proxy.String(), result.Err)
	} else {
		log.Printf("OK: \"%s\". IP: %s, CountryCode: %s, Country: %s", result.Proxy.String(), result.ExternalIp, result.CountryCode, result.Country)
	}
}
