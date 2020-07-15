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

type RoleCrWrapper struct {
	// Add real CR as type
	Cr *rbacv1.Role
}

func (wrapper RoleCrWrapper) GetName() string {
	return wrapper.Cr.Name
}
func (wrapper RoleCrWrapper) GetNamespace() string {
	return wrapper.Cr.Namespace
}
func (wrapper RoleCrWrapper) GetResourceVersion() string {
	return wrapper.Cr.ResourceVersion
}
func (wrapper RoleCrWrapper) GetKind() string {
	return wrapper.Cr.Kind
}
func (wrapper RoleCrWrapper) GetCr() runtime.Object {
	return wrapper.Cr
}
func (wrapper RoleCrWrapper) SetAnnotations(annotations map[string]string) {
	wrapper.Cr.SetAnnotations(annotations)
}
func (wrapper RoleCrWrapper) GetAnnotations() map[string]string {
	return wrapper.Cr.GetAnnotations()
}

type RoleReference struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
	Rules                []roleCrRole
}

type roleCrRole struct {
	APIGroup  []string
	Resources []string
	Verbs     []string
}

func (reference *RoleReference) CrType() string {

	return Rbacv1Role
}

func (reference *RoleReference) DeleteExistingValues(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(RoleCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.RoleReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	delete(ant, modify_crs.IntegreatlyName)
	delete(ant, modify_crs.IntegreatlyNamespace)
	cr.Cr.SetAnnotations(ant)
	cr.Cr.Rules = nil

	intContainer.Put(cr)
}

func (reference *RoleReference) CopyRequiredValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(RoleCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.RoleReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	reference.IntegreatlyName = ant[modify_crs.IntegreatlyName]
	reference.IntegreatlyNamespace = ant[modify_crs.IntegreatlyNamespace]

	for _, rule := range cr.Cr.Rules {
		reference.Rules = append(reference.Rules, roleCrRole{
			APIGroup:  rule.APIGroups,
			Resources: rule.Resources,
			Verbs:     rule.Verbs,
		})
	}

	intContainer.Put(cr)
}

func (reference *RoleReference) ChangeCRValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(RoleCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.RoleReference from intContainer", phase)
	}
	ant := cr.Cr.GetAnnotations()
	if ant != nil {
		ant[modify_crs.IntegreatlyName] = "Bad Value"
		ant[modify_crs.IntegreatlyNamespace] = "Bad Value"
		cr.Cr.SetAnnotations(ant)
	} else {
		t.Logf("cr: %s, is missing Integreatly Labels", cr.GetName())
	}

	for index, rule := range cr.Cr.Rules {
		for i := range rule.Resources {
			cr.Cr.Rules[index].Resources[i] = "Bad Value"
		}

		for i := range rule.Verbs {
			cr.Cr.Rules[index].Verbs[i] = "Bad Value"
		}

		for i := range rule.APIGroups {
			cr.Cr.Rules[index].APIGroups[i] = "Bad Value"
		}
	}

	intContainer.Put(cr)
}

func (reference *RoleReference) CompareValues(t *testing.T, intContainer *modify_crs.Container, phase string) *[]modify_crs.CompareResult {
	var values []modify_crs.CompareResult

	cr, ok := intContainer.Get().(RoleCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.RoleReference from intContainer", phase)
	}

	ant := cr.GetAnnotations()
	if ant[modify_crs.IntegreatlyName] != reference.IntegreatlyName {
		values = append(values, modify_crs.CompareResult{
			Type:  cr.Cr.Kind,
			Name:  cr.Cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[modify_crs.IntegreatlyName], reference.IntegreatlyName),
		})
	}

	if ant[modify_crs.IntegreatlyNamespace] != reference.IntegreatlyNamespace {
		values = append(values, modify_crs.CompareResult{
			Type:  cr.Cr.Kind,
			Name:  cr.Cr.Name,
			Key:   "metadata.annotations.integreatly-namespace",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[modify_crs.IntegreatlyNamespace], reference.IntegreatlyNamespace),
		})
	}

	for _, rule := range cr.Cr.Rules {
		for _, value := range rule.Resources {
			err := reference.compareResources(value)
			if err != nil {
				values = append(values, modify_crs.CompareResult{
					Type:  cr.Cr.Kind,
					Name:  cr.Cr.Name,
					Key:   "Roles.Resources",
					Error: err.Error(),
				})
			}
		}

		for _, value := range rule.Verbs {
			err := reference.compareVerbs(value)
			if err != nil {
				values = append(values, modify_crs.CompareResult{
					Type:  cr.Cr.Kind,
					Name:  cr.Cr.Name,
					Key:   "Roles.Verbs",
					Error: err.Error(),
				})
			}
		}

		for _, value := range rule.APIGroups {
			err := reference.compareAPIGroups(value)
			if err != nil {
				values = append(values, modify_crs.CompareResult{
					Type:  cr.Cr.Kind,
					Name:  cr.Cr.Name,
					Key:   "Roles.APIGroup",
					Error: err.Error(),
				})
			}
		}
	}

	intContainer.Put(cr)
	if len(values) > 0 {
		return &values
	} else {
		return nil
	}
}

func (reference *RoleReference) compareAPIGroups(value string) error {
	for _, item := range reference.Rules {
		for _, expected := range item.APIGroup {
			if value == expected {
				return nil
			}
		}
	}
	return fmt.Errorf("could not find %s in copied CR Roles.APIGroup", value)
}

func (reference *RoleReference) compareVerbs(value string) error {
	for _, item := range reference.Rules {
		for _, expected := range item.Verbs {
			if value == expected {
				return nil
			}
		}
	}
	return fmt.Errorf("could not find %s in copied CR Roles.Verbs", value)
}

func (reference *RoleReference) compareResources(value string) error {
	for _, item := range reference.Rules {
		for _, expected := range item.Resources {
			if value == expected {
				return nil
			}
		}
	}
	return fmt.Errorf("could not find %s in copied CR Roles.Resources", value)
}
