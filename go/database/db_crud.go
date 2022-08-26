package database

import (
	"gorm.io/gorm"
)

func (d *Database) Create(obj interface{}) (err error) {
	result := d.db.Create(obj)
	return result.Error
}

func (d *Database) Update(obj interface{}) (affected int, err error) {
	result := d.db.Save(obj)
	return int(result.RowsAffected), result.Error
}

func (d *Database) Delete(obj interface{}) (affected int, err error) {
	result := d.db.Delete(obj)
	return int(result.RowsAffected), result.Error
}

func (d *Database) FindById(id string, obj interface{}) (err error) {
	//_ = d.Query.Eq("id", id).FindRaw(obj)
	return
}

func (d *Database) Find(dest interface{}, conditions ...interface{}) (err error) {
	result := d.db.Find(dest, conditions...)
	return result.Error
}

func (d *Database) GormInstance() *gorm.DB {
	return d.db
}

//ExecSQL funcao para executar insert, update, delete e atualizar schema do banco, essa funcao NÃO FAZ QUERY
func (d *Database) ExecSQL(query string, args ...interface{}) (result *gorm.DB, err error) {
	if d.tx != nil {
		result = d.tx.Exec(query, args...)
	} else {
		result = d.db.Exec(query, args...)
	}
	err = result.Error
	return
}

//Select funcao para executar uma consulta(query), o parametro `dest` é usado para gravar o resultado da query
func (d *Database) Select(dest interface{}, query string, args ...interface{}) (err error) {
	var result *gorm.DB
	if d.tx != nil {
		result = d.tx.Raw(query, args...).Scan(dest)
	} else {
		result = d.db.Raw(query, args...).Scan(dest)
	}

	err = result.Error
	return
}

func (d *Database) Migrate(dst ...interface{}) (err error) {
	result := d.db.AutoMigrate(dst...)
	return result
}
