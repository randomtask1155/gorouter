// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"code.cloudfoundry.org/gorouter/route"
)

type FakeEndpointIterator struct {
	NextStub        func() *route.Endpoint
	nextMutex       sync.RWMutex
	nextArgsForCall []struct{}
	nextReturns     struct {
		result1 *route.Endpoint
	}
	nextReturnsOnCall map[int]struct {
		result1 *route.Endpoint
	}
	EndpointFailedStub        func(err error)
	endpointFailedMutex       sync.RWMutex
	endpointFailedArgsForCall []struct {
		err error
	}
	PreRequestStub        func(e *route.Endpoint)
	preRequestMutex       sync.RWMutex
	preRequestArgsForCall []struct {
		e *route.Endpoint
	}
	PostRequestStub        func(e *route.Endpoint)
	postRequestMutex       sync.RWMutex
	postRequestArgsForCall []struct {
		e *route.Endpoint
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeEndpointIterator) Next() *route.Endpoint {
	fake.nextMutex.Lock()
	ret, specificReturn := fake.nextReturnsOnCall[len(fake.nextArgsForCall)]
	fake.nextArgsForCall = append(fake.nextArgsForCall, struct{}{})
	fake.recordInvocation("Next", []interface{}{})
	fake.nextMutex.Unlock()
	if fake.NextStub != nil {
		return fake.NextStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.nextReturns.result1
}

func (fake *FakeEndpointIterator) NextCallCount() int {
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	return len(fake.nextArgsForCall)
}

func (fake *FakeEndpointIterator) NextReturns(result1 *route.Endpoint) {
	fake.NextStub = nil
	fake.nextReturns = struct {
		result1 *route.Endpoint
	}{result1}
}

func (fake *FakeEndpointIterator) NextReturnsOnCall(i int, result1 *route.Endpoint) {
	fake.NextStub = nil
	if fake.nextReturnsOnCall == nil {
		fake.nextReturnsOnCall = make(map[int]struct {
			result1 *route.Endpoint
		})
	}
	fake.nextReturnsOnCall[i] = struct {
		result1 *route.Endpoint
	}{result1}
}

func (fake *FakeEndpointIterator) EndpointFailed(err error) {
	fake.endpointFailedMutex.Lock()
	fake.endpointFailedArgsForCall = append(fake.endpointFailedArgsForCall, struct {
		err error
	}{err})
	fake.recordInvocation("EndpointFailed", []interface{}{err})
	fake.endpointFailedMutex.Unlock()
	if fake.EndpointFailedStub != nil {
		fake.EndpointFailedStub(err)
	}
}

func (fake *FakeEndpointIterator) EndpointFailedCallCount() int {
	fake.endpointFailedMutex.RLock()
	defer fake.endpointFailedMutex.RUnlock()
	return len(fake.endpointFailedArgsForCall)
}

func (fake *FakeEndpointIterator) EndpointFailedArgsForCall(i int) error {
	fake.endpointFailedMutex.RLock()
	defer fake.endpointFailedMutex.RUnlock()
	return fake.endpointFailedArgsForCall[i].err
}

func (fake *FakeEndpointIterator) PreRequest(e *route.Endpoint) {
	fake.preRequestMutex.Lock()
	fake.preRequestArgsForCall = append(fake.preRequestArgsForCall, struct {
		e *route.Endpoint
	}{e})
	fake.recordInvocation("PreRequest", []interface{}{e})
	fake.preRequestMutex.Unlock()
	if fake.PreRequestStub != nil {
		fake.PreRequestStub(e)
	}
}

func (fake *FakeEndpointIterator) PreRequestCallCount() int {
	fake.preRequestMutex.RLock()
	defer fake.preRequestMutex.RUnlock()
	return len(fake.preRequestArgsForCall)
}

func (fake *FakeEndpointIterator) PreRequestArgsForCall(i int) *route.Endpoint {
	fake.preRequestMutex.RLock()
	defer fake.preRequestMutex.RUnlock()
	return fake.preRequestArgsForCall[i].e
}

func (fake *FakeEndpointIterator) PostRequest(e *route.Endpoint) {
	fake.postRequestMutex.Lock()
	fake.postRequestArgsForCall = append(fake.postRequestArgsForCall, struct {
		e *route.Endpoint
	}{e})
	fake.recordInvocation("PostRequest", []interface{}{e})
	fake.postRequestMutex.Unlock()
	if fake.PostRequestStub != nil {
		fake.PostRequestStub(e)
	}
}

func (fake *FakeEndpointIterator) PostRequestCallCount() int {
	fake.postRequestMutex.RLock()
	defer fake.postRequestMutex.RUnlock()
	return len(fake.postRequestArgsForCall)
}

func (fake *FakeEndpointIterator) PostRequestArgsForCall(i int) *route.Endpoint {
	fake.postRequestMutex.RLock()
	defer fake.postRequestMutex.RUnlock()
	return fake.postRequestArgsForCall[i].e
}

func (fake *FakeEndpointIterator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.nextMutex.RLock()
	defer fake.nextMutex.RUnlock()
	fake.endpointFailedMutex.RLock()
	defer fake.endpointFailedMutex.RUnlock()
	fake.preRequestMutex.RLock()
	defer fake.preRequestMutex.RUnlock()
	fake.postRequestMutex.RLock()
	defer fake.postRequestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeEndpointIterator) recordInvocation(key string, args []interface{}) {
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

var _ route.EndpointIterator = new(FakeEndpointIterator)
