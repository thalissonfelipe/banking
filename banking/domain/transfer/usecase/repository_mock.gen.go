// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package usecase

import (
	"context"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"sync"
)

// Ensure, that RepositoryMock does implement transfer.Repository.
// If this is not the case, regenerate this file with moq.
var _ transfer.Repository = &RepositoryMock{}

// RepositoryMock is a mock implementation of transfer.Repository.
//
// 	func TestSomethingThatUsesRepository(t *testing.T) {
//
// 		// make and configure a mocked transfer.Repository
// 		mockedRepository := &RepositoryMock{
// 			ListTransfersFunc: func(contextMoqParam context.Context, accountID vos.AccountID) ([]entities.Transfer, error) {
// 				panic("mock out the ListTransfers method")
// 			},
// 			PerformTransferFunc: func(contextMoqParam context.Context, transfer *entities.Transfer) error {
// 				panic("mock out the PerformTransfer method")
// 			},
// 		}
//
// 		// use mockedRepository in code that requires transfer.Repository
// 		// and then make assertions.
//
// 	}
type RepositoryMock struct {
	// ListTransfersFunc mocks the ListTransfers method.
	ListTransfersFunc func(contextMoqParam context.Context, accountID vos.AccountID) ([]entities.Transfer, error)

	// PerformTransferFunc mocks the PerformTransfer method.
	PerformTransferFunc func(contextMoqParam context.Context, transfer *entities.Transfer) error

	// calls tracks calls to the methods.
	calls struct {
		// ListTransfers holds details about calls to the ListTransfers method.
		ListTransfers []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// AccountID is the accountID argument value.
			AccountID vos.AccountID
		}
		// PerformTransfer holds details about calls to the PerformTransfer method.
		PerformTransfer []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Transfer is the transfer argument value.
			Transfer *entities.Transfer
		}
	}
	lockListTransfers   sync.RWMutex
	lockPerformTransfer sync.RWMutex
}

// ListTransfers calls ListTransfersFunc.
func (mock *RepositoryMock) ListTransfers(contextMoqParam context.Context, accountID vos.AccountID) ([]entities.Transfer, error) {
	if mock.ListTransfersFunc == nil {
		panic("RepositoryMock.ListTransfersFunc: method is nil but Repository.ListTransfers was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		AccountID       vos.AccountID
	}{
		ContextMoqParam: contextMoqParam,
		AccountID:       accountID,
	}
	mock.lockListTransfers.Lock()
	mock.calls.ListTransfers = append(mock.calls.ListTransfers, callInfo)
	mock.lockListTransfers.Unlock()
	return mock.ListTransfersFunc(contextMoqParam, accountID)
}

// ListTransfersCalls gets all the calls that were made to ListTransfers.
// Check the length with:
//     len(mockedRepository.ListTransfersCalls())
func (mock *RepositoryMock) ListTransfersCalls() []struct {
	ContextMoqParam context.Context
	AccountID       vos.AccountID
} {
	var calls []struct {
		ContextMoqParam context.Context
		AccountID       vos.AccountID
	}
	mock.lockListTransfers.RLock()
	calls = mock.calls.ListTransfers
	mock.lockListTransfers.RUnlock()
	return calls
}

// PerformTransfer calls PerformTransferFunc.
func (mock *RepositoryMock) PerformTransfer(contextMoqParam context.Context, transfer *entities.Transfer) error {
	if mock.PerformTransferFunc == nil {
		panic("RepositoryMock.PerformTransferFunc: method is nil but Repository.PerformTransfer was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Transfer        *entities.Transfer
	}{
		ContextMoqParam: contextMoqParam,
		Transfer:        transfer,
	}
	mock.lockPerformTransfer.Lock()
	mock.calls.PerformTransfer = append(mock.calls.PerformTransfer, callInfo)
	mock.lockPerformTransfer.Unlock()
	return mock.PerformTransferFunc(contextMoqParam, transfer)
}

// PerformTransferCalls gets all the calls that were made to PerformTransfer.
// Check the length with:
//     len(mockedRepository.PerformTransferCalls())
func (mock *RepositoryMock) PerformTransferCalls() []struct {
	ContextMoqParam context.Context
	Transfer        *entities.Transfer
} {
	var calls []struct {
		ContextMoqParam context.Context
		Transfer        *entities.Transfer
	}
	mock.lockPerformTransfer.RLock()
	calls = mock.calls.PerformTransfer
	mock.lockPerformTransfer.RUnlock()
	return calls
}
