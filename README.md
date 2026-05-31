# Expense Tracker

A command-line application to manage your personal finances. Add, update, delete and summarise expenses stored as JSON.

## Build

```bash
go build -o tracker ./cmd/expense-tracker
```

## Usage

```
./tracker <command> [flags]
```

## Commands

### Add an expense

```bash
./tracker add --description "Lunch" --amount 20 --category Food --date 2026-05-28
```

Flags:
- `--description` — expense description
- `--amount` — expense amount
- `--category` — expense category (default: Unknown)
- `--date` — date in `YYYY-MM-DD` format (default: today)

### Update an expense

```bash
./tracker update --id 1 --description "Brunch" --amount 15 --category Food --date 2026-05-28
```

### Delete an expense

```bash
./tracker delete --id 1
```

### List all expenses

```bash
./tracker list
```

### List by category

```bash
./tracker list-category --category Food
```

### List by month

```bash
./tracker list-month --month 5 --year 2026
```

### List by year

```bash
./tracker list-year --year 2026
```

### Summary

Total of all expenses:

```bash
./tracker summary
```

Total for a specific month of the current year:

```bash
./tracker summary --month 5
```

### Export to CSV

```bash
./tracker export
```

Exports all expenses to `expenses.csv` in the current directory.

### Help

```bash
./tracker help
```

## Categories

| Category      |
|---------------|
| Business      |
| Clothing      |
| Education     |
| Entertainment |
| Food          |
| Gifts         |
| Health        |
| Housing       |
| Savings       |
| Social        |
| Subscriptions |
| Transport     |
| Travel        |
| Unknown       |
| Utilities     |

## Example session

```bash
$ ./tracker add --description "Lunch" --amount 20 --category Food --date 2026-05-28
Expense added successfully.

$ ./tracker add --description "Coursera" --amount 50 --category Education --date 2026-05-28
Expense added successfully.

$ ./tracker list
| ID: 1 | 2026-05-28 | Food      | Lunch    | 20.00 |
| ID: 2 | 2026-05-28 | Education | Coursera | 50.00 |

Total Amount: 70

$ ./tracker summary
Total expenses: €70.00

$ ./tracker summary --month 5
Total expenses for May: €70.00

$ ./tracker delete --id 1
Expense deleted successfully.

$ ./tracker summary
Total expenses: €50.00
```

## Running tests

```bash
# All tests
go test ./...

# With coverage
go test ./... -coverprofile=coverage.txt && go tool cover -html=coverage.txt
```

## Project

https://roadmap.sh/projects/expense-tracker
