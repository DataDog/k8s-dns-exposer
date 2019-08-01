// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

package service

import (
	"context"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"github.com/go-logr/logr"
	//apiequality "k8s.io/apimachinery/pkg/api/equality"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	//"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/DataDog/k8s-dns-exposer/pkg/controller/config"
	"github.com/DataDog/k8s-dns-exposer/pkg/controller/predicate"
	"github.com/DataDog/k8s-dns-exposer/pkg/utils"
)

func TestReconcileService_Reconcile(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Service{})
	s.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Endpoints{})

	service1 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "service1",
			Namespace:   "foo",
			Annotations: map[string]string{config.K8sDNSExposerAnnotationKey: "true"},
		},
		Spec: corev1.ServiceSpec{
			ExternalName: "foo.datadoghq.com",
			ClusterIP:    "Nonde",
		},
	}

	type fields struct {
		client              client.Client
		scheme              *runtime.Scheme
		dnsResolver         utils.DNSResolverIface
		updateEndpointsFunc utils.UpdateEndpointsFunc
		watcherPredicate    predicate.AnnotationPredicate
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     reconcile.Result
		wantErr  bool
		wantFunc func(client client.Client) error
	}{
		{
			name: "Service Not exist, return without requeue",
			fields: fields{
				scheme:              s,
				client:              fake.NewFakeClient(),
				updateEndpointsFunc: utils.UpdateEndpoints,
				dnsResolver:         &FakeResolver{},
				watcherPredicate:    predicate.AnnotationPredicate{Key: config.K8sDNSExposerAnnotationKey},
			},
			args: args{
				request: reconcile.Request{},
			},
			wantErr: false,
		},
		{
			name: "Service exist, Endpoints doesn't create endpoint",
			fields: fields{
				scheme:              s,
				client:              fake.NewFakeClient(service1),
				updateEndpointsFunc: utils.UpdateEndpoints,
				dnsResolver:         &FakeResolver{},
				watcherPredicate:    predicate.AnnotationPredicate{Key: config.K8sDNSExposerAnnotationKey},
			},
			args: args{
				request: newRequest(service1.Namespace, service1.Name),
			},
			want:    reconcile.Result{RequeueAfter: config.DefaultRequeueDuration},
			wantErr: false,
			wantFunc: func(c client.Client) error {
				endpoint := &corev1.Endpoints{}
				return c.Get(context.TODO(), newRequest(service1.Namespace, service1.Name).NamespacedName, endpoint)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReconcileService{
				client:              tt.fields.client,
				scheme:              tt.fields.scheme,
				dnsResolver:         tt.fields.dnsResolver,
				updateEndpointsFunc: tt.fields.updateEndpointsFunc,
				watcherPredicate:    tt.fields.watcherPredicate,
			}
			got, err := r.Reconcile(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileService.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconcileService.Reconcile() = %v, want %v", got, tt.want)
			}
			if tt.wantFunc != nil {
				if err = tt.wantFunc(tt.fields.client); err != nil {
					t.Errorf("ReconcileService.Reconcile() validation function return an error: %v", err)
				}
			}
		})
	}
}

type FakeResolver struct {
	ips []string
	err error
}

func (f *FakeResolver) Resolve(entry string) ([]string, error) {
	return f.ips, f.err
}

func newRequest(ns, name string) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: ns,
			Name:      name,
		},
	}
}
