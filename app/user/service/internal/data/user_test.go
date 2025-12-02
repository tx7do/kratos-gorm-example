package data

import (
	"context"
	"io"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	pagination "github.com/tx7do/go-curd/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	userV1 "kratos-gorm-example/api/gen/go/user/service/v1"
	"kratos-gorm-example/app/user/service/internal/data/models"
)

func setupInMemoryDB(t *testing.T) *Data {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite memory db: %v", err)
	}
	if err = db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return &Data{db: db}
}

func strPtr(s string) *string { return &s }

func TestUserRepo_CreateGetListDelete(t *testing.T) {
	data := setupInMemoryDB(t)
	repo := NewUserRepo(data, log.NewStdLogger(io.Discard))
	ctx := context.Background()

	// Create
	in := &userV1.User{
		Id:       1,
		UserName: strPtr("alice"),
		NickName: strPtr("Alice"),
		Password: strPtr("password123"),
	}
	created, err := repo.Create(ctx, &userV1.CreateUserRequest{User: in})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created == nil {
		t.Fatalf("Create returned nil user")
	}
	if created.GetUserName() != "alice" {
		t.Fatalf("created username mismatch: want alice got %s", created.GetUserName())
	}
	if created.Id == 0 {
		t.Fatalf("created id empty")
	}

	// Get by id
	got, err := repo.Get(ctx, &userV1.GetUserRequest{QueryBy: &userV1.GetUserRequest_Id{Id: created.GetId()}})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got == nil || got.GetUserName() != "alice" {
		t.Fatalf("Get returned unexpected user: %+v", got)
	}

	// List
	listReq := &pagination.PagingRequest{Page: trans.Ptr(uint32(1)), PageSize: trans.Ptr(uint32(10))}
	listResp, err := repo.List(ctx, listReq)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if listResp == nil || listResp.Total < 1 || len(listResp.Items) < 1 {
		t.Fatalf("List returned empty result: %+v", listResp)
	}

	// Delete
	deleted, err := repo.Delete(ctx, &userV1.DeleteUserRequest{Id: created.GetId()})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !deleted {
		t.Fatalf("Delete returned false")
	}

	// After delete, try get (expect error or nil depending on implementation)
	after, err := repo.Get(ctx, &userV1.GetUserRequest{QueryBy: &userV1.GetUserRequest_Id{Id: created.GetId()}})
	// 接受两种合理结果：返回 nil (repo 层保护) 或返回 gorm.ErrRecordNotFound 的错误
	if err == nil && after != nil {
		t.Fatalf("expected not found after delete, got: %+v", after)
	}
}
