package dashboard

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/roleTask"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SegmentDetailStruct struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Fullname string `json:"fullname"`
	Path     string `json:"path"`
	TaskId   int    `json:"taskId"`
	Image    string `json:"image"`
	Color    string `json:"color"`
	Status   string `json:"status"`
	// Available int    `json:"available"`
	// PreApply  int    `json:"preApply"`
	// Upcoming  int    `json:"upcoming"`
}

type DashBoardRespStruct struct {
	SegmentArr []SegmentDetailStruct `json:"segmentArr"`
	RouterArr  []roleTask.Task       `json:"routerArr"`
	Status     string                `json:"status"`
	ErrMsg     string                `json:"errMsg"`
}

type SetDashRespStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func GetDashboardDetail(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDashboardDetail (+)")
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			_, lErr := brokers.GetBrokerId(origin) // TO get brokerId
			log.Println(lErr, origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Path,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		var lDashResp DashBoardRespStruct
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lCurrentRouter := r.Header.Get("Path")
		log.Println("lCurrentRouter := ", lCurrentRouter)
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr1 != nil {
			log.Println("DGDD01", lErr1)
			lDashResp.Status = common.ErrorCode
			lDashResp.ErrMsg = "DGDD01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("DGDD01", "Error while fetching your details"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("DGDD02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lDashResp.Status = common.SuccessCode
		// lMasterArr, lErr3 := CalcMasterDetails(lClientId, lBrokerId)
		// if lErr3 != nil {
		// 	log.Println("DGDD03", lErr3)
		// 	lDashResp.Status = common.ErrorCode
		// 	lDashResp.ErrMsg = "DGDD03" + lErr3.Error()
		// 	fmt.Fprintf(w, helpers.GetErrorString("DGDD03", "Error while fetching your details"))
		// 	return
		// } else {

		lResultArr, lErr4 := GetSegmentDetails(lCurrentRouter)
		if lErr4 != nil {
			log.Println("DGDD04", lErr4)
			lDashResp.Status = common.ErrorCode
			lDashResp.ErrMsg = "DGDD04" + lErr4.Error()
			fmt.Fprintf(w, helpers.GetErrorString("DGDD04", "Error while fetching segment details"))
			return
		} else {
			lDashResp.SegmentArr = lResultArr
			// log.Println("lResultArr", lDashResp.SegmentArr)

			// This method only call when the api is called from the router "/dashboard/setup"
			if lCurrentRouter == "/dashboard/setup" {
				lRouterArr, lErr5 := roleTask.GetTaskList()
				if lErr5 != nil {
					log.Println("DGDD05", lErr5)
					lDashResp.Status = common.ErrorCode
					lDashResp.ErrMsg = "DGDD05" + lErr5.Error()
					fmt.Fprintf(w, helpers.GetErrorString("DGDD05", "Error while fetching segment details"))
					return
				} else {
					lDashResp.RouterArr = lRouterArr
				}
			}
		}
		// }

		// Marshaling the response structure to lData
		lData, lErr6 := json.Marshal(lDashResp)
		if lErr6 != nil {
			log.Println("DGDD06", lErr6)
			fmt.Fprintf(w, helpers.GetErrorString("DGDD06", "Issue in Getting Datas!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}
	log.Println("GetDashboardDetail (-)")
}

// func GetSegmentDetails(pCurrentPath string, pMasterDetail []SegmentDetailStruct) ([]SegmentDetailStruct, error) {
func GetSegmentDetails(pCurrentPath string) ([]SegmentDetailStruct, error) {
	log.Println("GetSegmentDetails (+)")
	var lDetail SegmentDetailStruct
	var lDetailArr []SegmentDetailStruct

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("DGSD01", lErr1)
		return lDetailArr, lErr1
	} else {
		defer lDb.Close()

		lAdditionalCond := ""
		lCoreString := `select ns.Id, nvl(ns.Name,''),nvl(ns.FullName,'') ,nvl(air.Router,'') ,air.Id,
		nvl(ns.Image,'') ,nvl(ns.CardColor ,''),nvl(ns.Status,'')
		from a_novo_segment ns,a_ipo_router air 
		where air.Id  = ns.RouterId`

		if pCurrentPath == "/dashboard" {
			lAdditionalCond = ` and ns.Status = 'Y'`
		}
		lCoreString1 := lCoreString + lAdditionalCond
		lRows2, lErr2 := lDb.Query(lCoreString1)
		if lErr2 != nil {
			log.Println("DGSD02", lErr2)
			return lDetailArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows2.Next() {
				lErr3 := lRows2.Scan(&lDetail.Id, &lDetail.Name, &lDetail.Fullname, &lDetail.Path, &lDetail.TaskId, &lDetail.Image, &lDetail.Color,
					&lDetail.Status)
				if lErr3 != nil {
					log.Println("DGSD03", lErr3)
					return lDetailArr, lErr3
				} else {
					// for _, lmaster := range pMasterDetail {
					// 	if lmaster.Name == lDetail.Name {
					// 		lDetail.Available = lmaster.Available
					// 		lDetail.PreApply = lmaster.PreApply
					// 		lDetail.Upcoming = lmaster.Upcoming

					lDetailArr = append(lDetailArr, lDetail)
					// 	}
					// }

				}
			}
		}
	}
	log.Println("GetSegmentDetails (-)")
	return lDetailArr, nil
}

// func CalcMasterDetails(pClientId string, pBrokerId int) ([]SegmentDetailStruct, error) {
// 	log.Println("CalcMasterDetails (+)")
// 	var lSegRec SegmentDetailStruct
// 	var lCalcArr []SegmentDetailStruct

// 	lIpoArr, lErr1 := localdetail.GetIpoMasterDetail(pClientId, pBrokerId)
// 	if lErr1 != nil {
// 		log.Println("DCMD01", lErr1)
// 		return lCalcArr, lErr1
// 	} else {
// 		lSegRec.Name = "Ipo"
// 		for _, lIpo := range lIpoArr {
// 			if lIpo.Upcoming == "C" {
// 				lSegRec.Available++
// 			} else if lIpo.Upcoming == "U" {
// 				lSegRec.Upcoming++
// 			} else if lIpo.PreApply == "pre" {
// 				lSegRec.PreApply++
// 			}
// 		}
// 		// append the calculated Ipo master information in the segmentArray
// 		lCalcArr = append(lCalcArr, lSegRec)
// 		lSegRec = SegmentDetailStruct{}

// 		lSgbArr, lErr2 := localdetails.GetSGBdetail(pClientId, pBrokerId)
// 		if lErr2 != nil {
// 			log.Println("DCMD02", lErr2)
// 			return lCalcArr, lErr2
// 		} else {
// 			lSegRec.Name = "Sgb"
// 			for _, _lSgb := range lSgbArr {
// 				if _lSgb.Upcoming == "C" {
// 					lSegRec.Available++
// 				} else if _lSgb.Upcoming == "U" {
// 					lSegRec.Upcoming++
// 				} else {
// 					lSegRec.PreApply++
// 				}
// 			}
// 			// append the calculated Sgb master information in the segmentArray
// 			lCalcArr = append(lCalcArr, lSegRec)
// 			lSegRec = SegmentDetailStruct{}

// 		}

// 		lSegRec.Name = "Gsecs"
// 		// append the calculated Gsecs master information in the segmentArray
// 		lCalcArr = append(lCalcArr, lSegRec)

// 		lSegRec.Name = "MutualFund"
// 		// append the calculated MutualFund master information in the segmentArray
// 		lCalcArr = append(lCalcArr, lSegRec)

// 		lSegRec.Name = "CorporateBond"
// 		// append the calculated CorporateBond master information in the segmentArray
// 		lCalcArr = append(lCalcArr, lSegRec)

// 	}

// 	log.Println("CalcMasterDetails (-)")
// 	return lCalcArr, nil
// }

func SetDashboardDetail(w http.ResponseWriter, r *http.Request) {
	log.Println("SetDashboardDetail (+)")
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "POST" {
		var lDashRec SegmentDetailStruct
		var lDashResp DashBoardRespStruct

		lDashResp.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/dashboard/setup")
		if lErr1 != nil {
			log.Println("DSDD01", lErr1)
			lDashResp.Status = common.ErrorCode
			lDashResp.ErrMsg = "DSDD01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("DSDD01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lErr2 := common.CustomError("Client Details Not Found")
				lDashResp.Status = common.ErrorCode
				lDashResp.ErrMsg = "DSDD02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("DSDD02", "Client Details Not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// read the body values from request
		lBody, lErr3 := ioutil.ReadAll(r.Body)
		if lErr3 != nil {
			lDashResp.Status = common.ErrorCode
			lDashResp.ErrMsg = "DSDD03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("DSDD03", "Reading Value Issue.Please try after sometime"))
			return
		} else {
			// unmarsal the body values to corresponding structure
			lErr4 := json.Unmarshal(lBody, &lDashRec)
			log.Println("lDashRec", lDashRec)
			if lErr4 != nil {
				lDashResp.Status = common.ErrorCode
				lDashResp.ErrMsg = "DSDD04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("DSDD04", "Unable to get Request,Please try after sometime"))
				return
			} else {
				// This method should insert or update the request body to the database
				lErr5 := insertAndUpdateSegmentDetails(lDashRec, lClientId)
				if lErr5 != nil {
					lDashResp.Status = common.ErrorCode
					lDashResp.ErrMsg = "DSDD05" + lErr5.Error()
					fmt.Fprintf(w, helpers.GetErrorString("DSDD05", "Unable to update your Request right now,Please try after sometime"))
					return
				}
			}
		}
		// Marshaling the response structure to lData
		lData, lErr6 := json.Marshal(lDashResp)
		if lErr6 != nil {
			log.Println("DSDD06", lErr6)
			fmt.Fprintf(w, helpers.GetErrorString("DSDD06", "Issue in Getting Datas!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("SetDashboardDetail (-)")
	}
}

func insertAndUpdateSegmentDetails(pSegment SegmentDetailStruct, pClientId string) error {
	log.Println("insertAndUpdateSegmentDetails (+)")
	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("DIAUSD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		//if the given id is not equal to  0 then update the details in databas
		if pSegment.Id == 0 {
			ExistFlag, lErr2 := CheckSegmentExist(pSegment.Name)
			if lErr2 != nil {
				log.Println("DIAUSD02", lErr2)
				return lErr2
			} else {
				if ExistFlag == "Y" {
					lErr3 := common.CustomError("Segment Already Exists")
					log.Println("DIAUSD03", lErr3)
					return lErr3

				} else if ExistFlag == "N" {

					sqlString := `insert into a_novo_segment (Name,FullName,RouterId,Image,CardColor,Status,CreatedBy,CreatedDate)
					values (?,?,?,?,?,?,?,now())`

					_, lErr4 := lDb.Exec(sqlString, pSegment.Name, pSegment.Fullname, pSegment.TaskId, pSegment.Image, pSegment.Color, pSegment.Status, pClientId)
					if lErr4 != nil {
						log.Println("DIAUSD03", lErr4)
						return lErr4
					}
				}
			}

		} else {
			sqlString := `update a_novo_segment ns 
					set ns.Name = ?, ns.FullName = ?,ns.RouterId = ?,ns.Image = ?,ns.CardColor = ?,ns.Status = ?,ns.UpdatedBy = ?,ns.UpdatedDate = now()
					where ns.Id = ?`
			_, lErr5 := lDb.Exec(sqlString, pSegment.Name, pSegment.Fullname, pSegment.TaskId, pSegment.Image, pSegment.Color, pSegment.Status, pClientId, pSegment.Id)
			if lErr5 != nil {
				log.Println("DIAUSD05", lErr5)
				return lErr5
			}
		}

	}
	log.Println("insertAndUpdateSegmentDetails (-)")
	return nil
}

func CheckSegmentExist(segment string) (string, error) {
	log.Println("CheckSegmentExist (+)")
	var lFlag string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("DCSE01", lErr1)
		return lFlag, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select (case when ns.Name = ? then 'Y' else 'N' end)
						from a_novo_segment ns,a_ipo_router air`

		lRows2, lErr2 := lDb.Query(lCoreString, segment)
		if lErr2 != nil {
			log.Println("DCSE02", lErr2)
			return lFlag, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in the variable
			for lRows2.Next() {
				lErr3 := lRows2.Scan(&lFlag)
				if lErr3 != nil {
					log.Println("DCSE03", lErr3)
					return lFlag, lErr3
				}
			}
		}
	}
	log.Println("CheckSegmentExist (-)")
	return lFlag, nil
}
