package httpapi

import (
	"encoding/json"
	"expense-tracker/internal/expense"
	"expense-tracker/internal/storage"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard) // silence request logging
	os.Exit(m.Run())
}

// newTestServer returns an API backed by a JSON file in a temp dir,
// pre-filled with the given expenses (IDs 1, 2, ...).
func newTestServer(t *testing.T, expenses ...expense.Expense) *httptest.Server {
	t.Helper()
	f, err := os.Create(filepath.Join(t.TempDir(), "expenses.json"))
	if err != nil {
		t.Fatalf("could not create file: %v", err)
	}
	t.Cleanup(func() { f.Close() })

	svc := expense.NewService(storage.NewJSONRepo(f))
	for _, e := range expenses {
		if err := svc.AddExpense(e.ExpenseAmount(), e.ExpenseDescription(), e.ExpenseCategory(), e.ExpenseDate()); err != nil {
			t.Fatalf("could not add expense: %v", err)
		}
	}

	srv := httptest.NewServer(newHandler(svc))
	t.Cleanup(srv.Close)
	return srv
}

func TestGetExpense(t *testing.T) {
	day := expense.Date(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	rent, _ := expense.NewExpense(expense.Amount(500*100), "Rent", expense.Housing, day)
	lunch, _ := expense.NewExpense(expense.Amount(12.5*100), "Lunch", expense.Food, day)
	srv := newTestServer(t, rent, lunch)

	resp, err := http.Get(srv.URL + "/expenses/2")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var got struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Date        string `json:"date"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("invalid JSON %s: %v", body, err)
	}
	if got.ID != 2 || got.Description != "Lunch" || got.Category != "Food" || got.Date != "2026-05-01" {
		t.Errorf("unexpected expense: %s", body)
	}

	var raw map[string]json.RawMessage
	json.Unmarshal(body, &raw)
	if string(raw["amount"]) != "12.50" {
		t.Errorf("expected amount 12.50, got %s", raw["amount"])
	}
}

func TestGetExpenseErrors(t *testing.T) {
	e, _ := expense.NewExpense(expense.Amount(1*100), "A", expense.Food, expense.Date(time.Now()))
	srv := newTestServer(t, e)

	cases := map[string]int{
		"/expenses/99":  http.StatusNotFound,
		"/expenses/abc": http.StatusBadRequest,
	}
	for path, want := range cases {
		resp, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != want {
			t.Errorf("GET %s: expected %d, got %d", path, want, resp.StatusCode)
		}
	}
}
