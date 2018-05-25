package main

import (
	"fmt"
	"sort"
)

func main() {

	//sort_by_interface()

	sort_by_key()
}

//인터페이스를 이용한 정렬
func sort_by_interface() {
	s := Students{
		{"Maria", 89.3},
		{"AceChan", 69.3},
		{"CarmoleKim", 99.3},
	}

	sort.Sort(s)
	fmt.Println(s)

	//점수를 기준으로 내림차순 정렬
	sort.Sort(sort.Reverse(ByScore{s}))

	fmt.Println(s)
}

type Student struct {
	name  string
	score float32
}
type Students []Student

func (s Students) Len() int {
	return len(s)
}
func (s Students) Less(i, j int) bool {
	return s[i].name < s[j].name
}
func (s Students) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ByScore struct {
	Students
}

//Students를 임베딩(is-a)시켰wlaks Less함수를 재정의 해서 스코어로 정렬이 되도록 구조가 바뀜.
func (s ByScore) Less(i, j int) bool {
	return s.Students[i].score < s.Students[j].score
}

////////////////////////////////////////////////////

//정렬 키를 이용한 정렬 처리.
func sort_by_key() {
	s := Students{
		{"Maria", 89.3},
		{"AceChan", 69.3},
		{"CarmoleKim", 99.3},
	}

	//ss := []Student(s)
	ss := interface{}(s)
	fmt.Println(ss)

	/*
		nameby := func(p1, p2 *Student) bool {
			return p1.name < p2.name
		}
		By(nameby).Sort(s)
		//*/

	Compair(func(p1, p2 interface{}) bool {
		x := p1.(Student)
		y := p2.(Student)
		return x.name < y.name
	}).Sort(s, len(s))

	fmt.Println(s)
}

type SortDatas interface{} //정렬하고자 하는 데이터

type Compair func(s1, s2 interface{}) bool //각 상황별 정렬 함수를 저장할 타입

type DataSorter struct {
	sortDatas interface{}
	size      int
	compair   Compair //func(s1, s2 *Student) bool
}

func (compair Compair) Sort(sortdatas interface{}, size int) {
	sorter := &DataSorter{sortdatas, size, compair}
	sort.Sort(sorter)
}

func (s *DataSorter) Len() int {
	return s.size //len(s.sortDatas)
}
func (s *DataSorter) Less(i, j int) bool {
	return s.compair(&s.sortDatas[i], &s.sortDatas[j])
}
func (s *DataSorter) Swap(i, j int) {
	td := s.sortDatas.(Students)
	s.sortDatas[i], s.sortDatas[j] = s.sortDatas[j], s.sortDatas[i]
}
