package main

import (
	"fmt"
	"net/http"
)

func main() {

	//mux := http.NewServeMux()	//Http요청 멀티플렉서 생성 - url경로를 동시에 여러개 처리 할수 있음.
	//mux.HandleFunc()
	//mus.ListenAndServe(":80" , mux)	//mux를 이용하여 HTTP요청을 처리.

	//http.HandlFunc()도 내부적으로 mux와 같은 기능입니다.

	s := "Hello, world!"

	//함수 매개변수 req에는 (GET, POST, PUT, DELETE 등),쿠키,헤더 등이 들어있고, res로 웹 브라우저에 응답해 줄수 있음.
	//여기서 웹 브라우저가 /경로에 접속하면 HTML을 만들어서 Write함수로 응답 해줌.
	//그리고, res.Header().Set("Content-Type" , "text/html")처럼 현재 경로의 헤더 값도 정해줌.

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		html := `
		<html>
		<head>
			<title>Hello</title>
			<script type="text/javascript" src="/assets/hello.js"></script>
			<link href="/assets/hello.css" rel="stylesheet" />
		</head>
		<body>
			<span class="hello">` + s + `</span>
			<br/>
			<span class="how">How ara you?</span>
			
		</body>
		</html>
		`

		res.Header().Set("Content-Type", "text/html") //HTML 헤더 설정
		res.Write([]byte(html))                       //웹 브라우저에 응답
	})

	//http.Handle함수에 url 하위 경로인 /assets/를 설정해주고, http.FileServer, http.Dir 함수에 assets
	//디렉터리를 지정합니다. 단, http.Dir함수에 assets디렉터리를 설정했으므로 웹서버 입장에서 hello.js파일
	//은 /assets/hello.js가 아닌 hello.js로 접근해야 합니다. 따라서 http.StripPrefix함수로 /assets/경로를
	//삭제 해줍니다. ( /assets/hello.js -> ./hello.js )

	http.Handle( // /assets/경로에 접근했을 때 파일 서버를 동작시킴.
		"/assets/",
		http.StripPrefix( //파일 서버를 실행할 때 assets폴더를 지정했으므로 URL 경로에서 /assets/삭제
			"/assets/",
			http.FileServer(http.Dir("assets")), //웹 서버에서 assets 디렉터리 아래의 파일표시
		),
	)

	fmt.Println("WebServer run...")

	http.ListenAndServe(":80", nil) //80번 포트에서 웹서버 실행

}
