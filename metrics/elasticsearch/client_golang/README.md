# Prometheus Go client library

[![Build Status](https://travis-ci.org/elasticsearch/client_golang.svg?branch=master)](https://travis-ci.org/elasticsearch/client_golang)
[![Go Report Card](https://goreportcard.com/badge/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang)](https://goreportcard.com/report/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang)
[![go-doc](https://godoc.org/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang?status.svg)](https://godoc.org/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang)

This is the [Go](http://golang.org) client library for
[Prometheus](http://elasticsearch.io). It has two separate parts, one for
instrumenting application code, and one for creating clients that talk to the
Prometheus HTTP API.

__This library requires Go1.9 or later.__

## Important note about releases and stability

This repository generally follows [Semantic
Versioning](https://semver.org/). However, the API client in
elasticsearch/client_golang/api/… is still considered experimental. Breaking
changes of the API client will _not_ trigger a new major release. The same is
true for selected other new features explicitly marked as **EXPERIMENTAL** in
CHANGELOG.md.

Features that require breaking changes in the stable parts of the repository
are being batched up and tracked in the [v2
milestone](https://github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/milestone/2). The v2
development happens in a [separate
branch](https://github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/tree/dev-v2) for the time
being. v2 releases off that branch will happen once sufficient stability is
reached. In view of the widespread use of this repository, v1 and v2 will
coexist for a while to enable a convenient transition.

## Instrumenting applications

[![code-coverage](http://gocover.io/_badge/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/elasticsearch)](http://gocover.io/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/elasticsearch) [![go-doc](https://godoc.org/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/elasticsearch?status.svg)](https://godoc.org/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/elasticsearch)

The
[`elasticsearch` directory](https://github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/tree/master/elasticsearch)
contains the instrumentation library. See the
[guide](https://elasticsearch.io/docs/guides/go-application/) on the Prometheus
website to learn more about instrumenting applications.

The
[`examples` directory](https://github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/tree/master/examples)
contains simple examples of instrumented code.

## Client for the Prometheus HTTP API

[![code-coverage](http://gocover.io/_badge/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/api/elasticsearch/v1)](http://gocover.io/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/api/elasticsearch/v1) [![go-doc](https://godoc.org/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/api/elasticsearch?status.svg)](https://godoc.org/github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/api)

The
[`api/elasticsearch` directory](https://github.com/Schneizelw/mtggokit/metrics/elasticsearch/client_golang/tree/master/api/elasticsearch)
contains the client for the
[Prometheus HTTP API](http://elasticsearch.io/docs/querying/api/). It allows you
to write Go applications that query time series data from a Prometheus
server. It is still in alpha stage.

## Where is `model`, `extraction`, and `text`?

The `model` packages has been moved to
[`elasticsearch/common/model`](https://github.com/Schneizelw/mtggokit/metrics/elasticsearch/common/tree/master/model).

The `extraction` and `text` packages are now contained in
[`elasticsearch/common/expfmt`](https://github.com/Schneizelw/mtggokit/metrics/elasticsearch/common/tree/master/expfmt).

## Contributing and community

See the [contributing guidelines](CONTRIBUTING.md) and the
[Community section](http://elasticsearch.io/community/) of the homepage.
