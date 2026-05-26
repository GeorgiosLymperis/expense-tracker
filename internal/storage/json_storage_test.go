package storage

import (
	"encoding/json"
	"expense-tracker/internal/expense"
	"io"
	"os"
	"testing"
	"time"
)

func compareRecExpense(t *testing.T, rec expenseRecord, e expense.Expense) {
	if rec.Amount != e.ExpenseAmount() {
		t.Errorf("Expense amount should be %v. Got %v", e.ExpenseAmount(), rec.Amount)
	}
	if rec.Date != e.ExpenseDate() {
		t.Errorf("Expense date should be %v. Got %v", e.ExpenseDate(), rec.Date)
	}
	if rec.Description != e.ExpenseDescription() {
		t.Errorf("Expense description should be %v. Got %v", e.ExpenseDescription(), rec.Description)
	}
	if rec.Category != e.ExpenseCategory() {
		t.Errorf("Expense category should be %v. Got %v", e.ExpenseCategory(), rec.Category)
	}
}

func tempJSONFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatal("could not create temporary file")
	}
	t.Cleanup(func() {
		f.Close()
		os.Remove(f.Name())
	})
	return f
}

func TestNewJSONRepo(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	if repo.file != f {
		t.Errorf("Repo file should be the same with temp")
	}
}

func TestJSONRepoSave(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)

	if err := repo.Save(&businessExpense); err != nil {
		t.Fatalf("Could not save expense %v", businessExpense)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[0], businessExpense)
}

func TestJSONRepoSaveOlderYearIsLast(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[1], educationExpense)

}

func TestJSONRepoSaveOlderMonthIsLast(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[1], educationExpense)
}

func TestJSONRepoSaveOlderDayIsLast(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[1], educationExpense)
}

func TestJSONRepoDelete(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)

	if err := repo.Delete(1); err != nil {
		t.Fatalf("Could not delete expense %v", educationExpense)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[0], educationExpense)
}

func TestJSONRepoUpdate(t *testing.T) {
	f := tempJSONFile(t)
	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	updatedExpense, _ := expense.NewEntertainmentExpense(
		25.00, "NETFLIX", expense.Date(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&businessExpense)
	if err := repo.Update(1, &updatedExpense); err != nil {{
		t.Errorf("Error in updating")
	}}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[0], updatedExpense)
}

func TestJSONRepoFindByDate(t *testing.T) {

}

func TestJSONRepoFindByMonth(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	businessExpense_2, _ := expense.NewBusinessExpense(
		110.50, "BANK",
		expense.Date(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)
	repo.Save(&businessExpense_2)


	expenses, err := repo.FindByMonth(2)
	if err != nil {
		t.Errorf("Error in finding by month")
	}

	compareRecExpense(t, NewExpenseRecord(&educationExpense), expenses[0])
	compareRecExpense(t, NewExpenseRecord(&businessExpense_2), expenses[1])
}

func TestJSONRepoFindByYear(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	businessExpense_2, _ := expense.NewBusinessExpense(
		110.50, "BANK",
		expense.Date(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)
	repo.Save(&businessExpense_2)


	expenses, err := repo.FindByYear(2023)
	if err != nil {
		t.Errorf("Error in finding by month")
	}

	compareRecExpense(t, NewExpenseRecord(&businessExpense), expenses[0])
	compareRecExpense(t, NewExpenseRecord(&educationExpense), expenses[1])
}

func TestJSONRepoFindByCategory(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	businessExpense_2, _ := expense.NewBusinessExpense(
		110.50, "BANK",
		expense.Date(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)
	repo.Save(&businessExpense_2)


	expenses, err := repo.FindByCategory(expense.BUSINESS)
	if err != nil {
		t.Errorf("Error in finding by category: %v", err)
	}

	compareRecExpense(t, NewExpenseRecord(&businessExpense), expenses[0])
	compareRecExpense(t, NewExpenseRecord(&businessExpense_2), expenses[1])

}

func TestJSONRepoFindByCategoryNotFound(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	businessExpense_2, _ := expense.NewBusinessExpense(
		110.50, "BANK",
		expense.Date(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)
	repo.Save(&businessExpense_2)


	expenses, err := repo.FindByCategory(expense.ENTERTAINMENT)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(expenses) != 0 {
		t.Errorf("Expected no expenses, got %v", len(expenses))
	}
}

func TestJSONRepoExport(t *testing.T) {

}
