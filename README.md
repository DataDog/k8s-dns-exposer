# K8s DNS Exposer

[![Build Status](https://travis-ci.com/DataDog/k8s-dns-exposer.svg?token=2DPCP7qYuuA4XUjbZaqq&branch=master)](https://travis-ci.com/DataDog/k8s-dns-exposer)

k8s-dns-exposer is a Kubernetes controller that helps create, and keep up to date Endpoints objects for Services that point to external domain name.

## Use cases

Some of the use cases we built this for are:

- enabling monitoring on every endpoint behind an external service (when you don't want the `/metrics` query to be load-balanced, because each replica behind the service returns different values)

## Config options

TODO

Made with :heart: @ Datadog
