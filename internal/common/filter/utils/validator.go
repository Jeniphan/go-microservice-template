package utils

import "errors"

// AdvanceFilterQuery validation struct
// Duplicated from filter package to avoid circular import
type AdvanceFilterQueryValidation struct {
	FilterBy                    []string
	Filter                      [][]interface{}
	FilterCondition             string
	FilterNestedBy              []string
	FilterNested                [][]interface{}
	FilterNestedCondition       string
	FilterNestedParentBy        []string
	FilterNestedParent          [][]interface{}
	FilterNestedParentCondition string
	FilterM2MBy                 []string
	FilterM2M                   [][]interface{}
	FilterM2MCondition          string
	FilterM2MJoinAlias          string
	SearchBy                    []string
	Search                      string
	StartBy                     string
	Start                       string
	EndBy                       string
	End                         string
	StartAndEndCondition        string
	SortBy                      []string
	Sort                        []string
	Page                        int
	PerPage                     int
	GroupBy                     []string
	GroupSortBy                 string
	GroupSort                   string
	Preload                     []string
	Limit                       int
}

// Error constants สำหรับ validation errors
var (
	ErrFilterLengthMismatch = errors.New("filter_by and filter must have same length")
	ErrSortLengthMismatch   = errors.New("sort_by and sort must have same length")
	ErrInvalidCondition     = errors.New("condition must be 'and' or 'or'")
	ErrInvalidSortDirection = errors.New("sort direction must be 'asc' or 'desc'")
	ErrPageLimitExceeded    = errors.New("per_page cannot exceed 100")
)

// ValidateQuery ตรวจสอบ query parameters ว่าถูกต้องหรือไม่
// ตรวจสอบ:
// - filter_by และ filter มี length เท่ากัน
// - sort_by และ sort มี length เท่ากัน (ถ้ามีทั้งสอง)
// - filter_condition เป็น "and" หรือ "or"
// - sort direction เป็น "asc" หรือ "desc"
// - per_page ไม่เกิน 100
func ValidateQuery(query *AdvanceFilterQueryValidation) error {
	// Validate filter_by and filter length match
	if len(query.FilterBy) != len(query.Filter) {
		return ErrFilterLengthMismatch
	}

	// Validate sort_by and sort length match
	if len(query.SortBy) > 0 && len(query.Sort) > 0 && len(query.SortBy) != len(query.Sort) {
		return ErrSortLengthMismatch
	}

	// Validate filter_condition
	if query.FilterCondition != "" && !IsValidCondition(query.FilterCondition) {
		return ErrInvalidCondition
	}

	// Validate filter_nested_condition
	if query.FilterNestedCondition != "" && !IsValidCondition(query.FilterNestedCondition) {
		return ErrInvalidCondition
	}

	// Validate filter_nested_parent_condition
	if query.FilterNestedParentCondition != "" && !IsValidCondition(query.FilterNestedParentCondition) {
		return ErrInvalidCondition
	}

	// Validate filter_m2m_condition
	if query.FilterM2MCondition != "" && !IsValidCondition(query.FilterM2MCondition) {
		return ErrInvalidCondition
	}

	// Validate start_and_end_condition
	if query.StartAndEndCondition != "" && !IsValidCondition(query.StartAndEndCondition) {
		return ErrInvalidCondition
	}

	// Validate sort directions
	for _, direction := range query.Sort {
		if !IsValidSortDirection(direction) {
			return ErrInvalidSortDirection
		}
	}

	// Validate group_sort
	if query.GroupSort != "" && query.GroupSort != "max" && query.GroupSort != "min" {
		return ErrInvalidSortDirection
	}

	// Validate per_page limit
	if query.PerPage > 100 {
		return ErrPageLimitExceeded
	}

	return nil
}

// IsValidCondition ตรวจสอบว่า condition เป็น "and" หรือ "or" หรือไม่
// เป็น case-insensitive
func IsValidCondition(condition string) bool {
	lowerCondition := condition
	return lowerCondition == "and" || lowerCondition == "or"
}

// IsValidSortDirection ตรวจสอบว่า direction เป็น "asc" หรือ "desc" หรือไม่
// เป็น case-insensitive
func IsValidSortDirection(direction string) bool {
	lowerDirection := direction
	return lowerDirection == "asc" || lowerDirection == "desc"
}
