package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

const messageFetchSize = 10

type Message struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	RoomId   bson.ObjectId `bson:"room_id" json:"room_id"`
	Content  string        `bson:"content" json:"content"`
	CreateAt time.Time     `bson:"create_at" json:"create_at"`
	User     *User         `bson:"user" json:"user"`
}

func (m *Message) create() error {
	session := mongoSession.Copy()
	defer session.Close()

	m.ID = bson.NewObjectId()
	m.CreateAt = time.Now()

	c := session.DB(DB_NAME).C(C_MESSAGES)
	if err := c.Insert(m); err != nil {
		return err
	}
	return nil
}

func retrieveMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	//쿼리 매개변수로 전달된 limit값 확인
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		//정상적인 limit값이 전달되지 않으면 limit를 messageFetchSize으로 셋팅
		limit = messageFetchSize
	}

	var message []Message
	//_id 역순으로 정렬하여 limit 수만큼 message조회
	err = session.DB(DB_NAME).C(C_MESSAGES).
		Find(bson.M{"room_id": bson.ObjectIdHex(ps.ByName("id"))}).
		Sort("-_id").Limit(limit).All(&message)

	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	//message 조회 결과 반환
	renderer.JSON(w, http.StatusOK, message)
}
