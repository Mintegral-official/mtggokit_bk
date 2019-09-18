package main

import (
	_ "fmt"
	"time"
    "./metrics/multi"
)

func main() {
    //lables := []string{"country", "province", "type"}
	//counters := multi.NewCounter("./config.yaml", lables)
	//counters.With("country","cn","province","guangdong","type","house").Add(1)
	//counters.With("country","cn","province","guangdong","type","house").Add(1.5)
	//time.Sleep(1000*time.Second)
    lables := []string{"type", "channel", "country", "campaign", "publisher"}
	gauges := multi.NewGauge("./config.yaml", lables)
	gauges.With("type","gauge","channel","clickadu","country","us","campaign","us|pop|xxx","publisher","1111").Set(100)
	time.Sleep(12*time.Second)
	gauges.With("type","gauge","channel","clickadu","country","us","campaign","us|pop|xxx","publisher","1111").Add(100)
	time.Sleep(13*time.Second)
	gauges.With("type","gauge","channel","clickadu","country","cn","campaign","cn|best|xxx","publisher","1111").Set(200)
	time.Sleep(13*time.Second)
	gauges.With("type","gauge","channel","clickadu","country","cn","campaign","cn|best|xxx","publisher","1111").Set(300)
	time.Sleep(13*time.Second)
	gauges.With("type","gauge","channel","popads","country","cn","campaign","cn|best2|xxx","publisher","2222").Set(10000)
	time.Sleep(13*time.Second)
	gauges.With("type","gauge","channel","popads","country","cn","campaign","cn|best2|xxx","publisher","1111").Set(100)
	time.Sleep(1000*time.Second)

}
