package amq_online

import (
	"fmt"
	"github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta1"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

//========================================================================================================
// enmasse StandardInfraConfig
//========================================================================================================

type BrokeredInfraConfigCrWrapper struct {
	// Add real CR as type
	Cr *v1beta1.BrokeredInfraConfig
}

func (wrapper BrokeredInfraConfigCrWrapper) GetName() string {
	return wrapper.Cr.Name
}
func (wrapper BrokeredInfraConfigCrWrapper) GetNamespace() string {
	return wrapper.Cr.Namespace
}
func (wrapper BrokeredInfraConfigCrWrapper) GetResourceVersion() string {
	return wrapper.Cr.ResourceVersion
}
func (wrapper BrokeredInfraConfigCrWrapper) GetKind() string {
	return wrapper.Cr.Kind
}
func (wrapper BrokeredInfraConfigCrWrapper) GetCr() runtime.Object {
	return wrapper.Cr
}
func (wrapper BrokeredInfraConfigCrWrapper) SetAnnotations(annotations map[string]string) {
	wrapper.Cr.SetAnnotations(annotations)
}
func (wrapper BrokeredInfraConfigCrWrapper) GetAnnotations() map[string]string {
	return wrapper.Cr.GetAnnotations()
}

type BrokeredInfraConfigReference struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
}

func (reference *BrokeredInfraConfigReference) CrType() string {

	return EnmasseBrokeredInfraConfig
}

func (reference *BrokeredInfraConfigReference) DeleteExistingValues(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(BrokeredInfraConfigCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.BrokeredInfraConfigReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	delete(ant, modify_crs.IntegreatlyName)
	delete(ant, modify_crs.IntegreatlyNamespace)
	cr.Cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func (reference *BrokeredInfraConfigReference) CopyRequiredValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(BrokeredInfraConfigCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.BrokeredInfraConfigReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	reference.IntegreatlyName = ant[modify_crs.IntegreatlyName]
	reference.IntegreatlyNamespace = ant[modify_crs.IntegreatlyNamespace]

	intContainer.Put(cr)
}

func (reference *BrokeredInfraConfigReference) ChangeCRValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(BrokeredInfraConfigCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.BrokeredInfraConfigReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	ant[modify_crs.IntegreatlyName] = "Bad Value"
	ant[modify_crs.IntegreatlyNamespace] = "Bad Value"
	cr.Cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func (reference *BrokeredInfraConfigReference) CompareValues(t *testing.T, intContainer *modify_crs.Container, phase string) *[]modify_crs.CompareResult {
	var values []modify_crs.CompareResult

	cr, ok := intContainer.Get().(BrokeredInfraConfigCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.BrokeredInfraConfigReference from intContainer", phase)
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
	intContainer.Put(cr)
	if len(values) > 0 {
		return &values
	} else {
		return nil
	}
}
