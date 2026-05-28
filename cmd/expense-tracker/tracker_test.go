package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var trackerBin string

func TestMain(m *testing.M) {
	bin, err := os.CreateTemp("", "tracker-test-*")
	if err != nil {
		panic("could not create temp file for binary: " + err.Error())
	}
	bin.Close()
	trackerBin = bin.Name()

	out, err := exec.Command("go", "build", "-o", trackerBin, "expense-tracker/cmd/expense-tracker").CombinedOutput()
	if err != nil {
		panic("build failed: " + string(out))
	}

	code := m.Run()
	os.Remove(trackerBin)
	os.Exit(code)
}

func tempWorkDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "tracker-work-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func run(t *testing.T, dir string, args ...string) (stdout string, ok bool) {
	t.Helper()
	cmd := exec.Command(trackerBin, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err == nil
}

func TestAddExpense(t *testing.T) {
	dir := tempWorkDir(t)
	out, ok := run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")
	if !ok {
		t.Fatal("add command failed")
	}
	if !strings.Contains(out, "added successfully") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestAddExpenseNegativeAmount(t *testing.T) {
	dir := tempWorkDir(t)
	_, ok := run(t, dir, "add", "--description", "Lunch", "--amount", "-5", "--date", "2026-05-28")
	if ok {
		t.Fatal("expected failure for negative amount, got success")
	}
}

func TestListAll(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")
	run(t, dir, "add", "--description", "Dinner", "--amount", "10", "--date", "2026-05-28")

	out, ok := run(t, dir, "list")
	if !ok {
		t.Fatal("list command failed")
	}
	if !strings.Contains(out, "Lunch") {
		t.Errorf("expected output to contain 'Lunch', got: %s", out)
	}
	if !strings.Contains(out, "Dinner") {
		t.Errorf("expected output to contain 'Dinner', got: %s", out)
	}
}

func TestSummaryTotal(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")
	run(t, dir, "add", "--description", "Dinner", "--amount", "10", "--date", "2026-05-28")

	out, ok := run(t, dir, "summary")
	if !ok {
		t.Fatal("summary command failed")
	}
	if !strings.Contains(out, "Total expenses: €30.00") {
		t.Errorf("expected 'Total expenses: €30.00', got: %s", out)
	}
}

func TestSummaryByMonth(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")
	run(t, dir, "add", "--description", "OldExpense", "--amount", "99", "--date", "2026-03-01")

	out, ok := run(t, dir, "summary", "--month", "5")
	if !ok {
		t.Fatal("summary --month command failed")
	}
	if !strings.Contains(out, "Total expenses for May: €20.00") {
		t.Errorf("expected 'Total expenses for May: €20.00', got: %s", out)
	}
}

func TestDeleteExpense(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")

	out, ok := run(t, dir, "delete", "--id", "1")
	if !ok {
		t.Fatal("delete command failed")
	}
	if !strings.Contains(out, "deleted successfully") {
		t.Errorf("expected success message, got: %s", out)
	}

	out, _ = run(t, dir, "summary")
	if !strings.Contains(out, "Total expenses: €0.00") {
		t.Errorf("expected €0.00 after delete, got: %s", out)
	}
}

func TestDeleteRequiresId(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")

	_, ok := run(t, dir, "delete")
	if ok {
		t.Fatal("expected failure when --id is omitted, got success")
	}
}

func TestUpdateExpense(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")

	_, ok := run(t, dir, "update", "--id", "1", "--description", "Brunch", "--amount", "15", "--date", "2026-05-28")
	if !ok {
		t.Fatal("update command failed")
	}

	out, _ := run(t, dir, "list")
	if !strings.Contains(out, "Brunch") {
		t.Errorf("expected 'Brunch' after update, got: %s", out)
	}
	if strings.Contains(out, "Lunch") {
		t.Errorf("expected 'Lunch' to be gone after update, got: %s", out)
	}
}

func TestUpdateRequiresId(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")

	_, ok := run(t, dir, "update", "--description", "Brunch", "--amount", "15", "--date", "2026-05-28")
	if ok {
		t.Fatal("expected failure when --id is omitted, got success")
	}
}

func TestExportCSV(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Lunch", "--amount", "20", "--date", "2026-05-28")

	_, ok := run(t, dir, "export")
	if !ok {
		t.Fatal("export command failed")
	}

	data, err := os.ReadFile(dir + "/expenses.csv")
	if err != nil {
		t.Fatalf("CSV file not created: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Lunch") {
		t.Errorf("expected CSV to contain 'Lunch', got: %s", content)
	}
}

func TestNoArgs(t *testing.T) {
	dir := tempWorkDir(t)
	_, ok := run(t, dir)
	if ok {
		t.Fatal("expected failure when no args, got success")
	}
}

func TestInvalidCommand(t *testing.T) {
	dir := tempWorkDir(t)
	_, ok := run(t, dir, "notacommand")
	if ok {
		t.Fatal("expected failure for invalid command, got success")
	}
}

func TestFilterByCategory(t *testing.T) {
	dir := tempWorkDir(t)
	run(t, dir, "add", "--description", "Coursera", "--amount", "50", "--category", "Education", "--date", "2026-05-28")
	run(t, dir, "add", "--description", "EFKA", "--amount", "100", "--category", "Business", "--date", "2026-05-28")

	out, ok := run(t, dir, "list-category", "--category", "Education")
	if !ok {
		t.Fatal("list-category command failed")
	}
	if !strings.Contains(out, "Coursera") {
		t.Errorf("expected 'Coursera' in output, got: %s", out)
	}
	if strings.Contains(out, "EFKA") {
		t.Errorf("expected 'EFKA' to be excluded, got: %s", out)
	}
}
