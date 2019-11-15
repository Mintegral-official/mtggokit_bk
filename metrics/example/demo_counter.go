package main

import (
    _ "fmt"
    "time"
    "github.com/Mintegral-official/mtggokit/metrics"
)

func main() {
    count()
}

func count() {
    labels := []string{"method", "code", "msg"}
    counter := metrics.NewCounter("./demo_counter.yaml", labels)
    counter.With("method","Get","code","200","msg","success").Add(1)
    counter.With("method","Post","code","200","msg","success").Add(4)
    time.Sleep(6*time.Second)
    counter.With("method","Post","code","200","msg","success").Add(4)
    counter.With("method","Get","code","501","msg","error").Add(2)
    counter.With("method","Post","code","404","msg","error").Add(3)
    time.Sleep(6*time.Second)
    counter.With("method","Post","code","404","msg","error").Add(3)
    time.Sleep(1000*time.Second)
}

