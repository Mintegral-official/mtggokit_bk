package main

import (
	_ "fmt"
	"time"
    "github.com/mtggokit/metrics/multi"
)

func main() {
    //lables := []string{"country", "province", "type"}
	//counters := multi.NewCounter("./config.yaml", lables)
	//counters.With("country","cn","province","guangdong","type","house").Add(1)
	//counters.With("country","cn","province","guangdong","type","house").Add(1.5)
	//time.Sleep(1000*time.Second)
    //lables := []string{"type", "channel", "country", "campaign", "publisher"}
	//gauges := multi.NewGauge("./config.yaml", lables)
	//gauges.With("type","gauge","channel","clickadu","country","us","campaign","us|pop|xxx","publisher","1111").Set(100)
	//time.Sleep(12*time.Second)
	//gauges.With("type","gauge","channel","clickadu","country","us","campaign","us|pop|xxx","publisher","1111").Add(100)
	//time.Sleep(13*time.Second)
	//gauges.With("type","gauge","channel","clickadu","country","cn","campaign","cn|best|xxx","publisher","1111").Set(200)
	//time.Sleep(13*time.Second)
	//gauges.With("type","gauge","channel","clickadu","country","cn","campaign","cn|best|xxx","publisher","1111").Set(300)
	//time.Sleep(13*time.Second)
	//gauges.With("type","gauge","channel","popads","country","cn","campaign","cn|best2|xxx","publisher","2222").Set(10000)
	//time.Sleep(13*time.Second)
	//gauges.With("type","gauge","channel","popads","country","cn","campaign","cn|best2|xxx","publisher","1111").Set(100)
	//time.Sleep(1000*time.Second)
    lables := []string{"score"}
	summaries := multi.NewSummary("./config.yaml", lables)
	summaries.With("score","math").Observe(93)
	summaries.With("score","math").Observe(93)
	summaries.With("score","math").Observe(99)
	summaries.With("score","math").Observe(100)
	summaries.With("score","math").Observe(100)
	summaries.With("score","math").Observe(30)
	summaries.With("score","math").Observe(30)
	summaries.With("score","math").Observe(30)
	summaries.With("score","math").Observe(30)
	summaries.With("score","math").Observe(30)
	summaries.With("score","math").Observe(40)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(50)
	summaries.With("score","math").Observe(51)
	summaries.With("score","math").Observe(60)
	summaries.With("score","math").Observe(90)
	summaries.With("score","math").Observe(91)
	summaries.With("score","math").Observe(92)
	summaries.With("score","math").Observe(93)
	summaries.With("score","math").Observe(100)
	summaries.With("score","math").Observe(98)
	time.Sleep(1200*time.Second)

}

