package types

// Money presents the amount of money in minimum units (cents, penny, dirams and others)
type Money int64

// PaymentCategory presents the category in which the payment was made(auto, pharmacy, restaraunts and others.)
type PaymentCategory string

// Predefined payment categories
const(
	PaymentCategoryAuto PaymentCategory = "Auto"
	PaymentCategoryFun PaymentCategory = "Fun"
	PaymentCategoryIT PaymentCategory = "IT"
)

// PaymentStatus presents the status of the payment
type PaymentStatus string

// Predefined payment statuses
const (
	PaymentStatusOk PaymentStatus = "OK"
	PaymentStatusFail PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment presents information about payment
type Payment struct {
	ID string
	AccountID int64
	Amount Money
	Category PaymentCategory
	Status PaymentStatus
}
// Phone presents a phone number
type Phone string

// Account presents information about the user's account
type Account struct {
	ID int64
	Phone Phone
	Balance Money
}

// Favorite presents information about Favorite payment
type Favorite struct {
	ID string
	AccountID int64
	Name string
	Amount Money
	Category PaymentCategory
}