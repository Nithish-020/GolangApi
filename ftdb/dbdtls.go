package ftdb

import (
	"fcs23pkg/genpkg"
	"fmt"
	"strconv"
)

const (
	MainDB     = "FTCDB"
	ClientDB   = "FTCDB"
	IPODB      = "IPODB"
	SSODB      = "SSODB"
	MariaFTPRD = "MARIAFTPRD"

	BrosePrefix     = "brose.dbo."
	TechExcelPrefix = "TECHEXCELPROD.CAPSFO.dbo."
)

// Initializing DB Details
func (d *AllUsedDatabases) Init() {
	dbconfig := genpkg.ReadTomlConfig("../dbconfig.toml")
	//Setting Techexcel UAT DB Connection Details
	d.TechExcelUAT.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATServer"])
	d.TechExcelUAT.Port, _ = strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATPort"]))
	d.TechExcelUAT.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATUser"])
	d.TechExcelUAT.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATPassword"])
	d.TechExcelUAT.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATDatabase"])
	d.TechExcelUAT.DBType = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATDBType"])
	d.TechExcelUAT.DB = MainDB
	//Setting Techexcel UAT DB Connection Details
	d.TechExcelUAT.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATServer"])
	d.TechExcelUAT.Port, _ = strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATPort"]))
	d.TechExcelUAT.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATUser"])
	d.TechExcelUAT.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATPassword"])
	d.TechExcelUAT.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATDatabase"])
	d.TechExcelUAT.DBType = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["TechExcelUATDBType"])
	d.TechExcelUAT.DB = ClientDB
	//setting SSO db connection details
	d.SSODB.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["SSODBServer"])
	d.SSODB.Port, _ = strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["SSODBPort"]))
	d.SSODB.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["SSODBUser"])
	d.SSODB.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["SSODBPassword"])
	d.SSODB.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["SSODBDatabase"])
	d.SSODB.DBType = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["SSODBDBType"])
	d.SSODB.DB = SSODB
	//setting IPO db connection details
	d.IPODB.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["IPODBServer"])
	d.IPODB.Port, _ = strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["IPODBPort"]))
	d.IPODB.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["IPODBUser"])
	d.IPODB.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["IPODBPassword"])
	d.IPODB.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["IPODBDatabase"])
	d.IPODB.DBType = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["IPODBDBType"])
	d.IPODB.DB = IPODB
	//setting Maria db connection details
	d.MariaDB.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaDBServer"])
	d.MariaDB.Port, _ = strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaDBPort"]))
	d.MariaDB.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaDBUser"])
	d.MariaDB.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaDBPassword"])
	d.MariaDB.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaDBDatabase"])
	d.MariaDB.DBType = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaDBDBType"])
	d.MariaDB.DB = MariaFTPRD
}
