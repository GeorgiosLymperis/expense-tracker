package httpapi

import (
	"encoding/json"
	"errors"
	"expense-tracker/internal/expense"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type api struct {
	svc expense.Service
	mu  sync.Mutex
}

type addExpenseRequest struct {
	Amount      expense.Amount   `json:"amount"`
	Description string           `json:"description"`
	Category    expense.Category `json:"category"`
	Date        string           `json:"date"`
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sr, r)

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sr.status, time.Since(start))
	})
}

func newHandler(svc expense.Service) http.Handler {
	a := &api{svc: svc}

	mux := http.NewServeMux()
	mux.Handle("GET /expenses", logging(http.HandlerFunc(a.handleList)))
	mux.Handle("GET /expenses/{id}", logging(http.HandlerFunc(a.handleGet)))
	mux.Handle("POST /expenses", logging(http.HandlerFunc(a.handleAdd)))
	mux.Handle("PUT /expenses/{id}", logging(http.HandlerFunc(a.handleUpdate)))
	mux.Handle("DELETE /expenses/{id}", logging(http.HandlerFunc(a.handleDelete)))
	return mux
}

func Run(svc expense.Service, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("serving on http://%s", ln.Addr())

	return http.Serve(ln, newHandler(svc))
}

func (a *api) handleList(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	expenses, err := a.svc.ListAll()
	if err != nil {
		http.Error(w, "error fetching expenses", http.StatusInternalServerError)
		return
	}
	if expenses == nil {
		expenses = []expense.Expense{}
	}
	if err := json.NewEncoder(w).Encode(expenses); err != nil {
		http.Error(w, "error encoding expenses", http.StatusInternalServerError)
		return
	}
}

func (a *api) handleGet(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	e, err := a.svc.GetExpense(id)
	if err != nil {
		if errors.Is(err, expense.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "error fetching expense", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(e); err != nil {
		http.Error(w, "error encoding expense", http.StatusInternalServerError)
		return
	}
}

func validateAddRequest(w http.ResponseWriter, req *addExpenseRequest) bool {
	if req.Amount < 0 {
		http.Error(w, "amount should be positive", http.StatusBadRequest)
		return false
	}

	if req.Description == "" {
		req.Description = "Unknown"
	}

	if req.Category == "" {
		req.Category = "Unknown"
	}
	if !req.Category.IsValid() {
		http.Error(w, "invalid category", http.StatusBadRequest)
		return false
	}

	if req.Date == "" {
		req.Date = expense.Date(time.Now()).String()
	}

	return true
}

func (a *api) handleAdd(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var req addExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ok := validateAddRequest(w, &req)
	if !ok {
		return
	}

	expenseDate, err := expense.ParseDate(req.Date)
	if err != nil {
		http.Error(w, "date should be yyyy-mm-dd", http.StatusBadRequest)
		return
	}
	if err := a.svc.AddExpense(req.Amount,
		req.Description,
		req.Category,
		expenseDate); err != nil {
		http.Error(w, "adding expense failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *api) handleUpdate(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req addExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ok := validateAddRequest(w, &req)
	if !ok {
		return
	}

	expenseDate, err := expense.ParseDate(req.Date)
	if err != nil {
		http.Error(w, "date should be yyyy-mm-dd", http.StatusBadRequest)
		return
	}

	if err := a.svc.UpdateExpense(id,
		req.Amount,
		req.Description,
		req.Category,
		expenseDate); err != nil {
		if errors.Is(err, expense.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "updating expense failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *api) handleDelete(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := a.svc.DeleteExpense(id); err != nil {
		if errors.Is(err, expense.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "deleting expense failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
