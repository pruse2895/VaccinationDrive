package daos

import (
	"github.com/FenixAra/go-util/log"

	"vaccinationDrive/models"

	"github.com/go-pg/pg"
)

type UserObj struct {
	l      *log.Logger
	dbConn *pg.DB
}

func NewUserData(l *log.Logger, dbConn *pg.DB) *UserObj {
	return &UserObj{
		l:      l,
		dbConn: dbConn,
	}
}

type UserDao interface {
	SaveUser(user models.User) error
}

func (u *UserObj) SaveUser(user models.User) error {

	if err := u.dbConn.Insert(&user); err != nil {
		u.l.Errorf("SaveUser Error ", err)
		return err
	}

	return nil

}
