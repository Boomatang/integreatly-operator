package amq_online

import (
	"fmt"
	"github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta2"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

//========================================================================================================
// enmasse addressPlan
//========================================================================================================

type AddressPlanCr struct {
	// Add real CR as type
	Cr *v1beta2.AddressPlan
}

func (cr AddressPlanCr) GetName() string {
	return cr.Cr.Name
}
func (cr AddressPlanCr) GetNamespace() string {
	return cr.Cr.Namespace
}
func (cr AddressPlanCr) GetResourceVersion() string {
	return cr.Cr.ResourceVersion
}
func (cr AddressPlanCr) GetKind() string {
	return cr.Cr.Kind
}
func (cr AddressPlanCr) GetCr() runtime.Object {
	return cr.Cr
}

type AddressPlan struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
}

func (i *AddressPlan) CrType() string {

	return "enmasse.AddressPlan"
}

func (i *AddressPlan) DeleteExistingValues(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(AddressPlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	delete(ant, modify_crs.IntegreatlyName)
	delete(ant, modify_crs.IntegreatlyNamespace)
	cr.Cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func (i *AddressPlan) CopyRequiredValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(AddressPlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	i.IntegreatlyName = ant[modify_crs.IntegreatlyName]
	i.IntegreatlyNamespace = ant[modify_crs.IntegreatlyNamespace]

	intContainer.Put(cr)
}

func (i *AddressPlan) ChangeCRValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(AddressPlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	ant[modify_crs.IntegreatlyName] = "Bad Value"
	ant[modify_crs.IntegreatlyNamespace] = "Bad Value"
	cr.Cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func (i *AddressPlan) CompareValues(t *testing.T, intContainer *modify_crs.Container, phase string) *[]modify_crs.CompareResult {
	var values []modify_crs.CompareResult

	cr, ok := intContainer.Get().(AddressPlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
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

func (i *AddressPlan) AddCRDummyValues(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(AddressPlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}
	ant := cr.Cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.Cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func (i *AddressPlan) CheckDummyValuesStillExist(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(AddressPlanCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}
	ant := cr.Cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}

	intContainer.Put(cr)
}
