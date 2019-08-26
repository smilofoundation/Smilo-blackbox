package types

var DBI DatabaseInstance

type DatabaseInstance interface {
	Close() error
	Delete(data interface{}) error
	Find(fieldname string, value interface{}, to interface{}) error
	Save(data interface{}) error
}
