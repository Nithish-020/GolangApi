package ncbplaceorder

// import (
// 	"fcs23pkg/common"
// 	"fcs23pkg/integration/nse/nsencb"
// 	"log"
// )

// type ErrorStruct struct {
// 	ErrCode string `json:"errCode"`
// 	ErrMsg  string `json:"errMsg"`
// }

// func NcbNsePlaceOrder(pExhangeReq nsencb.NcbAddReqStruct, pReqRec NcbReqStruct, pClientId string, pMailId string) (NcbRespStruct, nsencb.NcbAddResStruct, ErrorStruct) {
// 	log.Println("NcbNsePlaceOrder(+)")

// 	var lRespRec NcbRespStruct
// 	var lError ErrorStruct

// 	lResp, lErr1 := ProcessNseReq(pExhangeReq, pReqRec, pClientId, pMailId)
// 	if lErr1 != nil {
// 		log.Println("NNPO01", lErr1.Error())
// 		lRespRec.Status = common.ErrorCode
// 		lError.ErrCode = "NNPO01"
// 		lError.ErrMsg = "Exchange Server is Busy right now,Try After Sometime."
// 		return lRespRec, lResp, lError
// 	} else {
// 		log.Println("lRespRec.status", lResp.Status)
// 		if lResp.Status != "" {
// 			// for lRespIdx := 0; lRespIdx < len(lResp); lRespIdx++ {
// 			if lResp.Status == "failed" {
// 				lRespRec.OrderStatus = lResp.Status
// 				lRespRec.OrderReason = lResp.Reason
// 				lRespRec.Status = common.ErrorCode
// 				lError.ErrCode = "NNPO02"
// 				lError.ErrMsg = lResp.Reason + "\n" + "Application Failed"
// 				log.Println("NNPO02", lRespRec.OrderStatus, lRespRec.OrderReason)
// 				return lRespRec, lResp, lError
// 			} else if lResp.Status == "success" {
// 				lRespRec.OrderStatus = lResp.Status
// 				lRespRec.OrderReason = lResp.Reason
// 				lRespRec.Status = common.SuccessCode
// 			}
// 			// }
// 		} else {
// 			log.Println("lRespRec.status", lResp.Status)
// 			log.Println("NNPO03", "Unable to proceed Application!!!")
// 			lError.ErrCode = "NNPO03"
// 			lError.ErrMsg = "Unable to proceed Application!!!"
// 			log.Println("lRespRec", lRespRec)
// 			return lRespRec, lResp, lError
// 		}
// 	}
// 	log.Println("NcbNsePlaceOrder(-)")
// 	return lRespRec, lResp, lError
// }
