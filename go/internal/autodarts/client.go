package autodarts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/keycloak"
)

type AutodartsAPIClient struct {
	URL string
}

func NewAutodartsAPIClient() *AutodartsAPIClient {
	return &AutodartsAPIClient{URL: "https://api.autodarts.io"}
}

func (c *AutodartsAPIClient) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, c.URL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	tok, err := keycloak.AccessToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	return http.DefaultClient.Do(req)
}

func (c *AutodartsAPIClient) GetOnlineUsers() int {
	res, err := c.doRequest("GET", "/us/v0/users/online", nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return 0
	}
	var result OnlineUserMessage
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return 0
	}
	defer res.Body.Close()
	return result.Online
}

func (c *AutodartsAPIClient) GetNumMatches() int {
	res, err := c.doRequest("GET", "/gs/v0/matches/count", nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return 0
	}
	var result MatchesCountMessage
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return 0
	}
	defer res.Body.Close()
	return result.Count
}
