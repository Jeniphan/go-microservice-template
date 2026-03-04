package filter

// FilterResult ผลลัพธ์จาก AdvanceFilter
// ใช้ generics เพื่อรองรับ type ของ entity ที่หลากหลาย
type FilterResult[T any] struct {
	// Data ข้อมูลที่ได้จากการ query
	Data []T `json:"data"`

	// Total จำนวน record ทั้งหมด (ก่อน pagination)
	Total int64 `json:"total"`

	// TotalPage จำนวนหน้าทั้งหมด
	TotalPage int64 `json:"total_page"`

	// Page หมายเลขหน้าปัจจุบัน
	Page int `json:"page"`

	// PerPage จำนวนรายการต่อหน้า
	PerPage int `json:"per_page"`
}
