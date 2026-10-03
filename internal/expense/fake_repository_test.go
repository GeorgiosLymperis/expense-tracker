package expense

import (
	"fmt"
	"os"
	"slices"
	"time"
)

type fakeRepo struct {
	records []Expense
}

func (r *fakeRepo) Add(e *Expense) error {
	r.records = append(r.records, *e)
	return nil
}

func (r *fakeRepo) Get(id int) (Expense, error) {
	if id < 1 || id > len(r.records) {
		return Expense{}, fmt.Errorf("invalid id %d", id)
	}
	return r.records[id-1], nil
}

func (r *fakeRepo) Delete(id int) error {
	if id < 1 || id > len(r.records) {
		return fmt.Errorf("invalid id %d", id)
	}
	r.records = slices.Delete(r.records, id-1, id)
	return nil
}

func (r *fakeRepo) Update(id int, e *Expense) error {
	if id < 1 || id > len(r.records) {
		return fmt.Errorf("invalid id %d", id)
	}
	r.records[id-1] = *e
	return nil
}

func (r *fakeRepo) ListAll() ([]Expense, error) {
	return r.records, nil
}

func (r *fakeRepo) ListByDate(date Date) ([]Expense, error) {
	var out []Expense
	for _, e := range r.records {
		if e.date == date {
			out = append(out, e)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no expenses on %v", date)
	}
	return out, nil
}

func (r *fakeRepo) ListByMonth(month time.Month, year int) ([]Expense, error) {
	var out []Expense
	for _, e := range r.records {
		if e.date.Month() == month && e.date.Year() == year {
			out = append(out, e)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no expenses in %v %d", month, year)
	}
	return out, nil
}

func (r *fakeRepo) ListByYear(year int) ([]Expense, error) {
	var out []Expense
	for _, e := range r.records {
		if e.date.Year() == year {
			out = append(out, e)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no expenses in %d", year)
	}
	return out, nil
}

func (r *fakeRepo) ListByCategory(category Category) ([]Expense, error) {
	var out []Expense
	for _, e := range r.records {
		if e.category == category {
			out = append(out, e)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no expenses in category %v", category)
	}
	return out, nil
}

func (r *fakeRepo) ExportCSV(file *os.File) error {
	return nil
}
