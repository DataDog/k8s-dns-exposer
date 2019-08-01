// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

package predicate

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnnotationPredicate_isAnnotationKeyPresent(t *testing.T) {
	type fields struct {
		Key   string
		Value string
	}
	type args struct {
		obj v1.Object
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "annotation not present",
			fields: fields{
				Key:   "foo.key",
				Value: "",
			},
			args: args{
				obj: &corev1.Service{},
			},
			want: false,
		},
		{
			name: "annotation is present",
			fields: fields{
				Key:   "foo.key",
				Value: "",
			},
			args: args{
				obj: &corev1.Service{
					ObjectMeta: v1.ObjectMeta{
						Annotations: map[string]string{"foo.key": "bar"},
					},
				},
			},
			want: true,
		},
		{
			name: "annotation key and value is present",
			fields: fields{
				Key:   "foo.key",
				Value: "foo.value",
			},
			args: args{
				obj: &corev1.Service{
					ObjectMeta: v1.ObjectMeta{
						Annotations: map[string]string{"foo.key": "foo.value"},
					},
				},
			},
			want: true,
		},
		{
			name: "annotation value is present, but wrong value",
			fields: fields{
				Key:   "foo.key",
				Value: "bar.value",
			},
			args: args{
				obj: &corev1.Service{
					ObjectMeta: v1.ObjectMeta{
						Annotations: map[string]string{"foo.key": "foo.value"},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnnotationPredicate{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.isAnnotationKeyPresent(tt.args.obj); got != tt.want {
				t.Errorf("AnnotationPredicate.isAnnotationKeyPresent() = %v, want %v", got, tt.want)
			}
		})
	}
}
