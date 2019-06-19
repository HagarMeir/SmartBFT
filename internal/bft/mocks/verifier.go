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

// VerifyConsenterSig provides a mock function with given fields: signer, signature, prop
func (_m *Verifier) VerifyConsenterSig(signer uint64, signature []byte, prop types.Proposal) error {
	ret := _m.Called(signer, signature, prop)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, []byte, types.Proposal) error); ok {
		r0 = rf(signer, signature, prop)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyProposal provides a mock function with given fields: proposal, prevHeader
func (_m *Verifier) VerifyProposal(proposal types.Proposal, prevHeader []byte) error {
	ret := _m.Called(proposal, prevHeader)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Proposal, []byte) error); ok {
		r0 = rf(proposal, prevHeader)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyRequest provides a mock function with given fields: val
func (_m *Verifier) VerifyRequest(val []byte) error {
	ret := _m.Called(val)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(val)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}