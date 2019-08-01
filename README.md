# K8s DNS Exposer

[![Build Status](https://travis-ci.com/DataDog/k8s-dns-exposer.svg?token=2DPCP7qYuuA4XUjbZaqq&branch=master)](https://travis-ci.com/DataDog/k8s-dns-exposer)

k8s-dns-exposer is a Kubernetes controller that helps create, and keep up to date Endpoints objects for Services that point to external domain name.

## Use cases

The initial use case we built this for is for enabling monitoring on every endpoint behind an external service (when you don't want the `/metrics` query to be load-balanced, because each replica behind the service returns different values)

## Usage

Deploy the controller by running the following command from the root of this repository:

```bash
kubectl apply -f kubernetes/
```

Then, for any external service you want to expose inside your cluster, simply create a Kubernetes headless service with your service domain name as an `externalName`, and annotate it with `datadoghq.com/k8s-dns-exposer: "true"`.

## Config options

TODO

_Made with :heart: at Datadog_
