package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestUpdateEndpoints(t *testing.T) {
	for name, tc := range map[string]struct {
		svc    *corev1.Service
		eps    *corev1.Endpoints
		ips    []string
		result *corev1.Endpoints
	}{
		"nominal update case": {
			svc: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "My Service",
					UID:  types.UID("my-service-uid"),
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Service",
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Port: 80},
						{Port: 8080},
					},
				},
			},
			eps: &corev1.Endpoints{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "v1",
							Kind:       "Service",
							Name:       "My Service",
							UID:        types.UID("my-service-uid"),
						},
					},
					Name: "My Service",
					UID:  types.UID("my-endpoints-uid"),
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Endpoints",
				},
				Subsets: []corev1.EndpointSubset{
					{
						Addresses: []corev1.EndpointAddress{
							{
								IP: "10.0.0.1",
							},
							{
								IP: "10.0.0.2",
							},
						},
					},
				},
			},
			ips: []string{
				"10.0.0.3",
				"10.0.0.4",
			},
			result: &corev1.Endpoints{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "v1",
							Kind:       "Service",
							Name:       "My Service",
							UID:        types.UID("my-service-uid"),
						},
					},
					Name: "My Service",
					UID:  types.UID("my-endpoints-uid"),
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Endpoints",
				},
				Subsets: []corev1.EndpointSubset{
					{
						Addresses: []corev1.EndpointAddress{
							{
								IP: "10.0.0.3",
							},
							{
								IP: "10.0.0.4",
							},
						},
						Ports: []corev1.EndpointPort{
							{Port: 80},
							{Port: 8080},
						},
					},
				},
			},
		},
		"nominal creation case": {
			svc: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "My Service",
					Namespace: "myservicenamespace",
					UID:       types.UID("my-service-uid"),
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Service",
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{Port: 80},
						{Port: 8080},
					},
				},
			},
			eps: &corev1.Endpoints{
				ObjectMeta: metav1.ObjectMeta{},
				TypeMeta:   metav1.TypeMeta{},
				Subsets:    []corev1.EndpointSubset{},
			},
			ips: []string{
				"10.0.0.1",
				"10.0.0.2",
			},
			result: &corev1.Endpoints{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "v1",
							Kind:       "Service",
							Name:       "My Service",
							UID:        types.UID("my-service-uid"),
						},
					},
					Name:      "My Service",
					Namespace: "myservicenamespace",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Endpoints",
				},
				Subsets: []corev1.EndpointSubset{
					{
						Addresses: []corev1.EndpointAddress{
							{
								IP: "10.0.0.1",
							},
							{
								IP: "10.0.0.2",
							},
						},
						Ports: []corev1.EndpointPort{
							{Port: 80},
							{Port: 8080},
						},
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ec := EndpointController{}
			assert.EqualValues(t, tc.result, ec.UpdateEndpoints(tc.eps, tc.svc, tc.ips))
		})
	}
}
