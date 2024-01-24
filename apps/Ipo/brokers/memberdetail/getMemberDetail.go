package memberdetail

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
)

type ResultRespStruct struct {
	MemberDetail   memberDetailStruct  `json:"memberDetail"`
	SegmentDetails []SegmentDetailsRec `json:"segmentDetails"`
	Status         string              `json:"status"`
	ErrMsg         string              `json:"errMsg"`
}

func GetMemberDetail(w http.ResponseWriter, r *http.Request) {
	log.Println("GetMemberDetail (+)", r.Method)
	origin := r.Header.Get("Origin")
	var lBrokerId int
	var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
			log.Println(lErr, origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// This variable is used to store the records
		var lMemberRec memberDetailStruct

		var lCredRec CredentialStruct
		// This variable is used to send resp to frontend
		var lRespRec ResultRespStruct
		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/memberDetail")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BGBL01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BGBL01", "UserDetails not Found"))
			return
		} else {
			log.Println("GetMemberDetail lClientId", lClientId)
			if lClientId == "" {
				// var lErr2 error
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "BGBL02" + lErr1.Error()
				fmt.Fprintf(w, helpers.GetErrorString("BGBL02", "Access restricted"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lDb, lErr3 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr3 != nil {
			log.Println("BGBL03", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("BGBL03", "Unable to establish db connection"))
			return
		} else {
			defer lDb.Close()

			lCoreString2 := `select nvl(m.Flag,'N'),m.OrderPreference ,m.AllowedModules,nvl(d.Member,''),nvl(d.LoginId,'') ,nvl(d.Password,'') ,nvl(d.ibbsid,'') ,nvl(d.Stream ,'') 
			from a_ipo_directory d,a_ipo_memberdetails m
			where d.brokerMasterId = m.BrokerId and m.BrokerId = ?
			and d.Status = 'Y'`

			// and m.BiddingEndDate < curdate() Its goes befor client id in query
			lRows1, lErr4 := lDb.Query(lCoreString2, lBrokerId)
			if lErr4 != nil {
				log.Println("BGBL04", lErr4)
				fmt.Fprintf(w, helpers.GetErrorString("BGBL04", "UserDetails not Found"))
				return
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows1.Next() {
					lErr5 := lRows1.Scan(&lMemberRec.Flag, &lMemberRec.Order, &lMemberRec.Modules, &lCredRec.MemberId, &lCredRec.LoginId, &lCredRec.Password, &lCredRec.Ibbsib, &lCredRec.Exchange)
					if lErr5 != nil {
						log.Println("BGBL05", lErr5)
						fmt.Fprintf(w, helpers.GetErrorString("BGBL05", "UserDetails not Found"))
						return
					} else {
						lPassword, lErr6 := common.DecodeToString(lCredRec.Password)
						if lErr6 != nil {
							log.Println("BGBL05", lErr6)
							fmt.Fprintf(w, helpers.GetErrorString("BGBL06", "Error decoding password "))
							return
						} else {
							lCredRec.Password = lPassword
							// Append the Credentials Records into lCredRec Array
							lMemberRec.Credentials = append(lMemberRec.Credentials, lCredRec)
						}

					}
					// lRespRec.MemberDetail = lMemberRec
					// lRespRec.Status = common.SuccessCode
				}
				lSegmentArr, lErr7 := GetSegments(lBrokerId)
				if lErr7 != nil {
					log.Println("BGBL07", lErr7)
					fmt.Fprintf(w, helpers.GetErrorString("BGBL07", "Error in Geting Segments "))
					return
				} else {
					// This Temp Arr The Array of Data to customize
					lOrderedStruct, lErr8 := ConstructArrayOrder(lSegmentArr, lBrokerId)
					if lErr4 != nil {
						log.Println("BGBL08", lErr8)
						fmt.Fprintf(w, helpers.GetErrorString("BGBL08", "Error in Constructing Array order "))
						return
					} else {
						lRespRec.SegmentDetails = lOrderedStruct
						lRespRec.MemberDetail = lMemberRec
						lRespRec.Status = common.SuccessCode
					}

				}
			}
		}
		lData, lErr7 := json.Marshal(lRespRec)
		if lErr7 != nil {
			log.Println("BGBL06", lErr7)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BGBL06" + lErr7.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BGBL06", "Issue in Getting Your Datas.Try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetMemberDetail (-)", r.Method)
	}
}
