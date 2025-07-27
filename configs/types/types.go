package types

type (
	MemberRole  string
	Ordering    string
	AccountType string
)

const (
	OrderByAsc  Ordering = "asc"
	OrderByDesc Ordering = "desc"
)

func ValidateOrderBy(val Ordering) bool {
	if val == OrderByAsc {
		return true
	} else if val == OrderByDesc {
		return true
	}

	return false
}
