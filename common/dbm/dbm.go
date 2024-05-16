package dbm

import (
	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/influxdbm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8s.io/klog/v2"
	"time"
)

// GormInit init the database according to the config file.
func GormInit(config *config.DBConfig) *gorm.DB {
	switch config.Used {
	case "mysql":
		return GormMysql(config)
	case "postgre":
		return GormPostgreSQL(config)
	case "sqlite":
		return GormSQLite(config)
	default:
		return GormPostgreSQL(config)
	}
}

// GormMysql connect the mysql
func GormMysql(config *config.DBConfig) *gorm.DB {
	m := config.Mysql
	if m == nil || m.Dbname == "" {
		klog.Fatalf("Can't find Mysql database config information")
		return nil
	}

	klog.Infof("Connecting to mysql %s", m.Host)

	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // the default length of string type
		DisableDatetimePrecision:  true,    // Disable datetime precision, which is not supported in databases before MySQL 5.6
		DontSupportRenameIndex:    true,    // When renaming an index, delete and create a new one. Databases before MySQL 5.7 and MariaDB do not support renaming an index
		DontSupportRenameColumn:   true,    // Use 'change' to rename the column. Databases before MySQL 8 and MariaDB do not support renaming the column
		SkipInitializeWithVersion: false,   // Auto configuration based on version
	}
	ormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}

	db, err := gorm.Open(mysql.New(mysqlConfig), ormConfig)
	if err != nil {
		klog.Fatalf("Connect to mysql failed! %s", err.Error())
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	return db
}

// connect the postgresql
func GormPostgreSQL(config *config.DBConfig) *gorm.DB {
	m := config.Postgres
	if m == nil || m.Dbname == "" {
		klog.Fatalf("Can't find Mysql database config information")
		return nil
	}

	klog.Infof("Connecting to postgres %s, dsn= %s", m.Host, m.Dsn())

	db, err := gorm.Open(postgres.Open(m.Dsn()), &gorm.Config{})
	if err != nil {
		klog.Errorf("connect db failed %s", err.Error())
		return nil
	}

	//db.Debug()

	//create connection pool
	pdb, err := db.DB()
	pdb.SetMaxIdleConns(m.MaxIdleConns) //Set the maximum number of connections in the free connection pool
	//cpu cores * 2 + number of disk.
	pdb.SetMaxOpenConns(m.MaxOpenConns) //Set the maximum number of open database connections.
	pdb.SetConnMaxLifetime(time.Hour)
	return db
}

// Connect the sqlite3
func GormSQLite(config *config.DBConfig) *gorm.DB {
	m := config.SQLite3
	if m == nil || m.DbPath == "" {
		klog.Fatalf("Can't find Mysql database config information")
		return nil
	}

	klog.Infof("Connecting to sqlite3 %s", m.DbPath)

	db, err := gorm.Open(sqlite.Open(m.DbPath), &gorm.Config{})
	if err != nil {
		klog.Fatalf("Connect to sqlite3 failed! %s", err.Error())
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db
}

func init() {
	conf := config.GetDBConfig()
	if conf == nil {
		panic("DBConfig is missing")
	}

	global.DBAccess = GormInit(conf)
	if global.DBAccess == nil {
		panic("Oops, gorm init failed!")
	}

	//register tabels.
	err := model.RegisterTables(global.DBAccess)
	if err != nil {
		panic(err)
	}
	//init influxdb
	influxdbm.InitInfluxDb()

}
