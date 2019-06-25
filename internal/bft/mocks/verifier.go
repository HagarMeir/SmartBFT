// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import types "github.com/SmartBFT-Go/consensus/pkg/types"

// Verifier is an autogenerated mock type for the Verifier type
type Verifier struct {
	mock.Mock
}

// VerificationSequence provides a mock function with given fields:
func (_m *Verifier) VerificationSequence() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// VerifyConsenterSig provides a mock function with given fields: signature, prop
func (_m *Verifier) VerifyConsenterSig(signature types.Signature, prop types.Proposal) error {
	ret := _m.Called(signature, prop)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Signature, types.Proposal) error); ok {
		r0 = rf(signature, prop)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyProposal provides a mock function with given fields: proposal, prevHeader
func (_m *Verifier) VerifyProposal(proposal types.Proposal, prevHeader []byte) ([]types.RequestInfo, error) {
	ret := _m.Called(proposal, prevHeader)

	var r0 []types.RequestInfo
	if rf, ok := ret.Get(0).(func(types.Proposal, []byte) []types.RequestInfo); ok {
		r0 = rf(proposal, prevHeader)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.RequestInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.Proposal, []byte) error); ok {
		r1 = rf(proposal, prevHeader)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyRequest provides a mock function with given fields: val
func (_m *Verifier) VerifyRequest(val []byte) (types.RequestInfo, error) {
	ret := _m.Called(val)

	var r0 types.RequestInfo
	if rf, ok := ret.Get(0).(func([]byte) types.RequestInfo); ok {
		r0 = rf(val)
	} else {
		r0 = ret.Get(0).(types.RequestInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(val)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
