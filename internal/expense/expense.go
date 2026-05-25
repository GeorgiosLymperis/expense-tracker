package expense

import (
	"fmt"
	"time"
)

type Category string
type Date time.Time

const (
	EDUCATION     Category = "Education"
	ENTERTAINMENT Category = "Entertainment"
	BUSINESS      Category = "Business"
	UNKNOWN       Category = "Unknown"
)

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) MarshalJSON() ([]byte, error) {
	return fmt.Appendf([]byte{}, "\"%s\"", time.Time(d).Format("2006-01-02")), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2006-01-02\"", string(data))
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

type NegativeAmountError struct {
	amount float32
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf("Negative amount: %f", e.amount)
}

func validateAmount(amount float32) error {
	if amount < 0 {
		return &NegativeAmountError{amount: amount}
	}
	return nil
}

type Expense struct {
	amount      float32
	date        Date
	description string
	category    Category
}

func NewExpense(amount float32, description string, category Category, date Date) (Expense, error) {
	validAmount := validateAmount(amount)
	if validAmount != nil {
		return Expense{}, validAmount
	}
	return Expense{
		amount:      amount,
		date:        date,
		description: description,
		category:    category}, nil
}

func NewEducationExpense(amount float32, description string, date Date) (Expense, error) {
	return NewExpense(amount, description, EDUCATION, date)
}

func NewEntertainmentExpense(amount float32, description string, date Date) (Expense, error) {
	return NewExpense(amount, description, ENTERTAINMENT, date)
}

func NewBusinessExpense(amount float32, description string, date Date) (Expense, error) {
	return NewExpense(amount, description, BUSINESS, date)
}

func (e Expense) ExpenseAmount() float32 {
	return e.amount
}

func (e Expense) ExpenseDate() Date {
	return e.date
}

func (e Expense) ExpenseDescription() string {
	return e.description
}

func (e Expense) ExpenseCategory() Category {
	return e.category
}
