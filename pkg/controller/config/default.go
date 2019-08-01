// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

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
