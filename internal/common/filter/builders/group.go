package builders

import (
	"strings"

	"gorm.io/gorm"
)

// GroupFilterBuilder สร้าง group filter ด้วย DISTINCT ON (PostgreSQL)
type GroupFilterBuilder struct{}

// NewGroupFilterBuilder สร้าง GroupFilterBuilder instance ใหม่
func NewGroupFilterBuilder() *GroupFilterBuilder {
	return &GroupFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (g *GroupFilterBuilder) Name() string {
	return "GroupFilter"
}

// Priority คืนลำดับการ execute
func (g *GroupFilterBuilder) Priority() int {
	return 7
}

// Apply ประยุกต์ group filter ไปยัง GORM query
func (g *GroupFilterBuilder) Apply(
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

	if len(q.GroupBy) == 0 {
		return db, nil
	}

	alias := o.TableAlias
	if alias == "" {
		alias = "main"
	}

	// PostgreSQL DISTINCT ON
	groupColumns := make([]string, len(q.GroupBy))
	for i, field := range q.GroupBy {
		groupColumns[i] = alias + "." + field
	}

	// สร้าง subquery สำหรับ DISTINCT ON
	if q.GroupSortBy != "" {
		sortDirection := "DESC"
		if q.GroupSort == "min" {
			sortDirection = "ASC"
		}

		orderByClause := strings.Join(groupColumns, ", ") + ", " +
			alias + "." + q.GroupSortBy + " " + sortDirection

		subquery := db.Session(&gorm.Session{NewDB: true}).
			Table(alias + " AS " + alias).
			Select("DISTINCT ON (" + strings.Join(groupColumns, ", ") + ") id").
			Order(orderByClause)

		db = db.Where(alias+".id IN (?)", subquery)
	}

	return db, nil
}
