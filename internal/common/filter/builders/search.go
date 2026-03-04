package builders

import (
	"strings"

	"gorm.io/gorm"
)

// SearchFilterBuilder สร้าง search filter ด้วย LIKE/ILIKE
type SearchFilterBuilder struct{}

// NewSearchFilterBuilder สร้าง SearchFilterBuilder instance ใหม่
func NewSearchFilterBuilder() *SearchFilterBuilder {
	return &SearchFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (s *SearchFilterBuilder) Name() string {
	return "SearchFilter"
}

// Priority คืนลำดับการ execute
func (s *SearchFilterBuilder) Priority() int {
	return 5
}

// Apply ประยุกต์ search filter ไปยัง GORM query
func (s *SearchFilterBuilder) Apply(
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

	if q.Search == "" || len(q.SearchBy) == 0 {
		return db, nil
	}

	searchPattern := "%" + q.Search + "%"

	var conditions []string
	var args []interface{}

	for _, field := range q.SearchBy {
		column := s.buildSearchColumn(field, o.TableAlias)

		// ตรวจสอบว่าเป็น JSON field หรือไม่
		if strings.Contains(field, ".") && !strings.Contains(field, "->") {
			// JSON field: column.jsonKey
			parts := strings.SplitN(field, ".", 2)
			jsonColumn := parts[0]
			jsonKey := parts[1]

			alias := o.TableAlias
			if alias == "" {
				alias = "main"
			}
			column = alias + "." + jsonColumn + "->> '" + jsonKey + "'"
		}

		conditions = append(conditions, "CAST("+column+" AS TEXT) ILIKE ?")
		args = append(args, searchPattern)
	}

	if len(conditions) > 0 {
		db = db.Where("("+strings.Join(conditions, " OR ")+")", args...)
	}

	return db, nil
}

// buildSearchColumn สร้าง column expression สำหรับ search
func (s *SearchFilterBuilder) buildSearchColumn(field, alias string) string {
	if alias != "" && !strings.Contains(field, ".") {
		return alias + "." + field
	}
	return field
}
