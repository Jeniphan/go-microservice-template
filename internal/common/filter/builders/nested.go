package builders

import (
	"strings"

	"gorm.io/gorm"
)

// RelationFilter เก็บข้อมูล filter ของ relation
type RelationFilter struct {
	Column string
	Values []interface{}
}

// RelationMetadata metadata สำหรับ relation
type RelationMetadata struct {
	TableName  string
	ForeignKey string
}

// NestedFilterBuilder สร้าง nested relation filter ด้วย EXISTS subquery
type NestedFilterBuilder struct{}

// NewNestedFilterBuilder สร้าง NestedFilterBuilder instance ใหม่
func NewNestedFilterBuilder() *NestedFilterBuilder {
	return &NestedFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (n *NestedFilterBuilder) Name() string {
	return "NestedFilter"
}

// Priority คืนลำดับการ execute
func (n *NestedFilterBuilder) Priority() int {
	return 2
}

// Apply ประยุกต์ nested filter ไปยัง GORM query
func (n *NestedFilterBuilder) Apply(
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

	if len(q.FilterNestedBy) == 0 || len(q.FilterNested) == 0 {
		return db, nil
	}

	if len(q.FilterNestedBy) != len(q.FilterNested) {
		return nil, ErrFilterLengthMismatch
	}

	condition := q.FilterNestedCondition
	if condition == "" {
		condition = "and"
	}

	// Group by relation เพื่อสร้าง EXISTS ที่ถูกต้อง
	relationGroups := n.groupByRelation(q.FilterNestedBy, q.FilterNested)

	for relation, filters := range relationGroups {
		subquery := n.buildExistsSubquery(db, relation, filters, o)

		if condition == "or" {
			db = db.Or("EXISTS (?)", subquery)
		} else {
			db = db.Where("EXISTS (?)", subquery)
		}
	}

	return db, nil
}

// groupByRelation จัดกลุ่ม filters ตาม relation
func (n *NestedFilterBuilder) groupByRelation(
	filterNestedBy []string,
	filterNested [][]interface{},
) map[string][]RelationFilter {
	groups := make(map[string][]RelationFilter)

	for i, path := range filterNestedBy {
		parts := strings.SplitN(path, ".", 2)
		if len(parts) != 2 {
			continue
		}

		relation := parts[0]
		column := parts[1]

		groups[relation] = append(groups[relation], RelationFilter{
			Column: column,
			Values: filterNested[i],
		})
	}

	return groups
}

// buildExistsSubquery สร้าง EXISTS subquery
func (n *NestedFilterBuilder) buildExistsSubquery(
	db *gorm.DB,
	relation string,
	filters []RelationFilter,
	opts *FilterOptions,
) *gorm.DB {
	// ดึงข้อมูล relation metadata
	relationMeta := n.getRelationMetadata(relation)

	alias := opts.TableAlias
	if alias == "" {
		alias = "main"
	}

	subquery := db.Session(&gorm.Session{NewDB: true}).
		Table(relationMeta.TableName + " AS sub_" + relation).
		Select("1").
		Where("sub_" + relation + "." + relationMeta.ForeignKey + " = " + alias + ".id")

	// เพิ่ม soft delete filter
	if opts.SoftDelete {
		subquery = subquery.Where("sub_" + relation + ".deleted_at IS NULL")
	}

	// เพิ่ม filter conditions
	for _, f := range filters {
		subquery = subquery.Where("sub_"+relation+"."+f.Column+" IN ?", f.Values)
	}

	return subquery
}

// getRelationMetadata ดึง metadata ของ relation
// ใน production ควรใช้ relation registry เพื่อดึงข้อมูล
func (n *NestedFilterBuilder) getRelationMetadata(relation string) RelationMetadata {
	// Default implementation - ใน production ควรใช้ GORM schema
	return RelationMetadata{
		TableName:  relation + "s", // Simple pluralization
		ForeignKey: relation + "_id",
	}
}
