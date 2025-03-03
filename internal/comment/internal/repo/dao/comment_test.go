package dao

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"testing"
)

func openDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("postgres://root:root@localhost:15432/ink_flow?sslmode=disable"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func clearDB(db *gorm.DB, tbs ...any) error {
	// 清空所有表
	return db.Migrator().DropTable(tbs...)
}

func TestMultiLike(t *testing.T) {
	db, err := openDB()
	require.NoError(t, err)
	err = clearDB(db, &CommentLike{}, &CommentStatistic{}, &Comment{})
	require.NoError(t, err)
	Init(db)

	dao := NewGormCommentDAO(db)
	err = dao.Insert(context.Background(), Comment{
		UserId:      2,
		Content:     "test",
		Biz:         "ink",
		BizId:       1,
		Status:      CommentStatusPassed,
		ParentId:    -1,
		RootId:      -1,
		ReplyUserId: -1,
	})
	require.NoError(t, err)

	err = dao.Like(context.Background(), 1, 2)
	require.NoError(t, err)
	// 模拟重复调用
	err = dao.Like(context.Background(), 1, 2)
	require.NoError(t, err)

	// 并发
	wg := sync.WaitGroup{}
	for i := range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 30 {
				err = dao.Like(context.Background(), int64(i*10+j), 2)
				switch {
				case err == nil:
				case errors.Is(err, gorm.ErrDuplicatedKey):
				default:
					t.Log(err)
				}

			}
		}()
	}
	wg.Wait()
}
