# metrics

metrics为监控服务提供了统一的调用接口,主要包括counter,gauge,summary.而且为一些流行的metrics服务提供了适配器.

## usage

install:
```
go get -v -u github.com/Mintegral-official/mtggokit/metrics/metrics
```


counter usage:

```golang 
#path: /project/conf/counter.yaml
Open:
    Log: true
    Prometheus: false
    Elasticsearch: true
MonitorSystem:
    Default:
        Namespace: "Test"
        Subsystem: "testCount"
        Help: "just a test"
        Name: "test"
    Elasticsearch:
        Host: "xxxxx.com"
        Port: "8000"
        Index: "metric"
        Type: "metric_test"
        Interval: "10" #时间间隔
    Log:
        LogPath : "./.metricsLog"
        Interval: "10"
```

```golang 
import (
    "time"
    "github.com/Mintegral-official/mtggokit/metrics"
)
func main() {
    var logger *log.Logger
    var counter metrics.Counter
    lables := []string{"httpCode", "httpMethod"}
    counter = multi.NewCounter("/project/conf/counter.yaml", lables)
    counter.With({"httpCode":"200","httpMethod":"POST"}).Add(1)
    counter.With({"httpCode":"200","httpMethod":"GET"}).Add(2)
    counter.With({"httpCode":"200","httpMethod":"POST"}).Add(3)
    time.Sleep(1000*time.Second)
}
```


summary usage:

```golang
#path: /project/conf/summary.yaml
Open:
    Log: true
    Prometheus: false
    Elasticsearch: false
MonitorSystem:
    Default:
        Namespace: "summary"
        Subsystem: "summary test"
        Help: "just a test"
        Name: "test"
Metrics:
    Summary:
        Quantile50: 5 #对应分位数的误差
        Quantile90: 2 
        Quantile99: 1 
```

```golang 
import (
    "time"
    "github.com/Mintegral-official/mtggokit/metrics"
)
func main() {
    lables := []string{"score"}
    summaries := metrics.NewSummary("/project/conf/summary.yaml", lables)
    summaries.With("score","math").Observe(100)
    summaries.With("score","math").Observe(98)
    summaries.With("score","math").Observe(30)
    summaries.With("score","computer").Observe(30)
    summaries.With("score","computer").Observe(300)
    summaries.With("score","computer").Observe(200)
    time.Sleep(1000*time.Second)
}
```

