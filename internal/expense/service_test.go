package expense

import (
	"testing"
	"time"
)

var (
	day20200101 = Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	day20200102 = Date(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC))
	day20200202 = Date(time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC))
	day20210202 = Date(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC))
)

func newService() Service {
	return NewService(&fakeRepo{})
}

func TestAddExpense(t *testing.T) {
	s := newService()
	if err := s.AddExpense(100.0, "Coursera", Education, day20200101); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAddExpenseNegativeAmount(t *testing.T) {
	s := newService()
	if err := s.AddExpense(-10, "Coursera", Education, day20200101); err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

func TestAddExpenseInvalidCategory(t *testing.T) {
	s := newService()
	if err := s.AddExpense(10, "Coursera", Category("InvalidCategory"), day20200101); err == nil {
		t.Error("expected error for invalid category, got nil")
	}
}

func TestUpdateExpense(t *testing.T) {
	s := newService()
	s.AddExpense(100.0, "Coursera", Education, day20200101)
	if err := s.UpdateExpense(1, 200.0, "Coursera Pro", Education, day20200101); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDeleteExpense(t *testing.T) {
	s := newService()
	s.AddExpense(100.0, "Coursera", Education, day20200101)
	if err := s.DeleteExpense(1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestListAll(t *testing.T) {
	s := newService()
	s.AddExpense(100.0, "Coursera", Education, day20200101)
	s.AddExpense(12, "Datacamp", Education, day20200202)
	s.AddExpense(12, "EFKA", Business, day20210202)

	expenses, err := s.ListAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(expenses) != 3 {
		t.Errorf("expected 3 expenses, got %d", len(expenses))
	}
}

func TestListByCategory(t *testing.T) {
	s := newService()
	s.AddExpense(100.0, "Coursera", Education, day20200101)
	s.AddExpense(12, "Datacamp", Education, day20200102)
	s.AddExpense(12, "EFKA", Business, day20200101)

	expenses, err := s.ListByCategory(Education)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(expenses) != 2 {
		t.Errorf("expected 2 expenses, got %d", len(expenses))
	}
}

func TestListByMonth(t *testing.T) {
	s := newService()
	s.AddExpense(100.0, "Coursera", Education, day20200101)
	s.AddExpense(12, "Datacamp", Education, day20200202)
	s.AddExpense(12, "EFKA", Business, day20210202)

	expenses, err := s.ListByMonth(2, 2020)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(expenses) != 1 {
		t.Errorf("expected 1 expense, got %d", len(expenses))
	}
}

func TestListByYear(t *testing.T) {
	s := newService()
	s.AddExpense(100.0, "Coursera", Education, day20200101)
	s.AddExpense(12, "Datacamp", Education, day20200202)
	s.AddExpense(12, "EFKA", Business, day20210202)

	expenses, err := s.ListByYear(2020)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(expenses) != 2 {
		t.Errorf("expected 2 expenses, got %d", len(expenses))
	}
}

func TestTotalExpense(t *testing.T) {
	s := newService()
	s.AddExpense(100.0, "Coursera", Education, day20200101)
	s.AddExpense(12, "Datacamp", Education, day20200202)
	s.AddExpense(12, "EFKA", Business, day20210202)

	education, _ := s.ListByCategory(Education)
	if total := s.TotalExpense(education); total != 112 {
		t.Errorf("expected total 112, got %v", total)
	}
}
