package common

import (
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"testing"
	"time"
)

const (
	local_delay = (1 * time.Hour) + (15 * time.Minute)
)

type a struct {
	ResourceName string
	NameSpace    string
	Replicas     int32
	Resource     string
	Kind         string
	ApiVersion   string
	Uri          string
}

var (
	local_namespaces = []a{
		{"rhmi-operator", "redhat-rhmi-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"3scale-operator", "redhat-rhmi-3scale-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"enmasse-operator", "redhat-rhmi-amq-online", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"apicurito-operator", "redhat-rhmi-apicurito-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"codeready-operator", "redhat-rhmi-codeready-workspaces-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"syndesis-operator", "redhat-rhmi-fuse-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"keycloak-operator", "redhat-rhmi-rhsso-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"tutorial-web-app-operator", "redhat-rhmi-solution-explorer-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"unifiedpush-operator", "redhat-rhmi-ups-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"keycloak-operator", "redhat-rhmi-user-sso-operator", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
	}

	local_stateful_sets = []a{
		{"keycloak", "redhat-rhmi-rhsso", 2, "StatefulSets", "StatefulSet", "apps/v1", "/apis/apps/v1"},
		{"keycloak", "redhat-rhmi-user-sso", 2, "StatefulSets", "StatefulSet", "apps/v1", "/apis/apps/v1"},
	}

	local_product_pods = []a{
		{"apicast-production", "redhat-rhmi-3scale", 2, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"apicast-staging", "redhat-rhmi-3scale", 2, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"backend-cron", "redhat-rhmi-3scale", 2, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"backend-listener", "redhat-rhmi-3scale", 2, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"backend-worker", "redhat-rhmi-3scale", 2, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"system-app", "redhat-rhmi-3scale", 2, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"system-memcache", "redhat-rhmi-3scale", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"system-sidekiq", "redhat-rhmi-3scale", 2, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"system-sphinx", "redhat-rhmi-3scale", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"zync", "redhat-rhmi-3scale", 2, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"zync-database", "redhat-rhmi-3scale", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"zync-que", "redhat-rhmi-3scale", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"address-space-controller", "redhat-rhmi-amq-online", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"console", "redhat-rhmi-amq-online", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"none-authservice", "redhat-rhmi-amq-online", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"standard-authservice", "redhat-rhmi-amq-online", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"apicurito", "redhat-rhmi-apicurito", 2, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"fuse-apicurito-generator", "redhat-rhmi-apicurito", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"codeready", "redhat-rhmi-codeready-workspaces", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"devfile-registry", "redhat-rhmi-codeready-workspaces", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"plugin-registry", "redhat-rhmi-codeready-workspaces", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
		{"broker-amq", "redhat-rhmi-fuse", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"syndesis-meta", "redhat-rhmi-fuse", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"syndesis-oauthproxy", "redhat-rhmi-fuse", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"syndesis-prometheus", "redhat-rhmi-fuse", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"syndesis-server", "redhat-rhmi-fuse", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"syndesis-ui", "redhat-rhmi-fuse", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"tutorial-web-app", "redhat-rhmi-solution-explorer", 1, "DeploymentConfigs", "DeploymentConfig", "apps.openshift.io/v1", "/apis/apps.openshift.io/v1"},
		{"ups", "redhat-rhmi-ups", 1, "Deployments", "Deployment", "apps/v1", "/apis/apps/v1"},
	}
)

func TestWorkINProgress(t *testing.T, ctx *TestingContext) {
	defer tear_down(ctx, t)
	setup(ctx, t)

	t.Log("2 Check the dashboard Critical SLO summary after some time (~20min)")

	time.Sleep(local_delay)

	t.Fail()
}

func setup(ctx *TestingContext, t *testing.T) {
	for _, space := range local_namespaces {
		scale(ctx, t, 0, space.ResourceName, space.NameSpace, space.Resource, space.Kind, space.ApiVersion, space.Uri)
	}
	for _, stateful_set := range local_stateful_sets {
		scale(ctx, t, 0, stateful_set.ResourceName, stateful_set.NameSpace, stateful_set.Resource, stateful_set.Kind, stateful_set.ApiVersion, stateful_set.Uri)
	}
	for _, pod := range local_product_pods {
		scale(ctx, t, 0, pod.ResourceName, pod.NameSpace, pod.Resource, pod.Kind, pod.ApiVersion, pod.Uri)
	}
}

func tear_down(ctx *TestingContext, t *testing.T) {
	t.Logf("Starting the tear down")
	for _, space := range local_namespaces {
		scale(ctx, t, space.Replicas, space.ResourceName, space.NameSpace, space.Resource, space.Kind, space.ApiVersion, space.Uri)
	}
	for _, stateful_set := range local_stateful_sets {
		scale(ctx, t, stateful_set.Replicas, stateful_set.ResourceName, stateful_set.NameSpace, stateful_set.Resource, stateful_set.Kind, stateful_set.ApiVersion, stateful_set.Uri)
	}
	for _, pod := range local_product_pods {
		scale(ctx, t, pod.Replicas, pod.ResourceName, pod.NameSpace, pod.Resource, pod.Kind, pod.ApiVersion, pod.Uri)
	}
}

func scale(ctx *TestingContext, t *testing.T, replicas int32, resourceName string, namesSpace string, resource string, kind string, apiVersion string, uri string) {

	t.Logf("Scaling %s in %s namespace", resourceName, namesSpace)
	replica := fmt.Sprintf(`{
		"apiVersion": "%s",
		"kind": "%s",
		"spec": {
			"replicas": %v
		}
	}`, apiVersion, kind, replicas)
	replicaBytes := []byte(replica)

	request := ctx.ExtensionClient.RESTClient().Patch(types.MergePatchType).
		Resource(resource).
		Name(resourceName).
		Namespace(namesSpace).
		RequestURI(uri).Body(replicaBytes).Do()
	_, err := request.Raw()
	if err != nil {
		t.Logf(err.Error())
	}
}
