package gormx

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func CheckDuplicateErr(err error) (error, bool) {
	if err == nil {
		return nil, false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key") {
		return fmt.Errorf("%w, err: %s", gorm.ErrDuplicatedKey, err.Error()), true
	}
	return err, false
}
