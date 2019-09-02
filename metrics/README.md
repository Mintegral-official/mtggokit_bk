# metrics

metrics为监控服务提供了统一的调用接口,主要包括counter,gauge,summary,histograms.而且为一些流行的metrics服务提供了适配器.

## usage

```golang 
//path: /project/conf/config
es : true //开关设定
log : false
prometheus: true

esConfig:
    host: xxxx
    port: xxxx
    docId: xxxx
    docType: xxxx
    interval: 10s
    lables:
        httpCode
        httpMethod

logConfig:
    logFile: xxxx
    lables:
        httpCode

prometheusConfig:
    Namespace: xxx
    Subsystemp: xxxx
    Help: xxxx
    Name: xxxx
    Lables:
        httpCode
        httpMethod
```

```golang 
//use
func main() {
    var mulitCount multi.Counter
    multiCount = multi.NewCounter("/project/conf/config")
    multiCount.Add(1)
}

```

## 设计图

### Counter
![counter](img/Counter.png)
### Gauge
![gauge](img/Gauge.png)
### Summary
![summary](img/Summary.png)



