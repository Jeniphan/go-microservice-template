package builders

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// FilterBuilder interface สำหรับสร้าง filter components
type FilterBuilder interface {
	Apply(db *gorm.DB, query interface{}, opts interface{}) (*gorm.DB, error)
	Name() string
	Priority() int
}

// AdvanceFilterQuery - Query parameters struct
type AdvanceFilterQuery struct {
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

// FilterOptions - Filter options struct
type FilterOptions struct {
	TableAlias   string
	Preload      []string
	AppID        interface{}
	ParentTable  string
	SoftDelete   bool
	CustomScopes []func(*gorm.DB) *gorm.DB
}

// Error constants
var (
	ErrFilterLengthMismatch = errors.New("filter_by and filter must have same length")
)

// BasicFilterBuilder สร้าง basic filter สำหรับ IN/NOT IN filtering
type BasicFilterBuilder struct{}

// NewBasicFilterBuilder สร้าง BasicFilterBuilder instance ใหม่
func NewBasicFilterBuilder() *BasicFilterBuilder {
	return &BasicFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (b *BasicFilterBuilder) Name() string {
	return "BasicFilter"
}

// Priority คืนลำดับการ execute (ต่ำกว่าจะทำก่อน)
func (b *BasicFilterBuilder) Priority() int {
	return 1
}

// Apply ประยุกต์ basic filter ไปยัง GORM query
// รองรับ IN/NOT IN filtering พร้อม include/exclude values
func (b *BasicFilterBuilder) Apply(
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

	if len(q.FilterBy) == 0 || len(q.Filter) == 0 {
		return db, nil
	}

	if len(q.FilterBy) != len(q.Filter) {
		return nil, ErrFilterLengthMismatch
	}

	condition := q.FilterCondition
	if condition == "" {
		condition = "and"
	}

	for i, field := range q.FilterBy {
		values := q.Filter[i]
		if len(values) == 0 {
			continue
		}

		// แยกค่า include และ exclude (ค่าที่ขึ้นต้นด้วย !)
		includeValues, excludeValues := b.parseValues(values)

		column := b.buildColumn(field, o.TableAlias)

		if len(includeValues) > 0 {
			if condition == "or" && i > 0 {
				db = db.Or(column+" IN ?", includeValues)
			} else {
				db = db.Where(column+" IN ?", includeValues)
			}
		}

		if len(excludeValues) > 0 {
			db = db.Where(column+" NOT IN ?", excludeValues)
		}
	}

	return db, nil
}

// parseValues แยกค่า include และ exclude
// ค่าที่ขึ้นต้นด้วย "!" จะถือเป็น exclude value
func (b *BasicFilterBuilder) parseValues(values []interface{}) ([]interface{}, []interface{}) {
	var include, exclude []interface{}

	for _, v := range values {
		if str, ok := v.(string); ok && strings.HasPrefix(str, "!") {
			exclude = append(exclude, strings.TrimPrefix(str, "!"))
		} else {
			include = append(include, v)
		}
	}

	return include, exclude
}

// buildColumn สร้าง column expression พร้อม table alias
func (b *BasicFilterBuilder) buildColumn(field, alias string) string {
	if alias != "" {
		return alias + "." + field
	}
	return field
}
