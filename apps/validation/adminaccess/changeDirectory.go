package adminaccess

type StreamRespStruct struct {
	Stream string `json:"stream"`
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------
// func ChangeDirectory(w http.ResponseWriter, r *http.Request) {
// 	log.Println("ChangeDirectory (+)", r.Method)

// 	origin := r.Header.Get("Origin")
// 	for _, allowedOrigin := range common.ABHIAllowOrigin {
// 		if allowedOrigin == origin {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			log.Println(origin)
// 			break
// 		}
// 	}

// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// 	if r.Method == "PUT" {

// 		// create a instance of upiResponse struct
// 		var lRespRec StreamRespStruct
// 		lRespRec.Status = common.SuccessCode

// 		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
// 		lClientId := ""
// 		var lErr1 error
// 		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/config")
// 		if lErr1 != nil {
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "ACD01" + lErr1.Error()
// 			fmt.Fprintf(w, helpers.GetErrorString("ACD01", "UserDetails not Found"))
// 			return
// 		} else {
// 			if lClientId == "" {
// 				fmt.Fprintf(w, helpers.GetErrorString("ACD01", "UserDetails not Found"))
// 				return
// 			}
// 		}
// 		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

// 		// Calling LocalDbConect method in ftdb to estabish the database connection
// 		lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
// 		if lErr2 != nil {
// 			log.Println("ACD02", lErr2)
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "ACD02" + lErr2.Error()
// 			fmt.Fprintf(w, helpers.GetErrorString("ACD02", "Issue in Connecting DB.Please try after sometime"))
// 			return
// 		} else {
// 			defer lDb.Close()

// 			lCoreString := ``
// 			_, lErr3 := lDb.Exec(lCoreString)
// 			if lErr3 != nil {
// 				log.Println("ACD03", lErr3)
// 				lRespRec.Status = common.ErrorCode
// 				lRespRec.ErrMsg = "ACD03" + lErr3.Error()
// 				fmt.Fprintf(w, helpers.GetErrorString("ACD03", "Unable to Change Directory"))
// 				return
// 			} else {
// 				//This for loop is used to collect the record from the database and store them in structure
// 				// for lRows.Next() {
// 				// 	lErr4 := lRows.Scan(&lRespRec.Stream)
// 				// 	if lErr4 != nil {
// 				// 		log.Println("ACD04", lErr4)
// 				// 		lRespRec.Status = common.ErrorCode
// 				// 		lRespRec.ErrMsg = "ACD04" + lErr4.Error()
// 				// 		fmt.Fprintf(w, helpers.GetErrorString("ACD04", "Unable to get Directory. Please try after sometime"))
// 				// 		return
// 				// 	}
// 				// }
// 			}
// 		}
// 		// Marshall the structure into json
// 		lData, lErr5 := json.Marshal(lRespRec)
// 		if lErr5 != nil {
// 			log.Println("ACD05", lErr5)
// 			fmt.Fprintf(w, helpers.GetErrorString("ACD05", "Unable to process your request now. Please try after sometime"))
// 			return
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 		log.Println("ChangeDirectory (-)", r.Method)
// 	}
// }
