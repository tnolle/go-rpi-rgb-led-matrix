package keycloak

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	baseURL      = "https://login.autodarts.io"
	realm        = "autodarts"
	clientID     string
	clientSecret string
)

func Init(id, secret string) {
	clientID = id
	clientSecret = secret
}

func realmURL() string {
	return fmt.Sprintf("%s/realms/%s", baseURL, realm)
}

func tokenURL() string {
	return fmt.Sprintf("%s/protocol/openid-connect/token", realmURL())
}

var (
	token     TokenResponse
	publicKey *rsa.PublicKey
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func fetchToken() error {
	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)

	res, err := http.PostForm(tokenURL(), params)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return err
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		return err
	}

	_, err = verifyToken(token.AccessToken)
	if err != nil {
		return err
	}

	return nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	t, err := jwt.NewParser(
		jwt.WithIssuer(realmURL()),
		jwt.WithAudience("account"),
	).ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return fetchPublicKey(), nil
		},
	)
	if err != nil {
		return nil, err
	}
	exp, err := t.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	if time.Now().After(exp.Add(-1 * time.Minute)) {
		return nil, jwt.ErrTokenExpired
	}
	return t, nil
}

type Realm struct {
	Realm             string `json:"realm"`
	PublicKey         string `json:"public_key"`
	TokenServiceURL   string `json:"token-service"`
	AccountServiceURL string `json:"account-service"`
	TokensNotBefore   int    `json:"tokens-not-before"`
}

func fetchPublicKey() *rsa.PublicKey {
	if publicKey == nil {
		var r Realm
		res, err := http.Get(realmURL())
		if err != nil {
			return nil
		}
		json.NewDecoder(res.Body).Decode(&r)
		defer res.Body.Close()
		pemKey, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + r.PublicKey + "\n-----END PUBLIC KEY-----"))
		key, _ := x509.ParsePKIXPublicKey(pemKey.Bytes)
		publicKey, _ = key.(*rsa.PublicKey)
	}
	return publicKey
}

func AccessToken() (string, error) {
	_, err := verifyToken(token.AccessToken)
	if err != nil {
		err = fetchToken()
		if err != nil {
			return "", err
		}
	}
	return token.AccessToken, nil
}
