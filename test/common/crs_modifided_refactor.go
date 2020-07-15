package common

import (
	goctx "context"
	enmasseadminv1beta1 "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/admin/v1beta1"
	v1beta1 "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta1"
	enmasse "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta2"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs/amq-online"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"sync"
	"testing"
)

//========================================================================================================
// Setting up the test
//========================================================================================================

func TestResetCRs(t *testing.T, ctx *TestingContext) {
	var wg sync.WaitGroup
	//authenticationServiceTest(t, ctx, &wg) // Broken
	//addressSpacePlanTest(t, ctx, &wg)
	//addressPlanTest(t, ctx, &wg)
	//roleBindingTest(t, ctx, &wg)
	//roleTest(t, ctx, &wg)
	standardInfraConfigTest(t, ctx, &wg)
	wg.Wait()
}

//========================================================================================================
// enmasseadminv1beta1 AuthenticationService
//========================================================================================================

func authenticationServiceTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	asl := &enmasseadminv1beta1.AuthenticationServiceList{}
	err := ctx.Client.List(goctx.TODO(), asl, amq_online.ListOpts)
	if err != nil {
		t.Fatal("AuthenticationService : Failed to get a list of CR's from cluster")
	}
	var crNames []string
	for _, cr := range asl.Items {
		wg.Add(1)
		crNames = append(crNames, cr.Name)
		go authenticationServiceTestSetup(t, ctx, wg, cr)
		// Must check all cr's There is two or more configurations been checked
	}
	t.Logf("AuthenticationService CRs : %s", crNames)

}

func authenticationServiceTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr enmasseadminv1beta1.AuthenticationService) {
	defer wg.Done()
	as := amq_online.AuthenticationServiceReference{}
	ascr := amq_online.AuthenticationServiceCrWrapper{&cr}
	authenticationServiceContainer := &modify_crs.Container{}
	authenticationServiceContainer.Put(ascr)
	runCrTest(t, ctx, &as, authenticationServiceContainer)
}

//========================================================================================================
// enmasse addressSpacePlan
//========================================================================================================

func addressSpacePlanTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	aspl := &enmasse.AddressSpacePlanList{}
	err := ctx.Client.List(goctx.TODO(), aspl, amq_online.ListOpts)
	if err != nil {
		t.Fatal("addressSpacePlan : Failed to get a list of CR's from cluster")
	}
	var crNames []string
	for _, cr := range aspl.Items {
		wg.Add(1)
		crNames = append(crNames, cr.Name)
		go addressSpacePlanTestSetup(t, ctx, wg, cr)
	}
	t.Logf("addressSpacePlan CRs : %s", crNames)

}

func addressSpacePlanTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr enmasse.AddressSpacePlan) {
	defer wg.Done()
	asp := amq_online.AddressSpacePlanReference{}
	aspcr := amq_online.AddressSpacePlanCrWrapper{&cr}
	addressSpacePlanContainer := &modify_crs.Container{}
	addressSpacePlanContainer.Put(aspcr)
	runCrTest(t, ctx, &asp, addressSpacePlanContainer)
}

//========================================================================================================
// enmasse addressPlan
//========================================================================================================

func addressPlanTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	apl := &enmasse.AddressPlanList{}
	err := ctx.Client.List(goctx.TODO(), apl, amq_online.ListOpts)
	if err != nil {
		t.Fatal("addressPlan : Failed to get a list of address plan CR's from cluster")
	}
	var crNames []string
	for _, cr := range apl.Items {
		wg.Add(1)
		crNames = append(crNames, cr.Name)
		go addressPlanTestSetup(t, ctx, wg, cr)
	}
	t.Logf("addressPlan CRs : %s", crNames)

}

func addressPlanTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr enmasse.AddressPlan) {
	defer wg.Done()
	ap := amq_online.AddressPlanReference{}
	apcr := amq_online.AddressPlanCrWrapper{&cr}
	addressPlanContainer := &modify_crs.Container{}
	addressPlanContainer.Put(apcr)
	runCrTest(t, ctx, &ap, addressPlanContainer)
}

//========================================================================================================
// enmasse rbacv1.RoleBinding
//========================================================================================================

func roleBindingTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	rbl := &rbacv1.RoleBindingList{}
	err := ctx.Client.List(goctx.TODO(), rbl, amq_online.ListOpts)
	if err != nil {
		t.Fatal("RoleBinding : Failed to get a list of RoleBinding CR's from cluster")
	}

	var crNames []string
	var skipped []string
	for _, cr := range rbl.Items {
		if cr.Name == "dedicated-admins-service-admin" {
			wg.Add(1)
			crNames = append(crNames, cr.Name)
			go roleBindingTestSetup(t, ctx, wg, cr)
		} else {
			skipped = append(skipped, cr.Name)
		}
	}
	t.Logf("rbacv1.RoleBinding : The following CR's were skipped, %s", skipped)
	t.Logf("rbacv1.RoleBinding  CRs : %s", crNames)

}

func roleBindingTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr rbacv1.RoleBinding) {
	defer wg.Done()
	rb := amq_online.RoleBindingReference{}
	rbcr := amq_online.RoleBindingCrWrapper{&cr}
	roleBindingContainer := &modify_crs.Container{}
	roleBindingContainer.Put(rbcr)
	runCrTest(t, ctx, &rb, roleBindingContainer)
}

//========================================================================================================
// enmasse rbacv1.Role
//========================================================================================================

func roleTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	rbl := &rbacv1.RoleList{}
	err := ctx.Client.List(goctx.TODO(), rbl, amq_online.ListOpts)
	if err != nil {
		t.Fatal("Role : Failed to get a list of CR's from cluster")
	}

	var crNames []string
	var skippedCrs []string
	for _, cr := range rbl.Items {
		if strings.Contains(cr.Name, "amq-online") || cr.Name == "rhmi-registry-cs-configmap-reader" {
			skippedCrs = append(skippedCrs, cr.Name)
		} else {
			wg.Add(1)
			crNames = append(crNames, cr.Name)
			go roleTestSetup(t, ctx, wg, cr)
		}
	}
	t.Logf("rbacv1.Role CRs : %s", crNames)
	t.Logf("Skipped rbacv1.Role CRs : %s", skippedCrs)

}

func roleTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr rbacv1.Role) {
	defer wg.Done()
	rb := amq_online.RoleReference{}
	rbcr := amq_online.RoleCrWrapper{&cr}
	roleBindingContainer := &modify_crs.Container{}
	roleBindingContainer.Put(rbcr)
	runCrTest(t, ctx, &rb, roleBindingContainer)
}

//========================================================================================================
// enmasse StandardInfraConfig
//========================================================================================================

func standardInfraConfigTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	sicl := &v1beta1.StandardInfraConfigList{}
	err := ctx.Client.List(goctx.TODO(), sicl, amq_online.ListOpts)
	if err != nil {
		t.Fatal("Role : Failed to get a list of CR's from cluster")
	}

	var crNames []string
	for _, cr := range sicl.Items {
		wg.Add(1)
		crNames = append(crNames, cr.Name)
		go standardInfraConfigTestSetup(t, ctx, wg, cr)
	}
	t.Logf("rbacv1.Role CRs : %s", crNames)

}

func standardInfraConfigTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr v1beta1.StandardInfraConfig) {
	defer wg.Done()
	sic := amq_online.StandardInfraConfigReference{}
	siccr := amq_online.StandardInfraConfigCrWrapper{&cr}
	standardInfraConfigContainer := &modify_crs.Container{}
	standardInfraConfigContainer.Put(siccr)
	runCrTest(t, ctx, &sic, standardInfraConfigContainer)
}

//========================================================================================================
// generic functions
//========================================================================================================

func getCR(intContainer *modify_crs.Container, rt modify_crs.ResourceType) (modify_crs.CrInterface, bool) {
	switch rt.CrType() {
	case amq_online.EnmasseAddressPlan:
		cr, ok := intContainer.Get().(amq_online.AddressPlanCrWrapper)
		return cr, ok
	case amq_online.EnmasseAddressSpacePlan:
		cr, ok := intContainer.Get().(amq_online.AddressSpacePlanCrWrapper)
		return cr, ok
	case amq_online.EnmasseAuthenticationService:
		cr, ok := intContainer.Get().(amq_online.AuthenticationServiceCrWrapper)
		return cr, ok
	case amq_online.Rbacv1RoleBinding:
		cr, ok := intContainer.Get().(amq_online.RoleBindingCrWrapper)
		return cr, ok
	case amq_online.Rbacv1Role:
		cr, ok := intContainer.Get().(amq_online.RoleCrWrapper)
		return cr, ok
	case amq_online.EnmasseStandardInfraConfig:
		cr, ok := intContainer.Get().(amq_online.StandardInfraConfigCrWrapper)
		return cr, ok
	default:
		return nil, false
	}
}

func runCrTest(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, crData *modify_crs.Container) {
	modifyExistingValues(t, ctx, rt, crData)
	deleteExistingValues(t, ctx, rt, crData)
	addNewCRValues(t, ctx, rt, crData)
	tidyCr(t, ctx, rt, crData)
}

func modifyExistingValues(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, crData *modify_crs.Container) {
	phase := "Modify Existing CR Values"
	refreshCR(t, ctx, rt, crData, phase)

	rt.CopyRequiredValues(t, crData, phase)
	rt.ChangeCRValues(t, crData, phase)
	updateClusterCr(t, ctx, rt, crData, phase)
	results := compareResultsAfterReconcile(t, ctx, rt, crData, phase)
	if results != nil {
		for _, result := range *results {
			t.Logf("%s : %s: %s: %s: %s", phase, result.Type, result.Name, result.Key, result.Error)
		}
		t.Fatalf("%s : Failed to reset the CR values", phase)
	}
}

func deleteExistingValues(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, crData *modify_crs.Container) {
	phase := "Delete Existing CR Values"
	refreshCR(t, ctx, rt, crData, phase)

	rt.CopyRequiredValues(t, crData, phase)
	rt.DeleteExistingValues(t, crData, phase)
	updateClusterCr(t, ctx, rt, crData, phase)
	results := compareResultsAfterReconcile(t, ctx, rt, crData, phase)
	if results != nil {
		for _, result := range *results {
			t.Logf("%s : %s: %s: %s: %s", phase, result.Type, result.Name, result.Key, result.Error)
		}
		t.Fatalf("%s : Failed to reset the CR values", phase)
	}
}

func addNewCRValues(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, crData *modify_crs.Container) {
	phase := "Adding New CR Values"
	refreshCR(t, ctx, rt, crData, phase)

	AddDummyCRValues(t, rt, crData, phase)
	updateClusterCr(t, ctx, rt, crData, phase)
	compareAddedResultsAfterReconcile(t, ctx, rt, crData, phase)
}

func tidyCr(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, crData *modify_crs.Container) {
	phase := "Tidy New CR Values"
	refreshCR(t, ctx, rt, crData, phase)
	resetDummyValues(t, rt, crData, phase)
	updateClusterCr(t, ctx, rt, crData, phase)
}

func updateClusterCr(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, intContainer *modify_crs.Container, phase string) {
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Log(cr)
		t.Fatalf("%s : Unable to read CR from intContainer", phase)
	}

	retryCount := 3
	retry := true
	for retry {
		err := ctx.Client.Update(goctx.TODO(), cr.GetCr())
		if err != nil && retryCount == 0 {
			retry = false
			t.Log(cr)
			t.Log(err)
			t.Fatalf("%s : Failed to update CR on cluster", phase)
		} else if err != nil && retryCount > 0 {
			t.Logf("Update failed, retry count : %v", retryCount)
			retryCount -= 1
		} else {
			retry = false
		}
	}
	intContainer.Put(cr)
}

func compareResultsAfterReconcile(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, intContainer *modify_crs.Container, phase string) *[]modify_crs.CompareResult {
	var results *[]modify_crs.CompareResult
	retryCount := 3
	forceRetry := true

	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Log(cr)
		t.Fatalf("%s : Unable to read CR from intContainer", phase)
	}

	for forceRetry {
		// Force Retry is required to remove flaky test results after random updates
		err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.GetName(), Namespace: cr.GetNamespace()}, cr.GetCr())
		if err != nil {
			t.Fatalf("%s : Fail to refresh the cr", phase)
		}

		//t.Logf("%s : %s: count = %v, revison = %s", phase, cr.GetName(), retryCount, cr.GetResourceVersion())
		intContainer.Put(cr)
		_, err = waitReconcilingCR(t, ctx, rt, intContainer)
		if err != nil {
			t.Fatalf("%s : %s: %s:, %s", phase, cr.GetKind(), cr.GetName(), err)
		}

		results = rt.CompareValues(t, intContainer, phase)

		if results == nil {
			forceRetry = false
			retryCount -= 1
		}
		retryCount -= 1
		if retryCount < 0 {
			forceRetry = false
		}
	}
	return results
}

func compareAddedResultsAfterReconcile(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, intContainer *modify_crs.Container, phase string) {
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Log(cr)
		t.Fatalf("%s : Unable to read CR from intContainer", phase)
	}

	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.GetName(), Namespace: cr.GetNamespace()}, cr.GetCr())
	if err != nil {
		t.Fatalf("%s : Fail to refresh the cr", phase)
	}
	intContainer.Put(cr)

	_, err = waitReconcilingCR(t, ctx, rt, intContainer)
	if err != nil && err.Error() != "timed out waiting for the condition" {
		t.Fatalf("%s : %s: %s:, %s", phase, cr.GetKind(), cr.GetName(), err)
	}
	CheckDummyValuesStillExist(t, rt, intContainer, phase)
}

func waitReconcilingCR(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, intContainer *modify_crs.Container) (done bool, err error) {
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Log(cr)
		t.Fatalf("waitReconcilingCR : Unable to read CR from intContainer")
	}
	resourceVersion := cr.GetResourceVersion()
	intContainer.Put(cr)

	err = wait.Poll(modify_crs.RetryInterval, modify_crs.Timeout, func() (done bool, err error) {
		cr, ok := getCR(intContainer, rt)
		if !ok {
			t.Log(cr)
			t.Fatalf("waitReconcilingCR : Unable to read CR from intContainer")
		}
		defer intContainer.Put(cr)

		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.GetName(), Namespace: cr.GetNamespace()}, cr.GetCr())
		if err != nil {
			return false, err
		}
		//t.Logf("checking have %s : got %s", resourceVersion, cr.GetResourceVersion())
		if resourceVersion != cr.GetResourceVersion() {
			return true, nil
		} else {
			return false, nil
		}
	})
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func AddDummyCRValues(t *testing.T, rt modify_crs.ResourceType, intContainer *modify_crs.Container, phase string) {
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressSpacePlanReference from intContainer", phase)
	}
	ant := cr.GetAnnotations()
	if ant == nil {
		ant = map[string]string{}
	}
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func CheckDummyValuesStillExist(t *testing.T, rt modify_crs.ResourceType, intContainer *modify_crs.Container, phase string) {
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlanReference from intContainer", phase)
	}
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}

	intContainer.Put(cr)
}

func resetDummyValues(t *testing.T, rt modify_crs.ResourceType, intContainer *modify_crs.Container, phase string) {
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlanReference from intContainer", phase)
	}
	ant := cr.GetAnnotations()
	delete(ant, "dummy-value")
	cr.SetAnnotations(ant)

	intContainer.Put(cr)
}

func refreshCR(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, intContainer *modify_crs.Container, phase string) {
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Log(cr)
		t.Fatalf("%s : Unable to read CR from intContainer", phase)
	}

	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.GetName(), Namespace: cr.GetNamespace()}, cr.GetCr())
	if err != nil {
		t.Fatalf("%s : Fail to refresh the cr", phase)
	}
	intContainer.Put(cr)
}
