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
	if svc.addedAmount != expense.Amount(20*100) {
		t.Errorf("expected amount 20, got %v", svc.addedAmount)
	}
	if svc.addedDescription != "Lunch" {
		t.Errorf("expected description 'Lunch', got %v", svc.addedDescription)
	}
}

func TestAddCategory(t *testing.T) {
	svc := &fakeService{}
	Run(svc, []string{"add", "--description", "Coursera", "--amount", "50", "--category", "Education", "--date", "2026-05-28"})
	if svc.addedCategory != expense.Education {
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
	if svc.updatedAmount != expense.Amount(15*100) {
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
	lunch, _ := expense.NewExpense(expense.Amount(20*100), "Lunch", expense.Food,
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
	if svc.listByCat != expense.Food {
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

func TestUpdateKeepsFieldsNotGiven(t *testing.T) {
	day := expense.Date(time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC))
	rent, _ := expense.NewExpense(expense.Amount(500*100), "Rent", expense.Housing, day)
	svc := &fakeService{getResult: rent}

	if err := Run(svc, []string{"update", "--id", "1", "--amount", "450"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.updatedAmount != expense.Amount(450*100) {
		t.Errorf("expected amount 450.00, got %v", svc.updatedAmount)
	}
	if svc.updatedDescription != "Rent" || svc.updatedCategory != expense.Housing || svc.updatedDate != day {
		t.Errorf("expected other fields kept, got %q %v %v", svc.updatedDescription, svc.updatedCategory, svc.updatedDate)
	}
}

func TestUpdateUnknownIDFails(t *testing.T) {
	svc := &fakeService{getErr: expense.ErrNotFound}
	if err := Run(svc, []string{"update", "--id", "9", "--amount", "1"}); err == nil {
		t.Error("expected error for unknown id, got nil")
	}
	if svc.updatedID != 0 {
		t.Error("expected no update for unknown id")
	}
}

func TestAddRejectsBadAmount(t *testing.T) {
	svc := &fakeService{}
	if err := Run(svc, []string{"add", "--amount", "12.345"}); err == nil {
		t.Error("expected error for amount with 3 decimals, got nil")
	}
}

func TestListLastShowsMostRecent(t *testing.T) {
	var all []expense.Expense
	for d := 1; d <= 5; d++ {
		e, _ := expense.NewExpense(expense.Amount(d*100), "E", expense.Food,
			expense.Date(time.Date(2026, 5, d, 0, 0, 0, 0, time.UTC)))
		all = append(all, e)
	}

	for _, cmd := range []string{"list", "list-category", "list-month", "list-year"} {
		svc := &fakeService{listResult: all}
		if err := Run(svc, []string{cmd, "--last", "2"}); err != nil {
			t.Fatalf("%s: unexpected error: %v", cmd, err)
		}
		if len(svc.printedExpenses) != 2 {
			t.Fatalf("%s: expected 2 printed expenses, got %d", cmd, len(svc.printedExpenses))
		}
		if svc.printedExpenses[0].ExpenseDate() != all[3].ExpenseDate() ||
			svc.printedExpenses[1].ExpenseDate() != all[4].ExpenseDate() {
			t.Errorf("%s: expected the 2 most recent expenses", cmd)
		}
	}
}

func TestListLastLargerThanListShowsAll(t *testing.T) {
	e, _ := expense.NewExpense(expense.Amount(100), "E", expense.Food, expense.Date(time.Now()))
	svc := &fakeService{listResult: []expense.Expense{e}}
	if err := Run(svc, []string{"list", "--last", "20"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(svc.printedExpenses) != 1 {
		t.Errorf("expected 1 printed expense, got %d", len(svc.printedExpenses))
	}
}

func TestListLastNegativeFails(t *testing.T) {
	if err := Run(&fakeService{}, []string{"list", "--last", "-1"}); err == nil {
		t.Error("expected error for negative --last, got nil")
	}
}
