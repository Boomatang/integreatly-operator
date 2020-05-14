package modify_crs

import (
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
	"time"
)

const (
	RetryInterval        = time.Second * 5
	Timeout              = time.Second * 60 * 3
	IntegreatlyName      = "integreatly-name"
	IntegreatlyNamespace = "integreatly-namespace"
	AmqOnline            = "redhat-rhmi-amq-online"
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

type CompareResult struct {
	Type  string
	Name  string
	Key   string
	Error string
}

type ResourceType interface {
	CrType() string
	CopyRequiredValues(t *testing.T, intContainer *Container, phase string)
	DeleteExistingValues(t *testing.T, intContainer *Container, phase string)
	ChangeCRValues(t *testing.T, intContainer *Container, phase string)
	CompareValues(t *testing.T, intContainer *Container, phase string) *[]CompareResult
}

type CrInterface interface {
	// For any methods that are required in a common
	GetName() string
	GetNamespace() string
	GetResourceVersion() string
	GetKind() string
	GetCr() runtime.Object
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
}
