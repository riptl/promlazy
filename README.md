<div align="center">
  <h1>Go promlazy</h1>
  <p>
    <strong>Lazy promauto-style registration for Go Prometheus exporters</strong>
  </p>
</div>

## Summary

`promlazy` allows you to use the ergonomics of [`promauto`](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/promauto)
while being able to omit unused metrics. 

Lazily registered metrics will only show up in Prometheus exporters when they are written to.
`promauto` registered metrics on the other hand will always be exported, regardless if used or not.
`promlazy` is especially useful for larger applications where large pieces of code may be never be active depending on config.

## Example

Lazy metrics are defined via a `promlazy.Factory`.
`promlazy.New()` is a convenience method for constructing a factory against the default registerer.

```go
var (
    metrics = promlazy.New()
    myMetric1 = metrics.NewGauge(prometheus.GaugeOpts{
        Name: "my_metric_1",
    })
    myMetric2 = metrics.NewGauge(prometheus.GaugeOpts{
        Name: "my_metric_2",
    })
)
```

The metrics defined above start out without registration, meaning they won't be visible in gathered exports (in `/metrics`).

```go
myMetric1.Set(42)
```

As soon as they are written to, the factory registers all associated metrics (in this example `my_metric_1`, `my_metric_2`).

## Safety

`promlazy` is designed to delay metrics registration to when a metric is first written to.
This comes with drawbacks since registration can fail in the event of conflicting metrics:
The application can now **panic** in the `main()` phase (during a metric write),
whereas `promauto` will only panic in the `init()` phase.

Consider the following code:

```go
var (
    metrics = promlazy.New()
    myMetric1 = metrics.NewGauge(prometheus.GaugeOpts{
        Name: "my_metric",
    })
    myMetric2 = metrics.NewGauge(prometheus.GaugeOpts{
        Name: "my_metric",
    })
)
```

This code is clearly broken — Notice how `my_metric` gets registered twice.
Any write attempt will now panic, as metrics registration is delayed to first use.
Not failing fast is dangerous because it could make broken code appear seemingly healthy (i.e. not crashing) until it is actually used.

```go
func handleRequest() {
    myMetric2.Inc() // panics!
}
```

The equivalent `promauto` code would panic immediately.

```go
var (
    myMetric1 = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "my_metric",
    })
    myMetric2 = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "my_metric",
    })
)
```

To fail fast, call the `Register()` method on the factory as soon as you activate the code paths associated with your metrics.

```go
var (
	metrics = promlazy.New()
	...
)

func StartHandler() {
	metrics.Register() // panics fast
	...
}
```
