package testsort

import (
	"fmt"
	"testing"
)

func TestSortSlice(t *testing.T) {

	sort_data := []int{
		1,
		20,
		40,
		10,
		50,
	}
	sort_data = append(sort_data, 99)

	CustomSort(func() int {
		return len(sort_data)
	},
		func(x, y int) bool {
			return sort_data[x] < sort_data[y]
		},
		func(x, y int) {
			sort_data[x], sort_data[y] = sort_data[y], sort_data[x]
		})

	fmt.Printf("%+v\n", sort_data)
}

type Student struct {
	name  string
	id    uint32
	score float64
}
type Students []Student

func TestSortStruct(t *testing.T) {

	sortlist := Students{
		{"가나다", 1, 12.1},
		{"조대리", 2, 11.1},
		{"박과장", 3, 21.1},
	}

	parm1 := Student{"지화자", 10, 99.1}
	sortlist = append(sortlist, parm1)

	CustomSort(func() int {
		return len(sortlist)
	},
		func(x, y int) bool {
			//return sortlist[x].name < sortlist[y].name	//이름으로 소팅(오름차순)
			//return sortlist[x].name > sortlist[y].name //이름으로 소팅(내림차순)
			return sortlist[x].score < sortlist[y].score //점수로 소팅(오름차순)
		},
		func(x, y int) {
			sortlist[x], sortlist[y] = sortlist[y], sortlist[x]
		})

	fmt.Printf("%+v\n", sortlist)

}
