package utils

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const endpointsKind = "Endpoints"

// UpdateEndpointsFunc UpdateEndpoints function type
type UpdateEndpointsFunc func(ep *corev1.Endpoints, svc *corev1.Service, scheme *runtime.Scheme, ips []string) (*corev1.Endpoints, error)

// UpdateEndpoints used to update an Endpoints instance from a Service and a IP list
func UpdateEndpoints(ep *corev1.Endpoints, svc *corev1.Service, scheme *runtime.Scheme, ips []string) (*corev1.Endpoints, error) {
	ports := getPortsforEndpoints(svc)
	newEp := ep.DeepCopy()
	setMetaForEndpoints(newEp, svc)
	newEp.Subsets = getSubsetsForEndpoints(ips, ports)
	if err := controllerutil.SetControllerReference(svc, newEp, scheme); err != nil {
		return nil, err
	}
	return newEp, nil
}

// getSubsetsForEndpoints generates an EndpointSubset object from ips and ports slices
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

// getPortsforEndpoints retrieves ports from Service object
// prioritizes Target Ports over Ports
func getPortsforEndpoints(svc *corev1.Service) []int32 {
	ports := []int32{}
	for _, port := range svc.Spec.Ports {
		targetPort := port.TargetPort.IntValue()
		if targetPort > 0 {
			ports = append(ports, int32(targetPort))
		} else {
			ports = append(ports, port.Port)
		}
	}
	return ports
}

// setMetaForEndpoints defines the metadata of Endpoints object
func setMetaForEndpoints(ep *corev1.Endpoints, svc *corev1.Service) {
	ep.Name = svc.Name
	ep.Namespace = svc.Namespace
	ep.APIVersion = svc.APIVersion
	ep.Kind = endpointsKind
}
