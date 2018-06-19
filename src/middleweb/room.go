package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/mholt/binding" //HTTP요청 내용을 구조체로 변환하기 위해서 binding패키지 사용.
	"gopkg.in/mgo.v2/bson"
)

type Room struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson"name" json:"name"`
}

const (
	DB_NAME    = "test"
	C_ROOMS    = "rooms"
	C_MESSAGES = "messages"
)

// binding 패키지에서 Request 데이터를 Room 구조체로 변환하려면 Room 타입이 binding.FieldMapper 인터페이스여야 한다.
func (r *Room) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&r.Name: "name"}
}

func createRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//binding 패키지로 room 생성 요청 정보를 Room 타입 값으로 변환
	r := new(Room)
	errs := binding.Bind(req, r)
	//if errs.Handle(w) {
	if errs != nil {
		return
	}

	//mongoDB 세션 생성
	session := mongoSession.Copy()
	defer session.Close()

	//mongoDB ID 생성
	r.ID = bson.NewObjectId()
	//room 정보 저장을 위한 mongoDB 컬렉션 객체 생성
	c := session.DB(DB_NAME).C(C_ROOMS)

	//rooms 컬렉션에 room 정보 저장
	if err := c.Insert(r); err != nil {
		//오류 발생시 500에러 반환
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	//처리 결과 반환
	renderer.JSON(w, http.StatusCreated, r)
}

func retrieveRooms(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	var rooms []Room
	//모든 room 정보 조회
	err := session.DB(DB_NAME).C(C_ROOMS).Find(nil).All(&rooms)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	//room 조회 결과 반환
	renderer.JSON(w, http.StatusOK, rooms)
}
