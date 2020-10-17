package wallet

import (
	"io"
	"path/filepath"
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Eydzhpee08/wallet/pkg/types"
	"github.com/google/uuid"
)

// This is some predetermined errors.
var (
	ErrPhoneRegistered      = errors.New("phone is already registered")
	ErrAmountMustBePositive = errors.New("amount must be greater than 0")
	ErrAccountNotFound      = errors.New("account not found")
	ErrNotEnoughBalance     = errors.New("not enough money in your balance")
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrPaymentInFavorite    = errors.New("payment is already in favorite")
	ErrFavoriteNotFound     = errors.New("favorite payment not found")
)

// Service ...
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

// RegisterAccount is function for register new accounts. It was made by Alif and writed by _Muhammadkhon_
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
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

// Deposit is the function of account replenishment. It was made by Alif and writed by _Muhammadkhon_
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

// Pay is the function. It was made by Alif and writed by _Muhammadkhon_
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

// FindAccountByID is the function for find interested account by ID. It was made by _Muhammadkhon_
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound

}

// FindPaymentByID is the function for find interested payment by ID. It was made by _Muhammadkhon_
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if paymentID == payment.ID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

// Reject is the function for reject payment. It was made by __Muhammadkhon__
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}

// Repeat is function for repeat payment by ID. It was made by __Muhammadkhon__
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	newPayment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return newPayment, nil
}

// FavoritePayment ...
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	for _, favorite := range s.favorites {
		if favorite.ID == payment.ID {
			return nil, ErrPaymentInFavorite
		}
	}
	newFavorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, newFavorite)
	return newFavorite, nil
}

// FindFavoriteByID is the function for find favorite payment by ID. It was made by _Muhammadkhon_
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

// PayFromFavorite is function for repeat payment from favorite. It was made by __Muhammadkhon__
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}

// ExportToFile is function
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	result := ""
	for _, account := range s.accounts {
		result += strconv.Itoa(int(account.ID)) + ";"
		result += string(account.Phone) + ";"
		result += strconv.Itoa(int(account.Balance)) + "|"
	}

	_, err = file.Write([]byte(result))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ImportFromFile is function
func (s *Service) ImportFromFile(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	data := string(byteData)
	splitSlice := strings.Split(data, "|")
	for _, split := range splitSlice {
		if split != "" {
			datas := strings.Split(split, ";")
			id, err := strconv.Atoi(datas[0])
			if err != nil {
				log.Println(err)
				return err
			}

			balance, err := strconv.Atoi(datas[2])
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

// WriteToFile is function to write data to file
func WriteToFile(fileName string, data []byte) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			log.Print("WriteToFile. Could not create a folder. aaaa panic: ")
			panic(merr)
		}
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print("WriteToFile. Open file error: ", err)
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Print("WriteToFile. Close file error: ", closeErr)
		}
	}()
	_, err = file.Write(data)

	if err != nil {
		log.Print("WriteToFile. Write file error: ", err)
	}
	return nil
}

// Export is function for export data from services to file
func (s *Service) Export(dir string) error {
	log.Print("start exporting accounts entity, count of accounts: ", len(s.accounts))
	accExp := 0
	for _, account := range s.accounts {
		ID := strconv.FormatInt(account.ID, 10) + ";"
		phone := string(account.Phone) + ";"
		balance := strconv.FormatInt(int64(account.Balance), 10)
		err := WriteToFile(dir+"/accounts.dump", []byte(ID+phone+balance+"\n"))
		if err != nil {
			return err
		}
		accExp++
	}
	log.Print("end of exporting accounts entity, amount of exported accs: ", accExp)

	log.Print("start exporting payments entity, count of payments: ", len(s.payments))
	payExp := 0
	for _, payment := range s.payments {
		ID := payment.ID + ";"
		AccountID := strconv.FormatInt(payment.AccountID, 10) + ";"
		Amount := strconv.FormatInt(int64(payment.Amount), 10) + ";"
		Category := string(payment.Category) + ";"
		Status := string(payment.Status) + "\n"
		err := WriteToFile(dir+"/payments.dump", []byte(ID+AccountID+Amount+Category+Status))
		if err != nil {
			return err
		}
		payExp++
	}
	log.Print("end of exporting payments entity, amount of exported payments: ", payExp)

	log.Print("start exporting favorites entity, count of favorites: ", len(s.favorites))
	favExp := 0
	for _, favorite := range s.favorites {
		ID := favorite.ID + ";"
		AccountID := strconv.FormatInt(favorite.AccountID, 10) + ";"
		Name := favorite.Name + ";"
		Amount := strconv.FormatInt(int64(favorite.Amount), 10) + ";"
		Category := string(favorite.Category) + "\n"
		err := WriteToFile(dir+"/favorites.dump", []byte(ID+AccountID+Name+Amount+Category))
		favExp++
		if err != nil {
			return err
		}
	}
	log.Print("end of exporting favorites entity, amount of exported favorite: ", favExp)
	return nil
}

// Import is function
func (s *Service) Import(dir string) error {
	log.Print("account count in the start of import method: ", len(s.accounts))
	log.Print("Start Import method with param: " + dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return err
	}
	for _, file := range files {
		log.Print("files in Import->dir: " + file.Name())
		read, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if closeErr := read.Close(); closeErr != nil {
				log.Print(closeErr)
			}
		}()

		reader := bufio.NewReader(read)

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				log.Print("line in OEF: ", line)
				break
			}
			if err != nil {
				log.Print(err)
				return err
			}

			item := strings.Split(line, ";")
			switch file.Name() {
			case "accounts.dump":
				acc := s.convertToAccount(item)
				if acc != nil {
					s.accounts = append(s.accounts, acc)
				}
			case "favorites.dump":
				favorite := s.convertToFavorites(item)
				if favorite != nil {
					s.favorites = append(s.favorites, favorite)
				}
			case "payments.dump":
				payment := s.convertToPayments(item)
				if payment != nil {
					s.payments = append(s.payments, payment)
				}
			default:
				break
			}
		}

	}
	log.Print("account count in the end of import method: ", len(s.accounts))
	return nil
}

func (s *Service) convertToAccount(item []string) *types.Account {
	ID, _ := strconv.ParseInt(item[0], 10, 64)
	balance, _ := strconv.ParseInt(removeEndLine(item[2]), 10, 64)
	account, err := s.FindAccountByID(ID)
	if err != nil {
		s.nextAccountID++
		return &types.Account{
			ID:      ID,
			Phone:   types.Phone(item[1]),
			Balance: types.Money(balance),
		}
	}
	account.ID = ID
	account.Phone = types.Phone(item[1])
	account.Balance = types.Money(balance)
	return nil
}

func (s *Service) convertToFavorites(item []string) *types.Favorite {
	AccountID, _ := strconv.ParseInt(item[1], 10, 64)
	Amount, _ := strconv.ParseInt(item[3], 10, 64)

	favorite, err := s.FindFavoriteByID(item[0])
	if err != nil {
		return &types.Favorite{
			ID:        item[0],
			AccountID: AccountID,
			Name:      item[2],
			Amount:    types.Money(Amount),
			Category:  types.PaymentCategory(item[4]),
		}
	}
	favorite.ID = item[0]
	favorite.AccountID = AccountID
	favorite.Name = item[2]
	favorite.Amount = types.Money(Amount)
	favorite.Category = types.PaymentCategory(removeEndLine(item[4]))
	return nil
}

func (s *Service) convertToPayments(item []string) *types.Payment {
	AccountID, _ := strconv.ParseInt(item[1], 10, 64)
	Amount, _ := strconv.ParseInt(item[2], 10, 64)

	payment, err := s.FindPaymentByID(item[0])
	if err != nil {
		return &types.Payment{
			ID:        item[0],
			AccountID: AccountID,
			Amount:    types.Money(Amount),
			Category:  types.PaymentCategory(item[3]),
			Status:    types.PaymentStatus(removeEndLine(item[4])),
		}
	}
	payment.ID = item[0]
	payment.AccountID = AccountID
	payment.Amount = types.Money(Amount)
	payment.Category = types.PaymentCategory(item[3])
	payment.Status = types.PaymentStatus(item[4])
	return nil
}

func removeEndLine(balance string) string {
	return strings.TrimRightFunc(balance, func(c rune) bool {
		return c == '\r' || c == '\n'
	})
}