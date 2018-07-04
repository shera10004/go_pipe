package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Id        string
	AddressId string
}

const (
	PATH_ROOt          = "/"
	PATH_ABOUT         = "/about"
	PATH_LOGIN         = "/login"
	PATH_USERS_ID      = "/users/:id"
	PATH_USERS_ID_ADDR = "/users/:user_id/addresses/:address_id"
)

func main() {

	s := NewServer()

	s.HandleFunc("GET", PATH_ROOt, func(c *Context) {
		fmt.Println("[HandleFunc] -- index.html")
		c.RenderTemplate("/public/index.html", map[string]interface{}{"time": time.Now()})
	})

	s.HandleFunc("GET", PATH_ABOUT, func(c *Context) {
		fmt.Println("[HandleFunc] -- about")
		fmt.Fprintln(c.ResponseWriter, "about")
	})

	s.HandleFunc("GET", PATH_USERS_ID, func(c *Context) {
		fmt.Println("[HandleFunc] -- users/:id")
		u := User{Id: c.Params["id"].(string)}
		c.RenderXml(u)
	})

	s.HandleFunc("GET", PATH_USERS_ID_ADDR, func(c *Context) {
		fmt.Println("[HandleFunc] -- users/:user_id/addresses/:address_id")
		u := User{Id: c.Params["user_id"].(string),
			AddressId: c.Params["address_id"].(string),
		}
		c.RenderJson(u)
	})

	s.HandleFunc("GET", PATH_LOGIN, func(c *Context) {
		fmt.Println("[HandleFunc] -- login[GET]")
		c.RenderTemplate("/public/login.html", map[string]interface{}{"message": "로그인이필요합니다"})
	})

	s.HandleFunc("POST", PATH_LOGIN, func(c *Context) {
		fmt.Println("[HandleFunc] -- login[POST]")
		//로그인 정보를 확인하여 쿠키에 인증 토큰 값 기록
		fmt.Println("params len:", len(c.Params))
		for k, v := range c.Params {
			fmt.Printf("key:%v , value:%v \n", k, v)
		}

		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			fmt.Println("- 로그인 성공!")
			cookie := &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign(VerifyMessage),
				Path:  "/",
			}
			http.SetCookie(c.ResponseWriter, cookie)
			c.Redirect("/")
			return
		}
		c.RenderTemplate("/public/login.html",
			map[string]interface{}{"message": "id 또는 password가 일치하지 않습니다"})
	})

	s.Use(AuthHandler)

	s.Run(":8080")

}

const VerifyMessage = "verified"

func AuthHandler(next ContextFunc) ContextFunc {
	fmt.Println("<AuthHandler")
	ignore := []string{"/login", "public/index.html"}
	return func(c *Context) {
		log.Println("--AuthHandler URL_PATH:", c.Request.URL.Path)

		// url prefix가 "/login" , "public/index.html"이면 auth를 체크하지 않음.
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				fmt.Println("AuthHandler - HasPrefix", c.Request.URL.Path)
				next(c)
				return
			}
		} //for

		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			log.Println("AuthHandler - no cookie")
			//"X_AUTH"쿠키 값이 없으면 "/login"으로 이동
			c.Redirect("/login")
			return
		} else if err != nil {
			log.Println("AuthHandler - error")
			//에러처리
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify(VerifyMessage, v.Value) {
			log.Println("AuthHandler - verify")
			//쿠키값이 인증으로 확인되면 다음 핸들러로 넘어감
			next(c)
			return
		} else {
			fmt.Println("AuthHandler - Verify not :", v.Value)
		}

		fmt.Println("AuthHandler - login 인증 안됨.")
		c.Redirect(PATH_LOGIN)

	}
}

//인증 토큰 확인
func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}

func CheckLogin(username, password string) bool {
	//로그인처리
	const (
		USERNAME = "tester"
		PASSWORD = "1234"
	)
	return username == USERNAME && password == PASSWORD
}

const secretMessage = "golang-book-secret-key2"

//인증 토큰 생성
func Sign(message string) string {
	secretKey := []byte(secretMessage)
	if len(secretKey) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}
