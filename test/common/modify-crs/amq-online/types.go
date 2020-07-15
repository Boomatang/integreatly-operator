package amq_online

import (
	modify_crs "github.com/integr8ly/integreatly-operator/test/common/modify-crs"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	EnmasseAddressPlan           = "enmasse.AddressPlanReference"
	EnmasseAddressSpacePlan      = "enmasse.AddressSpacePlanReference"
	EnmasseStandardInfraConfig   = "enmasse.StandardInfraConfig"
	EnmasseAuthenticationService = "enmasseadminv1beta1.AuthenticationService"
	Rbacv1RoleBinding            = "rbacv1.RoleBinding"
	Rbacv1Role                   = "rbacv1.Role"
)

var ListOpts = &k8sclient.ListOptions{
	Namespace: modify_crs.AmqOnline,
}
