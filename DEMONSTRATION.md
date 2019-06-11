## 12 Factors

### I. Codebase

One codebase tracked in revision control, many deploys

This project is using GitHub to host it's code. Many deploys can be made from
this repository only using different CD tools like `Travis` or `Jenkins`

### II. Dependencies

Explicitly declare and isolate dependencies

Go modules are used to prevent dependencies from upgrading and break the whole
app, see `go.mod`

### III. Config

Store config in the environment

Viper is used to handle configuration and therefore adds the ability to handle
environment variables, config file and default values in a very easy and
comprehensive way. See `service/config.go`

### IV. Backing services

Treat backing services as attached resources

The app implements a TokenStore backed by Redis. It can connect to any
running instance of it by simply change the `Host` value of the Redis
configuration.

### V. Build, release, run

Strictly separate build and run stages

### VI. Processes

Execute the app as one or more stateless processes

Many replicas of the app can run at the same time and be balanced easily

### VII. Port binding

Export services via port binding

In addition to be a single binary that handles connections on it's own, it can
also run in a container

### VIII. Concurrency

Scale out via the process model

This app includes a main RPC server as well as an optional HTTP server. The
program is responsible for it's own multiplexing and handles concurrent requests
via go routines

### IX. Disposability

Maximize robustness with fast startup and graceful shutdown

The app can be stopped at any time without risking any data corruption thanks to
the graceful shutdown of the servers

### X. Dev/prod parity

Keep development, staging, and production as similar as possible

The code can be the same all the way.

### XI. Logs

Treat logs as event streams

All operations done by the app are all logged to `stdout` using the `logrus`
library that includes the time. This library is hookable to any event store such
as elastic or any other.

### XII. Admin processes

Run admin/management tasks as one-off processes

It is possible to use the same binary to execute custom commands that can be
easily added. For example, generating a new token in the server can be done
using the `cli`
