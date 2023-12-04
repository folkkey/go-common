package gocommon

import (
	"context"
	"errors"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IBaseService[Entity, ID, Dto, CreateDto, UpdateDto, FilterDto any] interface {
	Get(ctx context.Context, id ID) (*Dto, error)
	GetList(ctx context.Context, input FilterDto, query *PagingQuery) (PagedResultDto[Dto], error)
	Create(ctx context.Context, input CreateDto) (Dto, error)
	Update(ctx context.Context, id ID, input UpdateDto) (Dto, error)
	Delete(ctx context.Context, id ID) (bool, error)
}

type BaseService[Entity, ID, Dto, CreateDto, UpdateDto, FilterDto any] struct {
	BaseRepository IBaseRepository[Entity, ID]
}

func (s BaseService[Entity, ID, Dto, CreateDto, UpdateDto, FilterDto]) Get(ctx context.Context, id ID) (*Dto, error) {
	entity, err := s.BaseRepository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	dto, err := TypeConverter[Dto](entity)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (s BaseService[Entity, ID, Dto, CreateDto, UpdateDto, FilterDto]) GetList(ctx context.Context, input FilterDto, query *PagingQuery) (PagedResultDto[Dto], error) {
	rs := PagedResultDto[Dto]{
		Total: int64(0),
	}

	total, entities, err := s.BaseRepository.GetList(ctx, input, query)
	if err != nil {
		return rs, err
	}

	rs.Total = total
	if len(entities) > 0 {
		dto, err := TypeConverter[[]Dto](entities)
		if err != nil {
			return rs, err
		}
		rs.Items = *dto
	}
	return rs, nil
}

func (s BaseService[Entity, PkType, Dto, CreateDto, UpdateDto, FilterDto]) Create(ctx context.Context, input CreateDto) (Dto, error) {
	var entity Entity
	var dto Dto
	err := mapstructure.Decode(input, &entity)
	if err != nil {
		return dto, err
	}
	err = s.BaseRepository.Create(ctx, &entity)
	if err != nil {
		return dto, err
	}
	err = mapstructure.Decode(entity, &dto)
	if err != nil {
		return dto, err
	}
	return dto, nil
}

func (s BaseService[Entity, PkType, Dto, CreateDto, UpdateDto, FilterDto]) Update(ctx context.Context, id PkType, input UpdateDto) (Dto, error) {
	var dto Dto
	entity, err := s.BaseRepository.Get(ctx, id)
	if err != nil {
		return dto, err
	}

	err = mapstructure.Decode(input, &entity)
	if err != nil {
		return dto, err
	}
	err = s.BaseRepository.Update(ctx, entity)
	if err != nil {
		return dto, err
	}
	err = mapstructure.Decode(entity, &dto)
	if err != nil {
		return dto, err
	}
	return dto, nil
}

func (s BaseService[Entity, PkType, Dto, CreateDto, UpdateDto, FilterDto]) Delete(ctx context.Context, id PkType) (bool, error) {
	entity, err := s.BaseRepository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	_, err = s.BaseRepository.Delete(ctx, entity)
	if err != nil {
		return false, err
	}
	return true, nil
}
