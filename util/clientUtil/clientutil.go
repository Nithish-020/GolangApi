package clientUtil

import (
	"fcs23pkg/common"

	"fcs23pkg/ftdb"
	"log"
)

//------------------------------------------------------------
// function returns the mailid or mobilenumber
//Also return the encrypted values of (mailid or mobilenumber)
//------------------------------------------------------------
func GetEncryptedEmailPhone(clientId string, Type string, clientData string) (string, string, error) {
	log.Println("GetEncryptedEmailPhone(+)")
	var encryptedvalue string
	var err error
	//	clientData = Email
	log.Println(Type)
	log.Println(clientData)
	if clientData == "" {
		log.Println("insideif")
		// db, err := common.Localdbconnect(common.TechXL)
		// if err != nil {
		// 	common.LogDebug("DEBUG:GetEncryptedEmailPhone ", "(0MY099)", err.Error())
		// 	return clientData, encryptedvalue, err
		// } else {
		// defer db.Close()

		if Type == "EMAIL" {
			clientData, err = GetClientEmailId(clientId)
			if err != nil {
				common.LogDebug("clientUtil.GetEncryptedEmailPhone", clientId+":(CGEP01)", err.Error())
				return clientData, encryptedvalue, err
			}
		} else if Type == "MOBILE" {
			clientData, err = GetClientMobileNo(clientId)
			if err != nil {
				common.LogDebug("clientUtil.GetEncryptedEmailPhone", clientId+":(CGEP02)", err.Error())
				return clientData, encryptedvalue, err
			}
		}

		//}
	}

	if Type == "EMAIL" {
		// Get the encryptedvalue of given mailid
		encryptedvalue, err = common.GetEncryptedemail(clientData)
		log.Println("encryptedvalue:", encryptedvalue)
		if err != nil {
			common.LogDebug("clientUtil.GetEncryptedEmailPhone", clientId+":(CGEP03)", err.Error())
			return clientData, encryptedvalue, err
		}
	} else if Type == "MOBILE" {

		// Get the encryptedvalue of given mobilenumber
		encryptedvalue, err = common.GetEncryptedMobile(clientData)
		log.Println("encryptedvalue:", encryptedvalue)
		if err != nil {
			common.LogDebug("clientUtil.GetEncryptedEmailPhone", clientId+":(CGEP04)", err.Error())
			return clientData, encryptedvalue, err
		}
	}
	log.Println("clientData:", clientData)
	log.Println("GetEncryptedEmailPhone(-)")

	return clientData, encryptedvalue, nil
}

//-----------------------------------------------
// function fetch the company_code based on
// the given Clientid
//-----------------------------------------------
func GetClientSegment(ClientId string) (string, error) {
	log.Println("getClientSegment+")
	var Segment string
	var totalSegment string

	mssql, err := ftdb.LocalDbConnect(ftdb.ClientDB)
	if err != nil {
		common.LogDebug("clientUtil.GetClientSegment", ClientId+":(CGCS01)", err.Error())
		return Segment, err
	} else {
		defer mssql.Close()
		sqlString := `select distinct company_code from ` + common.TechExcelPrefix + `client_master
	where client_id = $1`
		rows, err := mssql.Query(sqlString, ClientId)
		if err != nil {
			common.LogDebug("clientUtil.GetClientSegment", ClientId+":(CGCS02)", err.Error())
			return Segment, err
		} else {
			for rows.Next() {
				err := rows.Scan(&Segment)
				if err != nil {
					common.LogDebug("clientUtil.GetClientSegment", ClientId+":(CGCS03)", err.Error())
					return Segment, err
				} else {
					Segment = Segment + " / "
					totalSegment = totalSegment + Segment
				}
			}
			totalSegment = totalSegment[0 : len(totalSegment)-3]

		}
	}
	log.Println("getClientSegment-")
	return totalSegment, nil

}

//------------------------------------------------
// Function returns the ClientName based on the
// given ClientId .
//------------------------------------------------
func GetClientName(ClientId string) (string, error) {
	log.Println("GetClientName (+)")
	var ClientName string
	//open a db connection
	db, err := ftdb.LocalDbConnect(ftdb.ClientDB)
	//if any error when opening db connection
	if err != nil {
		common.LogDebug("clientUtil.GetClientName", ClientId+":(CGCN01)", err.Error())
		return ClientName, err
	} else {
		defer db.Close()

		sqlString := `select distinct client_name from ` + common.TechExcelPrefix + `client_master
	where client_id = $1`
		// log.Println("sqlString := ", sqlString)
		rows, err := db.Query(sqlString, ClientId)
		if err != nil {
			common.LogDebug("clientUtil.GetClientName", ClientId+":(CGCN02)", err.Error())
			return ClientName, err
		} else {
			for rows.Next() {
				err := rows.Scan(&ClientName)
				if err != nil {
					common.LogDebug("clientUtil.GetClientName", ClientId+":(CGCN03)", err.Error())
					return ClientName, err
				}
			}
		}
	}
	//log.Println(filePath)
	log.Println("GetClientName (-)")
	return ClientName, nil
}

//-----------------------------------------------------
// function to get the Clientlocation for the given Clientid
//-----------------------------------------------------
func GetClientLocation(ClientId string) (string, error) {
	log.Println("getClientLocation+")
	var Location string
	var State string
	var Country string
	//open a db connection
	mssql, err := ftdb.LocalDbConnect(ftdb.ClientDB)
	//if any error when opening db connection
	if err != nil {
		common.LogDebug("clientUtil.GetClientLocation", ClientId+":(CGCL01)", err.Error())
		return Location, err
	} else {
		defer mssql.Close()
		Location = State + "," + Country
		sqlString := `select state,country from ` + common.TechExcelPrefix + `client_details
	where client_id = $1`
		//log.Println("sqlString", sqlString)
		rows, err := mssql.Query(sqlString, ClientId)
		if err != nil {
			common.LogDebug("clientUtil.GetClientLocation", ClientId+":(CGCL02)", err.Error())
			return Location, err
		} else {
			for rows.Next() {
				err := rows.Scan(&State, &Country)
				Location = State + "," + Country
				if err != nil {
					common.LogDebug("clientUtil.GetClientLocation", ClientId+":(CGCL03)", err.Error())
					return Location, err
				}
			}
		}

	}
	log.Println("getClientLocation-")
	return Location, nil
}

type ClientDetails struct {
	ClientName   string `json:"clientName"`
	Mobilenumber string `json:"mobilenumber"`
	EmailId      string `json:"emailId"`
}

//-------------------------------------------------------
//function gets client details for the given client id
//-------------------------------------------------------
func GetClientDetails(clientId string) (ClientDetails, error) {
	log.Println("GetClientDetails(+)")
	// log.Println(ID)
	var msg ClientDetails
	//open a db connection
	db, err := ftdb.LocalDbConnect(ftdb.ClientDB)
	//if any error when opening db connection
	if err != nil {
		common.LogDebug("clientUtil.GetClientDetails ", clientId+":(CGCD01)", err.Error())
		return msg, err
	} else {
		defer db.Close()
		sqlString := `select top 1 client_name,MOBILE_NO,EMAIL_ID from ` + common.TechExcelPrefix + `client_detail 
		where client_id = $1`
		rows, err := db.Query(sqlString, clientId)
		if err != nil {
			common.LogDebug("clientUtil.GetClientDetails ", clientId+":(CGCD02)", err.Error())
			return msg, err
		} else {
			for rows.Next() {
				err := rows.Scan(&msg.ClientName, &msg.Mobilenumber, &msg.EmailId)
				if err != nil {
					common.LogDebug("clientUtil.GetClientDetails ", clientId+":(CGCD03)", err.Error())
					return msg, err
				}
			}
		}
	}
	log.Println("GetClientDetails(-)")
	return msg, nil
}

//-----------------------------------------------------
// function to get the Clientlocation for the given Clientid
//-----------------------------------------------------
func GetClientEmailId(ClientId string) (string, error) {
	log.Println("GetClientEmailId+")
	var EmailId string
	//open a db connection
	mssql, err := ftdb.LocalDbConnect(ftdb.ClientDB)
	//if any error when opening db connection
	if err != nil {
		common.LogDebug("clientUtil.GetClientEmailId ", ClientId+":(CGCE01)", err.Error())
		return EmailId, err
	} else {
		defer mssql.Close()
		sqlString := `select EMAIL_ID from ` + common.TechExcelPrefix + `client_details
	where client_id = $1`
		log.Println("sqlString", sqlString)
		rows, err := mssql.Query(sqlString, ClientId)
		if err != nil {
			common.LogDebug("clientUtil.GetClientEmailId ", ClientId+":(CGCE02)", err.Error())
			return EmailId, err
		} else {
			for rows.Next() {
				err := rows.Scan(&EmailId)

				if err != nil {
					common.LogDebug("clientUtil.GetClientEmailId ", ClientId+":(CGCE03)", err.Error())
					return EmailId, err
				}
			}
		}

	}
	log.Println("GetClientEmailId-")
	return EmailId, nil
}

//-----------------------------------------------------
// function to get the Clientlocation for the given Clientid
//-----------------------------------------------------
func GetClientMobileNo(ClientId string) (string, error) {
	log.Println("GetClientMobileNo+")
	var MobileNumber string
	//open a db connection
	mssql, err := ftdb.LocalDbConnect(ftdb.ClientDB)
	//if any error when opening db connection
	if err != nil {
		common.LogDebug("clientUtil.GetClientMobileNo ", ClientId+":(CGCM01)", err.Error())
		return MobileNumber, err
	} else {
		defer mssql.Close()

		sqlString := `select MOBILE_NO from ` + common.TechExcelPrefix + `client_details
	where client_id = $1`
		log.Println("sqlString", sqlString)
		rows, err := mssql.Query(sqlString, ClientId)
		if err != nil {
			common.LogDebug("clientUtil.GetClientMobileNo ", ClientId+":(CGCM02)", err.Error())
			return MobileNumber, err
		} else {
			for rows.Next() {
				err := rows.Scan(&MobileNumber)

				if err != nil {
					common.LogDebug("clientUtil.GetClientMobileNo ", ClientId+":(CGCM03)", err.Error())
					return MobileNumber, err
				}
			}
		}

	}
	log.Println("GetClientMobileNo-")
	return MobileNumber, nil
}
