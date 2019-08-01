package config

import (
	"time"
)

const (
	// DefaultRequeueDuration default requeue duration
	DefaultRequeueDuration = 5 * time.Second
	// K8sDNSExposerAnnotationKey use to select Service only managed by the controller
	K8sDNSExposerAnnotationKey = "datadoghq.com/k8s-dns-exposer"
)
