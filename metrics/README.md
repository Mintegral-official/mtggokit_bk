# metrics

metrics为监控服务提供了统一的调用接口,主要包括counter,gauge,summary.而且为一些流行的metrics服务提供了适配器.

## usage

```golang 
//path: /project/conf/config
open:
    es : true //开关设定
    log : false
    prometheus: true

monitorSystem:
    es:
        Host: xxxx
        Port: xxxx
        DocId: xxxx
        DocType: xxxx
        Interval: 10s
    prometheus:
        Namespace: xxx
        Subsystemp: xxxx
        Help: xxxx
        Name: xxxx
lables:
    httpCode
    httpMethod
```

```

```golang 
//use
import (
    "log"
)
func main() {
    var logger *log.Logger
    var multiCount multi.Counter
    multiCount = multi.NewCounter("/project/conf/config", logger)
    multiCount.with("httpCode":"200","httpMethod":"POST").Add(1)
}

```

## 设计图

### Counter
![counter](img/Counter.png)
### Gauge
![gauge](img/Gauge.png)
### Summary
![summary](img/Summary.png)



