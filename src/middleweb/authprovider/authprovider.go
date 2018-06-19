package authprovider

import (
	"pipe/src/middleweb/authprovider/auth_google"
)

type iAuthClient interface {
	IsSet() bool
	GetClientID() string
	GetClientSecret() string
	GetRedirectURI() string
}

const (
	Google = iota
	Facebook
)

var authClient map[int]iAuthClient = make(map[int]iAuthClient)

func init() {

	googleData := auth_google.GetAuthData()

	authClient[Google] = googleData

}

func GetAuthClient(key int) iAuthClient {
	v, ok := authClient[key]
	if ok {
		return v
	}
	return nil
}
