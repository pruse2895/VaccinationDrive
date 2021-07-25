package userRegistration

import (
	"vaccinationDrive/internals/daos"
	"vaccinationDrive/models"

	"github.com/FenixAra/go-util/log"
	"github.com/go-pg/pg"
)

type UserData struct {
	dbConn  *pg.DB
	l       *log.Logger
	userDao daos.UserDao
}

func NewUserData(l *log.Logger, dbConn *pg.DB) *UserData {
	return &UserData{
		l:       l,
		dbConn:  dbConn,
		userDao: daos.NewUserData(l, dbConn),
	}

}

func (u *UserData) RegisterUser(user models.User) error {

	user.BeforeInsert()
	err := u.userDao.SaveUser(user)
	if err != nil {
		u.l.Errorf("RegisterUser Error -- ", err)
		return err
	}

	return nil

}
