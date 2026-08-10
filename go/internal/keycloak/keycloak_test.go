package keycloak

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	if os.Getenv("AUTODARTS_INTEGRATION_TEST") != "1" {
		t.Skip("set AUTODARTS_INTEGRATION_TEST=1 to run against Keycloak")
	}

	tok, err := AccessToken()
	assert.Nil(t, err)

	_, err = verifyToken(tok)
	fmt.Println(err)
	assert.NotNil(t, err)

	tok, err = AccessToken()
	assert.Nil(t, err)
}
