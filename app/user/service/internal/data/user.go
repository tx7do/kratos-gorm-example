package data

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/crypto"
	"github.com/tx7do/go-utils/mapper"

	pagination "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	gormCurd "github.com/tx7do/go-crud/gorm"

	"kratos-gorm-example/app/user/service/internal/data/models"

	userV1 "kratos-gorm-example/api/gen/go/user/service/v1"
)

type UserRepo struct {
	data *Data
	log  *log.Helper

	mapper     *mapper.CopierMapper[userV1.User, models.User]
	repository *gormCurd.Repository[userV1.User, models.User]
}

func NewUserRepo(data *Data, logger log.Logger) *UserRepo {
	l := log.NewHelper(log.With(logger, "module", "user/repo/user-service"))
	repo := &UserRepo{
		data:   data,
		log:    l,
		mapper: mapper.NewCopierMapper[userV1.User, models.User](),
	}

	repo.repository = gormCurd.NewRepository[userV1.User, models.User](
		repo.mapper,
	)

	repo.init()

	return repo
}

func (r *UserRepo) init() {
	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *UserRepo) List(ctx context.Context, req *pagination.PagingRequest) (*userV1.ListUserResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	ret, err := r.repository.ListWithPaging(ctx, r.data.db, req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &userV1.ListUserResponse{Total: 0, Items: nil}, nil
	}

	return &userV1.ListUserResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *UserRepo) Get(ctx context.Context, req *userV1.GetUserRequest) (*userV1.User, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}

	var whereCond *gorm.DB
	switch req.QueryBy.(type) {
	case *userV1.GetUserRequest_Id:
		whereCond = r.data.db.Where("id = ?", req.GetId())
	case *userV1.GetUserRequest_UserName:
		whereCond = r.data.db.Where("user_name = ?", req.GetUserName())
	default:
		whereCond = r.data.db.Where("id = ?", req.GetId())
	}

	dto, err := r.repository.Get(ctx, whereCond, req.GetViewMask())
	if err != nil {
		return nil, err
	}

	return dto, err
}

func (r *UserRepo) Create(ctx context.Context, req *userV1.CreateUserRequest) (*userV1.User, error) {
	if req == nil || req.Data == nil {
		return nil, errors.New("request is nil")
	}

	if req.Data.Password != nil && req.Data.GetPassword() != "" {
		cryptoPassword, err := crypto.HashPassword(req.Data.GetPassword())
		if err != nil {
			return nil, err
		}
		req.Data.Password = &cryptoPassword
	}

	result, err := r.repository.Create(ctx, r.data.db, req.Data, nil)

	return result, err
}

func (r *UserRepo) Update(ctx context.Context, req *userV1.UpdateUserRequest) (*userV1.User, error) {
	if req == nil || req.Data == nil {
		return nil, errors.New("request is nil")
	}

	if req.Data.Password != nil && req.Data.GetPassword() != "" {
		cryptoPassword, err := crypto.HashPassword(req.Data.GetPassword())
		if err != nil {
			return nil, err
		}
		req.Data.Password = &cryptoPassword
	}

	result, err := r.repository.Update(ctx, r.data.db, req.Data, req.GetUpdateMask())

	return result, err
}

func (r *UserRepo) Upsert(ctx context.Context, req *userV1.UpdateUserRequest) (*userV1.User, error) {
	if req == nil || req.Data == nil {
		return nil, errors.New("request is nil")
	}

	var err error

	if req.Data.Password != nil && req.Data.GetPassword() != "" {
		cryptoPassword, err := crypto.HashPassword(req.Data.GetPassword())
		if err != nil {
			return nil, err
		}
		req.Data.Password = &cryptoPassword
	}

	result, err := r.repository.Upsert(ctx, r.data.db, req.Data, req.GetUpdateMask())

	return result, err
}

func (r *UserRepo) Delete(ctx context.Context, req *userV1.DeleteUserRequest) (bool, error) {
	if req == nil {
		return false, errors.New("request is nil")
	}

	result, err := r.repository.Delete(ctx, r.data.db.Where("id = ?", req.GetId()))

	return result > 0, err
}
