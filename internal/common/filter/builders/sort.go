package builders

import (
	"strings"

	"gorm.io/gorm"
)

// SortFilterBuilder สร้าง sort filter ด้วย ORDER BY
type SortFilterBuilder struct{}

// NewSortFilterBuilder สร้าง SortFilterBuilder instance ใหม่
func NewSortFilterBuilder() *SortFilterBuilder {
	return &SortFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (s *SortFilterBuilder) Name() string {
	return "SortFilter"
}

// Priority คืนลำดับการ execute (หลังจาก group filter)
func (s *SortFilterBuilder) Priority() int {
	return 8
}

// Apply ประยุกต์ sort filter ไปยัง GORM query
func (s *SortFilterBuilder) Apply(
	db *gorm.DB,
	query interface{},
	opts interface{},
) (*gorm.DB, error) {
	q, ok := query.(*AdvanceFilterQuery)
	if !ok {
		return db, nil
	}

	o, ok := opts.(*FilterOptions)
	if !ok {
		o = &FilterOptions{}
	}

	if len(q.SortBy) == 0 {
		return db, nil
	}

	for i, field := range q.SortBy {
		direction := "ASC"
		if i < len(q.Sort) {
			if strings.ToLower(q.Sort[i]) == "desc" {
				direction = "DESC"
			}
		}

		column := s.buildSortColumn(field, o.TableAlias)
		db = db.Order(column + " " + direction)
	}

	return db, nil
}

// buildSortColumn สร้าง column expression สำหรับ sort
func (s *SortFilterBuilder) buildSortColumn(field, alias string) string {
	// รองรับ JSON field
	if strings.Contains(field, ".") && !strings.Contains(field, "->") {
		parts := strings.SplitN(field, ".", 2)
		if alias != "" {
			return alias + "." + parts[0] + "->> '" + parts[1] + "'"
		}
		return parts[0] + "->> '" + parts[1] + "'"
	}

	if alias != "" {
		return alias + "." + field
	}
	return field
}
