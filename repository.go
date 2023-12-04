package gocommon

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

type IBaseRepository[T, ID any] interface {
	Get(ctx context.Context, id ID) (*T, error)
	GetList(ctx context.Context, filters interface{}, paging *PagingQuery) (int64, []T, error)
	GetListWithQuery(ctx context.Context, query *gorm.DB, paging *PagingQuery) (int64, []T, error)
	Create(ctx context.Context, t *T) error
	CreateMany(ctx context.Context, t *[]T, size int) error
	Update(ctx context.Context, t *T) error
	Delete(ctx context.Context, t *T) (bool, error)
	QueryBuilder(query *gorm.DB, options interface{}) *gorm.DB
}

type BaseRepository[T, ID any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T, ID any](db *gorm.DB) *BaseRepository[T, ID] {
	return &BaseRepository[T, ID]{DB: db}
}

func (r BaseRepository[T, ID]) Get(ctx context.Context, id ID) (*T, error) {
	var model T
	err := r.DB.WithContext(ctx).Preload(clause.Associations).First(&model, id)
	if err != nil {
		return &model, err.Error
	}
	return &model, nil
}

func (r BaseRepository[T, ID]) GetList(ctx context.Context, filters interface{}, paging *PagingQuery) (int64, []T, error) {
	query := r.DB.Model(new(T)).Where(filters)
	r.QueryBuilder(query, nil)

	//Count data size
	total := int64(0)
	rs := query.Count(&total)
	if rs.Error != nil {
		//logger.Errorf("[BaseRepository] GetList %v error: %v", *new(T), rs.Error)
		return 0, nil, rs.Error
	}

	//Get items
	if paging != nil {
		query = query.Offset(paging.Page * paging.Size).Limit(paging.Size)
		if paging.OrderBy != nil {
			if len(*paging.OrderBy) > 0 {
				var sortBy = "asc"
				if len(*paging.SortBy) > 0 {
					sortBy = *paging.SortBy
				}
				//query = query.Order(sortBy + " " + *paging.OrderBy)
				query = query.Order(*paging.OrderBy + " " + sortBy)
			}
		}
	}
	var items []T
	rs = query.WithContext(ctx).Preload(clause.Associations).Find(&items)
	if rs.Error != nil {
		//logger.Errorf("[BaseRepository] GetList %v error: %v", *new(T), rs.Error)
		return 0, nil, rs.Error
	}

	return total, items, nil
}

func (r BaseRepository[T, ID]) GetListWithQuery(ctx context.Context, query *gorm.DB, paging *PagingQuery) (int64, []T, error) {
	r.QueryBuilder(query, nil)

	//Count data size
	total := int64(0)
	rs := query.Count(&total)
	if rs.Error != nil {
		//logger.Errorf("[BaseRepository] GetList %v error: %v", *new(T), rs.Error)
		return 0, nil, rs.Error
	}

	//Get items
	if paging != nil {
		query = query.Offset(paging.Page * paging.Size).Limit(paging.Size)
		if paging.OrderBy != nil {
			if len(*paging.OrderBy) > 0 {
				var sortBy = "asc"
				if len(*paging.SortBy) > 0 {
					sortBy = *paging.SortBy
				}
				//query = query.Order(sortBy + " " + *paging.OrderBy)
				query = query.Order(*paging.OrderBy + " " + sortBy)
			}
		}
	}
	var items []T
	rs = query.WithContext(ctx).Preload(clause.Associations).Find(&items)
	if rs.Error != nil {
		return 0, nil, rs.Error
	}

	return total, items, nil
}

//func (r BaseRepository[T, ID]) Find(ctx context.Context, specification common.Specification) ([]T, error) {
//	var models []T
//	err := r.DB.WithContext(ctx).Where(specification.GetQuery(), specification.GetValues()).Find(&models).Error
//	if err != nil {
//		return models, err
//	}
//
//	return models, nil
//}

func (r BaseRepository[T, ID]) Create(ctx context.Context, t *T) error {
	err := r.DB.WithContext(ctx).Create(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func (r BaseRepository[T, ID]) CreateMany(ctx context.Context, t *[]T, size int) error {
	result := r.DB.WithContext(ctx).CreateInBatches(t, size)
	if result.Error != nil {
		//log.Printf("Cannot save items %+v: %s", t, result.Error.Error())
		return result.Error
	}
	return nil
}

func (r BaseRepository[T, ID]) Update(ctx context.Context, t *T) error {
	rs := r.DB.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(t)
	if rs.Error != nil {
		//logger.Errorf("[BaseRepository] Update %v error: %v", *new(T), rs.Error)
		return rs.Error
	}
	return nil
}

func (r BaseRepository[T, ID]) Delete(ctx context.Context, t *T) (bool, error) {
	rs := r.DB.WithContext(ctx).Delete(t)
	if rs.Error != nil {
		//logger.Errorf("[BaseRepository] Delete %v error: %v", *new(T), rs.Error)
		return false, rs.Error
	}
	return true, nil
}

// ----------------------------------------
// ----------------------------------------
// QueryBuilder build additional query term
// ----------------------------------------
// ----------------------------------------
func (r BaseRepository[T, ID]) QueryBuilder(query *gorm.DB, options interface{}) *gorm.DB {
	if options != nil {
		optionsVal := reflect.ValueOf(options)

		if optionsVal.Kind() == reflect.Map {
			for _, key := range optionsVal.MapKeys() {
				switch key.Interface() {
				case "preload":
					preloadVal := optionsVal.MapIndex(key).Elem()

					if preloadVal.Kind() == reflect.Slice {
						for i := 0; i < preloadVal.Len(); i++ {
							query = query.Preload(preloadVal.Index(i).String())
						}
					}
				case "join":
					joinVal := optionsVal.MapIndex(key).Elem()

					if joinVal.Kind() == reflect.Slice {
						for i := 0; i < joinVal.Len(); i++ {
							query = query.Joins(joinVal.Index(i).String())
						}
					} else {
						query = query.Joins(joinVal.String())
					}
				}
			}
		}
	}

	return query
}
