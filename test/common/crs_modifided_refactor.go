package common

import (
	goctx "context"
	enmasse "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta2"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	"github.com/integr8ly/integreatly-operator/test/common/modify-crs/amq-online"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
	"testing"
)

//========================================================================================================
// Setting up the test
//========================================================================================================

func TestResetCRs(t *testing.T, ctx *TestingContext) {
	var wg sync.WaitGroup
	addressSpacePlanTest(t, ctx, &wg)
	addressPlanTest(t, ctx, &wg)
	wg.Wait()
}

//========================================================================================================
// enmasse addressSpacePlan
//========================================================================================================

func addressSpacePlanTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	apl := &enmasse.AddressSpacePlanList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: modify_crs.AmqOnline,
	}
	err := ctx.Client.List(goctx.TODO(), apl, listOpts)
	if err != nil {
		t.Fatal("addressSpacePlan : Failed to get a list of address plan CR's from cluster")
	}

	for _, cr := range apl.Items {
		wg.Add(1)
		go addressSpacePlanTestSetup(t, ctx, wg, &cr)
		break // This will need to be removed, issue with concurrence on writes
	}
}

func addressSpacePlanTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr *enmasse.AddressSpacePlan) {
	defer wg.Done()
	ap := amq_online.AddressSpacePlan{}
	apcr := amq_online.AddressSpacePlanCr{cr}
	addressPlanContainer := &modify_crs.Container{}
	addressPlanContainer.Put(apcr)
	modifyExistingValues(t, ctx, &ap, addressPlanContainer)
	deleteExistingValues(t, ctx, &ap, addressPlanContainer)
	addNewCRValues(t, ctx, &ap, addressPlanContainer)
}

//========================================================================================================
// enmasse addressPlan
//========================================================================================================

func addressPlanTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	apl := &enmasse.AddressPlanList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: modify_crs.AmqOnline,
	}
	err := ctx.Client.List(goctx.TODO(), apl, listOpts)
	if err != nil {
		t.Fatal("addressPlan : Failed to get a list of address plan CR's from cluster")
	}

	for _, cr := range apl.Items {
		wg.Add(1)
		go addressPlanTestSetup(t, ctx, wg, &cr)
		break // This will need to be removed, issue with concurrence on writes
	}
}

func addressPlanTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr *enmasse.AddressPlan) {
	defer wg.Done()
	ap := amq_online.AddressPlan{}
	apcr := amq_online.AddressPlanCr{cr}
	addressPlanContainer := &modify_crs.Container{}
	addressPlanContainer.Put(apcr)
	modifyExistingValues(t, ctx, &ap, addressPlanContainer)
	deleteExistingValues(t, ctx, &ap, addressPlanContainer)
	addNewCRValues(t, ctx, &ap, addressPlanContainer)
}

//========================================================================================================
// generic functions
//========================================================================================================

func modifyExistingValues(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, crData *modify_crs.Container) {
	phase := "Modify Existing CR Values"
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
	rt.AddCRDummyValues(t, crData, phase)
	updateClusterCr(t, ctx, rt, crData, phase)
	compareAddedResultsAfterReconcile(t, ctx, rt, crData, phase)
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

		t.Logf("%s : %s: count = %v, revison = %s", phase, cr.GetName(), retryCount, cr.GetResourceVersion())
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
	cr, ok := intContainer.Get().(modify_crs.CrInterface)
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
	rt.CheckDummyValuesStillExist(t, intContainer, phase)
}

func waitReconcilingCR(t *testing.T, ctx *TestingContext, rt modify_crs.ResourceType, intContainer *modify_crs.Container) (done bool, err error) {
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Log(cr)
		t.Fatalf("waitReconcilingCR : Unable to read CR from intContainer")
	}

	resourceVersion := cr.GetResourceVersion()
	err = wait.Poll(modify_crs.RetryInterval, modify_crs.Timeout, func() (done bool, err error) {
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.GetName(), Namespace: cr.GetNamespace()}, cr.GetCr())
		if err != nil {
			return false, err
		}
		if resourceVersion != cr.GetResourceVersion() {
			return true, nil
		} else {
			return false, nil
		}
	})
	intContainer.Put(cr)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func getCR(intContainer *modify_crs.Container, rt modify_crs.ResourceType) (modify_crs.CrInterface, bool) {
	switch {
	case rt.CrType() == amq_online.EnmasseAddressPlan:
		cr, ok := intContainer.Get().(amq_online.AddressPlanCr)
		return cr, ok
	case rt.CrType() == amq_online.EnmasseAddressSpacePlan:
		cr, ok := intContainer.Get().(amq_online.AddressSpacePlanCr)
		return cr, ok
	default:
		return nil, false
	}
}
