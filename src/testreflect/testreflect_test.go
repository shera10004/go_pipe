package testreflect

import (
	"container/list"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type Item struct {
	Name string `user name`
	Id   int    `user id`
}

func (i Item) Show() {
	fmt.Printf("Name:%v , Id:%v\n", i.Name, i.Id)
}

func (i Item) TLog(n int) string {
	return fmt.Sprintf("TLogValue : %v", n)
}

func Test_Reflect(t *testing.T) {

	{
		val1 := 69

		fmt.Printf("TypeOf:%v , ValueOf:%v \n", reflect.TypeOf(val1), reflect.ValueOf(val1))

		rv := reflect.ValueOf(val1)
		if rv.CanSet() {
			rv.SetInt(18)
		} else {
			fmt.Println("rv is not CanSet")
			rp := reflect.ValueOf(&val1)
			p := rp.Elem()
			fmt.Println(p.CanSet())
			p.SetInt(88)
			fmt.Println("val1 :", val1)
		}
	}

	fmt.Println("------------------------------")

	{
		langs := []string{"a", "b", "c"}
		sv := reflect.ValueOf(langs)
		v := sv.Index(1)
		if v.CanSet() {

			v.SetString("x")
			fmt.Println("langs :", langs)
		}

		rp := reflect.ValueOf(&langs).Elem()
		rp.Index(1).SetString("y")
		fmt.Println("xlangs :", langs)

	}

	fmt.Println("------------------------------")

	{
		u := Item{
			Name: "abc",
			Id:   18,
		}

		uType := reflect.TypeOf(u)
		uValue := reflect.ValueOf(u)

		fmt.Printf("[%v] \n", uType)

		if fName, ok := uType.FieldByName("Name"); ok {
			fmt.Println(fName.Type, ", ", fName.Name, ", ", fName.Tag)
		}

		if uValue.CanSet() {

		} else {
			fmt.Println("uValue is not CanSet")
			rp := reflect.ValueOf(&u).Elem()
			rp.FieldByName("Name").SetString("eee")
			rp.FieldByName("Id").SetInt(99)
			fmt.Printf("u : %#+v", u)
		}
	}

}

func TitleCase(s string, i int) (string, int) {
	return strings.Title(s), i * 10
}

func Test_FuncCall(t *testing.T) {
	caption := "go is an open souce programming language"
	caption2 := 18

	title, title2 := TitleCase(caption, caption2)
	fmt.Println(title, title2)

	fmt.Println("-----------------------------------------")

	titleFuncValue := reflect.ValueOf(TitleCase)

	fmt.Printf("%#+v\n", titleFuncValue)
	values := titleFuncValue.Call([]reflect.Value{reflect.ValueOf(caption), reflect.ValueOf(caption2)})

	fmt.Println(len(values))

	title = values[0].String()
	title2 = int(values[1].Int())
	fmt.Println(title, title2)

	fmt.Println("-----------------------------------------")

	{
		item := Item{Name: "kjs", Id: 21}

		itemType := reflect.TypeOf(item.Show)
		itemFunc := reflect.ValueOf(item.Show)
		fmt.Println("함수의 입력 파라미터 갯수 :", itemType.NumIn())
		fmt.Println("함수의 리턴 인자 갯수 :", itemType.NumOut())
		_ = itemFunc.Call(nil)
	}

	fmt.Println("-----------------------------------------")

	{
		item := Item{Name: "kjs", Id: 22}

		itemType := reflect.TypeOf(item)
		fmt.Println("리플렉션 String :", itemType.String())
		fmt.Println("리플렉션 Name :", itemType.Name())

		fmt.Println("함수 갯수 :", itemType.NumMethod())
		fmt.Println(">")
		for i := 0; i < itemType.NumMethod(); i++ {
			method := itemType.Method(i)
			fmt.Println("함수 이름 :", method.Name)

			mType := method.Type
			mCallValue := reflect.ValueOf(&item).MethodByName(method.Name)
			fmt.Println("함수의 입력 파라미터 갯수 :", mType.NumIn())
			fmt.Println("함수의 Func :", mCallValue, method.Func)

			inTypes := []reflect.Type{}
			_ = inTypes
			for i := 0; i < mType.NumIn(); i++ {
				fmt.Printf("%d번째 인자 :%v\n", i+1, mType.In(i))
				if mType.In(i) == itemType {
					continue
				}
				inTypes = append(inTypes, mType.In(i))
			}

			outTypes := []reflect.Type{}
			_ = outTypes
			fmt.Println("함수의 리턴 인자 갯수 :", mType.NumOut())
			for i := 0; i < mType.NumOut(); i++ {
				fmt.Printf("%d번째 인자 :%v\n", i+1, mType.Out(i))
				outTypes = append(outTypes, mType.Out(i))
			}

			if len(inTypes) > 0 {
				outValues := []reflect.Value{}
				inParams := []reflect.Value{}
				for _, t := range inTypes {
					switch t.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						inParams = append(inParams, reflect.ValueOf(18))
					case reflect.Float32, reflect.Float64:
						inParams = append(inParams, reflect.ValueOf(18.5))
					case reflect.String:
						inParams = append(inParams, reflect.ValueOf("a"))
					default:
						fmt.Println("???", t.Kind())
					}
				}

				fmt.Println("call ----inParams")
				outValues = mCallValue.Call(inParams)

				for _, v := range outValues {
					fmt.Println(v)
				}

			} else {
				if len(outTypes) == 0 {
					fmt.Println("call ----nil")
					_ = mCallValue.Call(nil)
				} else {
					fmt.Println("call ----nil")
					outValues := mCallValue.Call(nil)
					for _, v := range outValues {
						fmt.Println(v)
					}
				}
			}

			fmt.Println(">>")
		}

		/*
			itemFunc := reflect.ValueOf(item.TLog)

			fmt.Println("함수의 입력 파라미터 갯수 :", itemType.NumIn())
			for i := 0; i < itemType.NumIn(); i++ {
				fmt.Printf("%d번째 인자 :%v\n", i+1, itemType.In(i))
			}

			fmt.Println("함수의 리턴 인자 갯수 :", itemType.NumOut())
			for i := 0; i < itemType.NumOut(); i++ {
				fmt.Printf("%d번째 인자 :%v\n", i+1, itemType.Out(i))
			}

			rValues := itemFunc.Call([]reflect.Value{reflect.ValueOf(18)})
			for _, v := range rValues {
				fmt.Println(v)
			}
			//*/
	}

}

func Len(x interface{}) int {
	value := reflect.ValueOf(x)
	switch reflect.TypeOf(x).Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return value.Len()

	default:
		if method := value.MethodByName("Len"); method.IsValid() {
			values := method.Call(nil)
			//if len(values) == 1 {
			{
				return int(values[0].Int())
			}
		}
	}
	return -1
}

type cc int

func (c cc) Len() int {
	return int(c)
}

func Test_Len(t *testing.T) {
	a := list.New()
	_ = a

	b := list.New()
	b.PushBack(11)

	c := map[string]int{"a": 1, "b": 2}
	_ = c

	d := "one"
	_ = d

	e := []int{3, 2, 1, 0}
	_ = e

	var f cc
	f = 18

	fmt.Println(Len(a), Len(b), Len(c), Len(d), Len(e), Len(f))

}
