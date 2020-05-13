package common

import (
	goctx "context"
	"fmt"
	enmasseadminv1beta1 "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/admin/v1beta1"
	enmassev1beta1 "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta1"
	enmasse "github.com/integr8ly/integreatly-operator/pkg/apis-products/enmasse/v1beta2"
	rbacv1 "k8s.io/api/rbac/v1"
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

type compareResult struct {
	Type  string
	Name  string
	Key   string
	Error string
}

func TestResetCRs(t *testing.T, ctx *TestingContext) {
	var wg sync.WaitGroup
	testAMQOnline(t, ctx, &wg)

	wg.Wait()
}

//========================================================================================================
// enmasse
//========================================================================================================

func testAMQOnline(t *testing.T, ctx *TestingContext, wg *sync.WaitGroup) {
	testAddressSpacePlan(wg, t, ctx)
	testAddressPlan(wg, t, ctx)
	testAuthenticationServiceCr(wg, t, ctx)
	testBrokeredInfraConfigCr(wg, t, ctx)
	testStandardInfraConfigCr(wg, t, ctx)
	testRoleCr(wg, t, ctx)
	testRoleBindingCr(wg, t, ctx)
}

//========================================================================================================
// enmasse rbacv1.RoleBinding
// There are some CR that are been skipped. I do not know where these get created
//========================================================================================================

type roleBindingCr struct {
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

func testRoleBindingCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext) {
	crList := &rbacv1.RoleBindingList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: amqOnline,
	}

	err := ctx.Client.List(goctx.TODO(), crList, listOpts)
	if err != nil {
		t.Fatal("rbacv1.RoleBinding : Failed to get a list of CR's from cluster: ", err)
	}
	var skipped []string
	for _, cr := range crList.Items {
		if cr.Name == "dedicated-admins-service-admin" {
			wg.Add(1)
			go setUpRoleBindingCr(wg, t, ctx, cr)
		} else {
			skipped = append(skipped, cr.Name)
		}
	}
	t.Logf("rbacv1.RoleBinding : The following CR's were skipped, %s", skipped)

}

func setUpRoleBindingCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext, cr rbacv1.RoleBinding) {
	defer wg.Done()
	i := roleBindingCr{}
	i.runTests(t, ctx, cr)
}

func (i *roleBindingCr) runTests(t *testing.T, ctx *TestingContext, cr rbacv1.RoleBinding) {
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

func (i *roleBindingCr) modifyExistingValues(t *testing.T, ctx *TestingContext, cr rbacv1.RoleBinding) {
	i.copyRequiredValues(cr)
	i.changeCRValues(&cr)
	err := ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Modify Existing CR values : Failed to update CR on cluster: ", err)
	}

	var results *[]compareResult
	count := 3
	forceRetry := true
	// Force Retry is required to remove flaky test results after random updates
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

func (i *roleBindingCr) deleteExistingValues(t *testing.T, ctx *TestingContext, cr rbacv1.RoleBinding) {
	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
	if err != nil {
		t.Fatal("Deleting CR Values : Failed to refresh CR")
	}
	i.copyRequiredValues(cr)
	i.deleteCRValues(&cr)
	err = ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Log(err)
		t.Fatal("Deleting CR Values : Failed to update CR on cluster: ", err)
	}

	var results *[]compareResult
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

func (i *roleBindingCr) addNewValues(t *testing.T, ctx *TestingContext, cr rbacv1.RoleBinding) {
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

func (i *roleBindingCr) copyRequiredValues(cr rbacv1.RoleBinding) {
	ant := cr.GetAnnotations()
	i.IntegreatlyName = ant[integreatlyName]
	i.IntegreatlyNamespace = ant[integreatlyNamespace]
	i.RoleRefKind = cr.RoleRef.Kind
	i.RoleRefName = cr.RoleRef.Name
	for _, subject := range cr.Subjects {
		i.Subjects = append(i.Subjects, roleBindingSubject{
			SubjectName: subject.Name,
			SubjectKind: subject.Kind,
		})
	}

}

func (i *roleBindingCr) changeCRValues(cr *rbacv1.RoleBinding) {
	ant := cr.GetAnnotations()
	if ant == nil {
		ant = map[string]string{}
	}
	ant[integreatlyName] = "Bad Value"
	ant[integreatlyNamespace] = "Bad Value"
	cr.SetAnnotations(ant)
	//TODO Find a Role Kind that is allowed and is not Kind: Role
	//cr.RoleRef.Kind = "Bad Value"
	// Can not change role reference
	//cr.RoleRef.Name = "bad-value"
	for index := range cr.Subjects {
		cr.Subjects[index].Name = "bad-value"
		cr.Subjects[index].Kind = "ServiceAccount"
		cr.Subjects[index].APIGroup = ""
	}
}

func (i *roleBindingCr) waitReconcilingCR(ctx *TestingContext, cr rbacv1.RoleBinding) (done bool, err error) {
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

func (i *roleBindingCr) compareValues(cr *rbacv1.RoleBinding) *[]compareResult {
	var values []compareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-namespace",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyNamespace], i.IntegreatlyNamespace),
		})
	}

	for _, subject := range cr.Subjects {
		err := i.compareSubjectName(subject.Name)
		if err != nil {
			values = append(values, compareResult{
				Type:  cr.Kind,
				Name:  cr.Name,
				Key:   "subjects.[].name",
				Error: err.Error(),
			})
		}

		err = i.compareSubjectKind(subject.Kind)
		if err != nil {
			values = append(values, compareResult{
				Type:  cr.Kind,
				Name:  cr.Name,
				Key:   "subjects.[].Kind",
				Error: err.Error(),
			})
		}
	}

	if len(values) > 0 {
		return &values
	} else {
		return nil
	}
}

func (i *roleBindingCr) compareSubjectKind(value string) error {
	for _, item := range i.Subjects {
		if value == item.SubjectKind {
			return nil
		}
	}
	return fmt.Errorf("could not find %s in copied CR Subject.Kind", value)
}

func (i *roleBindingCr) compareSubjectName(value string) error {
	for _, item := range i.Subjects {
		if value == item.SubjectName {
			return nil
		}
	}
	return fmt.Errorf("could not find %s in copied CR Subject.Name", value)
}

func (i *roleBindingCr) deleteCRValues(cr *rbacv1.RoleBinding) {
	ant := cr.GetAnnotations()
	delete(ant, integreatlyName)
	delete(ant, integreatlyNamespace)
	cr.SetAnnotations(ant)
	cr.Subjects = nil
	//cr.RoleRef = rbacv1.RoleRef{}
}

func (i *roleBindingCr) addCRValue(cr rbacv1.RoleBinding) {
	ant := cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)
}

func (i *roleBindingCr) addedValuesStillExist(t *testing.T, cr rbacv1.RoleBinding) {
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}
}

//========================================================================================================
// enmasse rbacv1.Role
// There are some CR that are been skipped. I do not know where these get created
//========================================================================================================

type roleCr struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
	Roles                []roleCrRole
}

type roleCrRole struct {
	APIGroup  []string
	Resources []string
	Verbs     []string
}

func testRoleCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext) {
	crList := &rbacv1.RoleList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: amqOnline,
	}

	err := ctx.Client.List(goctx.TODO(), crList, listOpts)
	if err != nil {
		t.Fatal("rbacv1.Role : Failed to get a list of CR's from cluster: ", err)
	}
	var skipped []string
	for _, cr := range crList.Items {
		if cr.Name == "enmasse.io:service-admin" {
			wg.Add(1)
			go setUpRoleCr(wg, t, ctx, cr)
		}
		skipped = append(skipped, cr.Name)
	}
	t.Logf("rbacv1.Role : Skipping CR with name %s", skipped)

}

func setUpRoleCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
	defer wg.Done()
	i := roleCr{}
	i.runTests(t, ctx, cr)
}

func (i *roleCr) runTests(t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
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

func (i *roleCr) modifyExistingValues(t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
	i.copyRequiredValues(cr)
	i.changeCRValues(cr)
	err := ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Modify Existing CR values : Modify Existing CR values : Failed to update CR on cluster")
	}

	var results *[]compareResult
	count := 3
	forceRetry := true
	// Force Retry is required to remove flaky test results after random updates
	for forceRetry {
		err = ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
		if err != nil {
			t.Fatalf("Modify Existing CR values : Modify Existing CR values : Fail to refresh the cr")
		}

		t.Logf("Modify Existing CR values : Modify Existing CR values : %s: count = %v, revison = %s", cr.Name, count, cr.ResourceVersion)
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

func (i *roleCr) deleteExistingValues(t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
	err := ctx.Client.Get(goctx.TODO(), k8sclient.ObjectKey{Name: cr.Name, Namespace: cr.Namespace}, &cr)
	if err != nil {
		t.Fatal("Deleting CR Values : Failed to refresh CR")
	}
	i.copyRequiredValues(cr)
	i.deleteCRValues(&cr)
	err = ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Log(err)
		t.Fatal("Deleting CR Values : Failed to update CR on cluster")
	}

	var results *[]compareResult
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

func (i *roleCr) addNewValues(t *testing.T, ctx *TestingContext, cr rbacv1.Role) {
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

func (i *roleCr) copyRequiredValues(cr rbacv1.Role) {
	ant := cr.GetAnnotations()
	i.IntegreatlyName = ant[integreatlyName]
	i.IntegreatlyNamespace = ant[integreatlyNamespace]

	for _, rule := range cr.Rules {
		i.Roles = append(i.Roles, roleCrRole{
			APIGroup:  rule.APIGroups,
			Resources: rule.Resources,
			Verbs:     rule.Verbs,
		})
	}

}

func (i *roleCr) changeCRValues(cr rbacv1.Role) {
	ant := cr.GetAnnotations()
	if ant == nil {
		ant = map[string]string{}
	}
	ant[integreatlyName] = "Bad Value"
	ant[integreatlyNamespace] = "Bad Value"
	cr.SetAnnotations(ant)

	for index, rule := range cr.Rules {
		for i := range rule.Resources {
			cr.Rules[index].Resources[i] = "Bad Value"
		}

		for i := range rule.Verbs {
			cr.Rules[index].Verbs[i] = "Bad Value"
		}

		for i := range rule.APIGroups {
			cr.Rules[index].APIGroups[i] = "Bad Value"
		}
	}
}

func (i *roleCr) waitReconcilingCR(ctx *TestingContext, cr rbacv1.Role) (done bool, err error) {
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

func (i *roleCr) compareValues(cr *rbacv1.Role) *[]compareResult {
	var values []compareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-namespace",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyNamespace], i.IntegreatlyNamespace),
		})
	}

	for _, rule := range cr.Rules {
		for _, value := range rule.Resources {
			err := i.compareResources(value)
			if err != nil {
				values = append(values, compareResult{
					Type:  cr.Kind,
					Name:  cr.Name,
					Key:   "Roles.Resources",
					Error: err.Error(),
				})
			}
		}

		for _, value := range rule.Verbs {
			err := i.compareVerbs(value)
			if err != nil {
				values = append(values, compareResult{
					Type:  cr.Kind,
					Name:  cr.Name,
					Key:   "Roles.Verbs",
					Error: err.Error(),
				})
			}
		}

		for _, value := range rule.APIGroups {
			err := i.compareAPIGroups(value)
			if err != nil {
				values = append(values, compareResult{
					Type:  cr.Kind,
					Name:  cr.Name,
					Key:   "Roles.APIGroup",
					Error: err.Error(),
				})
			}
		}
	}

	if len(values) > 0 {
		return &values
	} else {
		return nil
	}
}

func (i *roleCr) compareAPIGroups(value string) error {
	for _, item := range i.Roles {
		for _, expected := range item.APIGroup {
			if value == expected {
				return nil
			}
		}
	}
	return fmt.Errorf("could not find %s in copied CR Roles.APIGroup", value)
}

func (i *roleCr) compareVerbs(value string) error {
	for _, item := range i.Roles {
		for _, expected := range item.Verbs {
			if value == expected {
				return nil
			}
		}
	}
	return fmt.Errorf("could not find %s in copied CR Roles.Verbs", value)
}

func (i *roleCr) compareResources(value string) error {
	for _, item := range i.Roles {
		for _, expected := range item.Resources {
			if value == expected {
				return nil
			}
		}
	}
	return fmt.Errorf("could not find %s in copied CR Roles.Resources", value)
}

func (i *roleCr) deleteCRValues(cr *rbacv1.Role) {
	ant := cr.GetAnnotations()
	delete(ant, integreatlyName)
	delete(ant, integreatlyNamespace)
	cr.SetAnnotations(ant)
	cr.Rules = nil
}

func (i *roleCr) addCRValue(cr rbacv1.Role) {
	ant := cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)
}

func (i *roleCr) addedValuesStillExist(t *testing.T, cr rbacv1.Role) {
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values go reset.")
	}
}

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

	var results *[]compareResult
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

	var results *[]compareResult
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

func (i *standardInfraConfig) compareValues(cr *enmassev1beta1.StandardInfraConfig) *[]compareResult {
	var values []compareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, compareResult{
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

	var results *[]compareResult
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

	var results *[]compareResult
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

func (i *brokeredInfraConfig) compareValues(cr *enmassev1beta1.BrokeredInfraConfig) *[]compareResult {
	var values []compareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, compareResult{
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

//========================================================================================================
// enmasse enmasseadminv1beta1 AuthenticationService
//========================================================================================================
const (
	NoneAuthservice     = "none-authservice"
	StandardAuthservice = "standard-authservice"
)

type authenticationService struct {
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

func testAuthenticationServiceCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext) {
	crList := &enmasseadminv1beta1.AuthenticationServiceList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: amqOnline,
	}
	err := ctx.Client.List(goctx.TODO(), crList, listOpts)
	if err != nil {
		t.Fatal("enmasseadminv1beta1.AuthenticationService : Failed to get a list of address space plan CR's from cluster")
	}

	for _, cr := range crList.Items {
		wg.Add(1)
		go setUpAuthenticationServiceCr(wg, t, ctx, cr)
	}

}

func setUpAuthenticationServiceCr(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext, cr enmasseadminv1beta1.AuthenticationService) {
	defer wg.Done()
	as := authenticationService{}
	as.runTests(t, ctx, cr)
}

func (i *authenticationService) runTests(t *testing.T, ctx *TestingContext, cr enmasseadminv1beta1.AuthenticationService) {
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

func (i *authenticationService) modifyExistingValues(t *testing.T, ctx *TestingContext, cr enmasseadminv1beta1.AuthenticationService) {
	i.copyRequiredValues(cr)
	i.changeCRValues(cr)
	err := ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Modify Existing CR values : Failed to update CR on cluster")
	}

	var results *[]compareResult
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

func (i *authenticationService) deleteExistingValues(t *testing.T, ctx *TestingContext, cr enmasseadminv1beta1.AuthenticationService) {
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

	var results *[]compareResult
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

func (i *authenticationService) addNewValues(t *testing.T, ctx *TestingContext, cr enmasseadminv1beta1.AuthenticationService) {
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

func (i *authenticationService) copyRequiredValues(cr enmasseadminv1beta1.AuthenticationService) {
	ant := cr.GetAnnotations()
	i.IntegreatlyName = ant[integreatlyName]
	i.IntegreatlyNamespace = ant[integreatlyNamespace]
	switch cr.Name {
	case NoneAuthservice:
		i.SpecType = cr.Spec.Type
	case StandardAuthservice:
		i.CredentialsSecretName = cr.Spec.Standard.Datasource.CredentialsSecret.Name
		i.CredentialsSecretNamespace = cr.Spec.Standard.Datasource.CredentialsSecret.Namespace
		i.DatasourceType = cr.Spec.Standard.Datasource.Type
		i.DatasourceDatabase = cr.Spec.Standard.Datasource.Database
		i.DatasourceHost = cr.Spec.Standard.Datasource.Host
		i.DatasourcePort = cr.Spec.Standard.Datasource.Port
	}
}

func (i *authenticationService) changeCRValues(cr enmasseadminv1beta1.AuthenticationService) {
	ant := cr.GetAnnotations()
	ant[integreatlyName] = "Bad Value"
	ant[integreatlyNamespace] = "Bad Value"
	cr.SetAnnotations(ant)

	switch cr.Name {
	case NoneAuthservice:
		cr.Spec.Type = "standard"
	case StandardAuthservice:
		cr.Spec.Type = "none"
		cr.Spec.Standard.Datasource.CredentialsSecret.Name = "bad value"
		cr.Spec.Standard.Datasource.CredentialsSecret.Namespace = "bad value"
		cr.Spec.Standard.Datasource.Type = "bad value"
		cr.Spec.Standard.Datasource.Database = "bad value"
		cr.Spec.Standard.Datasource.Host = "bad value"
		cr.Spec.Standard.Datasource.Port = 0
	}
}

func (i *authenticationService) waitReconcilingCR(ctx *TestingContext, cr enmasseadminv1beta1.AuthenticationService) (done bool, err error) {
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

func (i *authenticationService) compareValues(cr *enmasseadminv1beta1.AuthenticationService) *[]compareResult {
	var values []compareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, compareResult{
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

func (i *authenticationService) deleteCRValues(cr enmasseadminv1beta1.AuthenticationService) {
	ant := cr.GetAnnotations()
	delete(ant, integreatlyName)
	delete(ant, integreatlyNamespace)
	cr.SetAnnotations(ant)
	//TODO unable to delete cr.Spec.Type, do not know how to
}

func (i *authenticationService) addCRValue(cr enmasseadminv1beta1.AuthenticationService) {
	ant := cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)
}

func (i *authenticationService) addedValuesStillExist(t *testing.T, cr enmasseadminv1beta1.AuthenticationService) {
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}
}

//========================================================================================================
// enmasse addressSpacePlan
//========================================================================================================
type addressSpacePlan struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
}

func testAddressSpacePlan(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext) {
	aspl := &enmasse.AddressSpacePlanList{}
	listOpts := &k8sclient.ListOptions{
		Namespace: amqOnline,
	}
	err := ctx.Client.List(goctx.TODO(), aspl, listOpts)
	if err != nil {
		t.Fatal("addressSpacePlan : Failed to get a list of address space plan CR's from cluster")
	}

	for _, cr := range aspl.Items {
		wg.Add(1)
		go setUpAddressSpacePlan(wg, t, ctx, cr)
	}

}

func setUpAddressSpacePlan(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext, cr enmasse.AddressSpacePlan) {
	defer wg.Done()
	asp := addressSpacePlan{}
	asp.runTests(t, ctx, cr)
}

func (i *addressSpacePlan) runTests(t *testing.T, ctx *TestingContext, cr enmasse.AddressSpacePlan) {
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

func (i *addressSpacePlan) modifyExistingValues(t *testing.T, ctx *TestingContext, cr enmasse.AddressSpacePlan) {
	i.copyRequiredValues(cr)
	i.changeCRValues(cr)
	err := ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Modify Existing CR values : Failed to update CR on cluster")
	}

	var results *[]compareResult
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

func (i *addressSpacePlan) deleteExistingValues(t *testing.T, ctx *TestingContext, cr enmasse.AddressSpacePlan) {
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

	var results *[]compareResult
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

func (i *addressSpacePlan) addNewValues(t *testing.T, ctx *TestingContext, cr enmasse.AddressSpacePlan) {
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

func (i *addressSpacePlan) copyRequiredValues(cr enmasse.AddressSpacePlan) {
	ant := cr.GetAnnotations()
	i.IntegreatlyName = ant[integreatlyName]
	i.IntegreatlyNamespace = ant[integreatlyNamespace]
}

func (i *addressSpacePlan) changeCRValues(cr enmasse.AddressSpacePlan) {
	ant := cr.GetAnnotations()
	ant[integreatlyName] = "Bad Value"
	ant[integreatlyNamespace] = "Bad Value"
	cr.SetAnnotations(ant)
}

func (i *addressSpacePlan) waitReconcilingCR(ctx *TestingContext, cr enmasse.AddressSpacePlan) (done bool, err error) {
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

func (i *addressSpacePlan) compareValues(cr *enmasse.AddressSpacePlan) *[]compareResult {
	var values []compareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, compareResult{
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

func (i *addressSpacePlan) deleteCRValues(cr enmasse.AddressSpacePlan) {
	ant := cr.GetAnnotations()
	delete(ant, integreatlyName)
	delete(ant, integreatlyNamespace)
	cr.SetAnnotations(ant)
}

func (i *addressSpacePlan) addCRValue(cr enmasse.AddressSpacePlan) {
	ant := cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)
}

func (i *addressSpacePlan) addedValuesStillExist(t *testing.T, cr enmasse.AddressSpacePlan) {
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}
}

//========================================================================================================
// enmasse addressPlan
//========================================================================================================
type addressPlan struct {
	IntegreatlyName      string
	IntegreatlyNamespace string
}

func testAddressPlan(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext) {
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
		go setUpAddressPlan(wg, t, ctx, cr)
	}

}

func setUpAddressPlan(wg *sync.WaitGroup, t *testing.T, ctx *TestingContext, cr enmasse.AddressPlan) {
	defer wg.Done()
	ap := addressPlan{}
	ap.runTests(t, ctx, cr)
}

func (i *addressPlan) runTests(t *testing.T, ctx *TestingContext, cr enmasse.AddressPlan) {
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

func (i *addressPlan) modifyExistingValues(t *testing.T, ctx *TestingContext, cr enmasse.AddressPlan) {
	i.copyRequiredValues(cr)
	i.changeCRValues(cr)
	err := ctx.Client.Update(goctx.TODO(), &cr)
	if err != nil {
		t.Fatal("Modify Existing CR values : Failed to update CR on cluster")
	}

	var results *[]compareResult
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

func (i *addressPlan) deleteExistingValues(t *testing.T, ctx *TestingContext, cr enmasse.AddressPlan) {
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

	var results *[]compareResult
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

func (i *addressPlan) addNewValues(t *testing.T, ctx *TestingContext, cr enmasse.AddressPlan) {
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

func (i *addressPlan) waitReconcilingCR(ctx *TestingContext, cr enmasse.AddressPlan) (done bool, err error) {
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

func (i *addressPlan) compareValues(cr *enmasse.AddressPlan) *[]compareResult {
	var values []compareResult
	ant := cr.GetAnnotations()
	if ant[integreatlyName] != i.IntegreatlyName {
		values = append(values, compareResult{
			Type:  cr.Kind,
			Name:  cr.Name,
			Key:   "metadata.annotations.integreatly-name",
			Error: fmt.Sprintf("%s is not equal to expected %s", ant[integreatlyName], i.IntegreatlyName),
		})
	}

	if ant[integreatlyNamespace] != i.IntegreatlyNamespace {
		values = append(values, compareResult{
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

func (i *addressPlan) copyRequiredValues(cr enmasse.AddressPlan) {
	ant := cr.GetAnnotations()
	i.IntegreatlyName = ant[integreatlyName]
	i.IntegreatlyNamespace = ant[integreatlyNamespace]
}

func (i *addressPlan) changeCRValues(cr enmasse.AddressPlan) {
	ant := cr.GetAnnotations()
	ant[integreatlyName] = "Bad Value"
	ant[integreatlyNamespace] = "Bad Value"
	cr.SetAnnotations(ant)
}

func (i *addressPlan) deleteCRValues(cr enmasse.AddressPlan) {
	ant := cr.GetAnnotations()
	delete(ant, integreatlyName)
	delete(ant, integreatlyNamespace)
	cr.SetAnnotations(ant)
}

func (i *addressPlan) addCRValue(cr enmasse.AddressPlan) {
	ant := cr.GetAnnotations()
	ant["dummy-value"] = "dummy value"
	cr.SetAnnotations(ant)
}

func (i *addressPlan) addedValuesStillExist(t *testing.T, cr enmasse.AddressPlan) {
	ant := cr.GetAnnotations()
	if ant["dummy-value"] != "dummy value" {
		t.Fatal("Add New CR Values :  Added dummy values got reset.")
	}
}
