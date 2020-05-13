package amq_online

import (
	"fmt"
	"github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta2"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

//========================================================================================================
// enmasse addressSpacePlan
//========================================================================================================

type AddressSpacePlanCr struct {
	// Add real CR as type
	Cr *v1beta2.AddressSpacePlan
}

func (cr AddressSpacePlanCr) GetName() string {
	return cr.Cr.Name
}
func (cr AddressSpacePlanCr) GetNamespace() string {
	return cr.Cr.Namespace
}
func (cr AddressSpacePlanCr) GetResourceVersion() string {
	return cr.Cr.ResourceVersion
}
func (cr AddressSpacePlanCr) GetKind() string {
	return cr.Cr.Kind
}
func (cr AddressSpacePlanCr) GetCr() runtime.Object {
	return cr.Cr
}

type AddressSpacePlan struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
}

func (i *AddressSpacePlan) CrType() string {

	return EnmasseAddressSpacePlan
}

func (i *AddressSpacePlan) DeleteExistingValues(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(AddressSpacePlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressSpacePlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	delete(ant, modify_crs.IntegreatlyName)
	delete(ant, modify_crs.IntegreatlyNamespace)
	cr.Cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func (i *AddressSpacePlan) CopyRequiredValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(AddressSpacePlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressSpacePlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	i.IntegreatlyName = ant[modify_crs.IntegreatlyName]
	i.IntegreatlyNamespace = ant[modify_crs.IntegreatlyNamespace]

	intContainer.Put(cr)
}

func (i *AddressSpacePlan) ChangeCRValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(AddressSpacePlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressSpacePlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	ant[modify_crs.IntegreatlyName] = "Bad Value"
	ant[modify_crs.IntegreatlyNamespace] = "Bad Value"
	cr.Cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func (i *AddressSpacePlan) CompareValues(t *testing.T, intContainer *modify_crs.Container, phase string) *[]modify_crs.CompareResult {
	var values []modify_crs.CompareResult

	cr, ok := intContainer.Get().(AddressSpacePlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressSpacePlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	if ant[modify_crs.IntegreatlyName] != i.IntegreatlyName {
		values = append(values, modify_crs.CompareResult{
			Type:  cr.GetKind(),
			Name:  cr.GetName(),
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[modify_crs.IntegreatlyName], i.IntegreatlyName),
		})
	}

	if ant[modify_crs.IntegreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, modify_crs.CompareResult{
			Type:  cr.GetKind(),
			Name:  cr.GetName(),
			Key:   "metadata.annotations.integreatly-namespace",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[modify_crs.IntegreatlyNamespace], i.IntegreatlyNamespace),
		})
	}
	intContainer.Put(cr)
	if len(values) > 0 {
		return &values
	} else {
		return nil
	}
}

func (i *AddressSpacePlan) AddCRDummyValues(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(AddressSpacePlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressSpacePlan from intContainer", phase)
	}
	ant := cr.Cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.Cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func (i *AddressSpacePlan) CheckDummyValuesStillExist(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(AddressSpacePlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressSpacePlan from intContainer", phase)
	}
	ant := cr.Cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}

	intContainer.Put(cr)
}
