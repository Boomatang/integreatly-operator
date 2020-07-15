package common

import (
	goctx "context"
	"fmt"
	enmassev1beta1 "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta1"
	modify_crs "github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
	"testing"
	"time"
)

const (
	crRetryInterval      = time.Second * 5
	crTimeout            = time.Second * 60 * 3
	integreatlyName      = "integreatly-name"
	integreatlyNamespace = "integreatly-namespace"
	amqOnline            = "redhat-rhmi-amq-online"
	crFieldEdit          = true
	crFieldDelete        = true
	crFieldAdd           = true
)

func TestResetCRsold(t *testing.T, ctx *TestingContext) {
	var wg sync.WaitGroup
	testAMQOnline(t, ctx, &wg)

	wg.Wait()
}

//========================================================================================================
// enmasse
//========================================================================================================

func testAMQOnline(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	//testAddressSpacePlan(wg, t, ctx)
	//testAddressPlan(wg, t, ctx)
	//testAuthenticationServiceCr(wg, t, ctx)
	testBrokeredInfraConfigCr(wg, t, ctx)
	testStandardInfraConfigCr(wg, t, ctx)
	//testRoleCr(wg, t, ctx)
	//testRoleBindingCr(wg, t, ctx)
}

////========================================================================================================
//// enmasse rbacv1.Role
//// There are some CR that are been skipped. I do not know where these get created
////========================================================================================================
//
//type roleCr struct {
//	IntegreatlyName      string
//	IntegreatlyNamespace string
//	Roles                []roleCrRole
//}
//
//type roleCrRole struct {
//	APIGroup  []string
//	Resources []string
//	Verbs     []string
//}
//
//func testRoleCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext) {
//	crList := &rbacv1.RoleList{}
//	listOpts := &k8sclient.ListOptions{
//		Namespace: amqOnline,
//	}
//
//	err := ctx.Client.List(goctx.TODO(), crList, listOpts)
//	if err != nil {
//		t.Fatal("rbacv1.Role : Failed to get a list of CR's from cluster: ", err)
//	}
//	var skipped []string
//	for _, cr := range crList.Items {
//		if cr.Name == "enmasse.io:service-admin" {
//			wg.Add(1)
//			go setUpRoleCr(wg, t, ctx, cr)
//		}
//		skipped = append(skipped, cr.Name)
//	}
//	t.Logf("rbacv1.Role : Skipping CR with name %s", skipped)
//
//}
//
//func setUpRoleCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
//	defer wg.Done()
//	i := roleCr{}
//	i.runTests(t, ctx, cr)
//}
//
//func (i *roleCr) runTests(t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
//	if crFieldEdit {
//		i.modifyExistingValues(t, ctx, cr)
//	}
//	if crFieldDelete {
//		i.deleteExistingValues(t, ctx, cr)
//	}
//	if crFieldAdd {
//		i.addNewValues(t, ctx, cr)
//	}
//}
//
//func (i *roleCr) modifyExistingValues(t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
//	i.copyRequiredValues(cr)
//	i.changeCRValues(cr)
//	err := ctx.Client.Update(goctx.TODO(), &cr)
//	if err != nil {
//		t.Fatal("Modify Existing CR values : Modify Existing CR values : Failed to update CR on cluster")
//	}
//
//	var results *[]modify_crs.CompareResult
//	count := 3
//	forceRetry := true
//	// Force Retry is required to remove flaky test results after random updates
//	for forceRetry {
//		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
//		if err != nil {
//			t.Fatalf("Modify Existing CR values : Modify Existing CR values : Fail to refresh the cr")
//		}
//
//		t.Logf("Modify Existing CR values : Modify Existing CR values : %s: count = %v, revison = %s", cr.Name, count, cr.ResourceVersion)
//		_, err = i.waitReconcilingCR(ctx, cr)
//		if err != nil {
//			t.Fatalf("Modify Existing CR values : %s: %s:, %s", cr.Kind, cr.Name, err)
//		}
//		results = i.compareValues(&cr)
//
//		if results == nil {
//			forceRetry = false
//			count -= 1
//		}
//		count -= 1
//		if count < 0 {
//			forceRetry = false
//		}
//	}
//
//	if results != nil {
//		for _, result := range *results {
//			t.Logf("Modify Existing CR values : %s: %s: %s: %s", result.Type, result.Name, result.Key, result.Error)
//		}
//		t.Fatal("Modify Existing CR values : Failed to reset the CR values")
//	}
//}
//
//func (i *roleCr) deleteExistingValues(t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
//	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
//	if err != nil {
//		t.Fatal("Deleting CR Values : Failed to refresh CR")
//	}
//	i.copyRequiredValues(cr)
//	i.deleteCRValues(&cr)
//	err = ctx.Client.Update(goctx.TODO(), &cr)
//	if err != nil {
//		t.Log(err)
//		t.Fatal("Deleting CR Values : Failed to update CR on cluster")
//	}
//
//	var results *[]modify_crs.CompareResult
//	count := 3
//	forceRetry := true
//	// Force Retry is required to remove flaky test results after random updates
//	for forceRetry {
//		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
//		if err != nil {
//			t.Fatalf("Deleting CR Values : Fail to refresh the cr")
//		}
//
//		t.Logf("Deleting CR Values : %s: count = %v, revison = %s", cr.Name, count, cr.ResourceVersion)
//		_, err = i.waitReconcilingCR(ctx, cr)
//		if err != nil {
//			t.Fatalf("Deleting CR Values : %s: %s:, %s", cr.Kind, cr.Name, err)
//		}
//		results = i.compareValues(&cr)
//
//		if results == nil {
//			forceRetry = false
//			count -= 1
//		}
//		count -= 1
//		if count < 0 {
//			forceRetry = false
//		}
//	}
//
//	if results != nil {
//		for _, result := range *results {
//			t.Logf("Deleting CR Values : %s: %s: %s: %s", result.Type, result.Name, result.Key, result.Error)
//		}
//		t.Fatal("Deleting CR Values : Failed to reset the CR values")
//	}
//}
//
//func (i *roleCr) addNewValues(t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
//	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
//	if err != nil {
//		t.Fatal("Add New CR Values :  Failed to refresh CR")
//	}
//	i.addCRValue(cr)
//	err = ctx.Client.Update(goctx.TODO(), &cr)
//	if err != nil {
//		t.Fatal("Add New CR Values :  Failed to update CR on cluster")
//	}
//
//	// Refresh CR to get up-to-date version number
//	err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
//	if err != nil {
//		t.Fatalf("Add New CR Values :  Fail to refresh the cr")
//	}
//
//	_, err = i.waitReconcilingCR(ctx, cr)
//	if err != nil && err.Error() != "timed out waiting for the condition" {
//		t.Fatal(err)
//	} else {
//		i.addedValuesStillExist(t, cr)
//	}
//}
//
//func (i *roleCr) copyRequiredValues(cr rbacv1.Role) {
//	ant := cr.GetAnnotations()
//	i.IntegreatlyName = ant[integreatlyName]
//	i.IntegreatlyNamespace = ant[integreatlyNamespace]
//
//	for _, rule := range cr.Rules {
//		i.Roles = append(i.Roles, roleCrRole{
//			APIGroup:  rule.APIGroups,
//			Resources: rule.Resources,
//			Verbs:     rule.Verbs,
//		})
//	}
//
//}
//
//func (i *roleCr) changeCRValues(cr rbacv1.Role) {
//	ant := cr.GetAnnotations()
//	if ant == nil {
//		ant = map[string]string{}
//	}
//	ant[integreatlyName] = "Bad Value"
//	ant[integreatlyNamespace] = "Bad Value"
//	cr.SetAnnotations(ant)
//
//	for index, rule := range cr.Rules {
//		for i := range rule.Resources {
//			cr.Rules[index].Resources[i] = "Bad Value"
//		}
//
//		for i := range rule.Verbs {
//			cr.Rules[index].Verbs[i] = "Bad Value"
//		}
//
//		for i := range rule.APIGroups {
//			cr.Rules[index].APIGroups[i] = "Bad Value"
//		}
//	}
//}
//
//func (i *roleCr) waitReconcilingCR(ctx *TestingContext, cr rbacv1.Role) (done bool, err error) {
//	resourceVersion := cr.ResourceVersion
//	err = wait.Poll(crRetryInterval, crTimeout, func() (done bool, err error) {
//		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
//		if err != nil {
//			return false, err
//		}
//
//		if resourceVersion != cr.ResourceVersion {
//			return true, nil
//		} else {
//			return false, nil
//		}
//	})
//	if err != nil {
//		return false, err
//	} else {
//		return true, nil
//	}
//}
//
//func (i *roleCr) compareValues(cr *rbacv1.Role) *[]modify_crs.CompareResult {
//	var values []modify_crs.CompareResult
//	ant := cr.GetAnnotations()
//	if ant[integreatlyName] != i.IntegreatlyName {
//		values = append(values, modify_crs.CompareResult{
//			Type:  cr.Kind,
//			Name:  cr.Name,
//			Key:   "metadata.annotations.integreatly-name",
//			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
//		})
//	}
//
//	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
//		values = append(values, modify_crs.CompareResult{
//			Type:  cr.Kind,
//			Name:  cr.Name,
//			Key:   "metadata.annotations.integreatly-namespace",
//			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyNamespace], i.IntegreatlyNamespace),
//		})
//	}
//
//	for _, rule := range cr.Rules {
//		for _, value := range rule.Resources {
//			err := i.compareResources(value)
//			if err != nil {
//				values = append(values, modify_crs.CompareResult{
//					Type:  cr.Kind,
//					Name:  cr.Name,
//					Key:   "Roles.Resources",
//					Error: err.Error(),
//				})
//			}
//		}
//
//		for _, value := range rule.Verbs {
//			err := i.compareVerbs(value)
//			if err != nil {
//				values = append(values, modify_crs.CompareResult{
//					Type:  cr.Kind,
//					Name:  cr.Name,
//					Key:   "Roles.Verbs",
//					Error: err.Error(),
//				})
//			}
//		}
//
//		for _, value := range rule.APIGroups {
//			err := i.compareAPIGroups(value)
//			if err != nil {
//				values = append(values, modify_crs.CompareResult{
//					Type:  cr.Kind,
//					Name:  cr.Name,
//					Key:   "Roles.APIGroup",
//					Error: err.Error(),
//				})
//			}
//		}
//	}
//
//	if len(values) > 0 {
//		return &values
//	} else {
//		return nil
//	}
//}
//
//func (i *roleCr) compareAPIGroups(value string) error {
//	for _, item := range i.Roles {
//		for _, expected := range item.APIGroup {
//			if value == expected {
//				return nil
//			}
//		}
//	}
//	return fmt.Errorf("could not find %s in copied CR Roles.APIGroup", value)
//}
//
//func (i *roleCr) compareVerbs(value string) error {
//	for _, item := range i.Roles {
//		for _, expected := range item.Verbs {
//			if value == expected {
//				return nil
//			}
//		}
//	}
//	return fmt.Errorf("could not find %s in copied CR Roles.Verbs", value)
//}
//
//func (i *roleCr) compareResources(value string) error {
//	for _, item := range i.Roles {
//		for _, expected := range item.Resources {
//			if value == expected {
//				return nil
//			}
//		}
//	}
//	return fmt.Errorf("could not find %s in copied CR Roles.Resources", value)
//}
//
//func (i *roleCr) deleteCRValues(cr *rbacv1.Role) {
//	ant := cr.GetAnnotations()
//	delete(ant, integreatlyName)
//	delete(ant, integreatlyNamespace)
//	cr.SetAnnotations(ant)
//	cr.Rules = nil
//}
//
//func (i *roleCr) addCRValue(cr rbacv1.Role) {
//	ant := cr.GetAnnotations()
//	ant["dummy-value"] = "dummy value"
//	cr.SetAnnotations(ant)
//}
//
//func (i *roleCr) addedValuesStillExist(t *testing.T, cr rbacv1.Role) {
//	ant := cr.GetAnnotations()
//	if ant["dummy-value"] != "dummy value" {
//		t.Fatal("Add New CR Values :  Added dummy values go reset.")
//	}
//}

//========================================================================================================
// enmasse enmassev1beta1. StandardInfraConfig
//========================================================================================================

type standardInfraConfig struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
}

func testStandardInfraConfigCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext) {
	crList := &enmassev1beta1.StandardInfraConfigList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: amqOnline,
	}

	err := ctx.Client.List(goctx.TODO(), crList, listOpts)
	if err != nil {
		t.Fatal("enmassev1beta1.StandardInfraConfig : Failed to get a list of CR's from cluster: ", err)
	}

	for _, cr := range crList.Items {
		wg.Add(1)
		go setUpStandardInfraConfigCr(wg, t, ctx, cr)
	}

}

func setUpStandardInfraConfigCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext, cr enmassev1beta1.StandardInfraConfig) {
	defer wg.Done()
	i := standardInfraConfig{}
	i.runTests(t, ctx, cr)
}

func (i *standardInfraConfig) runTests(t *testing.T, ctx *TestingContext, cr enmassev1beta1.StandardInfraConfig) {
	if crFieldEdit {
		i.modifyExistingValues(t, ctx, cr)
	}
	if crFieldDelete {
		i.deleteExistingValues(t, ctx, cr)
	}
	if crFieldAdd {
		i.addNewValues(t, ctx, cr)
	}
}

func (i *standardInfraConfig) modifyExistingValues(t *testing.T, ctx *TestingContext, cr enmassev1beta1.StandardInfraConfig) {
	i.copyRequiredValues(cr)
	i.changeCRValues(cr)
	err := ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Modify Existing CR values : Failed to update CR on cluster")
	}

	var results *[]modify_crs.CompareResult
	count := 3
	forceRetry := true
	for forceRetry {
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
		if err != nil {
			t.Fatalf("Modify Existing CR values : Fail to refresh the cr")
		}
		t.Logf("Modify Existing CR values : %s: count = %v, revison = %s", cr.Name, count, cr.ResourceVersion)
		_, err = i.waitReconcilingCR(ctx, cr)
		if err != nil {
			t.Fatalf("Modify Existing CR values : %s: %s:, %s", cr.Kind, cr.Name, err)
		}
		results = i.compareValues(&cr)

		if results == nil {
			forceRetry = false
			count -= 1
		}
		count -= 1
		if count < 0 {
			forceRetry = false
		}
	}

	if results != nil {
		for _, result := range *results {
			t.Logf("Modify Existing CR values : %s: %s: %s: %s", result.Type, result.Name, result.Key, result.Error)
		}
		t.Fatal("Modify Existing CR values : Failed to reset the CR values")
	}
}

func (i *standardInfraConfig) deleteExistingValues(t *testing.T, ctx *TestingContext, cr enmassev1beta1.StandardInfraConfig) {
	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
	if err != nil {
		t.Fatal("Deleting CR Values : Failed to refresh CR")
	}
	i.copyRequiredValues(cr)
	i.deleteCRValues(cr)
	err = ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Log(err)
		t.Fatal("Deleting CR Values : Failed to update CR on cluster")
	}

	var results *[]modify_crs.CompareResult
	count := 3
	forceRetry := true
	// Force Retry is required to remove flaky test results after random updates
	for forceRetry {
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
		if err != nil {
			t.Fatalf("Deleting CR Values : Fail to refresh the cr")
		}

		t.Logf("Deleting CR Values : %s: count = %v, revison = %s", cr.Name, count, cr.ResourceVersion)
		_, err = i.waitReconcilingCR(ctx, cr)
		if err != nil {
			t.Fatalf("Deleting CR Values : %s: %s:, %s", cr.Kind, cr.Name, err)
		}
		results = i.compareValues(&cr)

		if results == nil {
			forceRetry = false
			count -= 1
		}
		count -= 1
		if count < 0 {
			forceRetry = false
		}
	}

	if results != nil {
		for _, result := range *results {
			t.Logf("Deleting CR Values : %s: %s: %s: %s", result.Type, result.Name, result.Key, result.Error)
		}
		t.Fatal("Deleting CR Values : Failed to reset the CR values")
	}
}

func (i *standardInfraConfig) addNewValues(t *testing.T, ctx *TestingContext, cr enmassev1beta1.StandardInfraConfig) {
	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
	if err != nil {
		t.Fatal("Add New CR Values :  Failed to refresh CR")
	}
	i.addCRValue(cr)
	err = ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Add New CR Values :  Failed to update CR on cluster")
	}

	// Refresh CR to get up-to-date version number
	err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
	if err != nil {
		t.Fatalf("Add New CR Values :  Fail to refresh the cr")
	}

	_, err = i.waitReconcilingCR(ctx, cr)
	if err != nil && err.Error() != "timed out waiting for the condition" {
		t.Fatal(err)
	} else {
		i.addedValuesStillExist(t, cr)
	}
}

func (i *standardInfraConfig) copyRequiredValues(cr enmassev1beta1.StandardInfraConfig) {
	ant := cr.GetAnnotations()
	i.IntegreatlyName = ant[integreatlyName]
	i.IntegreatlyNamespace = ant[integreatlyNamespace]
}

func (i *standardInfraConfig) changeCRValues(cr enmassev1beta1.StandardInfraConfig) {
	ant := cr.GetAnnotations()
	ant[integreatlyName] = "Bad Value"
	ant[integreatlyNamespace] = "Bad Value"
	cr.SetAnnotations(ant)
}

func (i *standardInfraConfig) waitReconcilingCR(ctx *TestingContext, cr enmassev1beta1.StandardInfraConfig) (done bool, err error) {
	resourceVersion := cr.ResourceVersion
	err = wait.Poll(crRetryInterval, crTimeout, func() (done bool, err error) {
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
		if err != nil {
			return false, err
		}

		if resourceVersion != cr.ResourceVersion {
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

func (i *standardInfraConfig) compareValues(cr *enmassev1beta1.StandardInfraConfig) *[]modify_crs.CompareResult {
	var values []modify_crs.CompareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, modify_crs.CompareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, modify_crs.CompareResult{
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

func (i *standardInfraConfig) deleteCRValues(cr enmassev1beta1.StandardInfraConfig) {
	ant := cr.GetAnnotations()
	delete(ant, integreatlyName)
	delete(ant, integreatlyNamespace)
	cr.SetAnnotations(ant)
	//TODO unable to delete cr.Spec.Type, do not know how to
}

func (i *standardInfraConfig) addCRValue(cr enmassev1beta1.StandardInfraConfig) {
	ant := cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)
}

func (i *standardInfraConfig) addedValuesStillExist(t *testing.T, cr enmassev1beta1.StandardInfraConfig) {
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}
}

//========================================================================================================
// enmasse enmassev1beta1.BrokeredInfraConfigList
//========================================================================================================

type brokeredInfraConfig struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
}

func testBrokeredInfraConfigCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext) {
	crList := &enmassev1beta1.BrokeredInfraConfigList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: amqOnline,
	}

	err := ctx.Client.List(goctx.TODO(), crList, listOpts)
	if err != nil {
		t.Fatal("enmassev1beta1.BrokeredInfraConfigList : Failed to get a list of CR's from cluster: ", err)
	}

	for _, cr := range crList.Items {
		wg.Add(1)
		go setUpBrokeredInfraConfigCr(wg, t, ctx, cr)
	}

}

func setUpBrokeredInfraConfigCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext, cr enmassev1beta1.BrokeredInfraConfig) {
	defer wg.Done()
	i := brokeredInfraConfig{}
	i.runTests(t, ctx, cr)
}

func (i *brokeredInfraConfig) runTests(t *testing.T, ctx *TestingContext, cr enmassev1beta1.BrokeredInfraConfig) {
	if crFieldEdit {
		i.modifyExistingValues(t, ctx, cr)
	}
	if crFieldDelete {
		i.deleteExistingValues(t, ctx, cr)
	}
	if crFieldAdd {
		i.addNewValues(t, ctx, cr)
	}
}

func (i *brokeredInfraConfig) modifyExistingValues(t *testing.T, ctx *TestingContext, cr enmassev1beta1.BrokeredInfraConfig) {
	i.copyRequiredValues(cr)
	i.changeCRValues(cr)
	err := ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Modify Existing CR values : Failed to update CR on cluster")
	}

	var results *[]modify_crs.CompareResult
	count := 3
	forceRetry := true
	for forceRetry {
		// Force Retry is required to remove flaky test results after random updates
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
		if err != nil {
			t.Fatalf("Modify Existing CR values : Fail to refresh the cr")
		}

		t.Logf("Modify Existing CR values : %s: count = %v, revison = %s", cr.Name, count, cr.ResourceVersion)
		_, err = i.waitReconcilingCR(ctx, cr)
		if err != nil {
			t.Fatalf("Modify Existing CR values : %s: %s:, %s", cr.Kind, cr.Name, err)
		}
		results = i.compareValues(&cr)

		if results == nil {
			forceRetry = false
			count -= 1
		}
		count -= 1
		if count < 0 {
			forceRetry = false
		}
	}

	if results != nil {
		for _, result := range *results {
			t.Logf("Modify Existing CR values : %s: %s: %s: %s", result.Type, result.Name, result.Key, result.Error)
		}
		t.Fatal("Modify Existing CR values : Failed to reset the CR values")
	}
}

func (i *brokeredInfraConfig) deleteExistingValues(t *testing.T, ctx *TestingContext, cr enmassev1beta1.BrokeredInfraConfig) {
	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
	if err != nil {
		t.Fatal("Deleting CR Values : Failed to refresh CR")
	}
	i.copyRequiredValues(cr)
	i.deleteCRValues(cr)
	err = ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Log(err)
		t.Fatal("Deleting CR Values : Failed to update CR on cluster")
	}

	var results *[]modify_crs.CompareResult
	count := 3
	forceRetry := true
	// Force Retry is required to remove flaky test results after random updates
	for forceRetry {
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
		if err != nil {
			t.Fatalf("Deleting CR Values : Fail to refresh the cr")
		}

		t.Logf("Deleting CR Values : %s: count = %v, revison = %s", cr.Name, count, cr.ResourceVersion)
		_, err = i.waitReconcilingCR(ctx, cr)
		if err != nil {
			t.Fatalf("Deleting CR Values : %s: %s:, %s", cr.Kind, cr.Name, err)
		}
		results = i.compareValues(&cr)

		if results == nil {
			forceRetry = false
			count -= 1
		}
		count -= 1
		if count < 0 {
			forceRetry = false
		}
	}

	if results != nil {
		for _, result := range *results {
			t.Logf("Deleting CR Values : %s: %s: %s: %s", result.Type, result.Name, result.Key, result.Error)
		}
		t.Fatal("Deleting CR Values : Failed to reset the CR values")
	}
}

func (i *brokeredInfraConfig) addNewValues(t *testing.T, ctx *TestingContext, cr enmassev1beta1.BrokeredInfraConfig) {
	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
	if err != nil {
		t.Fatal("Add New CR Values :  Failed to refresh CR")
	}
	i.addCRValue(cr)
	err = ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Add New CR Values :  Failed to update CR on cluster")
	}

	// Refresh CR to get up-to-date version number
	err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
	if err != nil {
		t.Fatalf("Add New CR Values :  Fail to refresh the cr")
	}

	_, err = i.waitReconcilingCR(ctx, cr)
	if err != nil && err.Error() != "timed out waiting for the condition" {
		t.Fatal(err)
	} else {
		i.addedValuesStillExist(t, cr)
	}
}

func (i *brokeredInfraConfig) copyRequiredValues(cr enmassev1beta1.BrokeredInfraConfig) {
	ant := cr.GetAnnotations()
	i.IntegreatlyName = ant[integreatlyName]
	i.IntegreatlyNamespace = ant[integreatlyNamespace]
}

func (i *brokeredInfraConfig) changeCRValues(cr enmassev1beta1.BrokeredInfraConfig) {
	ant := cr.GetAnnotations()
	ant[integreatlyName] = "Bad Value"
	ant[integreatlyNamespace] = "Bad Value"
	cr.SetAnnotations(ant)
}

func (i *brokeredInfraConfig) waitReconcilingCR(ctx *TestingContext, cr enmassev1beta1.BrokeredInfraConfig) (done bool, err error) {
	resourceVersion := cr.ResourceVersion
	err = wait.Poll(crRetryInterval, crTimeout, func() (done bool, err error) {
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
		if err != nil {
			return false, err
		}

		if resourceVersion != cr.ResourceVersion {
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

func (i *brokeredInfraConfig) compareValues(cr *enmassev1beta1.BrokeredInfraConfig) *[]modify_crs.CompareResult {
	var values []modify_crs.CompareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, modify_crs.CompareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, modify_crs.CompareResult{
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

func (i *brokeredInfraConfig) deleteCRValues(cr enmassev1beta1.BrokeredInfraConfig) {
	ant := cr.GetAnnotations()
	delete(ant, integreatlyName)
	delete(ant, integreatlyNamespace)
	cr.SetAnnotations(ant)
	//TODO unable to delete cr.Spec.Type, do not know how to
}

func (i *brokeredInfraConfig) addCRValue(cr enmassev1beta1.BrokeredInfraConfig) {
	ant := cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)
}

func (i *brokeredInfraConfig) addedValuesStillExist(t *testing.T, cr enmassev1beta1.BrokeredInfraConfig) {
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}
}
