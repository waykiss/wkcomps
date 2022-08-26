package database

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type Database struct {
	db *gorm.DB
	tx *gorm.DB // guardar conexao quando for uma transacao
}

var dbs map[string]*Database

func init() {
	dbs = map[string]*Database{}
}

//NewConnection return a new database connection
func NewConnection(host, port, user, password, dbName string) (r *Database, err error) {
	if dbs[dbName] == nil {
		dsn := getStringConnection(host, port, user, password, dbName)
		db, err2 := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err2 != nil {
			err = errors.New(err2.Error())
			logrus.Error(fmt.Sprintf("Could not connect to the database: msg %v dsn: %s", err, dsn))
			return
		}
		r = &Database{}
		r.db = db
		dbs[dbName] = r
	}
	r = dbs[dbName]
	return
}

//getStringConnection retorna a string de conexao para o banco dado
func getStringConnection(host, port, user, password, dbName string) string {
	if strings.Contains(host, "/") {
		return fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true", user, password, host, dbName)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}

//createDatabase função para criar um database, utilize o banco padrão para conectar ao servidor de banco de dados
//e então executar o comando sql para criar outro banco
func createDatabase(host, port, user, password, dbName string, ssl bool) error {
	dsn := getStringConnection(host, port, user, password, "postgres")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\";", dbName))
	return nil
}

type ViolationUniqKey struct {
	msg   string
	Field string
	Value string
}

func (v ViolationUniqKey) Error() string {
	return fmt.Sprintf("Já existe um registro com o valor %s para o campo %s ", v.Value, v.Field)
}

//StartTransaction start a database transaction
func (d *Database) StartTransaction() (err error) {
	d.tx = d.db.Begin()
	err = d.tx.Error
	return err
}

func (d *Database) CommitTransaction() (err error) {
	if d.tx == nil {
		err = errors.New("tentando fazer commit em uma transacao mas a mesma nao foi iniciada")
		return
	}

	d.tx.Commit()
	err = d.tx.Error
	d.tx = nil
	return err
}

func (d *Database) RollbackTransaction() (err error) {
	if d.tx == nil {
		err = errors.New("tentando fazer rollback em uma transacao mas a mesma nao foi iniciada")
		return
	}
	d.tx.Rollback()
	err = d.tx.Error
	d.tx = nil
	return err
}
