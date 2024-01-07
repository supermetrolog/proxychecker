## PROXY CHECKER

### Installation
```
git clone https://github.com/supermetrolog/proxychecker.git
```

### Launch
```
go run ./cmd/main.go
```

### Description

Список проксей по умолчанию находится в папке `resources` в файле `proxylist.txt`.
Вы можете изменить расположение передав флаг `-filepath=<value>`

Число воркеров по умолчанию равен `100`.
Так же его можно переопределить, передав число с флагом `-wc=<value>`

Таймаут соединения по умолчанию равен 5 секундам. Можно переопределить с помощью флага `-timeout=<value>`

Все флаги можно посмотреть так: `go run ./cmd/main.go -help`