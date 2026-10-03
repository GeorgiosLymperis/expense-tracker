package expense

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewExpense(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e, err := NewExpense(Amount(10*100), "Coursera", Education, Date(day))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.amount != Amount(10*100) {
		t.Errorf("expected amount 10, got %v", e.amount)
	}
	if e.description != "Coursera" {
		t.Errorf("expected description 'Coursera', got %v", e.description)
	}
	if e.date != Date(day) {
		t.Errorf("expected date 2020-01-01, got %v", e.date)
	}
	if e.category != Education {
		t.Errorf("expected category Education, got %v", e.category)
	}
}

func TestNewExpenseNegativeAmount(t *testing.T) {
	_, err := NewExpense(Amount(-10*100), "Coursera", Education, Date(time.Now()))
	if err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

func TestExpenseAmount(t *testing.T) {
	e, _ := NewExpense(Amount(10*100), "Amazon", Business, Date(time.Now()))
	if e.ExpenseAmount() != Amount(10*100) {
		t.Errorf("expected amount 10, got %v", e.ExpenseAmount())
	}
}

func TestExpenseDate(t *testing.T) {
	day := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e, _ := NewExpense(Amount(10*100), "Amazon", Business, Date(day))
	if e.ExpenseDate() != Date(day) {
		t.Errorf("expected date 2020-01-01, got %v", e.ExpenseDate())
	}
}

func TestExpenseDescription(t *testing.T) {
	e, _ := NewExpense(Amount(10*100), "Amazon", Business, Date(time.Now()))
	if e.ExpenseDescription() != "Amazon" {
		t.Errorf("expected description 'Amazon', got %v", e.ExpenseDescription())
	}
}

func TestExpenseCategory(t *testing.T) {
	e, _ := NewExpense(Amount(10*100), "Amazon", Business, Date(time.Now()))
	if e.ExpenseCategory() != Business {
		t.Errorf("expected category Business, got %v", e.ExpenseCategory())
	}
}

func TestCategoryIsValid(t *testing.T) {
	if !Education.IsValid() {
		t.Error("expected Education to be valid")
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

func TestParseAmount(t *testing.T) {
	valid := map[string]Amount{
		"0":      0,
		"12":     1200,
		"12.5":   1250,
		"12.50":  1250,
		"0.05":   5,
		"-3.40":  -340,
		" 7.1 ":  710,
		"160.46": 16046,
	}
	for in, want := range valid {
		got, err := ParseAmount(in)
		if err != nil || got != want {
			t.Errorf("ParseAmount(%q) = %v, %v; want %v", in, got, err, want)
		}
	}
	for _, in := range []string{"", "abc", "1.234", "1e3", ".5", "1.", "1,50", "--1", "99999999999999999999"} {
		if _, err := ParseAmount(in); err == nil {
			t.Errorf("ParseAmount(%q) should fail", in)
		}
	}
}

func TestAmountString(t *testing.T) {
	cases := map[Amount]string{0: "0.00", 5: "0.05", 1250: "12.50", -340: "-3.40", -5: "-0.05"}
	for in, want := range cases {
		if got := in.String(); got != want {
			t.Errorf("Amount(%d).String() = %q; want %q", int64(in), got, want)
		}
	}
}

func TestAmountSumIsExact(t *testing.T) {
	var sum Amount
	for range 1000 {
		sum += Amount(0.10 * 100)
	}
	if sum.String() != "100.00" {
		t.Errorf("1000 x 0.10 should be 100.00, got %s", sum)
	}
}
