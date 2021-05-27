# Progress [![PkgGoDev](https://pkg.go.dev/badge/github.com/aidenesco/progress)](https://pkg.go.dev/github.com/aidenesco/progress) [![Go Report Card](https://goreportcard.com/badge/github.com/aidenesco/progress)](https://goreportcard.com/report/github.com/aidenesco/progress)
This package uses a double sliding window to calculate item processing speed.

Due to the package's double sliding window design, fluctuations are smoothed out for a more accurate measurement of average speed. This design also means that the average isnt at its most stable until the window is full of data.

This code was inspired and derived from [average](https://github.com/prep/average) and [slidingwindow](https://github.com/bt/slidingwindow)

## Installation
```sh
go get -u github.com/aidenesco/progress
```

## Usage

```go
import "github.com/aidenesco/progress"

func main() {
    w := progress.NewWindow(time.Hour, time.Minute)
    defer w.End()
    
    for i := 0; i < 10; i++ {
        w.ItemCompleted()
    }
    
    fmt.Println(w.Average())
}
```
