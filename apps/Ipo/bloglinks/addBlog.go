package bloglinks

import (
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// response struct for AddBlogLinks API
type BlogRespStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func AddBlogLinks(w http.ResponseWriter, r *http.Request) {
	log.Println("AddBlogLinks (+)")
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "POST" {
		var lReqRec blogLinkStruct
		var lRespRec BlogRespStruct

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipoInfo")
		if lErr1 != nil {
			log.Println("BABL01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BABL01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BABL01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("BABL02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// read the body values from request
		lBody, lErr2 := ioutil.ReadAll(r.Body)
		if lErr2 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BABL03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BABL03", "Reading Value Issue.Please try after sometime"))
			return
		} else {
			// unmarsal the body values to corresponding structure
			lErr3 := json.Unmarshal(lBody, &lReqRec)
			log.Println("lReqRec", lReqRec)
			if lErr3 != nil {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "BABL04" + lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("BABL04", "Unable to get Request,Please try after sometime"))
				return
			} else {
				// call the add blog method to insert blog and drhp links in database
				lIsinFlag, lErr4 := AddBlog(lReqRec, lClientId)
				if lErr4 != nil {
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "BABL05" + helpers.ErrPrint(lErr4)
					fmt.Fprintf(w, helpers.GetErrorString("BABL05", "Issue in Updating datas,Please try after sometime"))
					return
				} else {
					if lIsinFlag != "" {
						lRespRec.Status = common.ErrorCode
						lRespRec.ErrMsg = lIsinFlag
					}

				}
			}
		}
		lData, lErr6 := json.Marshal(lRespRec)
		if lErr6 != nil {
			fmt.Fprintf(w, helpers.GetErrorString("BABL06", "Error in getting datas"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("AddBlogLinks (-)", r.Method)
	}
}

func AddBlog(pReq blogLinkStruct, pClientId string) (string, error) {
	log.Println("AddBlog (+)")
	var lExistflag string
	var lCategory int

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BAB01", lErr1)
		return lExistflag, helpers.ErrReturn(lErr1)
	} else {
		defer lDb.Close()

		if pReq.Sme == true {
			lCategory = 1
		} else {
			lCategory = 0
		}
		//if the given id is not equal to  0 then update the details in databas
		if pReq.Id == 0 {
			IsinFlag, lErr2 := CheckIsinExist(pReq)
			if lErr2 != nil {
				log.Println("BAB02", lErr2)
				return lExistflag, lErr2
			} else {
				log.Println("IsinFlag", IsinFlag)
				if IsinFlag == "Y" {

					lExistflag = "Isin Already Exist"
				} else if IsinFlag == "N" {

					sqlString := `insert into a_ipo_details (Symbol,Isin,detailsLink,drhpLink,category,allotmentFinal ,refundInitiate ,
						dematTransfer ,listingDate ,createdBy,createdDate,updatedBy,updatedDate)		
						values(?,?,?,?,?,?,?,?,?,?,now(),?,now())`

					_, lErr3 := lDb.Exec(sqlString, pReq.Symbol, pReq.Isin, pReq.BlogLink, pReq.DRHPLink, lCategory, pReq.AllotmentFinal, pReq.RefundInitiate, pReq.DematTransfer, pReq.ListingDate, pClientId, pClientId)
					if lErr3 != nil {
						log.Println("BAB03", lErr3)
						return lExistflag, lErr3
					}
				}
			}

		} else {
			// IsinFlag, lErr4 := CheckIsinExist(pReq)
			// if lErr4 != nil {
			// 	log.Println("BAB04", lErr4)
			// 	return lExistflag, lErr4
			// } else {
			// 	if IsinFlag == "Y" {
			// 		lExistflag = "Isin Already Exist"
			// 	} else if IsinFlag == "N" {

			sqlString := `update a_ipo_details d 
							set d.Symbol = ?, d.Isin = ?,d.detailsLink = ?,d.drhpLink = ?,d.category = ?,d.allotmentFinal = ?,d.refundInitiate = ?,
							d.dematTransfer = ?,d.listingDate = ?,d.updatedBy = ?,d.updatedDate = now()
							where d.Id = ?`
			_, lErr5 := lDb.Exec(sqlString, pReq.Symbol, pReq.Isin, pReq.BlogLink, pReq.DRHPLink, lCategory, pReq.AllotmentFinal, pReq.RefundInitiate, pReq.DematTransfer, pReq.ListingDate, pClientId, pReq.Id)
			if lErr5 != nil {
				log.Println("BAB05", lErr5)
				return lExistflag, lErr5
			}
		}
		// }
		// }
	}
	log.Println("AddBlog (-)")
	return lExistflag, nil
}

//  This method will Give The Flag For only inserting and isin
func CheckIsinExist(pBlogReq blogLinkStruct) (string, error) {
	log.Println("CheckIsinExist (+)")
	var lIsinFlag string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BCIE01", lErr1)
		return lIsinFlag, lErr1
	} else {
		defer lDb.Close()
		lCoreString2 := `select nvl((select  case when Isin  = ? then 'Y' else 'N'end
		from a_ipo_details aid 
		where isin = ?),'N') as flag`

		lRows1, lErr2 := lDb.Query(lCoreString2, pBlogReq.Isin, pBlogReq.Isin)
		if lErr2 != nil {
			log.Println("BCIE02", lErr2)
			return lIsinFlag, lErr2
		} else {
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lIsinFlag)
				if lErr3 != nil {
					log.Println("BCIE05", lErr3)
					return lIsinFlag, lErr3
				}
			}

		}
	}
	log.Println("CheckIsinExist (-)")
	return lIsinFlag, nil
}
