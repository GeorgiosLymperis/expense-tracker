# Expense Tracker

A command-line application to manage your personal finances. Add, update, delete and summarise expenses stored as JSON, from the terminal or through an [HTTP API](#http-api).

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
- `--amount` — expense amount, with at most 2 decimals (e.g. `12.50`)
- `--category` — expense category (default: Unknown)
- `--date` — date in `YYYY-MM-DD` format (default: today)

### Update an expense

```bash
./tracker update --id 1 --description "Brunch" --amount 15 --category Food --date 2026-05-28
```

Only the flags you pass are changed. For example, `./tracker update --id 1 --amount 15` changes the amount and keeps the description, category and date.

### Delete an expense

```bash
./tracker delete --id 1
```

### List all expenses

```bash
./tracker list
```

Show only the most recent expenses with `--last`. It works with every `list` command:

```bash
./tracker list --last 20
./tracker list-category --category Food --last 5
```

The total row then covers only the expenses shown.

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

## HTTP API

Start the server:

```bash
./tracker serve
```

It listens on `http://localhost:8080` and uses the same `expenses.json` file as the CLI. Requests and responses are JSON. Errors are returned as a plain-text message with the matching status code.

| Method   | Path             | Description           |
|----------|------------------|-----------------------|
| `GET`    | `/expenses`      | List all expenses     |
| `GET`    | `/expenses/{id}` | Get one expense       |
| `POST`   | `/expenses`      | Add an expense        |
| `PUT`    | `/expenses/{id}` | Replace an expense    |
| `DELETE` | `/expenses/{id}` | Delete an expense     |

### Expense IDs

Every expense has an `id`, assigned when it is added and saved in `expenses.json`. It is the same ID the CLI shows, and it doesn't change when other expenses are added, updated or deleted. A new expense gets the highest ID in use plus one.

### GET /expenses

Returns all expenses in the order they are stored: new expenses are inserted by date, and an update keeps an expense in place. Returns `[]` when there are none.

```bash
curl http://localhost:8080/expenses
```

```json
[
  {"id": 1, "amount": 12.50, "date": "2026-05-01", "description": "Lunch", "category": "Food"},
  {"id": 2, "amount": 50.00, "date": "2026-05-28", "description": "Coursera", "category": "Education"}
]
```

**Responses**
- `200 OK`: array of expenses
- `500 Internal Server Error`: the expenses file could not be read

### GET /expenses/{id}

Returns the expense with the given ID.

```bash
curl http://localhost:8080/expenses/1
```

```json
{"id": 1, "amount": 12.50, "date": "2026-05-01", "description": "Lunch", "category": "Food"}
```

**Responses**
- `200 OK`: the expense
- `400 Bad Request`: ID is not a number
- `404 Not Found`: no expense with that ID
- `500 Internal Server Error`: the expenses file could not be read

### POST /expenses

Adds an expense.

**Body**

| Field         | Type   | Required | Default   | Notes                                  |
|---------------|--------|----------|-----------|----------------------------------------|
| `amount`      | number | yes      |           | 0 or more, at most 2 decimals          |
| `description` | string | no       | `Unknown` |                                        |
| `category`    | string | no       | `Unknown` | One of the [categories](#categories)   |
| `date`        | string | no       | today     | `YYYY-MM-DD`                           |

```bash
curl -X POST http://localhost:8080/expenses \
  -H 'Content-Type: application/json' \
  -d '{"amount": 12.5, "description": "Lunch", "category": "Food", "date": "2026-05-01"}'
```

**Responses**
- `201 Created`: no body
- `400 Bad Request`: invalid JSON, amount negative or with more than 2 decimals, unknown category, or date not in `YYYY-MM-DD` format
- `500 Internal Server Error`: the expense could not be saved

### PUT /expenses/{id}

Replaces the expense with the given ID. The body is the same as for `POST /expenses`, with the same defaults: a field you leave out is reset to its default, not kept from the old expense.

```bash
curl -X PUT http://localhost:8080/expenses/1 \
  -H 'Content-Type: application/json' \
  -d '{"amount": 15, "description": "Brunch", "category": "Food", "date": "2026-05-01"}'
```

**Responses**
- `200 OK`: no body
- `400 Bad Request`: ID is not a number, or the body is invalid (same rules as `POST`)
- `404 Not Found`: no expense with that ID
- `500 Internal Server Error`: the expense could not be saved

### DELETE /expenses/{id}

Deletes the expense with the given ID.

```bash
curl -X DELETE http://localhost:8080/expenses/1
```

**Responses**
- `204 No Content`: no body
- `400 Bad Request`: ID is not a number
- `404 Not Found`: no expense with that ID
- `500 Internal Server Error`: the expense could not be deleted

## Categories

| Category      |
|---------------|
| Business      |
| Clothing      |
| Education     |
| Entertainment |
| Fitness       |
| Food          |
| Gifts         |
| Health        |
| Housing       |
| Savings       |
| Social        |
| Subscriptions |
| Supplements   |
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
