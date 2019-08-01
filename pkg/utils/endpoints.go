package utils

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EndpointControllerInterface interface {
	UpdateEndpoints(ep *corev1.Endpoints, svc *corev1.Service, ips []string) *corev1.Endpoints
}

type EndpointController struct {
}

func (ec *EndpointController) UpdateEndpoints(ep *corev1.Endpoints, svc *corev1.Service, ips []string) *corev1.Endpoints {
	newEp := ep.DeepCopy()
	ports := getPortsforEndpoints(svc)
	newEp.Subsets = getSubsetsForEndpoints(ips, ports)
	newEp.OwnerReferences = getOwnerRefForEndpoints(svc)
	return newEp
}

func getSubsetsForEndpoints(ips []string, ports []int32) []corev1.EndpointSubset {
	epAddresses := []corev1.EndpointAddress{}
	for _, ip := range ips {
		epAddresses = append(epAddresses, corev1.EndpointAddress{
			IP: ip,
		})
	}

	epPorts := []corev1.EndpointPort{}
	for _, port := range ports {
		epPorts = append(epPorts, corev1.EndpointPort{
			Port: port,
		})
	}

	return []corev1.EndpointSubset{
		{
			Addresses: epAddresses,
			Ports:     epPorts,
		},
	}
}

func getOwnerRefForEndpoints(svc *corev1.Service) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: svc.APIVersion,
			Kind:       svc.Kind,
			Name:       svc.Name,
			UID:        svc.UID,
		},
	}
}

func getPortsforEndpoints(svc *corev1.Service) []int32 {
	ports := []int32{}
	for _, port := range svc.Spec.Ports {
		if port.TargetPort.Size() > 0 {
			ports = append(ports, int32(port.TargetPort.IntValue()))
		} else {
			ports = append(ports, port.Port)
		}
	}
	return ports
}
