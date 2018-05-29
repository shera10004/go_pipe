package testsort

/*
CustomSort(func() int {
		return len(s)
	},
		func(x, y int) bool {
			return s[x].name < s[y].name
		},
		func(x, y int) {
			s[x], s[y] = s[y], s[x]
		},
	)
*/

import (
	"sort"
)

type fSize func() int
type fCompair func(s1, s2 int) bool //각 상황별 정렬 함수를 저장할 타입
type fChange func(s1, s2 int)

type stDataSorter struct {
	size    fSize
	compair fCompair //func(s1, s2 *Student) bool
	change  fChange
}

func CustomSort(l fSize, c fCompair, s fChange) {
	sorter := &stDataSorter{l, c, s}
	sort.Sort(sorter)
}

func (s *stDataSorter) Len() int {
	return s.size()
}
func (s *stDataSorter) Less(i, j int) bool {
	return s.compair(i, j)
}
func (s *stDataSorter) Swap(i, j int) {
	s.change(i, j)
}
