package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bufbuild/connect-go"

	pb "github.com/el-zacharoo/auth/internal/gen/auth/v1"
	pbcnn "github.com/el-zacharoo/auth/internal/gen/auth/v1/authv1connect"
)

type SignInServer struct {
	pbcnn.UnimplementedAuthServiceHandler
}

const auth0Domain = "https://react-messaging.au.auth0.com"

func (s *SignInServer) SignIn(ctx context.Context, req *connect.Request[pb.SignInRequest]) (*connect.Response[pb.SignInResponse], error) {
	reqMsg := req.Msg
	auth := reqMsg.AuthUserSignIn

	json_data, err := json.Marshal(auth)
	if err != nil {
		fmt.Println(err)
	}

	url := auth0Domain + "/oauth/token"
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
	}
	r.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Print(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	respBody := []byte(body)
	var rep map[string]interface{}
	json.Unmarshal(respBody, &rep)

	resp := &pb.SignInResponse{
		AccessToken: rep["access_token"].(string),
		Scope:       rep["scope"].(string),
		ExpiresIn:   int32(rep["expires_in"].(float64)),
		IdToken:     rep["id_token"].(string),
		TokenType:   rep["token_type"].(string),
	}

	return connect.NewResponse(resp), nil
}

func (s *SignInServer) GetAccount(ctx context.Context, req *connect.Request[pb.GetAccountRequest]) (*connect.Response[pb.GetAccountResponse], error) {
	reqMsg := req.Msg.AccessToken
	token := reqMsg

	url := auth0Domain + "/userinfo"
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Print(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	respBody := []byte(body)
	var rep map[string]interface{}
	json.Unmarshal(respBody, &rep)

	resp := &pb.GetAccountResponse{
		UserInfo: &pb.UserInfo{
			Sub:       rep["sub"].(string),
			Name:      rep["name"].(string),
			Nickname:  rep["nickname"].(string),
			Picture:   rep["picture"].(string),
			UpdatedAt: rep["updated_at"].(string),
			Email:     rep["email"].(string),
		},
	}

	return connect.NewResponse(resp), nil
}
