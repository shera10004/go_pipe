package auth_google

import (
	"custom/cpath"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Authweb struct {
	ClientID     string   `json:"client_id"`
	ProjectID    string   `json:"project_id"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
	AuthProvider string   `json:"auth_provider"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIS []string `json:"redirect_uris"`
}
type AuthData struct {
	WebData Authweb `json:"web"`
	isSet   bool
}

func (authdata AuthData) IsSet() bool {
	return authdata.isSet
}

func (authdata AuthData) GetClientID() string {
	return authdata.WebData.ClientID
}

func (authdata AuthData) GetClientSecret() string {
	return authdata.WebData.ClientSecret
}

func (authdata AuthData) GetRedirectURI() string {
	if len(authdata.WebData.RedirectURIS) < 1 {
		return ""
	}
	return authdata.WebData.RedirectURIS[0]
}

func GetAuthData() AuthData {

	authdata := AuthData{isSet: false}

	currentPath := cpath.GetCurrent()

	//auth_goole/auth.json 파일 로드 ( 인증서 파일 )
	bytes, err := ioutil.ReadFile(currentPath + "/authprovider/auth_google/auth.json")
	if err != nil {
		log.Println("auth file read err :", err)
	} else {
		fmt.Println(string(bytes))

		json.Unmarshal(bytes, &authdata)
		authdata.isSet = true

		fmt.Println("- google_auth 2.0 - ")
		fmt.Printf("%#+v \n", authdata)
		fmt.Println()
	}

	return authdata
}
