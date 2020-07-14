package amq_online

import (
	"fmt"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

//========================================================================================================
// enmasse rbacv1.RoleBinding
//========================================================================================================

type RoleBindingCrWrapper struct {
	// Add real CR as type
	Cr *rbacv1.RoleBinding
}

func (wrapper RoleBindingCrWrapper) GetName() string {
	return wrapper.Cr.Name
}
func (wrapper RoleBindingCrWrapper) GetNamespace() string {
	return wrapper.Cr.Namespace
}
func (wrapper RoleBindingCrWrapper) GetResourceVersion() string {
	return wrapper.Cr.ResourceVersion
}
func (wrapper RoleBindingCrWrapper) GetKind() string {
	return wrapper.Cr.Kind
}
func (wrapper RoleBindingCrWrapper) GetCr() runtime.Object {
	return wrapper.Cr
}
func (wrapper RoleBindingCrWrapper) SetAnnotations(annotations map[string]string) {
	wrapper.Cr.SetAnnotations(annotations)
}
func (wrapper RoleBindingCrWrapper) GetAnnotations() map[string]string {
	return wrapper.Cr.GetAnnotations()
}

type RoleBindingReference struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
	RoleRefName          string
	RoleRefKind          string
	Subjects             []roleBindingSubject
}

type roleBindingSubject struct {
	SubjectName string
	SubjectKind string
}

func (reference *RoleBindingReference) CrType() string {

	return Rbacv1RoleBinding
}

func (reference *RoleBindingReference) DeleteExistingValues(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(RoleBindingCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.RoleBindingReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	delete(ant, modify_crs.IntegreatlyName)
	delete(ant, modify_crs.IntegreatlyNamespace)
	cr.Cr.SetAnnotations(ant)
	cr.Cr.Subjects = nil

	intContainer.Put(cr)
}

func (reference *RoleBindingReference) CopyRequiredValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(RoleBindingCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.RoleBindingReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	reference.IntegreatlyName = ant[modify_crs.IntegreatlyName]
	reference.IntegreatlyNamespace = ant[modify_crs.IntegreatlyNamespace]
	reference.RoleRefKind = cr.Cr.RoleRef.Kind
	reference.RoleRefName = cr.Cr.RoleRef.Name

	for _, subject := range cr.Cr.Subjects {
		reference.Subjects = append(reference.Subjects, roleBindingSubject{
			SubjectKind: subject.Kind,
			SubjectName: subject.Name,
		})
	}

	intContainer.Put(cr)
}

func (reference *RoleBindingReference) ChangeCRValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(RoleBindingCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.RoleBindingReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	ant[modify_crs.IntegreatlyName] = "Bad Value"
	ant[modify_crs.IntegreatlyNamespace] = "Bad Value"
	cr.Cr.SetAnnotations(ant)

	for index := range cr.Cr.Subjects {
		cr.Cr.Subjects[index].Name = "bad-value"
		cr.Cr.Subjects[index].Kind = "ServiceAccount"
		cr.Cr.Subjects[index].APIGroup = ""
	}

	intContainer.Put(cr)
}

func (reference *RoleBindingReference) CompareValues(t *testing.T, intContainer *modify_crs.Container, phase string) *[]modify_crs.CompareResult {
	var values []modify_crs.CompareResult

	cr, ok := intContainer.Get().(RoleBindingCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.RoleBindingReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	if ant[modify_crs.IntegreatlyName] != reference.IntegreatlyName {
		values = append(values, modify_crs.CompareResult{
			Type:  cr.GetKind(),
			Name:  cr.GetName(),
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[modify_crs.IntegreatlyName], reference.IntegreatlyName),
		})
	}

	if ant[modify_crs.IntegreatlyNamespace] != reference.IntegreatlyNamespace {
		values = append(values, modify_crs.CompareResult{
			Type:  cr.GetKind(),
			Name:  cr.GetName(),
			Key:   "metadata.annotations.integreatly-namespace",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[modify_crs.IntegreatlyNamespace], reference.IntegreatlyNamespace),
		})
	}

	for _, subject := range cr.Cr.Subjects {
		err := reference.compareSubjectName(subject.Name)
		if err != nil {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.Cr.Kind,
				Name:  cr.Cr.Name,
				Key:   "subjects.[].name",
				Error: err.Error(),
			})
		}

		err = reference.compareSubjectKind(subject.Kind)
		if err != nil {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.Cr.Kind,
				Name:  cr.Cr.Name,
				Key:   "subjects.[].Kind",
				Error: err.Error(),
			})
		}
	}

	intContainer.Put(cr)
	if len(values) > 0 {
		return &values
	} else {
		return nil
	}
}

func (reference *RoleBindingReference) compareSubjectName(value string) error {
	for _, item := range reference.Subjects {
		if value == item.SubjectName {
			return nil
		}
	}
	return fmt.Errorf("could not find %s in copied CR Subject.Name", value)
}

func (reference *RoleBindingReference) compareSubjectKind(value string) error {
	for _, item := range reference.Subjects {
		if value == item.SubjectKind {
			return nil
		}
	}
	return fmt.Errorf("could not find %s in copied CR Subject.Kind", value)
}
