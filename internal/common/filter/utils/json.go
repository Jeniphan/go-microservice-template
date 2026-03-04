package utils

import "strings"

// ParseJSONPath แยก column และ key จาก JSON path
// เช่น "metadata.description" -> ("metadata", "description", true)
// ถ้า path ไม่มี "." จะ return (path, "", false)
func ParseJSONPath(path string) (column, key string, isJSON bool) {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}
	return path, "", false
}

// BuildJSONExpression สร้าง SQL expression สำหรับ JSON field
// เช่น BuildJSONExpression("a", "metadata", "description") -> "a.metadata->> 'description'"
func BuildJSONExpression(alias, column, key string) string {
	if alias != "" {
		return alias + "." + column + "->> '" + key + "'"
	}
	return column + "->> '" + key + "'"
}

// IsJSONField ตรวจสอบว่าเป็น JSON field หรือไม่
// JSON field คือ field ที่มี "." แต่ไม่มี "->" (เพราะ "->" คือ SQL operator ที่ใช้แล้ว)
// เช่น "metadata.description" -> true, "metadata->description" -> false, "status" -> false
func IsJSONField(field string) bool {
	// ถ้ามี "->" แสดงว่าเป็น SQL expression ที่ใช้แล้ว ไม่ใช่ JSON path
	if strings.Contains(field, "->") {
		return false
	}
	// ถ้ามี "." แสดงว่าเป็น JSON path (เช่น metadata.description)
	return strings.Contains(field, ".")
}
