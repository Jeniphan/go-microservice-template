package builders

import (
	"strings"

	"gorm.io/gorm"
)

// M2MMetadata metadata สำหรับ many-to-many relation
type M2MMetadata struct {
	JoinTable   string
	SourceFK    string
	TargetTable string
	TargetFK    string
}

// M2MFilterBuilder สร้าง many-to-many filter ด้วย subquery
type M2MFilterBuilder struct{}

// NewM2MFilterBuilder สร้าง M2MFilterBuilder instance ใหม่
func NewM2MFilterBuilder() *M2MFilterBuilder {
	return &M2MFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (m *M2MFilterBuilder) Name() string {
	return "M2MFilter"
}

// Priority คืนลำดับการ execute
func (m *M2MFilterBuilder) Priority() int {
	return 4
}

// Apply ประยุกต์ M2M filter ไปยัง GORM query
func (m *M2MFilterBuilder) Apply(
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

	if len(q.FilterM2MBy) == 0 || len(q.FilterM2M) == 0 {
		return db, nil
	}

	if len(q.FilterM2MBy) != len(q.FilterM2M) {
		return nil, ErrFilterLengthMismatch
	}

	condition := q.FilterM2MCondition
	if condition == "" {
		condition = "or"
	}

	for i, path := range q.FilterM2MBy {
		parts := strings.SplitN(path, ".", 2)
		if len(parts) != 2 {
			continue
		}

		relation := parts[0]
		column := parts[1]
		values := q.FilterM2M[i]

		// ดึง M2M metadata
		m2mMeta := m.getM2MMetadata(relation)

		alias := o.TableAlias
		if alias == "" {
			alias = "main"
		}

		// สร้าง subquery สำหรับ M2M
		subquery := db.Session(&gorm.Session{NewDB: true}).
			Table(m2mMeta.JoinTable+" AS jt").
			Select("jt."+m2mMeta.SourceFK).
			Joins("JOIN "+m2mMeta.TargetTable+" AS target ON target.id = jt."+m2mMeta.TargetFK).
			Where("target."+column+" IN ?", values)

		if condition == "or" && i > 0 {
			db = db.Or(alias+".id IN (?)", subquery)
		} else {
			db = db.Where(alias+".id IN (?)", subquery)
		}
	}

	return db, nil
}

// getM2MMetadata ดึง metadata ของ M2M relation
// ใน production ควรใช้ relation registry เพื่อดึงข้อมูล
func (m *M2MFilterBuilder) getM2MMetadata(relation string) M2MMetadata {
	// Default implementation - ใน production ควรใช้ GORM schema
	return M2MMetadata{
		JoinTable:   relation + "_tags", // Example: activities_tags
		SourceFK:    "activity_id",
		TargetTable: relation, // Example: tags
		TargetFK:    "tag_id",
	}
}
