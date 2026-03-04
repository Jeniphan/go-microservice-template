package builders

import "gorm.io/gorm"

// PaginationFilterBuilder สร้าง pagination filter ด้วย LIMIT/OFFSET
type PaginationFilterBuilder struct{}

// NewPaginationFilterBuilder สร้าง PaginationFilterBuilder instance ใหม่
func NewPaginationFilterBuilder() *PaginationFilterBuilder {
	return &PaginationFilterBuilder{}
}

// Name คืนชื่อของ filter builder
func (p *PaginationFilterBuilder) Name() string {
	return "PaginationFilter"
}

// Priority คืนลำดับการ execute (สุดท้ายสุด)
func (p *PaginationFilterBuilder) Priority() int {
	return 9
}

// Apply ประยุกต์ pagination filter ไปยัง GORM query
func (p *PaginationFilterBuilder) Apply(
	db *gorm.DB,
	query interface{},
	opts interface{},
) (*gorm.DB, error) {
	q, ok := query.(*AdvanceFilterQuery)
	if !ok {
		return db, nil
	}

	// ถ้ามี Limit ให้ใช้ Limit แทน pagination
	if q.Limit > 0 {
		return db.Limit(q.Limit), nil
	}

	page := q.Page
	if page <= 0 {
		page = 1
	}

	perPage := q.PerPage
	if perPage <= 0 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100 // Max limit
	}

	offset := (page - 1) * perPage

	db = db.Offset(offset).Limit(perPage)

	return db, nil
}

// CalculateTotalPage คำนวณจำนวนหน้าทั้งหมด
func CalculateTotalPage(total int64, perPage int) int64 {
	if perPage <= 0 {
		return 0
	}
	return (total + int64(perPage) - 1) / int64(perPage)
}
