package wallet

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"errors"
	"github.com/Eydzhpee08/wallet/pkg/types"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
	"sync"
)

type Error string

func (e Error) Error() string {

	return string(e)

}

var ErrPhoneNumberRegistred = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater that zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")
var ErrFileNotFound = errors.New("File Not found")

//Service -
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

//RegisterAccount создаем тут ак
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {

	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneNumberRegistred
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)

	return account, nil

}

// так, он находит по ID
// у нас есть слайст структур
// нам нужен отдельная структура, чтобы туда запихнуть данные
// вызываем цикл
// в аргументах например мы дали ID 1
// потом мы делаем проверку
// acc.ID == account.ID{
// account = acc
//

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive

	}

	var account *types.Account


	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
		}

	}
	if account == nil {
		return ErrAccountNotFound
	}
	account.Balance += amount

	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}

	}
	
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance

	}
	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID, 
		Amount:    amount,    
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil

}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil

}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			payment = pay
			break
		}
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}


func (s *Service) Reject(paymentID string) error {
	var targetPayment *types.Payment

	
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			targetPayment = payment
			break
		}

	}
	if targetPayment == nil {
		return ErrPaymentNotFound
	}
	var targetAccount *types.Account

	for _, account := range s.accounts {
		
		if account.ID == targetPayment.AccountID {
			targetAccount = account
			break
		}
	}
	if targetAccount == nil {
		return ErrAccountNotFound
	}
	targetPayment.Status = types.PaymentStatusFail
	targetAccount.Balance += targetPayment.Amount

	return nil

}



func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(pay.AccountID, pay.Amount, pay.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}




// он создает FavoritePayment 
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return nil, err
	}

	favoriteID := uuid.New().String()
	newFavorite := &types.Favorite{
		ID:        favoriteID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, newFavorite)
	return newFavorite, nil
}


func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}



//PayFromFavorite для совершения платежа в Избранное
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

//ExportToFile - для импорта данных
func (s *Service) ExportToFile(path string) error  {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}

		
	}()
	str := ""
	for _, account := range s.accounts {
		str += strconv.Itoa(int(account.ID)) + ";"
		str += string(account.Phone) + ";"
		str += strconv.Itoa(int(account.Balance)) + "|"
	}

	_, err = file.Write([]byte(str))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}


func (s *Service) ImportFromFile(path string) error {


	byteData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}
	
	// обязательно ли преобразовать слайс байтов преобразовывать в стринг?
	data := string(byteData)
	fmt.Println(data)


	splitSlice := strings.Split(data, "|")
	// здесь он должен показать  весь список, тажкже убрать |
	fmt.Println(splitSlice)
	// тут он покажет 2 номер
	fmt.Println(splitSlice[1])
	// тут 1
	fmt.Println(splitSlice[0])
	//splitSlice = strings.Split(data, ";")
	for _, split := range splitSlice {
		if split != "" {
			datas := strings.Split(split, ";")

			id, err := strconv.Atoi(datas[0])
			if err != nil {
				log.Println(err)
				return err
			}
			fmt.Println(datas[0])
			balance, err := strconv.Atoi(datas[2])
			fmt.Println(datas[2])
			if err != nil {
				log.Println(err)
				return err
			}

			newAccount := &types.Account{
				ID:      int64(id),
				Phone:   types.Phone(datas[1]),
				Balance: types.Money(balance),
			}

			s.accounts = append(s.accounts, newAccount)
		}
	}

	return nil
}

// у нас же данные есть ведь? - например RegisterAc - который в памяти хранится
// если она пустая то дай ошибку
// нужна отдельная функция для создания файлов, чтобы 3 раза не писать одно и тоже

func (s *Service) Export(dir string) error {
		// внутри него данные
		// будет тру
	if s.accounts != nil{
		acc := ""
		for _, account := range s.accounts{
			acc += strconv.Itoa(int(account.ID)) + ";"
			acc +=string(account.Phone) + ";"
			acc += strconv.Itoa(int(account.Balance)) +";"
			acc += string('\n')
		}
		err := WriteToFile(dir +"/accounts.dump", acc)
		if err != nil {
			log.Print(err)
			return err
		}
	}
	if  s.payments != nil{
		pay := ""

		for _, payment := range s.payments {

			pay += 	payment.ID + ";"
			pay +=  strconv.Itoa(int(payment.AccountID)) + ";"
			pay +=  strconv.Itoa(int(payment.Amount)) + ";"
			pay +=  string(payment.Category) + ";"
			pay += string(payment.Status) + ";"
			pay += "\n"
			}
		err := WriteToFile(dir + "/payments.dump", pay)
		if err != nil {
			log.Print(err)
			return err
		}

	}
	if s.favorites != nil{
		fav := ""
		for _, favorite := range s.favorites{
			fav += favorite.ID + ";"
			fav += strconv.Itoa(int(favorite.AccountID)) + ";"
			fav += favorite.Name + ";"
			fav += strconv.Itoa(int(favorite.Amount)) + ";"
			fav += string(favorite.Category) + ";"
			fav += "\n"
		}
		err := WriteToFile(dir + "/favorites.dump", fav)
		if err != nil {
			log.Print(err)
			return err
		}

	}







	return  nil

}




func WriteToFile(path string, data string)error  {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err =  file.Close()
		if err != nil {
			log.Print(err)
			return
		}
	}()
	// он возвращает кол-во байтов

	_, err = file.WriteString(data)
	if err != nil {
		log.Print(err)
		return  err
	}
	return nil
}
func (s *Service) Import(dir string) error {
	err := s.actionByAccounts(dir + "/accounts.dump")
	if err != nil {
		log.Println("err from actionByAccount")
		return err
	}

	err = s.actionByPayments(dir + "/payments.dump")
	if err != nil {
		log.Println("err from actionByPayments")
		return err
	}

	err = s.actionByFavorites(dir + "/favorites.dump")
	if err != nil {
		log.Println("err from actionByFavorites")
		return err
	}

	return nil
}

func (s *Service) actionByAccounts(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")

			id, err := strconv.Atoi(data[0])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			phone := types.Phone(data[1])

			balance, err := strconv.Atoi(data[2])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}
			// вот это зачем?
			account, err := s.FindAccountByID(int64(id))
			if err != nil {
				acc, err := s.RegisterAccount(phone)
				if err != nil {
					log.Println("err from register account")
					return err
				}

				acc.Balance = types.Money(balance)
			} else {
				account.Phone = phone
				account.Balance = types.Money(balance)
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}

func (s *Service) actionByPayments(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")
			id := data[0]

			accountID, err := strconv.Atoi(data[1])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			amount, err := strconv.Atoi(data[2])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			category := types.PaymentCategory(data[3])

			status := types.PaymentStatus(data[4])

			payment, err := s.FindPaymentByID(id)
			if err != nil {
				newPayment := &types.Payment{
					ID:        id,
					AccountID: int64(accountID),
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(category),
					Status:    types.PaymentStatus(status),
				}

				s.payments = append(s.payments, newPayment)
			} else {
				payment.AccountID = int64(accountID)
				payment.Amount = types.Money(amount)
				payment.Category = category
				payment.Status = status
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}

func (s *Service) actionByFavorites(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")
			id := data[0]

			accountID, err := strconv.Atoi(data[1])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			name := data[2]

			amount, err := strconv.Atoi(data[3])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			category := types.PaymentCategory(data[4])

			favorite, err := s.FindFavoriteByID(id)
			if err != nil {
				newFavorite := &types.Favorite{
					ID:        id,
					AccountID: int64(accountID),
					Name:      name,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(category),
				}

				s.favorites = append(s.favorites, newFavorite)
			} else {
				favorite.AccountID = int64(accountID)
				favorite.Name = name
				favorite.Amount = types.Money(amount)
				favorite.Category = category
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}


func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}
	accountPayments := []types.Payment{}

	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			accountPayments = append(accountPayments, types.Payment{
				ID:        payment.ID,
				AccountID: payment.AccountID,
				Amount:    payment.Amount,
				Category:  payment.Category,
				Status:    payment.Status,
			})
		}
	}

	return accountPayments, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if len(payments) == 0 {
		return nil
	}
	if len(payments) <= records {
		exportPayments(payments, dir+"/payments.dump")
		return nil
	}
	var iterator int = 1
	count := 0
	for i := 1; i <= len(payments); i++ {
		strIterator := strconv.Itoa(iterator)
		fileName := dir + "/payments" + strIterator + ".dump"
		if i == len(payments) && i-count != records {
			exportPayments(payments[count:i], fileName)
		}

		if i-count == records {
			exportPayments(payments[count:i], fileName)
			count += records
			iterator++
		}
	}
	return nil
}

func exportPayments(payments []types.Payment, path string) error {
	pay := ""

	for _, payment := range payments {

		pay += payment.ID + ";"
		pay += strconv.Itoa(int(payment.AccountID)) + ";"
		pay += strconv.Itoa(int(payment.Amount)) + ";"
		pay += string(payment.Category) + ";"
		pay += string(payment.Status) + ";"
		pay += "\n"
	}
	err := WriteToFile(path, pay)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}


func (s *Service) SumPayments(goroutines int) types.Money {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var summ types.Money = 0
	if goroutines == 0 || goroutines == 1 {
		wg.Add(1)
		go func(payments []*types.Payment) {
			defer wg.Done()
			for _, payment := range payments {
				summ += payment.Amount
			}
		}(s.payments)
	} else {
		from := 0
		count := len(s.payments) / goroutines
		for i := 1; i <= goroutines; i++ {
			wg.Add(1)
			last := len(s.payments) - i*count
			if i == goroutines {
				last = 0
			}
			to := len(s.payments) - last
			go func(payments []*types.Payment) {
				defer wg.Done()
				s := types.Money(0)
				for _, payment := range payments {
					s += payment.Amount
				}
				mu.Lock()
				defer mu.Unlock()
				summ += s
			}(s.payments[from:to])
			from += count
		}
	}

	wg.Wait()

	return summ
}
