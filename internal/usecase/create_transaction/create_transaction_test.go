package create_transaction

import (
	"context"
	"testing"

	"github.com/flaviomdutra/walletcore/internal/entity"
	"github.com/flaviomdutra/walletcore/internal/event"
	"github.com/flaviomdutra/walletcore/internal/usecase/mocks"
	"github.com/flaviomdutra/walletcore/pkg/events"
	"github.com/flaviomdutra/walletcore/pkg/uow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "john@doe.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("Jane Doe", "jane@doe.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	accountGateway := &mocks.AccountGatewayMock{}
	accountGateway.On("FindByID", account1.ID).Return(account1, nil).Once()
	accountGateway.On("FindByID", account2.ID).Return(account2, nil).Once()
	accountGateway.On("UpdateBalance", account1).Return(nil).Once()
	accountGateway.On("UpdateBalance", account2).Return(nil).Once()

	transactionGateway := &mocks.TransactionGatewayMock{}
	transactionGateway.On("Create", mock.Anything).Return(nil).Once()

	mockUow := &mocks.UowMock{}
	mockUow.On("GetRepository", "AccountDB").Return(accountGateway, nil).Once()
	mockUow.On("GetRepository", "TransactionDB").Return(transactionGateway, nil).Once()
	mockUow.On("Do", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(*uow.Uow) error)
		assert.NoError(t, fn(nil))
	}).Once()

	inputDTO := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100.0,
	}

	dispatcher := events.NewEventDispatcher()
	eventTransaction := event.NewTransactionCreated()
	eventBalance := event.NewBalanceUpdated()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(mockUow, dispatcher, eventTransaction, eventBalance)
	output, err := uc.Execute(ctx, inputDTO)
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 900.0, account1.Balance)
	assert.Equal(t, 1100.0, account2.Balance)
	assert.Same(t, output, eventTransaction.GetPayload())

	balanceOutput, ok := eventBalance.GetPayload().(*BalanceUpdatedOutputDTO)
	if assert.True(t, ok) {
		assert.Equal(t, account1.ID, balanceOutput.AccountIDFrom)
		assert.Equal(t, account2.ID, balanceOutput.AccountIDTo)
		assert.Equal(t, 900.0, balanceOutput.BalanceAccountIDFrom)
		assert.Equal(t, 1100.0, balanceOutput.BalanceAccountIDTo)
	}

	accountGateway.AssertExpectations(t)
	transactionGateway.AssertExpectations(t)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
