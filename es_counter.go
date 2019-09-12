package main

import (
    "./metrics/multi"
)

func main() {
    lables := []string{"country", "province", "type"}
	counters := multi.NewCounter("./config.yaml", lables)
	counters.With("country","cn","province","guangdong","type","house").Add(1)
}
