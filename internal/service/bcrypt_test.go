package service

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestBCrypt(t *testing.T) {
	r, err := bcrypt.GenerateFromPassword([]byte("hellosadadasdasdasdasda"), bcrypt.DefaultCost)
	require.NoError(t, err)
	fmt.Println(len(string(r)))
	fmt.Println(string(r))
}
