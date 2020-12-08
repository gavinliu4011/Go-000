package main

import (
	"database/sql"
	"errors"
	"fmt"

	pkgerr "github.com/pkg/errors"
)

type User struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"type:varchar(36);not null;"`
}

var (
	ErrDaoRecordNotFound = pkgerr.New("record not found")
)

func MockSqlErrNoRows() error {
	return sql.ErrNoRows
}

type Dao struct {
	db *sql.DB
}

func NewDao(db *sql.DB) *Dao {
	return &Dao{
		db: db,
	}
}

func (d *Dao) GetUserByID(id int) (*User, error) {
	user := &User{}
	// err := d.db.QueryRow("SELECT id, name FROM users WHERE id = ?", id).Scan(user)
	err := MockSqlErrNoRows()

	if errors.Is(err, sql.ErrNoRows) {
		// 这里看自己团队设计而定，在这里我先返回pkgerr
		err = ErrDaoRecordNotFound
	}
	return user, pkgerr.Wrapf(err, "query user %d failed", id)
}

type Service struct {
	dao *Dao
}

func NewService(d *Dao) *Service {
	return &Service{
		dao: d,
	}
}

func (s *Service) GetUserByID(id int) (*User, error) {
	user, err := s.dao.GetUserByID(id)
	if errors.Is(pkgerr.Cause(err), ErrDaoRecordNotFound) {
		// 这里可以处理没有查询到记录时的业务逻辑
		// err = nil
		return user, err
	}
	return user, err
}

func main() {
	db := &sql.DB{} // 假装这里成功初始化了一个db
	service := NewService(NewDao(db))
	_, err := service.GetUserByID(1)
	if err != nil {
		fmt.Printf("err: %+v, stack:%+v", pkgerr.Cause(err), pkgerr.WithStack(err))
	}

}
