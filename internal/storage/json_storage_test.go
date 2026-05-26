package storage

import (
	"encoding/csv"
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

func tempCSVFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "example-*.csv")
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
	if repo.GetFile() != f {
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
	if err := repo.Update(1, &updatedExpense); err != nil {
		{
			t.Errorf("Error in updating")
		}
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	compareRecExpense(t, expenses[0], updatedExpense)
}

func TestJSONRepoFindByDate(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	businessExpense_2, _ := expense.NewBusinessExpense(
		110.50, "BANK",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)
	repo.Save(&businessExpense_2)

	date, err := time.Parse("2006-01-02", "2023-02-01")
	if err != nil {
		t.Fatalf("Could not parse date: %v", err)
	}
	expenses, err := repo.FindByDate(expense.Date(date))
	if err != nil {
		t.Errorf("Error in finding by date: %v", err)
	}

	compareRecExpense(t, NewExpenseRecord(&educationExpense), expenses[0])
	compareRecExpense(t, NewExpenseRecord(&businessExpense_2), expenses[1])
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
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&businessExpense)

	c := tempCSVFile(t)
	if err := repo.ExportCSV(c); err != nil {
		t.Fatalf("Could not export to CSV: %v", err)
	}
	c.Seek(0, 0)
	r := csv.NewReader(c)

	header, err := r.Read()
	if err != nil {
		t.Fatalf("Could not read header: %v", err)
	}

	if header[0] != "Date" {
		t.Errorf("Expected header to be Date, got %v", header[0])
	}
	if header[1] != "Category" {
		t.Errorf("Expected header to be Category, got %v", header[1])
	}
	if header[2] != "Description" {
		t.Errorf("Expected header to be Description, got %v", header[2])
	}
	if header[3] != "Amount" {
		t.Errorf("Expected header to be Amount, got %v", header[3])
	}

	row, err := r.Read()
	if err != nil {
		t.Fatalf("Could not read row: %v", err)
	}

	if row[0] != "2023-01-01" {
		t.Errorf("Expected date to be 2023-01-01, got %v", row[0])
	}

	if row[1] != "Business" {
		t.Errorf("Expected category to be Business, got %v", row[1])
	}

	if row[2] != "EFKA" {
		t.Errorf("Expected description to be EFKA, got %v", row[2])
	}

	if row[3] != "100.5" {
		t.Errorf("Expected amount to be 100.5, got %v", row[3])
	}
}

func TestListAll(t *testing.T) {
	f := tempJSONFile(t)

	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	businessExpense_2, _ := expense.NewBusinessExpense(
		110.50, "BANK",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	repo := NewJSONRepo(f)
	repo.Save(&educationExpense)
	repo.Save(&businessExpense)
	repo.Save(&businessExpense_2)

	expenses, err := repo.ListAll()
	if err != nil {
		t.Errorf("Error in finding by month")
	}

	compareRecExpense(t, NewExpenseRecord(&businessExpense_2), expenses[0])
	compareRecExpense(t, NewExpenseRecord(&educationExpense), expenses[1])
	compareRecExpense(t, NewExpenseRecord(&businessExpense), expenses[2])
}