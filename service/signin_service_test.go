package service

import (
	"testing"

	"github.com/mishozz/Library/auth"
	"github.com/stretchr/testify/assert"
)

var sign = signInStruct{}

func TestSignInStruct_SignIn(t *testing.T) {
	var authD auth.AuthDetails
	authD.UserId = 1
	authD.AuthUuid = "83b09612-9dfc-4c1d-8f7d-a589acec7081"

	token, err := sign.SignIn(authD)
	assert.Nil(t, err)
	assert.NotNil(t, token)
}
