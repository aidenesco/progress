# Progress [![PkgGoDev](https://pkg.go.dev/badge/github.com/aidenesco/progress)](https://pkg.go.dev/github.com/aidenesco/progress) [![Go Report Card](https://goreportcard.com/badge/github.com/aidenesco/progress)](https://goreportcard.com/report/github.com/aidenesco/progress)
This package uses a double sliding window to help you calculate item processing speed within your application.

Because this package uses a double sliding window, fluctuations are smoothed out for a more accurate measurement of average speed.

This code was inspired and derived from [average](https://github.com/prep/average) and [slidingwindow](https://github.com/bt/slidingwindow)

## Installation
```sh
go get -u github.com/aidenesco/progress
```

## Usage

```go
import "github.com/aidenesco/progress"

func main() {
    prog := progress.NewWindow(time.Minute, time.Second)
    defer prog.End()
    
    for i := 0; i < 10; i++ {
        prog.ItemCompleted()
    }
    
    fmt.Println(prog.Average())
}
```
