package auth

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	au := AuthDetails{
		AuthUuid: "43b78a87-6bcf-439a-ab2e-940d50c4dc33", //this can be anything
		UserId:   1,
	}
	token, err := CreateToken(au)
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestVerifyToken(t *testing.T) {
	//In order to generate a request, let use the logout endpoint
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Error(err)
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiJjMmUxYjBjMy00ZGRjLTQ0NjUtYWVkNC1iNGE2NDM5NzI4M2MiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjoxfQ.FWbfdhEJeK7mjZ-lWvs9scuyUrSKPrC4xafUoEqkduc"

	tokenString := fmt.Sprintf("Bearer %v", token)
	req.Header.Set("Authorization", tokenString)

	jwtAns, err := VerifyToken(req)

	assert.Nil(t, err)
	assert.NotNil(t, jwtAns) //this is of type *jwt.Token
}

func TestExtractToken(t *testing.T) {
	//In order to generate a request, let use the logout endpoint
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Error(err)
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiJjMmUxYjBjMy00ZGRjLTQ0NjUtYWVkNC1iNGE2NDM5NzI4M2MiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjoxfQ.FWbfdhEJeK7mjZ-lWvs9scuyUrSKPrC4xafUoEqkduc"

	tokenString := fmt.Sprintf("Bearer %v", token)
	req.Header.Set("Authorization", tokenString)

	result := ExtractToken(req)
	assert.NotNil(t, result)
	assert.EqualValues(t, result, token)
}

//Check the auth details from the token:
func TestExtractTokenAuth(t *testing.T) {
	//In order to generate a request, let use the logout endpoint
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Error(err)
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiIzZTY0NjM0MS0wMzdhLTRmZjMtOGViYy0wNzRlNDdmZTNmMzAiLCJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTAxODg4NzIsInVzZXJfaWQiOjEsInVzZXJfcm9sZSI6IkFkbWluIn0.zxTRJEJhwABZXQ04bMxwwCJDvgvkgt8S1yyWQFmEE-s"

	tokenString := fmt.Sprintf("Bearer %v", token)
	req.Header.Set("Authorization", tokenString)

	result, err := ExtractTokenAuth(req)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.UserId)
	assert.NotNil(t, result.AuthUuid)
}

func TestTokenValid(t *testing.T) {
	//In order to generate a request, let use the logout endpoint
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Error(err)
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiJjMmUxYjBjMy00ZGRjLTQ0NjUtYWVkNC1iNGE2NDM5NzI4M2MiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjoxfQ.FWbfdhEJeK7mjZ-lWvs9scuyUrSKPrC4xafUoEqkduc"

	tokenString := fmt.Sprintf("Bearer %v", token)
	req.Header.Set("Authorization", tokenString)

	errToken := TokenValid(req)
	assert.Nil(t, errToken)
}

//i added garbage to the token, so is not valid
func TestToken_Invalid(t *testing.T) {
	//In order to generate a request, let use the logout endpoint
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Error(err)
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiJjMmUxYjBjMy00ZGRjLTQ0NjUtYWVkNC1iNGE2NDM5NzI4M2MiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjoxfQ.FWbfdhEJeK7mjZ-lWvs9scuyUrSKPrC4xafUoEqkducxx"

	tokenString := fmt.Sprintf("Bearer %v", token)
	req.Header.Set("Authorization", tokenString)

	errToken := TokenValid(req)
	assert.NotNil(t, errToken)
	assert.EqualValues(t, "illegal base64 data at input byte 45", errToken.Error())
}
