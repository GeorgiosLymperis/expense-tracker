package expense

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Amount is a sum of money in cents. Whole cents avoid the rounding errors
// of floating point, so totals are always exact.
type Amount int64

// ParseAmount parses a decimal string such as "12", "12.5" or "-3.40".
// It rejects more than two decimal places instead of rounding them away.
func ParseAmount(s string) (Amount, error) {
	str := strings.TrimSpace(s)
	neg := strings.HasPrefix(str, "-")
	str = strings.TrimPrefix(str, "-")

	whole, frac, hasDot := strings.Cut(str, ".")
	if whole == "" || (hasDot && frac == "") || !isDigits(whole) || !isDigits(frac) {
		return 0, fmt.Errorf("invalid amount %q", s)
	}
	if len(frac) > 2 {
		return 0, fmt.Errorf("invalid amount %q: at most 2 decimal places", s)
	}

	units, err := strconv.ParseInt(whole, 10, 64)
	if err != nil || units > math.MaxInt64/100-1 {
		return 0, fmt.Errorf("invalid amount %q: too large", s)
	}
	cents := int64(0)
	if frac != "" {
		frac += strings.Repeat("0", 2-len(frac))
		cents, _ = strconv.ParseInt(frac, 10, 64)
	}

	a := Amount(units*100 + cents)
	if neg {
		a = -a
	}
	return a, nil
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// String formats the amount with two decimals, e.g. "12.50".
func (a Amount) String() string {
	sign := ""
	if a < 0 {
		sign = "-"
		a = -a
	}
	return fmt.Sprintf("%s%d.%02d", sign, a/100, a%100)
}

func (a Amount) MarshalJSON() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Amount) UnmarshalJSON(data []byte) error {
	parsed, err := ParseAmount(string(data))
	if err != nil {
		return err
	}
	*a = parsed
	return nil
}
