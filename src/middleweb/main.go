package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/urfave/negroni"
	//"github.com/codegangsta/negroni"	[obsolute] -> github.com/urfave/negroni
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/websocket"
)

const (
	//애플리케이션에서 사용할 세션의 키 정보
	sessionKey    = "simple_chat_session"
	sessionSecret = "simple_chat_session_secret"

	socketBufferSize = 1024
)

var (
	renderer     *render.Render
	mongoSession *mgo.Session

	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

func init() {
	fmt.Println("main init()")

	//렌더러 생성
	renderer = render.New()

	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	mongoSession = s
}

func main() {
	//라우터 생성
	router := httprouter.New()

	//핸들러 정의
	//type Handle func(http.ResponseWriter, *http.Request, Params)
	/*
		 , func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			})
			//*/

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		//렌더러를 사용하여 템플릿 렌더링
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Simple Chat!"})
	})

	router.GET("/login", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		//로그인 페이지 렌더링
		renderer.HTML(w, http.StatusOK, "login", nil)
	})
	router.GET("/logout", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		//세션에서 사용자 정보 제거 후 로그인 페이지로 이동
		sessions.GetSession(req).Delete(currentUserKey)
		http.Redirect(w, req, "/login", http.StatusFound)
	})

	router.GET("/auth/:action/:provider", loginHandler)

	router.POST("/"+C_ROOMS, createRoom)
	router.GET("/"+C_ROOMS, retrieveRooms)

	router.GET("/"+C_ROOMS+"/:id/"+C_MESSAGES, retrieveMessage)

	router.GET("/ws/:room_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("ServeHTTP:", err)
			return
		}
		newClient(socket, ps.ByName("room_id"), GetCurrentUser(r))
	})

	//negroni 미들웨어 생성
	n := negroni.Classic()
	store := cookiestore.New([]byte(sessionSecret))

	n.Use(sessions.Sessions(sessionKey, store))

	n.Use(LoginRequired("/login", "/auth"))

	//negroni에 router를 핸들러로 등록
	n.UseHandler(router)

	log.Println("server run...")
	//웹서버 실행
	n.Run(":3000")

}
