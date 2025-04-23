package keycloak

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	tok, err := AccessToken()
	assert.Nil(t, err)

	_, err = verifyToken(tok)
	fmt.Println(err)
	assert.NotNil(t, err)

	tok, err = AccessToken()
	assert.Nil(t, err)
}
