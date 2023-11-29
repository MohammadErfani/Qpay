package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashingPassword(t *testing.T) {
	hp, err := HashPassword("123asdf@#BNM")
	fmt.Println(hp)
	assert.NoError(t, err)
	assert.NotEqual(t, "", hp)
}

func TestCheckPassword(t *testing.T) {
	hp, _ := HashPassword("123asdf@#BNM")
	assert.True(t, CheckPassword("123asdf@#BNM", hp))
	assert.False(t, CheckPassword("123asdf@#BN", hp))
	assert.False(t, CheckPassword("123asdf@#BNm", hp))
	assert.False(t, CheckPassword("123asdf#BNM", hp))
	assert.False(t, CheckPassword("123asdF@#BNM", hp))

}
