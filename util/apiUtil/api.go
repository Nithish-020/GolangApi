package apiUtil

import (
	"bytes"
	"io/ioutil"
	"log" //import newly
	"net/http"
	"strings"

	"fcs23pkg/apps/Ipo/Function"
	"fcs23pkg/common"
)

type HeaderDetails struct {
	Key   string
	Value string
}

func Api_call(url string, methodType string, jsonData string, header []HeaderDetails, Source string) (string, error) {
	log.Println("Api_call+")

	var lLogInputRec Function.ParameterStruct
	lLogInputRec.Request = jsonData
	lLogInputRec.EndPoint = url
	lLogInputRec.Flag = common.INSERT
	lLogInputRec.ClientId = "Program"
	lLogInputRec.Method = methodType
	// LogEntry method is used to store the Request in Database
	lId, lErr1 := Function.LogEntry(lLogInputRec)
	if lErr1 != nil {
		// log.Println("CVMP01", lErr1)
		return "", lErr1
	} else {

		//var resp KycApiResponse
		var body []byte
		var err error
		var request *http.Request
		//var endPoint string

		// StrArr := strings.Split(url, "/")
		// //alertSource := Source + " EndPoint: /" + StrArr[len(StrArr)-1]

		// log.Println(StrArr[len(StrArr)-1])
		// endPoint = StrArr[len(StrArr)-1]
		// if endPoint == "" {
		// 	endPoint = StrArr[len(StrArr)-2]
		// }

		//Call API
		log.Println("JsonData: ", jsonData)
		if methodType != "GET" {
			request, err = http.NewRequest(strings.ToUpper(methodType), url, bytes.NewBuffer([]byte(jsonData)))
		} else {
			request, err = http.NewRequest(strings.ToUpper(methodType), url, nil)
		}

		//request, err := http.NewRequest(strings.ToUpper(methodType), url, postJsonBody)
		if err != nil {
			common.LogDebug("apiUtil.Api_call", "(AAC01)", err.Error())
			// err1 := adminAlert.SendAlertMsg(Source, "(AAC01)", url)
			// if err1 != nil {
			// 	common.LogError("emailUtil.SendEmail", "(AAC02)", err1.Error())
			// }
			return "", err
		} else {

			if len(header) > 0 {
				for i := 0; i < len(header); i++ {
					request.Header.Set(header[i].Key, header[i].Value)
				}
			}

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				common.LogDebug("apiUtil.Api_call", "(AAC03)", err.Error())
				// err1 := adminAlert.SendAlertMsg(Source, "(AAC03)", url)
				// if err1 != nil {
				// 	common.LogError("emailUtil.SendEmail", "(AAC04)", err1.Error())
				// }
				log.Println("response: ", jsonData)
				return "", err

			} else {
				defer response.Body.Close()

				body, err = ioutil.ReadAll(response.Body)
				lLogInputRec.Response = string(body)
				lLogInputRec.LastId = lId
				lLogInputRec.Flag = common.UPDATE
				_, lErr4 := Function.LogEntry(lLogInputRec)
				log.Println("Error:", lErr4)
				if err != nil {
					common.LogDebug("apiUtil.Api_call", "(AAC05)", err.Error())
					// err1 := adminAlert.SendAlertMsg(Source, "(AAC05)", url)
					// if err1 != nil {
					// 	common.LogError("emailUtil.SendEmail", "(AAC06)", err1.Error())
					// }
					log.Println("response: ", body)
					return "", err

				}

			}
		}
		log.Println("Api_call-")
		return string(body), nil
	}
}
