package filter

import (
	"context"
	"sort"

	"gorm.io/gorm"

	"order-v2-microservice/internal/common/filter/builders"
)

// AdvanceFilter main struct สำหรับ advance filtering
// ใช้ generics เพื่อรองรับ type ของ entity ที่หลากหลาย
type AdvanceFilter[T any] struct {
	db       *gorm.DB
	builders []builders.FilterBuilder
	options  *builders.FilterOptions
}

// NewAdvanceFilter สร้าง AdvanceFilter instance ใหม่
// รับ db *gorm.DB และ opts *FilterOptions
// ถ้า opts เป็น nil จะใช้ค่า default
func NewAdvanceFilter[T any](db *gorm.DB, opts *builders.FilterOptions) *AdvanceFilter[T] {
	if opts == nil {
		opts = &builders.FilterOptions{
			SoftDelete: true,
		}
	}

	return &AdvanceFilter[T]{
		db:      db,
		options: opts,
		builders: []builders.FilterBuilder{
			builders.NewBasicFilterBuilder(),
			builders.NewNestedFilterBuilder(),
			builders.NewParentFilterBuilder(),
			builders.NewM2MFilterBuilder(),
			builders.NewSearchFilterBuilder(),
			builders.NewRangeFilterBuilder(),
			builders.NewGroupFilterBuilder(),
			builders.NewSortFilterBuilder(),
			builders.NewPaginationFilterBuilder(),
		},
	}
}

// Apply ประยุกต์ filter ทั้งหมด
// รับ ctx context.Context และ query *AdvanceFilterQuery
// คืนค่า *FilterResult[T] และ error (ถ้ามี)
func (af *AdvanceFilter[T]) Apply(ctx context.Context, query *AdvanceFilterQuery) (*FilterResult[T], error) {
	// Convert query to builders.QueryParams for validation
	queryForValidation := &builders.AdvanceFilterQuery{
		FilterBy:                    query.FilterBy,
		Filter:                      query.Filter,
		FilterCondition:             query.FilterCondition,
		FilterNestedBy:              query.FilterNestedBy,
		FilterNested:                query.FilterNested,
		FilterNestedCondition:       query.FilterNestedCondition,
		FilterNestedParentBy:        query.FilterNestedParentBy,
		FilterNestedParent:          query.FilterNestedParent,
		FilterNestedParentCondition: query.FilterNestedParentCondition,
		FilterM2MBy:                 query.FilterM2MBy,
		FilterM2M:                   query.FilterM2M,
		FilterM2MCondition:          query.FilterM2MCondition,
		FilterM2MJoinAlias:          query.FilterM2MJoinAlias,
		SearchBy:                    query.SearchBy,
		Search:                      query.Search,
		StartBy:                     query.StartBy,
		Start:                       query.Start,
		EndBy:                       query.EndBy,
		End:                         query.End,
		StartAndEndCondition:        query.StartAndEndCondition,
		SortBy:                      query.SortBy,
		Sort:                        query.Sort,
		Page:                        query.Page,
		PerPage:                     query.PerPage,
		GroupBy:                     query.GroupBy,
		GroupSortBy:                 query.GroupSortBy,
		GroupSort:                   query.GroupSort,
		Preload:                     query.Preload,
		Limit:                       query.Limit,
	}

	// Convert options
	optsForBuilders := &builders.FilterOptions{
		TableAlias:   af.options.TableAlias,
		Preload:      af.options.Preload,
		AppID:        af.options.AppID,
		ParentTable:  af.options.ParentTable,
		SoftDelete:   af.options.SoftDelete,
		CustomScopes: af.options.CustomScopes,
	}

	// Start query
	db := af.db.WithContext(ctx).Model(new(T))

	// Apply table alias
	if af.options.TableAlias != "" {
		db = db.Table(af.getTableName() + " " + af.options.TableAlias)
	}

	// Apply soft delete filter
	if af.options.SoftDelete {
		alias := af.options.TableAlias
		if alias == "" {
			alias = af.getTableName()
		}
		db = db.Where(alias + ".deleted_at IS NULL")
	}

	// Apply app_id constraint
	if af.options.AppID != nil {
		alias := af.options.TableAlias
		if alias == "" {
			alias = af.getTableName()
		}
		db = db.Where(alias+".app_id = ?", af.options.AppID)
	}

	// Sort builders by priority
	sortedBuilders := make([]builders.FilterBuilder, len(af.builders))
	copy(sortedBuilders, af.builders)
	sort.Slice(sortedBuilders, func(i, j int) bool {
		return sortedBuilders[i].Priority() < sortedBuilders[j].Priority()
	})

	// Apply builders in priority order
	for _, builder := range sortedBuilders {
		var err error
		db, err = builder.Apply(db, queryForValidation, optsForBuilders)
		if err != nil {
			return nil, err
		}
	}

	// Apply preload relations from query
	for _, relation := range query.Preload {
		db = db.Preload(relation)
	}

	// Apply preload from options
	for _, relation := range af.options.Preload {
		db = db.Preload(relation)
	}

	// Get total count before pagination
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply distinct
	db = db.Distinct()

	// Execute query
	var results []T
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	// Build result
	perPage := query.PerPage
	if perPage <= 0 {
		perPage = 10
	}

	page := query.Page
	if page <= 0 {
		page = 1
	}

	return &FilterResult[T]{
		Data:      results,
		Total:     total,
		TotalPage: CalculateTotalPage(total, perPage),
		Page:      page,
		PerPage:   perPage,
	}, nil
}

// getTableName ดึงชื่อ table จาก generic type T
func (af *AdvanceFilter[T]) getTableName() string {
	// ใช้ reflect เพื่อดึงชื่อ table จาก type
	// สำหรับ GORM models ที่ implement TableName() string
	var t T
	if filterable, ok := any(t).(Filterable); ok {
		return filterable.TableName()
	}
	// Fallback: ใช้ชื่อ type โดยลบ s ตัวสุดท้ายถ้าเป็น plural
	return ""
}

// CalculateTotalPage คำนวณจำนวนหน้าทั้งหมด
// รับ total จำนวน record ทั้งหมด และ perPage จำนวน record ต่อหน้า
// คืนค่าจำนวนหน้าทั้งหมด
func CalculateTotalPage(total int64, perPage int) int64 {
	if perPage <= 0 {
		return 0
	}
	return (total + int64(perPage) - 1) / int64(perPage)
}
