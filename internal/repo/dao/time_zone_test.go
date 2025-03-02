package dao

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

type TestModel struct {
	Id        int64 `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

func openDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("postgres://root:root@localhost:15432/ink_flow?sslmode=disable"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestInsertDB(t *testing.T) {
	db, err := openDB()
	require.NoError(t, err)

	db.Save(&TestModel{
		Id:        1,
		Name:      "default",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	db.Save(&TestModel{
		Id:        2,
		Name:      "specify utc",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	db.Save(&TestModel{
		Id:        3,
		Name:      "utc now",
		CreatedAt: time.Now().Add(-time.Hour * 8).UTC(),
		UpdatedAt: time.Now().Add(-time.Hour * 8).UTC(),
	})
}

func TestGet(t *testing.T) {
	db, err := openDB()
	require.NoError(t, err)

	m := TestModel{}
	err = db.Where("name = ?", "default").First(&m).Error
	require.NoError(t, err)
	fmt.Println(m)

	m = TestModel{}
	err = db.Where("name = ?", "specify utc").First(&m).Error
	require.NoError(t, err)
	fmt.Println(m)

	m = TestModel{}
	err = db.Where("name = ?", "utc now").First(&m).Error
	require.NoError(t, err)
	fmt.Println(m)
}
