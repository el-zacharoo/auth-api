package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
)

type SignInServer struct{}

type Credentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	GrantType    string `json:"grant_type"`
	Audience     string `json:"audience"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (s *SignInServer) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	defer r.Body.Close()
	reqByt, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	}

	var cred Credentials
	json.Unmarshal(reqByt, &cred)
	json_data, _ := json.Marshal(cred)

	url := "https://react-messaging.au.auth0.com/oauth/token"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(body)

}

type UserToken struct {
	AccessToken string `json:"access_token"`
}

func (s *SignInServer) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	defer r.Body.Close()
	reqByt, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v", err)))
	}

	var tkn UserToken
	json.Unmarshal(reqByt, &tkn)

	url := "https://react-messaging.au.auth0.com/userinfo"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Authorization", "Bearer "+tkn.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(body)
}
