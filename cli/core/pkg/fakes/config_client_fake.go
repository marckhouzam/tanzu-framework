// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/vmware-tanzu/tanzu-framework/cli/core/pkg/interfaces"
)

type FakeConfigClient struct {
	GetEnvConfigurationsStub        func() map[string]string
	getEnvConfigurationsMutex       sync.RWMutex
	getEnvConfigurationsArgsForCall []struct {
	}
	getEnvConfigurationsReturns struct {
		result1 map[string]string
	}
	getEnvConfigurationsReturnsOnCall map[int]struct {
		result1 map[string]string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConfigClient) GetEnvConfigurations() map[string]string {
	fake.getEnvConfigurationsMutex.Lock()
	ret, specificReturn := fake.getEnvConfigurationsReturnsOnCall[len(fake.getEnvConfigurationsArgsForCall)]
	fake.getEnvConfigurationsArgsForCall = append(fake.getEnvConfigurationsArgsForCall, struct {
	}{})
	stub := fake.GetEnvConfigurationsStub
	fakeReturns := fake.getEnvConfigurationsReturns
	fake.recordInvocation("GetEnvConfigurations", []interface{}{})
	fake.getEnvConfigurationsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConfigClient) GetEnvConfigurationsCallCount() int {
	fake.getEnvConfigurationsMutex.RLock()
	defer fake.getEnvConfigurationsMutex.RUnlock()
	return len(fake.getEnvConfigurationsArgsForCall)
}

func (fake *FakeConfigClient) GetEnvConfigurationsCalls(stub func() map[string]string) {
	fake.getEnvConfigurationsMutex.Lock()
	defer fake.getEnvConfigurationsMutex.Unlock()
	fake.GetEnvConfigurationsStub = stub
}

func (fake *FakeConfigClient) GetEnvConfigurationsReturns(result1 map[string]string) {
	fake.getEnvConfigurationsMutex.Lock()
	defer fake.getEnvConfigurationsMutex.Unlock()
	fake.GetEnvConfigurationsStub = nil
	fake.getEnvConfigurationsReturns = struct {
		result1 map[string]string
	}{result1}
}

func (fake *FakeConfigClient) GetEnvConfigurationsReturnsOnCall(i int, result1 map[string]string) {
	fake.getEnvConfigurationsMutex.Lock()
	defer fake.getEnvConfigurationsMutex.Unlock()
	fake.GetEnvConfigurationsStub = nil
	if fake.getEnvConfigurationsReturnsOnCall == nil {
		fake.getEnvConfigurationsReturnsOnCall = make(map[int]struct {
			result1 map[string]string
		})
	}
	fake.getEnvConfigurationsReturnsOnCall[i] = struct {
		result1 map[string]string
	}{result1}
}

func (fake *FakeConfigClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getEnvConfigurationsMutex.RLock()
	defer fake.getEnvConfigurationsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeConfigClient) recordInvocation(key string, args []interface{}) {
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

var _ interfaces.ConfigClient = new(FakeConfigClient)
