package helpers

import (
	"encoding/json"
	"fcs23pkg/common"
	"log"
)

type Error_Response struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

type Msg_Response struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetMsg_String(Msg_Title string, Msg_Description string) string {

	var Msg_Res Msg_Response

	Msg_Res.Title = Msg_Title
	Msg_Res.Description = Msg_Description

	result, err := json.Marshal(Msg_Res)

	if err != nil {
		log.Println(err)
	}

	return string(result)

}

func GetErrorString(ErrCode string, ErrDescription string) string {
	// log.Println(Err_Title, Err_Description)
	var Err_Response Error_Response
	Err_Response.Status = common.ErrorCode
	Err_Response.ErrMsg = ErrCode + "/" + ErrDescription

	lResult, err := json.Marshal(Err_Response)

	if err != nil {
		log.Println(err)
	}

	return string(lResult)

}
