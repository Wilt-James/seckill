package repositories

import (
	"database/sql"
	"errors"
	"seckill/common"
	"seckill/datamodels"
)

type IUser interface {
	Conn() error
	Insert(*datamodels.User)(int64, error)
	Select(string)(*datamodels.User, error)
}

type UserManager struct {
	table string
	mysqlConn *sql.DB
}

func NewUserManager(table string, db *sql.DB) IUser {
	return &UserManager{table:table, mysqlConn:db}
}

// 数据连接
func (u *UserManager) Conn()(err error) {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "users"
	}
	return
}

// 插入
func (u *UserManager) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "INSERT " + u.table + " SET nickName = ?, userName = ?, passWord = ?"
	stmt, err := u.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}


func (u *UserManager) Select(userName string) (userResult *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空！")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "SELECT * FROM " + u.table + " WHERE userName = ?"
	rows, err := u.mysqlConn.Query(sql, userName)
	defer rows.Close()
	if err != nil {
		return &datamodels.User{}, err
	}

	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, nil
	}

	userResult = &datamodels.User{}
	common.DataToStructByTagSql(result, userResult)
	return
}
