package expense

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCreationEducationExpense(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e, err := NewEducationExpense(10, "Coursera", Date(day))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if e.amount != 10 {
		t.Errorf("Expected amount was 10, got %v", e.amount)
	}
	if e.description != "Coursera" {
		t.Errorf("Expected description was Coursera, got %v", e.description)
	}
	if e.date != Date(day) {
		t.Errorf("Expected date was 2020-01-01, got %v", e.date)
	}
	if e.category != EDUCATION {
		t.Errorf("Expected category was Education, got %v", e.category)
	}
}

func TestCreationEducationExpenseWithNegativeAmount(t *testing.T) {
	_, err := NewEducationExpense(-10, "Coursera", Date(time.Now()))
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreationEntertainmentExpense(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e, err := NewEntertainmentExpense(10, "Netflix", Date(day))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if e.amount != 10 {
		t.Errorf("Expected amount was 10, got %v", e.amount)
	}
	if e.description != "Netflix" {
		t.Errorf("Expected description was Netflix, got %v", e.description)
	}
	if e.date != Date(day) {
		t.Errorf("Expected date was 2020-01-01, got %v", e.date)
	}
	if e.category != ENTERTAINMENT {
		t.Errorf("Expected category was Entertainment, got %v", e.category)
	}
}

func TestCreationEntertainmentExpenseWithNegativeAmount(t *testing.T) {
	_, err := NewEntertainmentExpense(-10, "Netflix", Date(time.Now()))
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreationBusinessExpense(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e, err := NewBusinessExpense(10, "Amazon", Date(day))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if e.amount != 10 {
		t.Errorf("Expected amount was 10, got %v", e.amount)
	}
	if e.description != "Amazon" {
		t.Errorf("Expected description was Amazon, got %v", e.description)
	}
	if e.date != Date(day) {
		t.Errorf("Expected date was 2020-01-01, got %v", e.date)
	}
	if e.category != BUSINESS {
		t.Errorf("Expected category was Business, got %v", e.category)
	}
}

func TestCreationBusinessExpenseWithNegativeAmount(t *testing.T) {
	_, err := NewBusinessExpense(-10, "Amazon", Date(time.Now()))
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestExpenseAmount(t *testing.T) {
	e, _ := NewBusinessExpense(10, "Amazon", Date(time.Now()))
	if e.ExpenseAmount() != 10 {
		t.Errorf("Expected amount was 10, got %v", e.ExpenseAmount())
	}
}

func TestExpenseDate(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e, _ := NewBusinessExpense(10, "Amazon", Date(day))
	if e.ExpenseDate() != Date(day) {
		t.Errorf("Expected date was 2020-01-01, got %v", e.ExpenseDate())
	}
}

func TestExpenseDescription(t *testing.T) {
	e, _ := NewBusinessExpense(10, "Amazon", Date(time.Now()))
	if e.ExpenseDescription() != "Amazon" {
		t.Errorf("Expected description was Amazon, got %v", e.ExpenseDescription())
	}
}

func TestExpenseCategory(t *testing.T) {
	e, _ := NewBusinessExpense(10, "Amazon", Date(time.Now()))
	if e.ExpenseCategory() != BUSINESS {
		t.Errorf("Expected category was Business, got %v", e.ExpenseCategory())
	}
}

func TestDateMarshallJSON(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	data, err := json.Marshal(Date(day))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(data) != `"2020-01-01"` {
		t.Errorf("Expected date was 2020-01-01, got %v", string(data))
	}
}

func TestDateUnmarshallJSON(t *testing.T) {
	data := []byte(`"2020-01-01"`)
	var d Date
	err := json.Unmarshal(data, &d)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if d != Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected date was 2020-01-01, got %v", d)
	}
}