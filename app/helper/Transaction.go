package helper

import "gorm.io/gorm"

func CommitOrRollback(db *gorm.DB) {
	err := recover()

	if err != nil {
		db.Rollback()
		panic(err)
	} else {
		db.Commit()
	}
}
