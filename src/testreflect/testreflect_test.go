package testreflect

import (
	"container/list"
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
)

type ctType struct {
	Tname string
}

type Item struct {
	Name string `user name`
	Id   int    `user id`
	Ctt  ctType
}

func (i Item) SetCt(ct ctType) Item {
	i.Ctt = ct

	return i
}

func (i Item) Show() {
	fmt.Printf("Name:%v , Id:%v\n", i.Name, i.Id)
}

func (i Item) ShowId() {
	fmt.Println("id:", i.Id)
}

func (i Item) TLog(n int) string {
	return fmt.Sprintf("TLogValue : %v", n)
}

// 기본적으로 reflect에서 private함수는 지원하지 않는다. but github에 private method를 지원하는 reflect패키지가 있당!!!!
func (i Item) privateFunc() {
	fmt.Println("Item.privateFunc")
}

//본 함수를 reflect로 읽기 위해서는 reflect.TypeOf(&item) 주소값을 넘겨 주어야 한다.
func (i *Item) SetID(id int) {
	i.Id = id
}

func Test_Reflect(t *testing.T) {
	fmt.Println("************* Test_Reflect ****************")
	run := func() { fmt.Println("run") }
	defer func() {
		run()
	}()
	_ = func() {
		val1 := 69

		fmt.Printf("TypeOf:%v , ValueOf:%v %n\n", reflect.TypeOf(val1), reflect.ValueOf(val1), &val1)

		rv := reflect.ValueOf(val1)
		typeStr := rv.Type().Kind().String()
		if typeStr != "ptr" {
			fmt.Println("this reflect value is not ptr.", "this type is ", typeStr)
		}

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

		//return
	}
	_ = func() {
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
	run = func() {
		type Temp struct {
			Name string
			Id   int
			age  int
		}
		u := Temp{
			Name: "abc",
			Id:   18,
			age:  20,
		}

		var uu interface{}
		uu = &u

		uValue := reflect.ValueOf(uu)
		uType := uValue.Type()

		if uType.Kind() == reflect.Ptr {
			fmt.Println("Pointer????")
			t := uType.Elem()

			if t.Kind() == reflect.Struct {
				fmt.Println("struct")
				uType = t
			} else {
				fmt.Println("struct not ", t.Kind())

				return
			}

		}

		fmt.Printf("[%v] \n", uType)
		for i := 0; i < uType.NumField(); i++ {
			sf := uType.Field(i)
			fmt.Printf("%v : %v - %v \n", sf.Type, sf.Name, sf.Tag)

			v := []byte(sf.Name)
			if v[0] >= 65 && v[0] <= 90 {
				fmt.Println("True")
			} else {
				fmt.Println("False")
			}

		}
		/*
			if fName, ok := uType.FieldByName("Name"); ok {
				fmt.Println(fName.Type, ", ", fName.Name, ", ", fName.Tag)
			}
		*/

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

func Test_FuncCallA(t *testing.T) {
	item := Item{Name: "kjs", Id: 22}

	var stType = reflect.TypeOf((*Item)(nil)).Elem()
	//var stValue = reflect.ValueOf((*Item)(nil))

	fmt.Println("stType:", stType)

	iType := reflect.TypeOf(&item) // &타입으로 보내면 *형 함수도 처리할 수 있다.
	fmt.Println("iType:", iType)

	iValue := reflect.ValueOf(&item)
	fmt.Println("iValue:", iValue)

	reflectValue := iType

	for i := 0; i < reflectValue.NumMethod(); i++ {

		method := reflectValue.Method(i)

		if method.Name == "Show" {
			fmt.Println("name:", method.Name)
			fmt.Println("funcType", reflect.Indirect(method.Func))

			for j := 0; j < method.Type.NumIn(); j++ {
				fmt.Println("-", method.Type.In(j))
			}

			mFunc1 := method.Func                                      //Server
			mFunc2 := reflect.ValueOf(&item).MethodByName(method.Name) //

			fmt.Println(mFunc1, mFunc2)

			rFunc := mFunc1

			fmt.Println("Call by Method")
			retVals := rFunc.Call([]reflect.Value{reflect.ValueOf(&item)}) //case mFunc1
			//retVals := rFunc.Call(nil)									//case mFunc2
			for _, v := range retVals {
				fmt.Println(v)
			}
		} else {
			fmt.Println(method.Func)
		}
		fmt.Println("-------")
	}

}

func Test_FuncCall(t *testing.T) {

	fmt.Println("-----------------------------------------")

	{
		caption := "go is an open souce programming language"
		caption2 := 18

		titleFuncValue := reflect.ValueOf(TitleCase)

		fmt.Printf("%#+v\n", titleFuncValue)
		values := titleFuncValue.Call([]reflect.Value{reflect.ValueOf(caption), reflect.ValueOf(caption2)})

		fmt.Println(len(values))

		title := values[0].String()
		title2 := int(values[1].Int())
		fmt.Println(title, title2)
	}

	fmt.Println("===============================================")

	{
		var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
		_ = contextType

		var stType = reflect.TypeOf((*ctType)(nil)).Elem()
		_ = stType

		item := Item{Name: "kjs", Id: 22}

		itemType := reflect.TypeOf(&item)
		fmt.Println("리플렉션 String :", itemType.String())
		fmt.Println("리플렉션 Name :", itemType.Name())

		fmt.Println("함수 갯수 :", itemType.NumMethod())
		fmt.Println(">\n")

		for i := 0; i < itemType.NumMethod(); i++ {
			method := itemType.Method(i)
			fmt.Println("함수 이름 :", method.Name)
			fmt.Println("-----------------------------------------")

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
				if mType.In(i) == stType {
					fmt.Println("ctType struct !!")
				}
			}

			outTypes := []reflect.Type{}
			_ = outTypes
			fmt.Println("함수의 리턴 인자 갯수 :", mType.NumOut())
			for i := 0; i < mType.NumOut(); i++ {
				fmt.Printf("%d번째 인자 :%v\n", i+1, mType.Out(i))
				outTypes = append(outTypes, mType.Out(i))
			}

			if len(inTypes) > 0 {
				inParams := []reflect.Value{}
				for _, t := range inTypes {
					switch t.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						inParams = append(inParams, reflect.ValueOf(18))
					case reflect.Float32, reflect.Float64:
						inParams = append(inParams, reflect.ValueOf(18.5))
					case reflect.String:
						inParams = append(inParams, reflect.ValueOf("a"))
					case reflect.Struct:
						if t == stType {
							inParams = append(inParams, reflect.ValueOf(ctType{Tname: "a1818"}))
						} else {
							fmt.Println("What the fuck? struct!")
						}
					default:
						fmt.Println("???", t.Kind())
					}
				}

				fmt.Println("call ----inParams")
				outValues := mCallValue.Call(inParams)

				for _, v := range outValues {
					fmt.Println(v)
				}
				fmt.Println(item)

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

			fmt.Println(">>\n")
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

func Test_Loop(t *testing.T) {

	var loop int32
	loop = 1
	cnt := 0

	for i := atomic.LoadInt32(&loop); i == 1; {
		cnt++
		if cnt > 10 {
			break
		}
	}
	fmt.Println(cnt)

}
