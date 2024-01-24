package coresettings

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"log"
)

// --------------------------------------------------------------------
// function to get value from core setting for a given Key
// --------------------------------------------------------------------
func GetCoreSettingValue(dbName string, key string) string {
	method := "coresettings.GetCoreSettingValue"
	var value string
	db, err := ftdb.LocalDbConnect(dbName)
	if err != nil {
		return value
	} else {
		defer db.Close()
		sqlString := "select valuev from CoreSettings where keyv ='" + key + "'"
		rows, err := db.Query(sqlString)
		if err != nil {
			log.Println(common.NoPanic, method, err.Error())
		}
		for rows.Next() {
			err := rows.Scan(&value)
			if err != nil {
				log.Println(common.NoPanic, method, err.Error())
			}
		}
		return value
	}

}
