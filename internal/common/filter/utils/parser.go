package utils

import (
	"strconv"
	"strings"
	"time"
)

// ParseRangeValue แปลงค่าสำหรับ range filter
// ตรวจสอบว่าเป็น timestamp field หรือตัวเลข
// ถ้าเป็น timestamp field จะ parse เป็น time.Time
// ถ้าเป็นตัวเลขจะ parse เป็น float64
// ถ้าเป็น string ธรรมดาจะ return ค่าเดิม
func ParseRangeValue(value, field string) (interface{}, error) {
	// ตรวจสอบว่าเป็น timestamp field หรือไม่
	if IsTimestampField(field) {
		return time.Parse(time.RFC3339, value)
	}

	// ตรวจสอบว่าเป็นตัวเลขหรือไม่
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, nil
	}

	// ถ้าไม่ใช่ทั้งสองอย่าง return ค่าเดิมเป็น string
	return value, nil
}

// IsTimestampField ตรวจสอบว่า field เป็น timestamp หรือไม่
// โดยดูจากชื่อ field ที่มี "at" หรือ "date" (case-insensitive)
// เช่น "created_at", "updated_at", "occurred_at", "start_date", "end_date"
func IsTimestampField(field string) bool {
	lowerField := strings.ToLower(field)
	return strings.Contains(lowerField, "at") || strings.Contains(lowerField, "date")
}

// ParseBool แปลง string เป็น bool อย่างปลอดภัย
// รองรับค่า "true", "1", "yes", "on" (case-insensitive) เป็น true
// รองรับค่า "false", "0", "no", "off" (case-insensitive) เป็น false
// ถ้าไม่รู้จักจะ return error
func ParseBool(value string) (bool, error) {
	lowerValue := strings.ToLower(value)

	switch lowerValue {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, ErrInvalidBoolValue
	}
}

// ErrInvalidBoolValue error สำหรับค่า bool ที่ไม่ถูกต้อง
var ErrInvalidBoolValue = &invalidBoolError{}

// invalidBoolError เป็น custom error สำหรับ bool parsing
type invalidBoolError struct{}

func (e *invalidBoolError) Error() string {
	return "invalid boolean value"
}
