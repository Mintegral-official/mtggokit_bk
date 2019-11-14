# metrics

metrics为监控服务提供了统一的调用接口,主要包括counter,gauge,summary.而且为一些流行的metrics服务提供了适配器.

## usage

install:
```
go get -v -u github.com/Mintegral-official/mtggokit/metrics/metrics
```

counter usage:

```golang 
#FILENAME: demo_counter.yaml
Open:
    Log: true
    Prometheus: false
MonitorSystem:
    Default:
        Namespace: "Demo"
        Subsystem: "prometheus"
        Name: "httpRequest"
        Help: "test"
    Log:
        LogPath: "./temp" 
        Interval: "5" 

```

```golang 
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
    counter.With("method","Post","code","200","msg","success").Add(4)
    counter.With("method","Get","code","501","msg","error").Add(2)
    counter.With("method","Post","code","404","msg","error").Add(3)
    counter.With("method","Post","code","404","msg","error").Add(3)
    time.Sleep(10*time.Second)
}

```

log输出json数据格式:
```json
{
    "Data":[
        {
            "Value":1,
            "code":"200",
            "method":"Get",
            "msg":"success"
        },
        {
            "Value":8,
            "code":"200",
            "method":"Post",
            "msg":"success"
        },
        {
            "Value":2,
            "code":"501",
            "method":"Get",
            "msg":"error"
        },
        {
            "Value":6,
            "code":"404",
            "method":"Post",
            "msg":"error"
        }
    ],
    "FqName":"Demo_prometheus_httpRequest",
    "Help":"test",
    "Timestamp":"2019-11-14T03:25:43Z"
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


