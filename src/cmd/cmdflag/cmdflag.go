package main

import (
	"flag"
	"fmt"
)

func main() {

	title := flag.String("title", "", "영화 이름")
	runtime := flag.Int("runtime", 0, "상영 시간")
	rating := flag.Float64("rating", 0.0, "평점")
	release := flag.Bool("release", false, "개봉 여부")

	flag.Parse() //명령줄 옵션의 내용을 각 자료형별로 분석
	if flag.NFlag() == 0 {
		flag.Usage() //명령줄 기본사용법 출력
		return
	}

	fmt.Println("인자 갯수:", flag.NFlag())

	fmt.Printf(
		" 영화 이름:%s\n 상영 시간:%d\n 평점:%f\n",
		*title,
		*runtime,
		*rating,
	)

	if *release == true {
		fmt.Println("개봉 여부:개봉")
	} else {
		fmt.Println("개봉 여부:미개봉")
	}

}
