package ftdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

// Structure to hold database connection details
type DatabaseType struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
	DBType   string
	DB       string
}

// structure to hold all db connection details used in this program
type AllUsedDatabases struct {
	TechExcelUAT  DatabaseType
	Kyc           DatabaseType
	MariaDB       DatabaseType
	KamabalaDB    DatabaseType
	KamabalaApiDB DatabaseType
	MariaEkyc     DatabaseType
	SSODB         DatabaseType
	IPODB         DatabaseType
}

// ---------------------------------------------------------------------------------
// function opens the db connection and return connection variable
// ---------------------------------------------------------------------------------
func LocalDbConnect(DBtype string) (*sql.DB, error) {
	DbDetails := new(AllUsedDatabases)
	DbDetails.Init()

	connString := ""
	localDBtype := ""

	var db *sql.DB
	var err error
	var dataBaseConnection DatabaseType
	// get connection details
	if DBtype == DbDetails.TechExcelUAT.DB {
		dataBaseConnection = DbDetails.TechExcelUAT
		localDBtype = DbDetails.TechExcelUAT.DBType
	} else if DBtype == DbDetails.KamabalaDB.DB {
		dataBaseConnection = DbDetails.KamabalaDB
		localDBtype = DbDetails.KamabalaDB.DBType
	} else if DBtype == DbDetails.MariaDB.DB {
		dataBaseConnection = DbDetails.MariaDB
		localDBtype = DbDetails.MariaDB.DBType
	} else if DBtype == DbDetails.Kyc.DB {
		dataBaseConnection = DbDetails.Kyc
		localDBtype = DbDetails.Kyc.DBType
	} else if DBtype == DbDetails.KamabalaApiDB.DB {
		dataBaseConnection = DbDetails.KamabalaApiDB
		localDBtype = DbDetails.KamabalaApiDB.DBType
	} else if DBtype == DbDetails.MariaEkyc.DB {
		dataBaseConnection = DbDetails.MariaEkyc
		localDBtype = DbDetails.MariaEkyc.DBType
	} else if DBtype == DbDetails.SSODB.DB {
		dataBaseConnection = DbDetails.SSODB
		localDBtype = DbDetails.SSODB.DBType
	} else if DBtype == DbDetails.IPODB.DB {
		dataBaseConnection = DbDetails.IPODB
		localDBtype = DbDetails.SSODB.DBType
	}
	// log.Println("localDBtype", localDBtype)
	// log.Println("dataBaseConnection", dataBaseConnection)

	// Prepare connection string
	if localDBtype == "mssql" {
		connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", dataBaseConnection.Server, dataBaseConnection.User, dataBaseConnection.Password, dataBaseConnection.Port, dataBaseConnection.Database)
	} else if localDBtype == "mysql" {
		connString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dataBaseConnection.User, dataBaseConnection.Password, dataBaseConnection.Server, dataBaseConnection.Port, dataBaseConnection.Database)
	}

	//make a connection to db
	if localDBtype != "" {
		db, err = sql.Open(localDBtype, connString)
		//db, err := util.Getdb(localDBtype, connString)
		if err != nil {
			log.Println("Open connection failed:", err.Error())
		} else {
			if DBtype == SSODB || DBtype == IPODB {
				db.SetConnMaxIdleTime(time.Second * 5)
			}
		}
	} else {
		return db, fmt.Errorf("Invalid DB Details")
	}

	return db, err
}

// --------------------------------------------------------------------
//
//	execute bulk inserts
//
// --------------------------------------------------------------------
func ExecuteBulkStatement(db *sql.DB, sqlStringValues string, sqlString string) error {
	log.Println("ExecuteBulkStatement+")
	//trim the last ,
	sqlStringValues = sqlStringValues[0 : len(sqlStringValues)-1]
	_, err := db.Exec(sqlString + sqlStringValues)
	if err != nil {
		log.Println(err)
		log.Println("ExecuteBulkStatement-")
		return err
	} else {
		log.Println("inserted Sucessfully")
	}
	log.Println("ExecuteBulkStatement-")
	return nil
}
