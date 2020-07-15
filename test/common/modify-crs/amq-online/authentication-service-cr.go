package amq_online

import (
	"fmt"
	enmasseadminv1beta1 "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/admin/v1beta1"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

//========================================================================================================
// enmasseadminv1beta1 AuthenticationService
//========================================================================================================
const (
	NoneAuthservice     = "none-authservice"
	StandardAuthservice = "standard-authservice"
)

type AuthenticationServiceCrWrapper struct {
	// Add real CR as type
	Cr *enmasseadminv1beta1.AuthenticationService
}

func (wrapper AuthenticationServiceCrWrapper) GetName() string {
	return wrapper.Cr.Name
}
func (wrapper AuthenticationServiceCrWrapper) GetNamespace() string {
	return wrapper.Cr.Namespace
}
func (wrapper AuthenticationServiceCrWrapper) GetResourceVersion() string {
	return wrapper.Cr.ResourceVersion
}
func (wrapper AuthenticationServiceCrWrapper) GetKind() string {
	return wrapper.Cr.Kind
}
func (wrapper AuthenticationServiceCrWrapper) GetCr() runtime.Object {
	return wrapper.Cr
}
func (wrapper AuthenticationServiceCrWrapper) SetAnnotations(annotations map[string]string) {
	wrapper.Cr.SetAnnotations(annotations)
}
func (wrapper AuthenticationServiceCrWrapper) GetAnnotations() map[string]string {
	return wrapper.Cr.GetAnnotations()
}

type AuthenticationServiceReference struct {
	IntegreatlyName            string
	IntegreatlyNamespace       string
	SpecType                   enmasseadminv1beta1.AuthenticationServiceType
	CredentialsSecretName      string
	CredentialsSecretNamespace string
	DatasourceType             enmasseadminv1beta1.DatasourceType
	DatasourceDatabase         string
	DatasourceHost             string
	DatasourcePort             int
}

func (reference *AuthenticationServiceReference) CrType() string {

	return EnmasseAuthenticationService
}

func (reference *AuthenticationServiceReference) DeleteExistingValues(t *testing.T, intContainer *modify_crs.Container, phase string) {
	cr, ok := intContainer.Get().(AuthenticationServiceCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlanReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	delete(ant, modify_crs.IntegreatlyName)
	delete(ant, modify_crs.IntegreatlyNamespace)
	cr.Cr.SetAnnotations(ant)
	//TODO delete cr.Spec.Type, do not know how to do this currently

	intContainer.Put(cr)
}

func (reference *AuthenticationServiceReference) CopyRequiredValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(AuthenticationServiceCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlanReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	reference.IntegreatlyName = ant[modify_crs.IntegreatlyName]
	reference.IntegreatlyNamespace = ant[modify_crs.IntegreatlyNamespace]
	reference.SpecType = cr.Cr.Spec.Type

	switch cr.GetName() {
	case StandardAuthservice:
		reference.CredentialsSecretName = cr.Cr.Spec.Standard.Datasource.CredentialsSecret.Name
		reference.CredentialsSecretNamespace = cr.Cr.Spec.Standard.Datasource.CredentialsSecret.Namespace
		reference.DatasourceType = cr.Cr.Spec.Standard.Datasource.Type
		reference.DatasourceDatabase = cr.Cr.Spec.Standard.Datasource.Database
		reference.DatasourceHost = cr.Cr.Spec.Standard.Datasource.Host
		reference.DatasourcePort = cr.Cr.Spec.Standard.Datasource.Port
	}
	intContainer.Put(cr)
}

func (reference *AuthenticationServiceReference) ChangeCRValues(t *testing.T, intContainer *modify_crs.Container, phase string) {

	cr, ok := intContainer.Get().(AuthenticationServiceCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlanReference from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	ant[modify_crs.IntegreatlyName] = "Bad Value"
	ant[modify_crs.IntegreatlyNamespace] = "Bad Value"
	cr.Cr.SetAnnotations(ant)

	switch cr.GetName() {
	case NoneAuthservice:
		cr.Cr.Spec.Type = "standard"
	case StandardAuthservice:
		cr.Cr.Spec.Type = "none"
		cr.Cr.Spec.Standard.Datasource.CredentialsSecret.Name = "bad value"
		cr.Cr.Spec.Standard.Datasource.CredentialsSecret.Namespace = "bad value"
		cr.Cr.Spec.Standard.Datasource.Type = "bad value"
		cr.Cr.Spec.Standard.Datasource.Database = "bad value"
		cr.Cr.Spec.Standard.Datasource.Host = "bad value"
		cr.Cr.Spec.Standard.Datasource.Port = 0
	}

	intContainer.Put(cr)
}

func (reference *AuthenticationServiceReference) CompareValues(t *testing.T, intContainer *modify_crs.Container, phase string) *[]modify_crs.CompareResult {
	var values []modify_crs.CompareResult

	cr, ok := intContainer.Get().(AuthenticationServiceCrWrapper)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlanReference from intContainer", phase)
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

	switch cr.GetName() {
	case NoneAuthservice:
		if cr.Cr.Spec.Type != reference.SpecType {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.GetKind(),
				Name:  cr.GetName(),
				Key:   "spec.type",
				Error: fmt.Sprintf("%s is not equal to expected %s", cr.Cr.Spec.Type, reference.SpecType),
			})
		}
	case StandardAuthservice:
		if cr.Cr.Spec.Type != reference.SpecType {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.GetKind(),
				Name:  cr.GetName(),
				Key:   "spec.type",
				Error: fmt.Sprintf("%s is not equal to expected %s", cr.Cr.Spec.Type, reference.SpecType),
			})
		}

		if cr.Cr.Spec.Standard.Datasource.CredentialsSecret.Name != reference.CredentialsSecretName {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.GetKind(),
				Name:  cr.GetName(),
				Key:   "spec.standard.datasource.credentialssecret.name",
				Error: fmt.Sprintf("%s is not equal to expected %s", cr.Cr.Spec.Standard.Datasource.CredentialsSecret.Name, reference.CredentialsSecretName),
			})
		}

		if cr.Cr.Spec.Standard.Datasource.CredentialsSecret.Namespace != reference.CredentialsSecretNamespace {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.GetKind(),
				Name:  cr.GetName(),
				Key:   "spec.standard.datasource.credentialssecret.namespace",
				Error: fmt.Sprintf("%s is not equal to expected %s", cr.Cr.Spec.Standard.Datasource.CredentialsSecret.Namespace, reference.CredentialsSecretNamespace),
			})
		}

		if cr.Cr.Spec.Standard.Datasource.Type != reference.DatasourceType {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.GetKind(),
				Name:  cr.GetName(),
				Key:   "spec.standard.datasource.type",
				Error: fmt.Sprintf("%s is not equal to expected %s", cr.Cr.Spec.Standard.Datasource.Type, reference.DatasourceType),
			})
		}

		if cr.Cr.Spec.Standard.Datasource.Database != reference.DatasourceDatabase {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.GetKind(),
				Name:  cr.GetName(),
				Key:   "spec.standard.datasource.database",
				Error: fmt.Sprintf("%s is not equal to expected %s", cr.Cr.Spec.Standard.Datasource.Database, reference.DatasourceDatabase),
			})
		}

		if cr.Cr.Spec.Standard.Datasource.Host != reference.DatasourceHost {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.GetKind(),
				Name:  cr.GetName(),
				Key:   "spec.standard.datasource.host",
				Error: fmt.Sprintf("%s is not equal to expected %s", cr.Cr.Spec.Standard.Datasource.Host, reference.DatasourceHost),
			})
		}

		if cr.Cr.Spec.Standard.Datasource.Port != reference.DatasourcePort {
			values = append(values, modify_crs.CompareResult{
				Type:  cr.GetKind(),
				Name:  cr.GetName(),
				Key:   "spec.standard.datasource.port",
				Error: fmt.Sprintf("%s is not equal to expected %s", cr.Cr.Spec.Standard.Datasource.Port, reference.DatasourcePort),
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
