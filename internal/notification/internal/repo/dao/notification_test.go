package dao

import (
	"context"
	"fmt"
	"github.com/KNICEX/InkFlow/pkg/snowflakex"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
	"time"
)

func TestGormNotificationDAO_FindLikeMerge(t *testing.T) {
	db, err := gorm.Open(postgres.Open("host=localhost user=root password=root dbname=ink_flow port=15432"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&Notification{}))

	defer db.Exec("truncate table notification")

	dao := NewGormNotificationDAO(db, snowflakex.NewNode(snowflakex.DefaultStartTime, 0))

	for i := range 1000 {
		err = dao.Insert(context.Background(), Notification{
			NotificationType: "like",
			SubjectType:      "ink",
			SubjectId:        int64(i % 20),
			RecipientId:      1,
			SenderId:         int64(i),
			Content:          "",
			Read:             false,
		})
		require.NoError(t, err)
	}

	res, err := dao.FindLikeMerge(context.Background(), 1, time.Now().UnixMilli(), 0, 10)
	require.NoError(t, err)

	fmt.Printf("res: %+v\n", res)

}
