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
}

type Album struct {
	Name   string `json:"name"`
	Total  int    `json:"total"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Release    string `json:"release_date"`
	Label      string `json:"label"`
	Popularity int    `json:"popularity"`
	Tracks     []struct {
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Name   string `json:"name"`
		Number int    `json:"track_number"`
		Disc   int    `json:"disc_number"`
	} `json:"items"`
}

// this should be in some kind of config but eh
var (
	tokenEndpoint string = "https://accounts.spotify.com/api/token"
)

// Gets auth token and sets in User struct
func (u *User) SetAuth() error {
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
	} else if resp.StatusCode != http.StatusOK {
		return &ErrBadRequest{"auth", resp.StatusCode}
	}

	d := json.NewDecoder(resp.Body)
	return d.Decode(&u)
}

// if something's missing then one should set market in request
func (u *User) GetAlbum(id string) (Album, error) {
	var a Album
	if u.AuthToken == "" {
		return Album{}, &ErrNoAuth{}
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Album{}, err
	}

	req.Header.Add("Authorization", "Bearer "+u.AuthToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Album{}, err
	} else if resp.StatusCode != http.StatusOK {
		return Album{}, &ErrBadRequest{"album", resp.StatusCode}
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&a)
	if err != nil {
		return Album{}, err
	}

	url = fmt.Sprintf("https://api.spotify.com/v1/albums/%s", id)
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Album{}, err
	}

	req.Header.Add("Authorization", "Bearer "+u.AuthToken)
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		return Album{}, err
	} else if resp.StatusCode != http.StatusOK {
		return Album{}, &ErrBadRequest{"album", resp.StatusCode}
	}

	d = json.NewDecoder(resp.Body)
	err = d.Decode(&a)
	return a, err
}
