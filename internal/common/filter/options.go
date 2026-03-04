package filter

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
)

// FilterOptions ตัวเลือกเพิ่มเติมสำหรับ filter
// ใช้กำหนดค่า config ต่างๆ ของ AdvanceFilter
type FilterOptions struct {
	// TableAlias alias สำหรับตารางหลัก
	// เช่น "a" สำหรับ activities, "o" สำหรับ orders
	// ใช้ในการ prefix column เพื่อหลีกเลี่ยง conflict
	TableAlias string

	// Preload relations ที่ต้องการ preload อัตโนมัติ
	// ใช้หลีกเลี่ยง N+1 query problem
	Preload []string

	// AppID app ID constraint สำหรับ multi-tenant
	// ถ้ามีค่า จะถูกเพิ่มเป็น WHERE condition ทุก query
	AppID *uuid.UUID

	// ParentTable ชื่อตารางแม่ (สำหรับ parent filter)
	// ใช้เมื่อต้องการ filter จากตารางที่เชื่อมโยงกัน
	ParentTable string

	// SoftDelete เปิดใช้งาน soft delete filter
	// ถ้าเป็น true จะกรอง record ที่มี deleted_at ไม่เป็น NULL ออก
	// default: true
	SoftDelete bool

	// CustomScopes custom GORM scopes ที่ต้องการ apply
	// ใช้สำหรับ logic พิเศษที่ต้องการเพิ่มเข้าไปทุก query
	CustomScopes []func(*gorm.DB) *gorm.DB
}
