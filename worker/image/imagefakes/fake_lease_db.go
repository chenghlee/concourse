// This file was generated by counterfeiter
package imagefakes

import (
	"sync"
	"time"

	"github.com/concourse/atc/db"
	"github.com/concourse/atc/worker/image"
	"github.com/pivotal-golang/lager"
)

type FakeLeaseDB struct {
	GetLeaseStub        func(logger lager.Logger, leaseName string, interval time.Duration) (db.Lease, bool, error)
	getLeaseMutex       sync.RWMutex
	getLeaseArgsForCall []struct {
		logger    lager.Logger
		leaseName string
		interval  time.Duration
	}
	getLeaseReturns struct {
		result1 db.Lease
		result2 bool
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLeaseDB) GetLease(logger lager.Logger, leaseName string, interval time.Duration) (db.Lease, bool, error) {
	fake.getLeaseMutex.Lock()
	fake.getLeaseArgsForCall = append(fake.getLeaseArgsForCall, struct {
		logger    lager.Logger
		leaseName string
		interval  time.Duration
	}{logger, leaseName, interval})
	fake.recordInvocation("GetLease", []interface{}{logger, leaseName, interval})
	fake.getLeaseMutex.Unlock()
	if fake.GetLeaseStub != nil {
		return fake.GetLeaseStub(logger, leaseName, interval)
	} else {
		return fake.getLeaseReturns.result1, fake.getLeaseReturns.result2, fake.getLeaseReturns.result3
	}
}

func (fake *FakeLeaseDB) GetLeaseCallCount() int {
	fake.getLeaseMutex.RLock()
	defer fake.getLeaseMutex.RUnlock()
	return len(fake.getLeaseArgsForCall)
}

func (fake *FakeLeaseDB) GetLeaseArgsForCall(i int) (lager.Logger, string, time.Duration) {
	fake.getLeaseMutex.RLock()
	defer fake.getLeaseMutex.RUnlock()
	return fake.getLeaseArgsForCall[i].logger, fake.getLeaseArgsForCall[i].leaseName, fake.getLeaseArgsForCall[i].interval
}

func (fake *FakeLeaseDB) GetLeaseReturns(result1 db.Lease, result2 bool, result3 error) {
	fake.GetLeaseStub = nil
	fake.getLeaseReturns = struct {
		result1 db.Lease
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeLeaseDB) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getLeaseMutex.RLock()
	defer fake.getLeaseMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeLeaseDB) recordInvocation(key string, args []interface{}) {
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

var _ image.LeaseDB = new(FakeLeaseDB)
