package common

import (
	goctx "context"
	"fmt"
	enmasse "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
	"testing"
)

type Container []interface{}

func (c *Container) Put(elem interface{}) {
	*c = append(*c, elem)
}
func (c *Container) Get() interface{} {
	elem := (*c)[0]
	*c = (*c)[1:]
	return elem
}

type compareResult2BR struct {
	Type  string
	Name  string
	Key   string
	Error string
}

type resourceType interface {
	crType() string
	copyRequiredValues(t *testing.T, intContainer Container, phase string)
	deleteExistingValues(t *testing.T, intContainer Container, phase string)
	changeCRValues(t *testing.T, intContainer Container, phase string)
	compareValues(t *testing.T, intContainer Container, phase string) *[]compareResult2BR
	addCRDummyValues(t *testing.T, intContainer Container, phase string)
	checkDummyValuesStillExist(t *testing.T, intContainer Container, phase string)
}

//========================================================================================================
// Working on
//========================================================================================================

type crInterface interface {
	// For any methods that are required in a common
	getName() string
	getNamespace() string
	getResourceVersion() string
	getKind() string
	getCr() runtime.Object
}

type AddressSpaceCr struct {
	// Add real CR as type
	Cr enmasse.AddressPlan
}

func (cr AddressSpaceCr) getName() string {
	return cr.Cr.Name
}
func (cr AddressSpaceCr) getNamespace() string {
	return cr.Cr.Namespace
}
func (cr AddressSpaceCr) getResourceVersion() string {
	return cr.Cr.ResourceVersion
}
func (cr AddressSpaceCr) getKind() string {
	return cr.Cr.Kind
}
func (cr AddressSpaceCr) getCr() runtime.Object {
	return &cr.Cr
}

//========================================================================================================
// enmasse addressPlan
//========================================================================================================

type addressPlan2BR struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
}

func (i *addressPlan2BR) crType() string {

	return "enmasse.AddressPlan"
}

func (i *addressPlan2BR) deleteExistingValues(t *testing.T, intContainer Container, phase string) {
	cr, ok := intContainer.Get().(enmasse.AddressPlan)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}

	ant := cr.GetAnnotations()
	delete(ant, integreatlyName)
	delete(ant, integreatlyNamespace)
	cr.SetAnnotations(ant)
}

func (i *addressPlan2BR) copyRequiredValues(t *testing.T, intContainer Container, phase string) {

	cr, ok := intContainer.Get().(AddressSpaceCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	i.IntegreatlyName = ant[integreatlyName]
	i.IntegreatlyNamespace = ant[integreatlyNamespace]
}

func (i *addressPlan2BR) changeCRValues(t *testing.T, intContainer Container, phase string) {

	cr, ok := intContainer.Get().(AddressSpaceCr)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}

	ant := cr.Cr.GetAnnotations()
	ant[integreatlyName] = "Bad Value"
	ant[integreatlyNamespace] = "Bad Value"
	cr.Cr.SetAnnotations(ant)

	//intContainer.Put(cr)
}

func (i *addressPlan2BR) compareValues(t *testing.T, intContainer Container, phase string) *[]compareResult2BR {
	var values []compareResult2BR

	cr, ok := intContainer.Get().(enmasse.AddressPlan)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}

	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, compareResult2BR{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, compareResult2BR{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-namespace",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyNamespace], i.IntegreatlyNamespace),
		})
	}

	if len(values) > 0 {
		return &values
	} else {
		return nil
	}
}

func (i *addressPlan2BR) addCRDummyValues(t *testing.T, intContainer Container, phase string) {
	cr, ok := intContainer.Get().(enmasse.AddressPlan)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}
	ant := cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)
}

func (i *addressPlan2BR) checkDummyValuesStillExist(t *testing.T, intContainer Container, phase string) {
	cr, ok := intContainer.Get().(enmasse.AddressPlan)
	if !ok {
		t.Fatalf("%s : Unable to read enmasse.AddressPlan from intContainer", phase)
	}
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}
}

//========================================================================================================
// Setting up the test
//========================================================================================================

func TestResetCRs2BR(t *testing.T, ctx *TestingContext) {
	var wg sync.WaitGroup
	addressPlanTest(t, ctx, &wg)
	wg.Wait()
}

func addressPlanTest(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	apl := &enmasse.AddressPlanList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: amqOnline,
	}
	err := ctx.Client.List(goctx.TODO(), apl, listOpts)
	if err != nil {
		t.Fatal("addressPlan : Failed to get a list of address plan CR's from cluster")
	}

	for _, cr := range apl.Items {
		wg.Add(1)
		go addressPlanTestSetup(t, ctx, wg, cr)
		break
	}
}

func addressPlanTestSetup(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup, cr enmasse.AddressPlan) {
	defer wg.Done()
	ap := addressPlan2BR{}
	ascr := AddressSpaceCr{cr}
	addressPlanContainer := &Container{}
	addressPlanContainer.Put(ascr)
	modifyExistingValues(t, ctx, &ap, *addressPlanContainer)
	//deleteExistingValues(t, ctx, &ap, *addressPlanContainer)
	//addNewCRValues(t, ctx, &ap, *addressPlanContainer)
}

//========================================================================================================
// generic functions
//========================================================================================================

func modifyExistingValues(t *testing.T, ctx *TestingContext, rt resourceType, crData Container) {
	phase := "Modify Existing CR Values"
	rt.copyRequiredValues(t, crData, phase)
	rt.changeCRValues(t, crData, phase)
	updateClusterCr(t, ctx, rt, crData, phase)
	results := compareResultsAfterReconcile(t, ctx, rt, crData, phase)
	if results != nil {
		for _, result := range *results {
			t.Logf("%s : %s: %s: %s: %s", phase, result.Type, result.Name, result.Key, result.Error)
		}
		t.Fatalf("%s : Failed to reset the CR values", phase)
	}
}

func deleteExistingValues(t *testing.T, ctx *TestingContext, rt resourceType, crData Container) {
	phase := "Delete Existing CR Values"
	rt.copyRequiredValues(t, crData, phase)
	rt.deleteExistingValues(t, crData, phase)
	updateClusterCr(t, ctx, rt, crData, phase)
	results := compareResultsAfterReconcile(t, ctx, rt, crData, phase)
	if results != nil {
		for _, result := range *results {
			t.Logf("%s : %s: %s: %s: %s", phase, result.Type, result.Name, result.Key, result.Error)
		}
		t.Fatalf("%s : Failed to reset the CR values", phase)
	}
}

//func addNewCRValues(t *testing.T, ctx *TestingContext, rt resourceType, crData Container) {
//	phase := "Adding New CR Values"
//	rt.addCRDummyValues(t, crData, phase)
//	updateClusterCr(t, ctx, rt, crData, phase)
//	compareAddedResultsAfterReconcile(t, ctx, rt, crData, phase)
//}

func updateClusterCr(t *testing.T, ctx *TestingContext, rt resourceType, intContainer Container, phase string) {
	t.Log("Brake Point")
	//cr := intContainer.Get().(crInterface)
	//cr, ok := intContainer.Get().(crInterface)
	cr, ok := getCR(intContainer, rt)
	if !ok {
		t.Log(cr)
		t.Fatalf("%s : Unable to read CR from intContainer", phase)
	}

	retryCount := 3
	retry := true
	for retry {
		err := ctx.Client.Update(goctx.TODO(), cr.getCr())
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
}

func compareResultsAfterReconcile(t *testing.T, ctx *TestingContext, rt resourceType, intContainer Container, phase string) *[]compareResult2BR {
	var results *[]compareResult2BR
	retryCount := 3
	forceRetry := true

	cr, ok := getCR(intContainer, rt)
	//cr, ok := intContainer.Get().(crInterface)
	t.Log("Checking cr")
	if !ok {
		t.Log(cr)
		t.Fatalf("%s : Unable to read CR from intContainer", phase)
	}

	for forceRetry {
		// Force Retry is required to remove flaky test results after random updates
		err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.getName(), Namespace: cr.getNamespace()}, cr.getCr())
		if err != nil {
			t.Fatalf("%s : Fail to refresh the cr", phase)
		}

		t.Logf("%s : %s: count = %v, revison = %s", phase, cr.getName(), retryCount, cr.getResourceVersion())
		_, err = waitReconcilingCR(t, ctx, rt, intContainer)
		if err != nil {
			t.Fatalf("%s : %s: %s:, %s", phase, cr.getKind(), cr.getName(), err)
		}
		if len(intContainer) == 0 {
			intContainer.Put(cr)
		}
		results = rt.compareValues(t, intContainer, phase)

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

//func compareAddedResultsAfterReconcile(t *testing.T, ctx *TestingContext, rt resourceType, intContainer Container, phase string) {
//	//TODO refactor to be generic
//	cr, ok := intContainer.Get().(crInterface)
//	if !ok {
//		t.Log(cr)
//		t.Fatalf("%s : Unable to read CR from intContainer", phase)
//	}
//
//	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.getName(), Namespace: cr.getNamespace()}, cr.getCr())
//	if err != nil {
//		t.Fatalf("%s : Fail to refresh the cr", phase)
//	}
//
//	_, err = waitReconcilingCR(t, ctx, cr)
//	if err != nil {
//		t.Fatalf("%s : %s: %s:, %s",phase, cr.getKind(), cr.getName(), err)
//	}
//	if len(intContainer) == 0 {
//		intContainer.Put(cr)
//	}
//	rt.checkDummyValuesStillExist(t, intContainer, phase)
//}

//TODO refactor to be generic
func waitReconcilingCR(t *testing.T, ctx *TestingContext, rt resourceType, intContainer Container) (done bool, err error) {
	cr, ok := getCR(intContainer, rt)
	t.Log("Checking cr")
	if !ok {
		t.Log(cr)
		//t.Fatalf("%s : Unable to read CR from intContainer", phase)
	}

	resourceVersion := cr.getResourceVersion()
	err = wait.Poll(crRetryInterval, crTimeout, func() (done bool, err error) {
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.getName(), Namespace: cr.getNamespace()}, cr.getCr())
		if err != nil {
			return false, err
		}
		t.Log(cr.getResourceVersion())
		if resourceVersion != cr.getResourceVersion() {
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

func getCR(intContainer Container, rt resourceType) (crInterface, bool) {
	switch {
	case rt.crType() == "enmasse.AddressPlan":
		cr, ok := intContainer.Get().(AddressSpaceCr)
		return cr, ok
	default:
		return nil, false
	}
}
