package predicate

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = AnnotationPredicate{}

// AnnotationPredicate Predicate interface implementation that is looking
// for a specific Annotation key in the object metadata
// if Value is also set, validate that the Value is also equal
type AnnotationPredicate struct {
	Key   string
	Value string
}

// Create returns true if the Create event should be processed
func (a AnnotationPredicate) Create(e event.CreateEvent) bool {
	return a.isAnnotationKeyPresent(e.Meta)
}

// Delete returns true if the Delete event should be processed
func (a AnnotationPredicate) Delete(e event.DeleteEvent) bool {
	return a.isAnnotationKeyPresent(e.Meta)
}

// Update returns true if the Update event should be processed
func (a AnnotationPredicate) Update(e event.UpdateEvent) bool {
	return a.isAnnotationKeyPresent(e.MetaNew) || a.isAnnotationKeyPresent(e.MetaOld)
}

// Generic returns true if the Generic event should be processed
func (a AnnotationPredicate) Generic(e event.GenericEvent) bool {
	return a.isAnnotationKeyPresent(e.Meta)
}

func (a AnnotationPredicate) isAnnotationKeyPresent(obj v1.Object) bool {
	if val, ok := obj.GetAnnotations()[a.Key]; ok {
		if a.Value != "" && val != a.Value {
			return false
		}
		return true
	}
	return false
}
