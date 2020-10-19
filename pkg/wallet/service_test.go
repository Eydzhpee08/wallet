package wallet

import (
	"github.com/google/uuid"
	"reflect"
	"fmt"
	"testing"
	"github.com/Eydzhpee08/wallet/pkg/types"
)

type testService struct {
	*Service
}

type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
		amount types.Money
		category types.PaymentCategory
	}
}
var defaultTestACcount = testAccount {
	phone: "+992937870880",
	balance: 10_000_00,
	payments: []struct {
		amount types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}
func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *Service) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cant't register account, error = %v", err)
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}
	return account, payments, nil
}

func (s *Service) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	// Register user
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account, error = %v", err)
	}
	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account, error = %v", err)
	}
	return account, nil
}

func TestService_FindAccountByID_success(t *testing.T) {
	var service Service
	service.RegisterAccount("992937870880")

	account, err := service.FindAccountByID(1)

	if err != nil {
		t.Errorf("account => %v", account)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	var service Service
	service.RegisterAccount("992937870880")

	account, err := service.FindAccountByID(2)

	if err == nil {
		t.Errorf("method returned nil error, account => %v", account)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	// Creating service
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestACcount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	// Creating service
	s := newTestService()
	_, _, err := s.addAccount(defaultTestACcount)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = s.FindPaymentByID(uuid.New().String())

	if err == nil {
		t.Errorf("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestACcount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
	}
	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}
	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestACcount.balance {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedPayment)
		return
	}
}


func TestService_Reject_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestACcount)
	if err != nil {
		t.Error(err)
		return
	}
	// payment := payments[0]
	err = s.Reject(uuid.New().String())
	if err == nil {
		t.Errorf("Reject(): must return error, returned nil")
	}
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("Reject(): must return error, returned nil")
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()

	account, err := s.addAccountWithBalance("+992902616961", 100)
	if err != nil {
		t.Errorf("account => %v", account)
		return
	}

	payment, err := s.Pay(account.ID, 10, "Food")
	if err != nil {
		t.Errorf("payment => %v", payment)
		return
	}

	newPayment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("newPayment => %v", newPayment)
		return
	}
}
func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()

	account, err := s.addAccountWithBalance("+992902616961", 100)
	if err != nil {
		t.Errorf("account => %v", account)
		return
	}

	payment, err := s.Pay(account.ID, 10, "Food")
	if err != nil {
		t.Errorf("payment => %v", payment)
		return
	}
	_, err = s.FavoritePayment(payment.ID, "Food")
	if err != nil {
		t.Errorf("error => %v", err)
		return
	}
}

func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTestService()

	account, err := s.addAccountWithBalance("+992902616961", 100)
	if err != nil {
		t.Errorf("account => %v", account)
		return
	}

	payment, err := s.Pay(account.ID, 10, "Food")
	if err != nil {
		t.Errorf("payment => %v", payment)
		return
	}
	favorite, err := s.FavoritePayment(uuid.New().String(), "Food")
	if err == nil {
		t.Errorf("error => %v", favorite)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()

	account, err := s.addAccountWithBalance("+992902616961", 100)
	if err != nil {
		t.Errorf("account => %v", account)
		return
	}

	payment, err := s.Pay(account.ID, 10, "Food")
	if err != nil {
		t.Errorf("payment => %v", payment)
		return
	}
	favorite, err := s.FavoritePayment(payment.ID, "Food")
	if err != nil {
		t.Errorf("error => %v", err)
		return
	}
	payFromFavorite, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("payment => %v", payFromFavorite)
	}
}
func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()

	account, err := s.addAccountWithBalance("+992902616961", 100)
	if err != nil {
		t.Errorf("account => %v", account)
		return
	}

	payment, err := s.Pay(account.ID, 10, "Food")
	if err != nil {
		t.Errorf("payment => %v", payment)
		return
	}
	favorite, err := s.FavoritePayment(payment.ID, "Food")
	if err != nil {
		t.Errorf("error => %v", favorite)
		return
	}
	favoritePayment, err := s.PayFromFavorite(uuid.New().String())
	if err == nil {
		t.Errorf("payment => %v", favoritePayment)
	}
}

