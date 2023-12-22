package Function

type TrackParaStruct struct {
	HeaderId     int     `json:"header"` // Last inserted id of Header
	ClientId     string  `json:"clientId"`
	BidRefNo     string  `json:"bidRefNo"`
	ActivityType string  `json:"activityType"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
	LastId       int     `json:"lastId"`
}

/*
Pupose: This method is used to store the data for Trackbids datatable
Parameters:
		send TraackParaStruct and Flag as a parameter to this method
Response:
    *On Sucess
    =========
    In case of a successful execution of this method, you will get the dpStructArr data
	from the a_ipo_oder_header Data Table

    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error
Author: Nithish Kumar
Date: 07JUNE2023
*/

// func TrackBid(pTrack TrackParaStruct, pFlag string) (int, error) {
// 	log.Println("TrackBid +")
// 	// This variable is use to store the last inserted id of the TrackBid method
// 	var lId int
// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr != nil {
// 		log.Println("FTB01", lErr)
// 		return lId, lErr
// 	} else {
// 		defer lDb.Close()
// 		// check is the flag is Insert
// 		if pFlag == "INSERT" {
// 			lSqlString1 := `INSERT INTO a_bidtracking_table (headerId,bidRefNo,activityType,quantity,price,CreatedBy,CreatedDate)
// 		values (?,?,?,?,?,?,now())`

// 			lInserted, lErr := lDb.Exec(lSqlString1, &pTrack.HeaderId, &pTrack.BidRefNo, &pTrack.ActivityType, &pTrack.Quantity, &pTrack.Price,
// 				&pTrack.ClientId)

// 			if lErr != nil {
// 				log.Println("FTB02", lErr)
// 				return lId, lErr
// 			} else {
// 				lTrack, _ := lInserted.LastInsertId()
// 				lId = int(lTrack)
// 			}
// 			// Check if the flag is Update
// 		} else if pFlag == "UPDATE" {
// 			lSqlString2 := `UPDATE a_bidtracking_table bt SET bt.status = ?,bt.UpdatedBy = ?,bt.UpdatedDate = now()
// 			WHERE bt.Id = ?`

// 			_, lErr := lDb.Exec(lSqlString2, &pTrack.Status, &pTrack.ClientId, &pTrack.LastId)
// 			if lErr != nil {
// 				log.Println("FTB03", lErr)
// 				return lId, lErr
// 			}
// 		}
// 	}
// 	log.Println("TrackBid -")
// 	return lId, nil
// }
