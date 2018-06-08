package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	*router
	middlewares  []Middleware
	startHandler ContextFunc
}

func NewServer() *Server {

	r := &router{make(map[string]map[string]ContextFunc)}

	s := &Server{router: r}

	s.middlewares = []Middleware{
		logHandler,
		recoverHandler,
		staticHandler,
		parseFormHandler,
		parseJSONBodyHandler,
	}

	fmt.Println("> NewServer")
	return s
}

func (s *Server) Run(addr string) {

	//startHandler를 라우터 핸들러 함수로 지정
	s.startHandler = s.router.getContextFunc()

	//등록된 미들웨어를 라우터 핸들러 앞에 하나씩 추가
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.startHandler = s.middlewares[i](s.startHandler)
	}

	fmt.Println("> web server Run...")
	//웹 서버 시작
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ServeHTTP")

	//Context 생성
	c := &Context{
		Params:         make(map[string]interface{}),
		ResponseWriter: w,
		Request:        r,
	}
	for k, v := range r.URL.Query() {
		fmt.Printf("query ,k:%v , v:%+v \n", k, v)
		c.Params[k] = v[0]
	}

	s.startHandler(c)
}

func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}
