package corero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"lscor/config"
	"net/http"
)

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Cookie *http.Cookie

func FetchToken() error {
	cfg := config.GlobalConfig.Corero
	credentials := login{
		Username: cfg.User,
		Password: cfg.Pass,
	}
	payload, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal login credentials: %s", err)
	}

	res, err := http.Post(cfg.URL+"auth", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to post auth: %s", err)
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "Authentication" {
			Cookie = cookie
		}
	}

	log.Println("Finished logging in")
	return nil
}

func Logout() error {
	req, err := http.NewRequest("POST", config.GlobalConfig.Corero.URL+"log_out", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("content-type", "application/json")
	if err != nil {
		return fmt.Errorf("failed to format postt: %s", err)
	}

	req.AddCookie(Cookie)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("got non 200 status code: %s", res.Status)
	}

	log.Println("Finished logging out")
	return nil
}
