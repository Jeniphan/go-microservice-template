package filter

// AdvanceFilterQuery โครงสร้างหลักสำหรับ query parameters
// รองรับการ filter แบบต่างๆ ผ่าน query parameters
type AdvanceFilterQuery struct {
	// ============ Basic Filter ============

	// FilterBy รายชื่อ fields ที่ต้องการ filter (เช่น ["status", "type"])
	FilterBy []string `form:"filter_by" json:"filter_by"`

	// Filter ค่าที่ต้องการ filter สำหรับแต่ละ field
	// เป็น 2D array เพื่อรองรับหลายค่าต่อ field (เช่น [["active", "pending"], ["email", "sms"]])
	Filter [][]interface{} `form:"filter" json:"filter"`

	// FilterCondition เงื่อนไขระหว่าง filters ( "and" | "or" )
	// default: "and"
	FilterCondition string `form:"filter_condition" json:"filter_condition"`

	// ============ Nested Relation Filter ============

	// FilterNestedBy รายชื่อ nested relation fields ที่ต้องการ filter
	// format: "relation.column" (เช่น "items.status")
	FilterNestedBy []string `form:"filter_nested_by" json:"filter_nested_by"`

	// FilterNested ค่าที่ต้องการ filter สำหรับ nested relation
	FilterNested [][]interface{} `form:"filter_nested" json:"filter_nested"`

	// FilterNestedCondition เงื่อนไขระหว่าง nested filters ( "and" | "or" )
	// default: "and"
	FilterNestedCondition string `form:"filter_nested_condition" json:"filter_nested_condition"`

	// ============ Parent Relation Filter ============

	// FilterNestedParentBy รายชื่อ parent relation fields ที่ต้องการ filter
	// format: "parent_alias.column" (เช่น "activity.app_id")
	FilterNestedParentBy []string `form:"filter_nested_parent_by" json:"filter_nested_parent_by"`

	// FilterNestedParent ค่าที่ต้องการ filter สำหรับ parent relation
	FilterNestedParent [][]interface{} `form:"filter_nested_parent" json:"filter_nested_parent"`

	// FilterNestedParentCondition เงื่อนไขระหว่าง parent filters ( "and" | "or" )
	// default: "and"
	FilterNestedParentCondition string `form:"filter_nested_parent_condition" json:"filter_nested_parent_condition"`

	// ============ Many-to-Many Filter ============

	// FilterM2MBy รายชื่อ M2M relation fields ที่ต้องการ filter
	// format: "relation.column" (เช่น "tags.id")
	FilterM2MBy []string `form:"filter_m2m_by" json:"filter_m2m_by"`

	// FilterM2M ค่าที่ต้องการ filter สำหรับ M2M relation
	FilterM2M [][]interface{} `form:"filter_m2m" json:"filter_m2m"`

	// FilterM2MCondition เงื่อนไขระหว่าง M2M filters ( "and" | "or" )
	// default: "or"
	FilterM2MCondition string `form:"filter_m2m_condition" json:"filter_m2m_condition"`

	// FilterM2MJoinAlias alias สำหรับ join table ใน M2M filter
	FilterM2MJoinAlias string `form:"filter_m2m_join_alias" json:"filter_m2m_join_alias"`

	// ============ Search ============

	// SearchBy รายชื่อ fields ที่ต้องการค้นหา
	SearchBy []string `form:"search_by" json:"search_by"`

	// Search คำค้นหา (support partial match ด้วย LIKE/ILIKE)
	Search string `form:"search" json:"search"`

	// ============ Range Filter ============

	// StartBy ชื่อ field สำหรับ start value (ใช้กับ >=)
	StartBy string `form:"start_by" json:"start_by"`

	// Start ค่าเริ่มต้น (สำหรับ range filter)
	Start string `form:"start" json:"start"`

	// EndBy ชื่อ field สำหรับ end value (ใช้กับ <=)
	EndBy string `form:"end_by" json:"end_by"`

	// End ค่าสิ้นสุด (สำหรับ range filter)
	End string `form:"end" json:"end"`

	// StartAndEndCondition เงื่อนไขระหว่าง start และ end ( "and" | "or" )
	// default: "and"
	StartAndEndCondition string `form:"start_and_end_condition" json:"start_and_end_condition"`

	// ============ Sorting ============

	// SortBy รายชื่อ fields ที่ต้องการเรียงลำดับ
	SortBy []string `form:"sort_by" json:"sort_by"`

	// Sort ทิศทางการเรียงลำดับ ( "asc" | "desc" )
	// ต้องมี length เท่ากับ SortBy
	Sort []string `form:"sort" json:"sort"`

	// ============ Pagination ============

	// Page หมายเลขหน้าปัจจุบัน (เริ่มที่ 1)
	// default: 1
	Page int `form:"page" json:"page"`

	// PerPage จำนวนรายการต่อหน้า
	// default: 10, max: 100
	PerPage int `form:"per_page" json:"per_page"`

	// ============ Grouping ============

	// GroupBy รายชื่อ fields สำหรับจัดกลุ่ม (DISTINCT ON ใน PostgreSQL)
	GroupBy []string `form:"group_by" json:"group_by"`

	// GroupSortBy ชื่อ field สำหรับเรียงลำดับก่อนจับกลุ่ม
	// ใช้กับ DISTINCT ON เพื่อเลือก record ที่ต้องการ
	GroupSortBy string `form:"group_sort_by" json:"group_sort_by"`

	// GroupSort ทิศทางการเรียงลำดับสำหรับ group ( "max" | "min" )
	// "max" = เลือก record ที่มีค่ามากที่สุด
	// "min" = เลือก record ที่มีค่าน้อยที่สุด
	GroupSort string `form:"group_sort" json:"group_sort"`

	// ============ Preload ============

	// Preload รายชื่อ relations ที่ต้องการ preload
	// ใช้หลีกเลี่ยง N+1 query problem
	Preload []string `form:"preload" json:"preload"`

	// ============ Limit ============

	// Limit จำนวน record สูงสุดที่ต้องการ (ไม่รวม pagination)
	// ใช้สำหรับ query ที่ไม่ต้องการ pagination
	Limit int `form:"limit" json:"limit"`
}
