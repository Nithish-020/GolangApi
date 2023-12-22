package memberdetail

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/apps/validation/menu"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type CredentialStruct struct {
	Exchange string `json:"exchange"`
	MemberId string `json:"memberId"`
	LoginId  string `json:"login"`
	Password string `json:"password"`
	Ibbsib   string `json:"ibbsid"`
}

type memberDetailStruct struct {
	Flag          string              `json:"flag"`
	Modules       string              `json:"selectedModules"`
	Order         string              `json:"selectedOrder"`
	Credentials   []CredentialStruct  `json:"credentials"`
	SegmentShares []SegmentDetailsRec `json:"SegmentDetails"`
}
type memberDetailResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func SetMemberDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("SetMemberDetails (+)")
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
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {

		var lReqRec memberDetailStruct
		var lRespRec memberDetailResp

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/report")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BSB01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BSB01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "BSB02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("BSB02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lBody, lErr3 := ioutil.ReadAll(r.Body)
		if lErr3 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BSB03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BSB03", "Request cannot be empty.Please try after sometime"))
			return
		} else {
			lErr4 := json.Unmarshal(lBody, &lReqRec)
			log.Println("Request SetMemberDetail", lReqRec)
			if lErr4 != nil {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "BSB04" + lErr4.Error()
				log.Println("BSB04", lErr4)
				fmt.Fprintf(w, helpers.GetErrorString("BSB04", "Unable to proccess your request right now.Please try after sometime"))
				return
			} else {

				lDb, lErr5 := ftdb.LocalDbConnect(ftdb.IPODB)
				if lErr5 != nil {
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "BSB05" + lErr5.Error()
					fmt.Fprintf(w, helpers.GetErrorString("BSB05", "Error in connecting server.Please try after sometime"))
					return
				} else {
					defer lDb.Close()
					// ! Get 0 if MemberDetailsId is not present,If present get MemberDetailsId
					lMemberDetailId, lErr6 := brokers.GetMemberDetailId(lBrokerId)
					if lErr6 != nil {
						lRespRec.Status = common.ErrorCode
						lRespRec.ErrMsg = "BSB06" + lErr6.Error()
						fmt.Fprintf(w, helpers.GetErrorString("BSB06", "Error in connecting server.Please try after sometime"))
						return
					} else {
						lExchangeStatus, lErr7 := InsertAndUpdateDetails(lMemberDetailId, lReqRec, lClientId, lBrokerId)
						if lErr7 != nil {
							log.Println("BAB07", lErr7)
							fmt.Fprintf(w, helpers.GetErrorString("BSB07", "Error in Setting Details .Please try after sometime"))
							return
						} else {
							if lExchangeStatus == "" {
								lErr8 := UpdateSegmentType(lBrokerId)
								if lErr8 != nil {
									log.Println("BAB08", lErr8)
									fmt.Fprintf(w, helpers.GetErrorString("BSB08", "Error in Disabling Segments .Please try after sometime"))
									return
								} else {
									lSharelimit, lErr9 := InsertAndUpdateSegments(lBrokerId, lReqRec)
									if lErr9 != nil {
										log.Println("BAB09", lErr9)
										fmt.Fprintf(w, helpers.GetErrorString("BSB09", "Error in Inseting And Updating Segments .Please try after sometime"))
										return
									} else {
										if lSharelimit != nil {
											var lSharelimitMsg string
											for _, Segments := range lSharelimit {
												lSharelimitMsg += Segments + " "

											}
											lSharelimitMsg += "Total Share Limit Lesser or Greater than 100 Not Acceptable"
											lRespRec.Status = common.ErrorCode
											lRespRec.ErrMsg = lSharelimitMsg

										}
									}
								}

							} else if lExchangeStatus != "" {
								lRespRec.Status = common.ErrorCode
								lRespRec.ErrMsg = lExchangeStatus

							}
						}
					}
				}
			}
			// Marshal the response structure into lData
			lData, lErr8 := json.Marshal(lRespRec)
			if lErr8 != nil {
				log.Println("BSB08", lErr8)
				fmt.Fprintf(w, helpers.GetErrorString("BSB08", "Error in Responsing ,Please try again!"))
				return
			} else {
				fmt.Fprintf(w, string(lData))
			}
		}
		log.Println("SetMemberDetails (-)")
	}
}

func GetAllowOrigin() []string {
	log.Println("GetAllowOrigin (+)")

	var lAllowOrigin []string
	var lDomain string
	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr != nil {
		log.Println("BGAO01", lErr)
		return lAllowOrigin
	} else {
		defer lDb.Close()

		// ! To get the active DomainNames from the DB
		lCoreString2 := `select bm.DomainName
			from a_ipo_brokerMaster bm`

		lRows1, lErr := lDb.Query(lCoreString2)
		if lErr != nil {
			log.Println("BGAO02", lErr)
			return lAllowOrigin
		} else {
			//This for loop is used to collect the records from the database and assign them to lDomain
			for lRows1.Next() {
				lErr := lRows1.Scan(&lDomain)
				if lErr != nil {
					log.Println("BGAO03", lErr)
					return lAllowOrigin
				} else {
					// Append the domains into the lArrowedOrigin array[]
					lAllowOrigin = append(lAllowOrigin, lDomain)
				}
			}
		}
	}
	log.Println("GetAllowOrigin (-)")
	return lAllowOrigin
}

func InsertAndUpdateDetails(pMemberDetailId int, pMemberDetail memberDetailStruct, pCliendId string, pBrokerId int) (string, error) {
	log.Println("InsertAndUpdateDetails (+)")
	var ExchangeStatus string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BIAUD01")
		return ExchangeStatus, lErr1
	} else {
		defer lDb.Close()
		if pMemberDetailId == 0 {
			lSqlString1 := `insert into a_ipo_memberdetails (BrokerId, OrderPreference ,AllowedModules ,Flag,CreatedBy ,CreatedDate ,UpdatedBy,UpdatedDate)
			values (?,?,?,?,?,now(),?,now())`
			_, lErr2 := lDb.Exec(lSqlString1, pBrokerId, pMemberDetail.Order, pMemberDetail.Modules, "Y", pCliendId, pCliendId)
			if lErr2 != nil {
				log.Println("BIAUD02", lErr2)
				return ExchangeStatus, lErr2
			}
		}
		for _, Credential := range pMemberDetail.Credentials {
			lpassword := common.EncodeToString(Credential.Password)
			//
			lSqlString2 := `Insert into a_ipo_directory (Member,LoginId,Password,ibbsid,brokerMasterId,Stream,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
									select ?,?,?,?,?,?,?,now(),?,now()
									from a_ipo_directory 
									where not exists(select * from a_ipo_directory where  brokerMasterId = ? and Stream= ?)
									limit 1`
			_, lErr3 := lDb.Exec(lSqlString2, Credential.MemberId, Credential.LoginId, lpassword, Credential.Ibbsib, pBrokerId, Credential.Exchange, pCliendId, pCliendId, pBrokerId, Credential.Exchange)
			if lErr3 != nil {
				log.Println("BIAUD03", lErr3)
				return ExchangeStatus, lErr3
			}
		}
		if pMemberDetailId != 0 {

			lFlag, lErr4 := AllowOrderPreference(pMemberDetail.Order, pMemberDetail.Modules, pBrokerId)
			if lErr4 != nil {
				log.Println("BIAUD04", lErr4)
				return ExchangeStatus, lErr4
			} else {
				log.Println("lFlag", lFlag)
				if lFlag == "N" {
					ExchangeStatus = "Unable To update the Order Preference Right Now"
				} else if lFlag == "Y" {

					lSqlString3 := `update a_ipo_memberdetails set OrderPreference = ?,AllowedModules = ?,UpdatedBy = ?,UpdatedDate = now() 
									where BrokerId = ?`
					_, lErr5 := lDb.Exec(lSqlString3, pMemberDetail.Order, pMemberDetail.Modules, pCliendId, pBrokerId)
					if lErr5 != nil {
						log.Println("BIAUD05", lErr5)
						return ExchangeStatus, lErr5
					}
				}

				for _, Credential := range pMemberDetail.Credentials {
					lpassword := common.EncodeToString(Credential.Password)
					lSqlString4 := `update a_ipo_directory set Member = ?,LoginId = ?,Password = ?,ibbsid = ?,UpdatedBy = ?,UpdatedDate = now() 
									where Stream = ? and brokerMasterId = ?`
					_, lErr6 := lDb.Exec(lSqlString4, Credential.MemberId, Credential.LoginId, lpassword, Credential.Ibbsib, pCliendId, Credential.Exchange, pBrokerId)
					if lErr6 != nil {
						log.Println("BIAUD06", lErr6)
						return ExchangeStatus, lErr6
					}

				}
			}
		}
		// }
	}
	log.Println("InsertAndUpdateDetails (-)")
	return ExchangeStatus, nil
}

func AllowOrderPreference(pExchange string, pModules string, pBrokerId int) (string, error) {
	log.Println("AllowOrderPreference (+)")
	// log.Println("ABhiBroker", common.ABHIBrokerId)
	var lcount int
	var lFlag string
	lFlag = "Y"
	var lExchange string
	if pExchange == common.NSE {
		lExchange = common.BSE
	} else {
		lExchange = common.NSE
	}
	lAllowedModules, lErr1 := menu.GetAllowModules(pBrokerId)
	if lErr1 != nil {
		log.Println("MAOP01", lErr1)
		return lFlag, lErr1
	} else {
		var UIstringArray []string
		UIstringArray = strings.Split(pModules, "/")

		lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr2 != nil {
			log.Println("MAOP02", lErr2)
			return lFlag, lErr2
		} else {
			defer lDb.Close()
			var DBModuleArray []string
			for _, lAllow := range lAllowedModules {
				Inticap := common.CapitalizeText(lAllow)
				DBModuleArray = append(DBModuleArray, Inticap)
			}
			var areEqual bool
			areEqual = reflect.DeepEqual(DBModuleArray, UIstringArray)
			if areEqual == false {

				count := 0

				// DB
				for _, Inticap := range DBModuleArray {
					count = 0
					//ui
					for _, Segments := range UIstringArray {
						if Inticap == Segments {
							count++
						}
					}
					if count == 0 {
						areEqual = false

						lFrom := ""
						lWhere := ""

						lIpoFrom := "a_ipo_master M,a_ipo_order_header IH"

						lIpoWhere := "M.Id = IH.MasterId and M.BiddingEndDate > Date(now()) and IH.cancelFlag !='N' and IH.Exchange = ? and BM.Id = IH.BrokerId and BM.Id = ?"

						lSgbFrom := "a_sgb_master S,a_sgb_orderheader SH "

						lSgbWhere := "SH.MasterId = S.id  and S.BiddingEndDate  > Date(now()) and  sh.cancelFlag != 'N' and SH.Exchange  = ? and BM.Id = SH.brokerId  and BM.Id = ? "
						if Inticap == "Ipo" {
							lFrom = lIpoFrom
							lWhere = lIpoWhere

						} else if Inticap == "Sgb" {
							lFrom = lSgbFrom
							lWhere = lSgbWhere
						} else if Inticap == "Gsec" {
							//  Need To Write for Gsec
							lFlag = "Y"
							return lFlag, nil
						}
						lCoreString1 := `select count(1)
					from a_ipo_brokermaster BM, ` + lFrom + `
					Where ` + lWhere + ``
						log.Println("lCoreString1", lCoreString1)
						lRows, lErr3 := lDb.Query(lCoreString1, lExchange, pBrokerId)
						if lErr3 != nil {
							log.Println("MAOP03", lErr3)
							return lFlag, lErr3
						} else {
							//This for loop is used to collect the records from the database and assign them to lDomain
							for lRows.Next() {
								lErr4 := lRows.Scan(&lcount)
								if lErr4 != nil {
									log.Println("MAOP03", lErr4)
									return lFlag, lErr4
								}
							}
						}
						if lcount > 0 {
							lFlag = "N"
							break
						} else {
							lFlag = "Y"
						}
					}
				}
			}
		}
		log.Println("AllowOrderPreference (-)")
		return lFlag, nil
	}
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var dummy []string
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		} else {
			dummy = append(dummy, item)
		}
	}
	return list
}
