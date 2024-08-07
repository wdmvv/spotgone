package spoton

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ClientID     string
	ClientSecret string
	AuthToken    string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// this should be in some kind of config but eh
var tokenEndpoint string = "https://accounts.spotify.com/api/token"

// Gets auth token and places in User struct
func (u *User) GetAuth() error {
	auth := fmt.Sprintf("%s:%s", u.ClientID, u.ClientSecret)
	authb64 := b64.StdEncoding.EncodeToString([]byte(auth))

	data := []byte("grant_type=client_credentials")

	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+authb64)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return &ErrBadAuth{resp.StatusCode}
	}

	d := json.NewDecoder(resp.Body)
	return d.Decode(&u)
}
