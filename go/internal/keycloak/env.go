//go:build with_env

package keycloak

import "os"

func init() {
	if url, ok := os.LookupEnv("AUTODARTS_KEYCLOAK_URL"); ok {
		baseURL = url
	}
	if r, ok := os.LookupEnv("AUTODARTS_KEYCLOAK_REALM"); ok {
		realm = r
	}
	if client, ok := os.LookupEnv("AUTODARTS_KEYCLOAK_CLIENT_ID"); ok {
		clientID = client
	}
	if secret, ok := os.LookupEnv("AUTODARTS_KEYCLOAK_CLIENT_SECRET"); ok {
		clientSecret = secret
	}
}
