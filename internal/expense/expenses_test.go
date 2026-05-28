package expense

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewExpense(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e, err := NewExpense(10, "Coursera", EDUCATION, Date(day))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.amount != 10 {
		t.Errorf("expected amount 10, got %v", e.amount)
	}
	if e.description != "Coursera" {
		t.Errorf("expected description 'Coursera', got %v", e.description)
	}
	if e.date != Date(day) {
		t.Errorf("expected date 2020-01-01, got %v", e.date)
	}
	if e.category != EDUCATION {
		t.Errorf("expected category Education, got %v", e.category)
	}
}

func TestNewExpenseNegativeAmount(t *testing.T) {
	_, err := NewExpense(-10, "Coursera", EDUCATION, Date(time.Now()))
	if err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

func TestExpenseAmount(t *testing.T) {
	e, _ := NewExpense(10, "Amazon", BUSINESS, Date(time.Now()))
	if e.ExpenseAmount() != 10 {
		t.Errorf("expected amount 10, got %v", e.ExpenseAmount())
	}
}

func TestExpenseDate(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e, _ := NewExpense(10, "Amazon", BUSINESS, Date(day))
	if e.ExpenseDate() != Date(day) {
		t.Errorf("expected date 2020-01-01, got %v", e.ExpenseDate())
	}
}

func TestExpenseDescription(t *testing.T) {
	e, _ := NewExpense(10, "Amazon", BUSINESS, Date(time.Now()))
	if e.ExpenseDescription() != "Amazon" {
		t.Errorf("expected description 'Amazon', got %v", e.ExpenseDescription())
	}
}

func TestExpenseCategory(t *testing.T) {
	e, _ := NewExpense(10, "Amazon", BUSINESS, Date(time.Now()))
	if e.ExpenseCategory() != BUSINESS {
		t.Errorf("expected category Business, got %v", e.ExpenseCategory())
	}
}

func TestCategoryIsValid(t *testing.T) {
	if !EDUCATION.IsValid() {
		t.Error("expected EDUCATION to be valid")
	}
	if Category("InvalidCategory").IsValid() {
		t.Error("expected 'InvalidCategory' to be invalid")
	}
}

func TestDateMarshalJSON(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	data, err := json.Marshal(Date(day))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `"2020-01-01"` {
		t.Errorf("expected '2020-01-01', got %v", string(data))
	}
}

func TestDateUnmarshalJSON(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(`"2020-01-01"`), &d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected 2020-01-01, got %v", d)
	}
}
