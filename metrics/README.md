# metrics

metrics为监控服务提供了统一的调用接口,主要包括counter,gauge,summary.而且为一些流行的metrics服务提供了适配器.

## usage

```golang 
//path: /project/conf/config.yaml
open:
    es: true //开关设定
    log: false
    prometheus: true

monitorSystem:
    default:
        Namespace: "counter_test"
        Subsystemp: "test1"
        Help: "http request"
        Name: "request"
    es:
        Host: "xxxx"
        Port: "xxxx"
        Index: "metric_test"
        Type: "test"
        Interval: 10
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
	lables := []string{"httpCode", "httpMethod"}
    multiCount = multi.NewCounter("/project/conf/config.yaml", lables)
    multiCount.with({"httpCode":"200","httpMethod":"POST"}).Add(1)
    multiCount.with({"httpCode":"200","httpMethod":"POST"}).Add(2)
    multiCount.with({"httpCode":"200","httpMethod":"POST"}).Add(3)
}
```

## 总体框架

![counter](img/总体框架.png)

## 设计图

### Counter
![counter](img/Counter.png)
### Gauge
![gauge](img/Gauge.png)
### Summary
![summary](img/Summary.png)


