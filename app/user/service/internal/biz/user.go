package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	pagination "github.com/tx7do/kratos-bootstrap/api/gen/go/pagination/v1"
	userV1 "kratos-gorm-example/api/gen/go/user/service/v1"
)

type UserRepo interface {
	List(ctx context.Context, req *pagination.PagingRequest) (*userV1.ListUserResponse, error)
	Get(ctx context.Context, req *userV1.GetUserRequest) (*userV1.User, error)
	Create(ctx context.Context, req *userV1.CreateUserRequest) (*userV1.User, error)
	Update(ctx context.Context, req *userV1.UpdateUserRequest) (*userV1.User, error)
	Delete(ctx context.Context, req *userV1.DeleteUserRequest) (bool, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	l := log.NewHelper(log.With(logger, "module", "user/usecase/user-service"))
	return &UserUseCase{
		repo: repo,
		log:  l,
	}
}

func (uc *UserUseCase) Get(ctx context.Context, req *userV1.GetUserRequest) (*userV1.User, error) {
	user, err := uc.repo.Get(ctx, req)
	if user != nil {
		user.Password = nil
	}
	return user, err
}

func (uc *UserUseCase) List(ctx context.Context, req *pagination.PagingRequest) (*userV1.ListUserResponse, error) {
	resp, err := uc.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(resp.Items); i++ {
		resp.Items[i].Password = nil
	}

	return resp, err
}

func (uc *UserUseCase) Create(ctx context.Context, req *userV1.CreateUserRequest) (*userV1.User, error) {
	user, err := uc.repo.Create(ctx, req)
	if user != nil {
		user.Password = nil
	}
	return user, err
}

func (uc *UserUseCase) Update(ctx context.Context, req *userV1.UpdateUserRequest) (*userV1.User, error) {
	user, err := uc.repo.Update(ctx, req)
	if user != nil {
		user.Password = nil
	}
	return user, err
}

func (uc *UserUseCase) Delete(ctx context.Context, req *userV1.DeleteUserRequest) (bool, error) {
	return uc.repo.Delete(ctx, req)
}
