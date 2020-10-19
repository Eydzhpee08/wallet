package types

// Money представляет собой денежную сумму в минимальных единицах (центы, копейки, дирамы и т.д.)
type Money int64





// Category представлет собой категорию, в которой был совершен платеж (авто, аптеки, рестораны и т.д)
type PaymentCategory string

// PaymentStatus - представляет собой статус платежа 
type PaymentStatus string




// PAN представляет номер карты
type PAN string

// Status представляет собой статус платежа


// представляют статусы 
const  (
	PaymentStatusOk PaymentStatus = "OK"
	PaymentStatusFail PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)





// Payment представляет информацию о платеже.
type Payment struct {
	ID string
	AccountID int64
	Amount Money
	Category PaymentCategory
	Status PaymentStatus
}


type Phone string

// Account - представляет информацию о счете пользователя
type Account struct{
	ID int64
	Phone Phone
	Balance Money
}

type Favorite struct {
	ID        string
	AccountID int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}