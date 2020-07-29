package common

import "testing"

func TestWorkINProgress(t *testing.T, ctx *TestingContext) {
	t.Log("Work in progress")
	t.Log("1 Make sure that there is at least one active alert for every panel in SLO dasboard (e.g. pod_down)")

	t.Log("1.1 Make sure rhmi-operator pod is scaled down to 0 pods")
	t.Log("1.2 Make sure all rhmi product operator pods are scaled down to 0, you can use code #1 below")
	t.Log("1.3 Make sure all keycloak stateful sets are scaled down to 0, you can use code #2 and #3 below")
	t.Log("1.4 Make sure all product pods are scaled down to 0, you can use code #4 below")

	t.Log("2 Check the dashboard Critical SLO summary after some time (~20min)")

	t.Log("3 Bring back up the pods in the 3scale, CodeReady, Solution Explorer, and UPS namespaces")

	t.Fail()
}
