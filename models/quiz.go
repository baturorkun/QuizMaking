package models

type Quiz struct {
	Model
	Name 		string    `schema:"name" gorm:"type:varchar(50)"`
	Details		string 	  `schema:"details" gorm:"type:text"`
}


func (quiz Quiz) GetAll() (res []Quiz) {

	db.Find(&res)
	return

}