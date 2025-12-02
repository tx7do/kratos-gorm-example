package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/crypto"
	"github.com/tx7do/go-utils/mapper"

	pagination "github.com/tx7do/go-curd/api/gen/go/pagination/v1"

	"gorm.io/gorm/clause"

	"kratos-gorm-example/app/user/service/internal/data/models"

	userV1 "kratos-gorm-example/api/gen/go/user/service/v1"
)

type UserRepo struct {
	data *Data
	log  *log.Helper

	mapper *mapper.CopierMapper[userV1.User, models.User]
}

func NewUserRepo(data *Data, logger log.Logger) *UserRepo {
	l := log.NewHelper(log.With(logger, "module", "user/repo/user-service"))
	repo := &UserRepo{
		data: data,
		log:  l,
	}

	repo.init()

	return repo
}

func (r *UserRepo) init() {
	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *UserRepo) List(_ context.Context, req *pagination.PagingRequest) (*userV1.ListUserResponse, error) {
	var entities []*models.User
	result := r.data.db.
		Limit(int(req.GetPageSize())).
		Offset(int(req.GetPageSize() * (req.GetPage() - 1))).
		Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}

	dtos := make([]*userV1.User, 0, len(entities))
	for _, entity := range entities {
		dto := r.mapper.ToDTO(entity)
		dtos = append(dtos, dto)
	}

	var count int64
	result = r.data.db.Model(&models.User{}).
		Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}

	return &userV1.ListUserResponse{
		Total: int32(count),
		Items: dtos,
	}, nil
}

func (r *UserRepo) Get(_ context.Context, req *userV1.GetUserRequest) (*userV1.User, error) {
	entity := &models.User{}
	err := r.data.db.First(entity, "id = ?", req.GetId()).Error
	return r.mapper.ToDTO(entity), err
}

func (r *UserRepo) Create(_ context.Context, req *userV1.CreateUserRequest) (*userV1.User, error) {
	cryptoPassword, err := crypto.HashPassword(req.User.GetPassword())
	if err != nil {
		return nil, err
	}

	entity := &models.User{
		UserName: req.User.GetUserName(),
		NickName: req.User.GetNickName(),
		Password: cryptoPassword,
	}

	result := r.data.db.Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *UserRepo) Update(_ context.Context, req *userV1.UpdateUserRequest) (*userV1.User, error) {
	var cryptoPassword string
	var err error
	if req.User.Password != nil {
		cryptoPassword, err = crypto.HashPassword(req.User.GetPassword())
		if err != nil {
			return nil, err
		}
	}

	entity := &models.User{
		UserName: req.User.GetUserName(),
		NickName: req.User.GetNickName(),
		Password: cryptoPassword,
	}

	result := r.data.db.Model(entity).Updates(entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *UserRepo) Upsert(_ context.Context, req *userV1.UpdateUserRequest) (*userV1.User, error) {
	var err error

	var cryptoPassword string
	if req.User.Password != nil {
		cryptoPassword, err = crypto.HashPassword(req.User.GetPassword())
		if err != nil {
			return nil, err
		}
	}

	entity := &models.User{
		UserName: req.User.GetUserName(),
		NickName: req.User.GetNickName(),
		Password: cryptoPassword,
	}

	result := r.data.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *UserRepo) Delete(_ context.Context, req *userV1.DeleteUserRequest) (bool, error) {
	result := r.data.db.Delete(&models.User{}, req.GetId())
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
