package models

import (
	"log"
	"net/url"
	"reflect"
)

type User struct {
	Model
	Name 		string    `schema:"name" gorm:"type:varchar(50)"`
	Email		string    `schema:"email" gorm:"type:varchar(150);unique"`
	Password	string    `schema:"password" gorm:"type:varchar(15)"`
	RePassword  string    `schema:"repassword" gorm:"-"`
}


func (user *User) Add(v url.Values) {

	err := decoder.Decode(&user, v)

	if err != nil {
		log.Fatal("Decode error: %+v", err)
	}

	//log.Printf(">>> %+v", u)

	db.Create(&user)

}

func (user *User) Get(v interface{}) {

	//log.Printf(">>> %v", reflect.TypeOf(v).Kind())

	if reflect.TypeOf(v).Kind() ==  reflect.Int {
		db.Where("id = ?", v.(int)).First(&user)
	} else if reflect.TypeOf(v).Kind() ==  reflect.String {
		db.Where("email = ?", v.(string)).First(&user)
		log.Println("*", user)
	} else {
		log.Fatal("Invalid type")
	}

}