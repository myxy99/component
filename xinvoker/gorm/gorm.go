package xgorm

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xlog"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (i *dbInvoker) newDatabaseClient(o *options) (db *gorm.DB) {
	var err error
	db, err = gorm.Open(o.getDSN(), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   o.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		xlog.Panic("NewDatabaseClient OpenDB", xlog.FieldErr(err))
	}
	if o.Debug {
		db = db.Debug()
	}
	d, err := db.DB()
	if err != nil {
		xlog.Panic("Application Starting",
			xlog.FieldComponentName("XInvoker"),
			xlog.FieldMethod("XInvoker.XGorm.NewDatabaseClient"),
			xlog.FieldDescription("NewDatabaseClient db.DB() Error"),
			xlog.FieldErr(err),
		)
	}
	d.SetMaxOpenConns(o.MaxOpenConnections)
	d.SetMaxIdleConns(o.MaxIdleConn)
	d.SetConnMaxLifetime(o.MaxConnectionLifeTime)
	//d.SetConnMaxIdleTime(o.MaxConnMaxIdleTime)
	return db
}

func (i *dbInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)

	prefix := i.key
	for name := range xcfg.GetStringMap(prefix) {
		cfg := xcfg.UnmarshalWithExpect(prefix+"."+name, newDatabaseOptions()).(*options)
		conf[name] = cfg
	}
	return conf
}
