// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	sync "sync"

	v1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	install "github.com/operator-framework/operator-lifecycle-manager/pkg/controller/install"
	operatorclient "github.com/operator-framework/operator-lifecycle-manager/pkg/lib/operatorclient"
	operatorlister "github.com/operator-framework/operator-lifecycle-manager/pkg/lib/operatorlister"
	ownerutil "github.com/operator-framework/operator-lifecycle-manager/pkg/lib/ownerutil"
)

type FakeStrategyResolverInterface struct {
	InstallerForStrategyStub        func(string, operatorclient.ClientInterface, operatorlister.OperatorLister, ownerutil.Owner, map[string]string, install.Strategy) install.StrategyInstaller
	installerForStrategyMutex       sync.RWMutex
	installerForStrategyArgsForCall []struct {
		arg1 string
		arg2 operatorclient.ClientInterface
		arg3 operatorlister.OperatorLister
		arg4 ownerutil.Owner
		arg5 map[string]string
		arg6 install.Strategy
	}
	installerForStrategyReturns struct {
		result1 install.StrategyInstaller
	}
	installerForStrategyReturnsOnCall map[int]struct {
		result1 install.StrategyInstaller
	}
	UnmarshalStrategyStub        func(v1alpha1.NamedInstallStrategy) (install.Strategy, error)
	unmarshalStrategyMutex       sync.RWMutex
	unmarshalStrategyArgsForCall []struct {
		arg1 v1alpha1.NamedInstallStrategy
	}
	unmarshalStrategyReturns struct {
		result1 install.Strategy
		result2 error
	}
	unmarshalStrategyReturnsOnCall map[int]struct {
		result1 install.Strategy
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStrategyResolverInterface) InstallerForStrategy(arg1 string, arg2 operatorclient.ClientInterface, arg3 operatorlister.OperatorLister, arg4 ownerutil.Owner, arg5 map[string]string, arg6 install.Strategy) install.StrategyInstaller {
	fake.installerForStrategyMutex.Lock()
	ret, specificReturn := fake.installerForStrategyReturnsOnCall[len(fake.installerForStrategyArgsForCall)]
	fake.installerForStrategyArgsForCall = append(fake.installerForStrategyArgsForCall, struct {
		arg1 string
		arg2 operatorclient.ClientInterface
		arg3 operatorlister.OperatorLister
		arg4 ownerutil.Owner
		arg5 map[string]string
		arg6 install.Strategy
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.recordInvocation("InstallerForStrategy", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.installerForStrategyMutex.Unlock()
	if fake.InstallerForStrategyStub != nil {
		return fake.InstallerForStrategyStub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.installerForStrategyReturns
	return fakeReturns.result1
}

func (fake *FakeStrategyResolverInterface) InstallerForStrategyCallCount() int {
	fake.installerForStrategyMutex.RLock()
	defer fake.installerForStrategyMutex.RUnlock()
	return len(fake.installerForStrategyArgsForCall)
}

func (fake *FakeStrategyResolverInterface) InstallerForStrategyCalls(stub func(string, operatorclient.ClientInterface, operatorlister.OperatorLister, ownerutil.Owner, map[string]string, install.Strategy) install.StrategyInstaller) {
	fake.installerForStrategyMutex.Lock()
	defer fake.installerForStrategyMutex.Unlock()
	fake.InstallerForStrategyStub = stub
}

func (fake *FakeStrategyResolverInterface) InstallerForStrategyArgsForCall(i int) (string, operatorclient.ClientInterface, operatorlister.OperatorLister, ownerutil.Owner, map[string]string, install.Strategy) {
	fake.installerForStrategyMutex.RLock()
	defer fake.installerForStrategyMutex.RUnlock()
	argsForCall := fake.installerForStrategyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeStrategyResolverInterface) InstallerForStrategyReturns(result1 install.StrategyInstaller) {
	fake.installerForStrategyMutex.Lock()
	defer fake.installerForStrategyMutex.Unlock()
	fake.InstallerForStrategyStub = nil
	fake.installerForStrategyReturns = struct {
		result1 install.StrategyInstaller
	}{result1}
}

func (fake *FakeStrategyResolverInterface) InstallerForStrategyReturnsOnCall(i int, result1 install.StrategyInstaller) {
	fake.installerForStrategyMutex.Lock()
	defer fake.installerForStrategyMutex.Unlock()
	fake.InstallerForStrategyStub = nil
	if fake.installerForStrategyReturnsOnCall == nil {
		fake.installerForStrategyReturnsOnCall = make(map[int]struct {
			result1 install.StrategyInstaller
		})
	}
	fake.installerForStrategyReturnsOnCall[i] = struct {
		result1 install.StrategyInstaller
	}{result1}
}

func (fake *FakeStrategyResolverInterface) UnmarshalStrategy(arg1 v1alpha1.NamedInstallStrategy) (install.Strategy, error) {
	fake.unmarshalStrategyMutex.Lock()
	ret, specificReturn := fake.unmarshalStrategyReturnsOnCall[len(fake.unmarshalStrategyArgsForCall)]
	fake.unmarshalStrategyArgsForCall = append(fake.unmarshalStrategyArgsForCall, struct {
		arg1 v1alpha1.NamedInstallStrategy
	}{arg1})
	fake.recordInvocation("UnmarshalStrategy", []interface{}{arg1})
	fake.unmarshalStrategyMutex.Unlock()
	if fake.UnmarshalStrategyStub != nil {
		return fake.UnmarshalStrategyStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.unmarshalStrategyReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStrategyResolverInterface) UnmarshalStrategyCallCount() int {
	fake.unmarshalStrategyMutex.RLock()
	defer fake.unmarshalStrategyMutex.RUnlock()
	return len(fake.unmarshalStrategyArgsForCall)
}

func (fake *FakeStrategyResolverInterface) UnmarshalStrategyCalls(stub func(v1alpha1.NamedInstallStrategy) (install.Strategy, error)) {
	fake.unmarshalStrategyMutex.Lock()
	defer fake.unmarshalStrategyMutex.Unlock()
	fake.UnmarshalStrategyStub = stub
}

func (fake *FakeStrategyResolverInterface) UnmarshalStrategyArgsForCall(i int) v1alpha1.NamedInstallStrategy {
	fake.unmarshalStrategyMutex.RLock()
	defer fake.unmarshalStrategyMutex.RUnlock()
	argsForCall := fake.unmarshalStrategyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStrategyResolverInterface) UnmarshalStrategyReturns(result1 install.Strategy, result2 error) {
	fake.unmarshalStrategyMutex.Lock()
	defer fake.unmarshalStrategyMutex.Unlock()
	fake.UnmarshalStrategyStub = nil
	fake.unmarshalStrategyReturns = struct {
		result1 install.Strategy
		result2 error
	}{result1, result2}
}

func (fake *FakeStrategyResolverInterface) UnmarshalStrategyReturnsOnCall(i int, result1 install.Strategy, result2 error) {
	fake.unmarshalStrategyMutex.Lock()
	defer fake.unmarshalStrategyMutex.Unlock()
	fake.UnmarshalStrategyStub = nil
	if fake.unmarshalStrategyReturnsOnCall == nil {
		fake.unmarshalStrategyReturnsOnCall = make(map[int]struct {
			result1 install.Strategy
			result2 error
		})
	}
	fake.unmarshalStrategyReturnsOnCall[i] = struct {
		result1 install.Strategy
		result2 error
	}{result1, result2}
}

func (fake *FakeStrategyResolverInterface) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.installerForStrategyMutex.RLock()
	defer fake.installerForStrategyMutex.RUnlock()
	fake.unmarshalStrategyMutex.RLock()
	defer fake.unmarshalStrategyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStrategyResolverInterface) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ install.StrategyResolverInterface = new(FakeStrategyResolverInterface)