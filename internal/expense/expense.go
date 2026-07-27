package expense

import (
	"fmt"
	"time"
)

type Category string
type Date time.Time

const (
	Education     Category = "Education"
	Entertainment Category = "Entertainment"
	Business      Category = "Business"
	Food          Category = "Food"
	Transport     Category = "Transport"
	Housing       Category = "Housing"
	Health        Category = "Health"
	Clothing      Category = "Clothing"
	Utilities     Category = "Utilities"
	Travel        Category = "Travel"
	Subscriptions Category = "Subscriptions"
	Savings       Category = "Savings"
	Gifts         Category = "Gifts"
	Unknown       Category = "Unknown"
	Social        Category = "Social"
	Fitness       Category = "Fitness"
	Supplements   Category = "Supplements"
)

var validCategories = map[Category]bool{
	Education:     true,
	Entertainment: true,
	Business:      true,
	Food:          true,
	Transport:     true,
	Housing:       true,
	Health:        true,
	Clothing:      true,
	Utilities:     true,
	Travel:        true,
	Subscriptions: true,
	Savings:       true,
	Gifts:         true,
	Unknown:       true,
	Social:        true,
	Fitness:       true,
	Supplements:   true,
}

func (c Category) IsValid() bool {
	return validCategories[c]
}

func ValidCategories() []Category {
	cats := make([]Category, 0, len(validCategories))
	for c := range validCategories {
		cats = append(cats, c)
	}
	return cats
}

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

func (d Date) Month() time.Month {
	return time.Time(d).Month()
}

func (d Date) Year() int {
	return time.Time(d).Year()
}

func ParseDate(str string) (Date, error) {
	strTime, err := time.Parse("2006-01-02", str)
	if err != nil {
		return Date{}, fmt.Errorf("Error parsing date: %v", err)
	}
	return Date(strTime), nil
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
	err := validateAmount(amount)
	if err != nil {
		return Expense{}, err
	}
	return Expense{
		amount:      amount,
		date:        date,
		description: description,
		category:    category}, nil
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
