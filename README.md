# Will it serve?

This microservice will confirm it can actually connect to--and do something with--
all services connected to it. More than just making a connection, it will try
to do something, like run a SQL query, to make sure that the firewall is not
blocking communications.

## Using

You must run this on Cloud Foundry (what's the point otherwise?). Simply bind the
services you wish to test, then hit the root endpoint.

## Building

Either run `go build` on an Ubuntu 14.04 machine and `cf push` with the
`binary_buildpack`, or build `Godeps` and push with the `go_buildpack`
