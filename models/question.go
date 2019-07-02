package models

type Question struct {
	Model
	Quiz		int64			`gorm:"foreignkey:ID"`
	Content 	string			`gorm:"type:text"`
	Answer 		string			`gorm:"type:char(1)"`
	Options 	[]Option		`gorm:"-"`
}

type Option struct {
	Model
	Question		int64		`gorm:"foreignkey:ID"`
	Content 		string		`gorm:"type:text"`
	Letter 			string		`gorm:"type:char(1)"`
}


func (question Question) Get() (res []Question) {

	db.Where("quiz = ?", question.Quiz).Find(&res)

	for i, q := range res {

		db.Model(Option{}).Where("question = ?", q.ID).Find(&res[i].Options)

	}

	return
}
