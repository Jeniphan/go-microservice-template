package filter

import (
	"gorm.io/gorm"
)

// FilterBuilder interface สำหรับสร้าง filter components
// แต่ละ builder จะ implement logic สำหรับ filter type ที่แตกต่างกัน
type FilterBuilder interface {
	// Apply ประยุกต์ filter ไปยัง GORM query
	// รับ db *gorm.DB, query *AdvanceFilterQuery, และ opts *FilterOptions
	// คืนค่า *gorm.DB ที่ถูก modify และ error (ถ้ามี)
	Apply(db *gorm.DB, query *AdvanceFilterQuery, opts *FilterOptions) (*gorm.DB, error)

	// Name คืนชื่อของ filter builder
	Name() string

	// Priority คืนลำดับการ execute (ต่ำกว่าจะทำก่อน)
	Priority() int
}

// Filterable interface ที่ entity ต้อง implement
// เพื่อให้ AdvanceFilter รู้จัก table และ fields ที่สามารถ filter ได้
type Filterable interface {
	// TableName คืนชื่อ table ของ entity
	TableName() string

	// FilterableFields คืนรายชื่อ fields ที่อนุญาตให้ filter
	// ใช้สำหรับ validate ว่า user ไม่ได้ filter ด้วย field ที่ไม่อนุญาต
	FilterableFields() []string
}
