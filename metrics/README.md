# metrics

metrics为监控服务提供了统一的调用接口,主要包括counter,gauge,summary.而且为一些流行的metrics服务提供了适配器.

## usage

install:
```
go get -v -u github.com/Mintegral-official/mtggokit/metrics/multi
```


```golang 
#path: /project/conf/config.yaml
Open:
    Log: true
    Prometheus: false
    Elasticsearch: false
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
        Interval: "10"
Metrics:
    Summary:
        Quantile50: 5 #配合对应分位数的误差
        Quantile90: 2 
        Quantile99: 1 

```

usage:
```golang 
import (
    "time"
    "github.com/mtggokit/metrics/multi"
)
func main() {
    var logger *log.Logger
    var multiCount multi.Counter
    lables := []string{"httpCode", "httpMethod"}
    multiCount = multi.NewCounter("/project/conf/config.yaml", lables)
    multiCount.With({"httpCode":"200","httpMethod":"POST"}).Add(1)
    multiCount.With({"httpCode":"200","httpMethod":"GET"}).Add(2)
    multiCount.With({"httpCode":"200","httpMethod":"POST"}).Add(3)
    time.Sleep(1000*time.Second)
}
```

