package builders

import (
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// RangeFilterBuilder สร้าง range filter ด้วย >= และ <=
type RangeFilterBuilder struct{}

// NewRangeFilterBuilder สร้าง RangeFilterBuilder instance ใหม่
func NewRangeFilterBuilder() *RangeFilterBuilder {
	return &RangeFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (r *RangeFilterBuilder) Name() string {
	return "RangeFilter"
}

// Priority คืนลำดับการ execute
func (r *RangeFilterBuilder) Priority() int {
	return 6
}

// Apply ประยุกต์ range filter ไปยัง GORM query
func (r *RangeFilterBuilder) Apply(
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

	condition := q.StartAndEndCondition
	if condition == "" {
		condition = "and"
	}

	// Start filter (>=)
	if q.StartBy != "" && q.Start != "" {
		column := r.buildRangeColumn(q.StartBy, o.TableAlias)
		value, err := r.parseRangeValue(q.Start, q.StartBy)
		if err != nil {
			return nil, err
		}
		db = db.Where(column+" >= ?", value)
	}

	// End filter (<=)
	if q.EndBy != "" && q.End != "" {
		column := r.buildRangeColumn(q.EndBy, o.TableAlias)
		value, err := r.parseRangeValue(q.End, q.EndBy)
		if err != nil {
			return nil, err
		}
		db = db.Where(column+" <= ?", value)
	}

	_ = condition // Used for future extension
	return db, nil
}

// buildRangeColumn สร้าง column expression สำหรับ range
func (r *RangeFilterBuilder) buildRangeColumn(field, alias string) string {
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

// parseRangeValue แปลงค่าสำหรับ range filter
func (r *RangeFilterBuilder) parseRangeValue(value, field string) (interface{}, error) {
	// ตรวจสอบว่าเป็น timestamp field หรือไม่
	if IsTimestampField(field) {
		return time.Parse(time.RFC3339, value)
	}

	// ตรวจสอบว่าเป็นตัวเลขหรือไม่
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, nil
	}

	return value, nil
}

// IsTimestampField ตรวจสอบว่า field เป็น timestamp หรือไม่
func IsTimestampField(field string) bool {
	lowerField := strings.ToLower(field)
	return strings.Contains(lowerField, "at") || strings.Contains(lowerField, "date")
}
