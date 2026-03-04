package builders

import (
	"strings"

	"gorm.io/gorm"
)

// ParentMetadata metadata สำหรับ parent table
type ParentMetadata struct {
	TableName  string
	ForeignKey string
}

// ParentFilterBuilder สร้าง parent relation filter ด้วย JOIN
type ParentFilterBuilder struct{}

// NewParentFilterBuilder สร้าง ParentFilterBuilder instance ใหม่
func NewParentFilterBuilder() *ParentFilterBuilder {
	return &ParentFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (p *ParentFilterBuilder) Name() string {
	return "ParentFilter"
}

// Priority คืนลำดับการ execute
func (p *ParentFilterBuilder) Priority() int {
	return 3
}

// Apply ประยุกต์ parent filter ไปยัง GORM query
func (p *ParentFilterBuilder) Apply(
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

	if len(q.FilterNestedParentBy) == 0 || len(q.FilterNestedParent) == 0 {
		return db, nil
	}

	if len(q.FilterNestedParentBy) != len(q.FilterNestedParent) {
		return nil, ErrFilterLengthMismatch
	}

	condition := q.FilterNestedParentCondition
	if condition == "" {
		condition = "and"
	}

	// Group by parent alias
	parentGroups := p.groupByParent(q.FilterNestedParentBy, q.FilterNestedParent)

	for parentAlias, filters := range parentGroups {
		// Join กับ parent table
		parentMeta := p.getParentMetadata(parentAlias, o.ParentTable)

		alias := o.TableAlias
		if alias == "" {
			alias = "main"
		}

		db = db.Joins(
			"JOIN " + parentMeta.TableName + " " + parentAlias +
				" ON " + parentAlias + ".id = " + alias + "." + parentMeta.ForeignKey,
		)

		// Apply filters
		for _, f := range filters {
			column := p.buildParentColumn(parentAlias, f.Column)

			if condition == "or" {
				db = db.Or(column+" IN ?", f.Values)
			} else {
				db = db.Where(column+" IN ?", f.Values)
			}
		}
	}

	return db, nil
}

// groupByParent จัดกลุ่ม filters ตาม parent alias
func (p *ParentFilterBuilder) groupByParent(
	filterNestedParentBy []string,
	filterNestedParent [][]interface{},
) map[string][]RelationFilter {
	groups := make(map[string][]RelationFilter)

	for i, path := range filterNestedParentBy {
		parts := strings.SplitN(path, ".", 2)
		if len(parts) != 2 {
			continue
		}

		parentAlias := parts[0]
		column := parts[1]

		groups[parentAlias] = append(groups[parentAlias], RelationFilter{
			Column: column,
			Values: filterNestedParent[i],
		})
	}

	return groups
}

// buildParentColumn สร้าง column expression สำหรับ parent
func (p *ParentFilterBuilder) buildParentColumn(parentAlias, column string) string {
	return parentAlias + "." + column
}

// getParentMetadata ดึง metadata ของ parent table
func (p *ParentFilterBuilder) getParentMetadata(parentAlias, parentTable string) ParentMetadata {
	// Default implementation - ใน production ควรใช้ relation registry
	return ParentMetadata{
		TableName:  parentTable,
		ForeignKey: parentAlias + "_id",
	}
}
