package storage

import (
	"encoding/csv"
	"expense-tracker/internal/expense"
	"testing"
	"time"
)

var day_2020_01_01 = expense.Date(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
var day_2020_01_02 = expense.Date(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC))
var day_2020_02_02 = expense.Date(time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC))
var day_2021_02_02 = expense.Date(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC))

func TestSaveExpense(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	err := s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	if err != nil {
		t.Errorf("Error in saving: %v", err)
	}
}

func TestUpdateExpense(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	err := s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	if err != nil {
		t.Errorf("Error in saving: %v", err)
	}
	err = s.UpdateExpense(1, 200.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	if err != nil {
		t.Errorf("Error in updating: %v", err)
	}
}

func TestDelete(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	err := s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	if err != nil {
		t.Errorf("Error in saving: %v", err)
	}
	err = s.DeleteExpense(1)
	if err != nil {
		t.Errorf("Error in deleting: %v", err)
	}
}

func TestListByCategory(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	s.SaveExpense(12, "Datacamp", expense.EDUCATION, day_2020_01_01)
	s.SaveExpense(12, "EFKA", expense.BUSINESS, day_2020_01_01)

	expenses, err := s.ListByCategory(expense.EDUCATION)
	if err != nil {
		t.Errorf("Error in listing by category: %v", err)
	}
	if len(expenses) != 2 {
		t.Errorf("Expected 2 expenses, got %v", len(expenses))
	}
}

func TestListByMonth(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	s.SaveExpense(12, "Datacamp", expense.EDUCATION, day_2020_02_02)
	s.SaveExpense(12, "EFKA", expense.BUSINESS, day_2021_02_02)

	expenses, err := s.ListByMonth(2)
	if err != nil {
		t.Errorf("Error in listing by month: %v", err)
	}
	if len(expenses) != 2 {
		t.Errorf("Expected 1 expense, got %v", len(expenses))
	}
}

func TestListByYear(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	s.SaveExpense(12, "Datacamp", expense.EDUCATION, day_2020_02_02)
	s.SaveExpense(12, "EFKA", expense.BUSINESS, day_2021_02_02)

	expenses, err := s.ListByYear(2020)
	if err != nil {
		t.Errorf("Error in listing by year: %v", err)
	}
	if len(expenses) != 2 {
		t.Errorf("Expected 2 expenses, got %v", len(expenses))
	}
}

func TestTotalExpense(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	s.SaveExpense(12, "Datacamp", expense.EDUCATION, day_2020_02_02)
	s.SaveExpense(12, "EFKA", expense.BUSINESS, day_2021_02_02)
	educationList, _ := s.ListByCategory(expense.EDUCATION)
	total := s.TotalExpense(&educationList)

	if total != 112 {
		t.Errorf("Expected total was 112, got %v", total)
	}
}

func TestServiceListAll(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	s.SaveExpense(12, "Datacamp", expense.EDUCATION, day_2020_02_02)
	s.SaveExpense(12, "EFKA", expense.BUSINESS, day_2021_02_02)
	expenses, err := s.ListAll()
	if err != nil {
		t.Errorf("Error in listing all: %v", err)
	}
	if len(expenses) != 3 {
		t.Errorf("Expected 3 expenses, got %v", len(expenses))
	}
}

func TestServiceExportCSV(t *testing.T) {
	f := tempJSONFile(t)
	repo := NewJSONRepo(f)
	s := expense.NewService(repo)
	s.SaveExpense(100.0, "Coursera", expense.EDUCATION, day_2020_01_01)
	s.SaveExpense(12, "Datacamp", expense.EDUCATION, day_2020_02_02)
	s.SaveExpense(12, "EFKA", expense.BUSINESS, day_2021_02_02)
	c := tempCSVFile(t)
	err := s.ExportCSV(c)
	if err != nil {
		t.Errorf("Error in exporting to CSV: %v", err)
	}

	c.Seek(0, 0)
	r := csv.NewReader(c)
	header, err := r.Read()
	if err != nil {
		t.Errorf("Error in reading header: %v", err)
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
}
