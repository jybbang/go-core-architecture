[![Go build and test](https://github.com/jybbang/go-core-architecture/actions/workflows/go_build_and_test.yml/badge.svg)](https://github.com/jybbang/go-core-architecture/actions/workflows/go_build_and_test.yml)

<div align='center'>
    <h3>🧿</h3>
    <h3>Go Core Architecture</h3>
    <em>keep clean and hexagonal</em>
</div>

<br>

<p align='center'>
</p>

<p align="center">
    <table>
        <tbody>
            <td align="center">
                <img width="2000" height="0"><br>
                Status: <b>beta 🚧</b><br>
                <a href="https://github.com/jybbang/go-nexinterface">Example</a><br>
                <img width="2000" height="0">
            </td>
        </tbody>
    </table>
</p>

<br>

## Features

- 📦 [Clean Architecture](https://github.com/jasontaylordev/CleanArchitecture) - mostly inspired

- 💾 CQRS

- ⚡️ Event Sourcing

- 🔥 Middlewares
  - [logging every requests](https://github.com/uber-go/zap)
  - [validation check](https://github.com/go-playground/validator)
  - check long running requests > 500 ms
  - panic recovery
  - ...and yours

- 📜 [Open tracing](https://github.com/openzipkin-contrib/zipkin-go-opentracing)

<br>

## Adapters

#### Repository adapters
| Adapter  | Status        |
|:----------|:------------|
| [MongoDB](https://github.com/mongodb/mongo-go-driver) | beta
| [PostgreSQL](https://gorm.io/) | beta
| [MySQL](https://gorm.io/) | alpha
| [SQL Server](https://gorm.io/) | alpha
| [SQLite](https://gorm.io/) | alpha
| Oracle | scheduled

#### State adapters
| Adapter  | Status        |
|:----------|:------------|
| [redis](https://github.com/go-redis/redis) | beta
| [etcd](https://github.com/etcd-io/etcd) | beta
| [LevelDB](https://github.com/syndtr/goleveldb) | beta
| [dapr](https://github.com/dapr/go-sdk) | alpha

#### Messaing adapters
| Adapter  | Status        |
|:----------|:------------|
| [redis](https://github.com/go-redis/redis) | beta
| [dapr](https://github.com/dapr/go-sdk) | alpha
| [NATS](https://github.com/nats-io/nats.go) | alpha
| AMQP | scheduled

<br>

## Getting Started

Use go get 🧿

	go get github.com/jybbang/go-core-architecture

Then import the 🧿 package into your own code.

	import "github.com/jybbang/go-core-architecture"
    
