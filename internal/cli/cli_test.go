package cli

import (
	"errors"
	"expense-tracker/internal/expense"
	"testing"
	"time"
)

func TestRunNoArgs(t *testing.T) {
	if err := Run(&fakeService{}, []string{}); err == nil {
		t.Error("expected error for no args, got nil")
	}
}

func TestRunInvalidCommand(t *testing.T) {
	if err := Run(&fakeService{}, []string{"notacommand"}); err == nil {
		t.Error("expected error for invalid command, got nil")
	}
}

func TestHelpCommand(t *testing.T) {
	if err := Run(&fakeService{}, []string{"help"}); err != nil {
		t.Errorf("expected no error for help, got %v", err)
	}
}

func TestAddRouting(t *testing.T) {
	svc := &fakeService{}
	err := Run(svc, []string{"add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// success message is printed to stdout; routing is verified via recorded fields below
	if svc.addedAmount != 20 {
		t.Errorf("expected amount 20, got %v", svc.addedAmount)
	}
	if svc.addedDescription != "Lunch" {
		t.Errorf("expected description 'Lunch', got %v", svc.addedDescription)
	}
}

func TestAddCategory(t *testing.T) {
	svc := &fakeService{}
	Run(svc, []string{"add", "--description", "Coursera", "--amount", "50", "--category", "Education", "--date", "2026-05-28"})
	if svc.addedCategory != expense.EDUCATION {
		t.Errorf("expected category Education, got %v", svc.addedCategory)
	}
}

func TestAddDate(t *testing.T) {
	svc := &fakeService{}
	Run(svc, []string{"add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28"})
	want, _ := expense.ParseDate("2026-05-28")
	if svc.addedDate != want {
		t.Errorf("expected date 2026-05-28, got %v", svc.addedDate)
	}
}

func TestDeleteRouting(t *testing.T) {
	svc := &fakeService{}
	if err := Run(svc, []string{"delete", "--id", "3"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.deletedID != 3 {
		t.Errorf("expected deleted id 3, got %v", svc.deletedID)
	}
}

func TestDeleteRequiresId(t *testing.T) {
	if err := Run(&fakeService{}, []string{"delete"}); err == nil {
		t.Error("expected error when --id omitted, got nil")
	}
}

func TestUpdateRouting(t *testing.T) {
	svc := &fakeService{}
	if err := Run(svc, []string{"update", "--id", "2", "--description", "Brunch", "--amount", "15", "--date", "2026-05-28"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.updatedID != 2 {
		t.Errorf("expected updated id 2, got %v", svc.updatedID)
	}
	if svc.updatedAmount != 15 {
		t.Errorf("expected updated amount 15, got %v", svc.updatedAmount)
	}
	if svc.updatedDescription != "Brunch" {
		t.Errorf("expected updated description 'Brunch', got %v", svc.updatedDescription)
	}
}

func TestUpdateRequiresId(t *testing.T) {
	if err := Run(&fakeService{}, []string{"update", "--description", "Brunch", "--amount", "15", "--date", "2026-05-28"}); err == nil {
		t.Error("expected error when --id omitted, got nil")
	}
}

func TestListCallsPrint(t *testing.T) {
	lunch, _ := expense.NewExpense(20, "Lunch", expense.FOOD,
		expense.Date(time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC)))
	svc := &fakeService{listResult: []expense.Expense{lunch}}

	if err := Run(svc, []string{"list"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(svc.printedExpenses) != 1 {
		t.Errorf("expected 1 printed expense, got %d", len(svc.printedExpenses))
	}
}

func TestListByMonthPassesCorrectMonth(t *testing.T) {
	svc := &fakeService{}
	if err := Run(svc, []string{"list-month", "--month", "3", "--year", "2026"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.listByMonthValue != time.March {
		t.Errorf("expected month March, got %v", svc.listByMonthValue)
	}
	if svc.listByMonthYear != 2026 {
		t.Errorf("expected year 2026, got %v", svc.listByMonthYear)
	}
}

func TestListByCategoryPassesCorrectCategory(t *testing.T) {
	svc := &fakeService{}
	if err := Run(svc, []string{"list-category", "--category", "Food"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.listByCat != expense.FOOD {
		t.Errorf("expected category Food, got %v", svc.listByCat)
	}
}

func TestListByYearPassesCorrectYear(t *testing.T) {
	svc := &fakeService{}
	if err := Run(svc, []string{"list-year", "--year", "2025"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.listByYear != 2025 {
		t.Errorf("expected year 2025, got %v", svc.listByYear)
	}
}

func TestSummaryCallsListAll(t *testing.T) {
	svc := &fakeService{}
	if err := Run(svc, []string{"summary"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSummaryByMonthPassesMonth(t *testing.T) {
	svc := &fakeService{}
	if err := Run(svc, []string{"summary", "--month", "5"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.listByMonthValue != time.May {
		t.Errorf("expected month May, got %v", svc.listByMonthValue)
	}
}

func TestServiceErrorPropagates(t *testing.T) {
	svc := &fakeService{cmdErr: errors.New("service error")}
	if err := Run(svc, []string{"add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28"}); err == nil {
		t.Error("expected service error to propagate, got nil")
	}
}
