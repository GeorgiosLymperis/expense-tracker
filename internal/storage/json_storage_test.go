package storage

import (
	"encoding/json"
	"expense-tracker/internal/expense"
	"io"
	"os"
	"testing"
	"time"
)

func TestNewJSONRepo(t *testing.T) {
	f, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatalf("Couldnt open temporary file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	repo := NewJSONRepo(f)
	if repo.file != f {
		t.Errorf("Repo file should be the same with temp")
	}
}

func TestJSONRepoSave(t *testing.T) {
	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))

	f, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatal("Couldnt open temporary file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	repo := NewJSONRepo(f)
	err = repo.Save(&businessExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", businessExpense)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	if amount := expenses[0].Amount; amount != 100.5 {
		t.Errorf("Expense amount should be 100.5. Got %v", amount)
	}

	if description := expenses[0].Description; description != "EFKA" {
		t.Errorf("Expense description should be EFKA. Got %v", description)
	}

	if category := expenses[0].Category; category != "Business" {
		t.Errorf("Expense category should be Business. Got %v", category)
	}

	if date := expenses[0].Date; date.String() != "2020-01-01" {
		t.Errorf("Expense date should be 2020-01-01. Got %v", date)
	}
}

func TestJSONRepoSaveOlderYearIsLast(t *testing.T) {
	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	f, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatal("Couldnt open temporary file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	repo := NewJSONRepo(f)

	err = repo.Save(&educationExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", educationExpense)
	}

	err = repo.Save(&businessExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", businessExpense)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	if amount := expenses[1].Amount; amount != 60.5 {
		t.Errorf("Expense amount should be 60.5. Got %v", amount)
	}

	if description := expenses[1].Description; description != "Coursera" {
		t.Errorf("Expense description should be Coursera. Got %v", description)
	}

	if category := expenses[1].Category; category != "Education" {
		t.Errorf("Expense category should be Education. Got %v", category)
	}

	if date := expenses[1].Date; date.String() != "2023-02-03" {
		t.Errorf("Expense date should be 2020-01-01. Got %v", date)
	}
}

func TestJSONRepoSaveOlderMonthIsLast(t *testing.T) {
	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	f, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatal("Couldnt open temporary file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	repo := NewJSONRepo(f)

	err = repo.Save(&educationExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", educationExpense)
	}

	err = repo.Save(&businessExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", businessExpense)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	if amount := expenses[1].Amount; amount != 60.5 {
		t.Errorf("Expense amount should be 60.5. Got %v", amount)
	}

	if description := expenses[1].Description; description != "Coursera" {
		t.Errorf("Expense description should be Coursera. Got %v", description)
	}

	if category := expenses[1].Category; category != "Education" {
		t.Errorf("Expense category should be Education. Got %v", category)
	}

	if date := expenses[1].Date; date.String() != "2023-02-03" {
		t.Errorf("Expense date should be 2020-01-01. Got %v", date)
	}
}

func TestJSONRepoSaveOlderDayIsLast(t *testing.T) {
	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	f, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatal("Couldnt open temporary file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	repo := NewJSONRepo(f)

	err = repo.Save(&educationExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", educationExpense)
	}

	err = repo.Save(&businessExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", businessExpense)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	if amount := expenses[1].Amount; amount != 60.5 {
		t.Errorf("Expense amount should be 60.5. Got %v", amount)
	}

	if description := expenses[1].Description; description != "Coursera" {
		t.Errorf("Expense description should be Coursera. Got %v", description)
	}

	if category := expenses[1].Category; category != "Education" {
		t.Errorf("Expense category should be Education. Got %v", category)
	}

	if date := expenses[1].Date; date.String() != "2023-02-03" {
		t.Errorf("Expense date should be 2020-01-01. Got %v", date)
	}
}

func TestJSONRepoDelete(t *testing.T) {
	businessExpense, _ := expense.NewBusinessExpense(
		100.50, "EFKA",
		expense.Date(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)))

	educationExpense, _ := expense.NewEducationExpense(
		60.50, "Coursera",
		expense.Date(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))

	f, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatal("Couldnt open temporary file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	repo := NewJSONRepo(f)

	err = repo.Save(&educationExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", educationExpense)
	}

	err = repo.Save(&businessExpense)
	if err != nil {
		t.Fatalf("Could not save expense %v", businessExpense)
	}

	err = repo.Delete(1)
	if err != nil {
		t.Fatalf("Could not delete expense %v", educationExpense)
	}

	byteValue, _ := io.ReadAll(f)
	var expenses []expenseRecord
	json.Unmarshal(byteValue, &expenses)

	if amount := expenses[0].Amount; amount != 60.5 {
		t.Errorf("Expense amount should be 60.5. Got %v", amount)
	}

	if description := expenses[0].Description; description != "Coursera" {
		t.Errorf("Expense description should be Coursera. Got %v", description)
	}

	if category := expenses[0].Category; category != "Education" {
		t.Errorf("Expense category should be Education. Got %v", category)
	}

	if date := expenses[0].Date; date.String() != "2023-02-03" {
		t.Errorf("Expense date should be 2020-01-01. Got %v", date)
	}
}

func TestJSONRepoUpdate(t *testing.T) {

}

func TestJSONRepoFindByDate(t *testing.T) {

}

func TestJSONRepoFindByMonth(t *testing.T) {

}

func TestJSONRepoFindByCategory(t *testing.T) {

}

func TestJSONRepoExport(t *testing.T) {

}
