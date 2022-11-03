package repos

import (
	"math"

	"m/v2/app/models"

	"gorm.io/gorm"
)

// returns a pagination scope.
// where is optional, use nil for none
func paginate(value interface{}, where *gorm.DB, pagination *models.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	if where == nil {
		db.Model(value).Count(&totalRows)
	} else {
		db.Model(value).Where(where).Count(&totalRows)
	}

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	if where == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort()).Where(where)
	}
}
