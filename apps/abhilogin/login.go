package abhilogin

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fcs23pkg/apigate"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/appsso"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type LoginMRequest struct {
	ClientId string `json:"clientId"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

type logout struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

type RespStruct struct {
	ClientId string `json:"clientId"`
	Stat     string `json:"stat"`
	Errmsg   string `json:"errmsg"`
}

type LoginRespStruct struct {
	ClientId string `json:"clientId"`
	Status   string `json:"status"`
	Errmsg   string `json:"errmsg"`
}

type JData struct {
	Apkversion string `json:"apkversion"`
	Uid        string `json:"uid"`
	Pwd        string `json:"pwd"`
	Factor2    string `json:"factor2"`
	Imei       string `json:"imei"`
	Source     string `json:"source"`
	Vc         string `json:"vc"`
	Appkey     string `json:"appkey"`
}

// TOKEN VALIDATION
type tokenValidation struct {
	ClientId string `json:"clientId"`
	Status   string `json:"status"`
	ErrMsg   string `json:"errMsg"`
}

type JDataRespStruct struct {
	Request_time   string   `json:"request_time"`
	Actid          string   `json:"actid"`
	Uname          string   `json:"uname"`
	Prarr          []JVprd  `json:"prarr"`
	Stat           string   `json:"stat"`
	EMail          string   `json:"email"`
	Exarr          []string `json:"exarr"`
	Brkname        string   `json:"FTC"`
	Lastaccesstime string   `json:"lastaccesstime"`
	Emsg           string   `json:"emsg"`
}
type JVprd struct {
	Prd        string `json:"lastaccesstime"`
	S_prdt_ali string `json:"s_prdt_ali"`
}

func ValidateAbhiToken(w http.ResponseWriter, req *http.Request) {
	log.Println("ValidateABHIToken(+) " + req.Method)
	origin := req.Header.Get("Origin")
	// var lBrokerId int
	// var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			//commented by nithish
			// the below method has no use in this function
			// lBrokerId, lErr := brokers.GetBrokerId(origin)
			// log.Println(origin)
			// log.Println("If Error in GetBrokerId : ", lErr)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)

	switch req.Method {
	case "GET":

		var validation tokenValidation

		publicTokenCookie, err := req.Cookie(common.ABHICookieName)
		// log.Println("Cookie: ", publicTokenCookie, lBrokerId)
		if err != nil {
			log.Println("error reading publicToken cookie")
			log.Println(err)
			validation.Status = "E"
			validation.ErrMsg = err.Error()
		} else {
			if publicTokenCookie.Value != "" {
				clientId := appsso.CheckPageTokenValidity2(publicTokenCookie.Value, common.ABHIAppName)
				if clientId != "" {
					validation.ClientId = clientId
					validation.Status = "S"
				} else {
					validation.Status = "I"
				}
				log.Println("Validation status 1: ", validation.Status)
			} else {
				validation.Status = "I"
				log.Println("Validation status 2: ", validation.Status)
			}
		}
		data, err := json.Marshal(validation)
		if err != nil {
			fmt.Fprintf(w, "Error taking data"+err.Error())
		} else {
			fmt.Fprintf(w, string(data))
		}

		log.Println("ValidateABHIToken(-) ")

	}

}

type abhiInputStruct struct {
	Token    string `json:"token"`
	ClientID string `json:"clientID"`
	Key      string `json:"key"`
}

type abhiReturnStruct struct {
	Token  string `json:"token"`
	Status string `json:"status"`
	Emsg   string `json:"emsg"`
}

//--------------------------------------------------------------------
//  function to insert sso token
//--------------------------------------------------------------------

func GenerateAbhiSessionToken(db *sql.DB, clientid string, appKey string, reqDtl apigate.RequestorDetails, pBrokerId int) string {

	b := make([]byte, 32)
	rand.Read(b)
	token := fmt.Sprintf("%x."+"%x", time.Now().UnixNano(), b)
	//insert token
	insertString := "insert into xxapi_ssotokens(app,clientid,token,createdtime,expiretime,realip,forwardedip,method,path,host,remoteaddr) values (?,?,?,now() ,ADDTIME(now(), '06:00:00.999998'),?,?,?,?,?,?)"
	_, err := db.Exec(insertString, appKey, clientid, token, reqDtl.RealIP, reqDtl.ForwardedIP, reqDtl.Method, reqDtl.Path, reqDtl.Host, reqDtl.RemoteAddr)
	if err != nil {
		log.Println("GenerateSSOToken insert error", err.Error())
	} else {
		// TO store the token in a_ipo_CookieDe
		lFlag, lerr := brokers.InsertCookieDetails(token, pBrokerId, clientid)
		if lerr != nil {
			log.Println("GenerateSSOToken insert error", lerr.Error())
		} else {
			common.ABHIFlag = lFlag
			log.Println("common.ABHIFlag", common.ABHIFlag)
			return token
		}
	}

	return ""
}

// --------------------------------------------------------------------
//
//	function validate HS token and genereate new token for thirdparty app
//
// --------------------------------------------------------------------
func ValidateAbhiPageHSToken(w http.ResponseWriter, req *http.Request) {

	origin := req.Header.Get("Origin")
	log.Println("origin", origin)
	var lBrokerId int
	var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		log.Println("allowedOrigin == origin", allowedOrigin, origin)
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
			log.Println(lErr, origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	//bodyMsg, _ := ioutil.ReadAll(req.Body)
	reqDtl := apigate.GetRequestorDetail(req)
	//reqDtl.Body = string(bodyMsg)

	log.Println("ValidateAbhiPageHSToken(+) " + req.Method)

	switch req.Method {
	case "GET":
		var authInputVar abhiInputStruct
		var wallreturnVar abhiReturnStruct
		// err := json.Unmarshal(bodyMsg, &authInputVar)
		// if err != nil {
		// 	log.Println("Unable to Unmarshal request")
		// 	log.Println(err)
		// 	wallreturnVar.Status = "E"
		// 	wallreturnVar.Emsg = err.Error()
		// } else {

		fullpath := req.URL.Path + "?" + req.URL.RawQuery
		//parse parameter values
		u, err := url.Parse(fullpath)
		if err != nil {
			wallreturnVar.Status = "E"
			wallreturnVar.Emsg = "Error: " + err.Error()
			log.Fatal(err)
		}
		//get parameter values

		q := u.Query()
		if q.Get("token") != "undefined" && q.Get("token") != "" {
			authInputVar.Token = q.Get("token")

		}
		if q.Get("clientId") != "undefined" && q.Get("clientId") != "" {
			authInputVar.ClientID = q.Get("clientId")

		}
		//-------------------------------------------------------------------
		//Login AUTH API method
		//-------------------------------------------------------------------
		// log.Println("AuthInputVar Response: ", authInputVar)
		lLoginResp, lErr := ValidateLoginForEmail(authInputVar)
		// log.Println("lLoginResp", lLoginResp)
		if lErr != nil {
			wallreturnVar.Status = "E"
			wallreturnVar.Emsg = lErr.Error()
		} else {
			// log.Println("ValidateLoginEmail Response: ", lLoginResp)
			if lLoginResp.Status == "SUCCESS" {
				lAppName, lErr := brokers.GetAppName(lBrokerId)
				if lErr != nil {
					wallreturnVar.Status = "E"
					wallreturnVar.Emsg = lErr.Error()
				} else {

					// validatedMsg := "TRUE" //ValidateHSToken(authInputVar.ClientID, authInputVar.Key)
					log.Println("App Name: ", lAppName)
					// validatedMsg := appsso.IsAuthTokenValid(lAppName, authInputVar.Token, authInputVar.ClientID)
					// log.Println("validatedMsg", validatedMsg)
					//pending
					// validatedMsg = "Y"
					// if validatedMsg == "Y" {
					db, err := ftdb.LocalDbConnect(ftdb.SSODB)
					//db, err := util.Getdb(config.Database.DbType, config.Database)
					if err != nil {
						log.Println(err)
						wallreturnVar.Status = "E"
						wallreturnVar.Emsg = err.Error()
					} else {
						defer db.Close()

						lDomain, lErr := brokers.GetDomain(lBrokerId)
						if lErr != nil {
							wallreturnVar.Status = "E"
							wallreturnVar.Emsg = err.Error()
						} else {

							//generate session
							wallreturnVar.Token = GenerateAbhiSessionToken(db, authInputVar.ClientID, lAppName, reqDtl, lBrokerId)
							// log.Println("token", wallreturnVar.Token)
							if wallreturnVar.Token == "" {
								//wallreturnVar.Emsg = "Something went wrong (0x133). please try after sometime."
								wallreturnVar.Status = "I"
							} else {
								// ABHIDomain := "novoadmin.flattrade.in"

								//cookie := &http.Cookie{Name: "ftac_pt", Value: wallreturnVar.Token, MaxAge: 3000, Domain: ".flattrade.in", Path: "/", HttpOnly: true, Secure: true}
								//cookie := http.Cookie{Name: common.WallCookieName, Value: wallreturnVar.Token, MaxAge: 21600, HttpOnly: true, Secure: true, Path: "/", Domain: common.WallDomain, SameSite: http.SameSiteNoneMode}
								// cookie := http.Cookie{Name: common.ABHICookieName, Value: wallreturnVar.Token, MaxAge: 21600, HttpOnly: true, Secure: false, Path: "/", Domain: common.ABHIDomain, SameSite: 0}
								cookie := http.Cookie{Name: common.ABHICookieName, Value: wallreturnVar.Token, MaxAge: 21600, HttpOnly: true, Secure: false, Path: "/", Domain: common.ABHIDomain, SameSite: 0}

								log.Println("setting cookie", cookie, lDomain)

								http.SetCookie(w, &cookie)
								wallreturnVar.Status = "S"

								authInputVar.ClientID = lLoginResp.LoginId

								//added by naveen for insert token with source (new code)
								//This method inserting the token in db
								//lErr := InsertToken(lLoginResp, wallreturnVar.Token)
								source := common.Web //from mobile or web
								lErr := InsertToken(lLoginResp, wallreturnVar.Token, source)
								log.Println("Insert token", lErr)
								if lErr != nil {
									wallreturnVar.Status = "E"
									wallreturnVar.Emsg = err.Error()
								}
							}
						}
					}
				}
			} else {
				wallreturnVar.Status = "I"
			}
		}
		// } else {
		// 	// wallreturnVar.Emsg = "Something went wrong (0x135). please try after sometime."
		// 	wallreturnVar.Status = "I"
		// }
		//here
		data, err := json.Marshal(wallreturnVar)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(data))
	}
	// w.WriteHeader(200)
	log.Println("ValidateAbhiPageHSToken(-)")

}

//---------------------------------------------------------------------------------
// function to log out from ledger
// --------------------------------------------------------------------------------

func DeleteTokenFromCookies(w http.ResponseWriter, req *http.Request) {
	lBrokerId := 0
	var logRec logout

	origin := req.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			BrokerId, lErr1 := brokers.GetBrokerId(origin)
			if lErr1 != nil {
				logRec.Status = "E"
			} else {
				lBrokerId = BrokerId
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	log.Println("DeleteTokenFromCookies(+) " + req.Method)

	switch req.Method {
	case "GET":

		logRec.Status = "E"
		lDomain, err := brokers.GetDomain(lBrokerId)
		if err != nil {
			logRec.Status = "E"
		} else {
			log.Println("lDomain", lDomain)
			// cookie := http.Cookie{Name: common.ABHICookieName, Value: "", MaxAge: -1, HttpOnly: true, Secure: true, Path: "/", Domain: lDomain, SameSite: 0}

			cookie := http.Cookie{Name: common.ABHICookieName, Value: "", MaxAge: -1, HttpOnly: true, Secure: false, Path: "/", Domain: common.ABHIDomain, SameSite: 0}
			http.SetCookie(w, &cookie)

			// Clientcookie := http.Cookie{Name: common.ABHIClientCookieName, Value: "", MaxAge: -1, HttpOnly: true, Secure: true, Path: "/", Domain: lDomain, SameSite: 0}
			Clientcookie := http.Cookie{Name: common.ABHIClientCookieName, Value: "", MaxAge: -1, HttpOnly: true, Secure: false, Path: "/", Domain: common.ABHIDomain, SameSite: 0}

			http.SetCookie(w, &Clientcookie)

			logRec.Status = "S"
		}
		data, err := json.Marshal(logRec)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(data))

		log.Println("DeleteTokenFromCookies(-)")
	}
}

type Request struct {
	App_Key    string `json:"app_key"`
	LoginId    string `json:"LoginId"`
	Token      string `json:"token"`
	Secret_Key string `json:"secret_key"`
}

type Response struct {
	LoginId string `json:"LoginId"`
	Name    string `json:"Name"`
	Email   string `json:"Email"`
	Phone   string `json:"Phone"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

//---------------------------------------------------------------------------------
// LOGIN AUTH API
// --------------------------------------------------------------------------------

func ValidateLoginForEmail(authInput abhiInputStruct) (Response, error) {
	log.Println("ValidateLoginForEmail (+)")
	var lHeaderRec []apiUtil.HeaderDetails
	var lLoginResp Response

	lConfigFile := common.ReadTomlConfig("toml/loginCredentials.toml")
	Url := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NovoAuthLogin"])
	lKey := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NovoAuthSecretKey"])

	lApiKey := ""
	lDecoded, lErr1 := common.DecodeToString(lKey)
	if lErr1 != nil {
		log.Println("Error in Decoding KEY", lErr1)
		return lLoginResp, lErr1
	} else {
		lApiKey = lDecoded
	}

	// Url := "https://authapi.flattrade.in/oauth/authorize"
	// lApiKey := "2022.764f78bfaa3ddsf84j5a5acfg678desdfdsf548c23432sd46"

	inputString := common.ABHIAppName + authInput.Token + lApiKey

	// convert the string to sha256 format
	lEncodedString := common.EncodeToSHA256(inputString)

	var lApiReq Request
	lApiReq.Token = authInput.Token
	lApiReq.LoginId = authInput.ClientID
	lApiReq.App_Key = common.ABHIAppName
	lApiReq.Secret_Key = lEncodedString

	lRequest, lErr2 := json.Marshal(lApiReq)
	if lErr2 != nil {
		log.Println("ValidateLoginForEmail (-)", lErr2)
		return lLoginResp, lErr2
	}
	lString := string(lRequest)

	Response, lErr3 := apiUtil.Api_call(Url, "POST", lString, lHeaderRec, "novodev")
	if lErr3 != nil {
		log.Println("ValidateLoginForEmail (-)", lErr3)
		return lLoginResp, lErr3

	} else {
		lErr4 := json.Unmarshal([]byte(Response), &lLoginResp)
		if lErr4 != nil {
			log.Println("ValidateLoginForEmail (-)", lErr4)
			return lLoginResp, lErr4
		}
	}
	log.Println("ValidateLoginForEmail (-)")
	return lLoginResp, nil
}

//---------------------------------------------------------------------------------
// Inserting The Created token and client details in novoToken Database
// --------------------------------------------------------------------------------

//added by naveen with parameter pSource
// func InsertToken(pRec Response, pToken string) error {
func InsertToken(pRec Response, pToken string, pSource string) error {
	log.Println("InsertToken (+)")

	db, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	//db, lErr1 := util.Getdb(config.Database.DbType, config.Database)
	if lErr1 != nil {
		log.Println(lErr1)
		return lErr1

	} else {
		defer db.Close()

		lSqlString := `insert into novo_token (Token,UserId,UserName,UserMailId,UserMobileNo,source,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
						values (?,?,?,?,?,?,?,now(),?,now())`

		_, lErr2 := db.Exec(lSqlString, pToken, pRec.LoginId, pRec.Name, pRec.Email, pRec.Phone, pSource, pRec.LoginId, pRec.LoginId)
		if lErr2 != nil {
			log.Println(lErr2)
			return lErr2
		}
	}
	log.Println("InsertToken (-)")
	return nil
}

/*
Purpose:This API is used to login the mobile apps of NOVO
Parameters:

    not Applicable

Response:

	==========
	*On Sucess
	==========
	{
		"stat": "OK"
	}
	=========
	!On Error
	=========
	{
		"stat": "NOT OK",
		"errMsg": "Error : Invalid UserId"
	}

Author: KAVYA DHARSHANI
Date: 09SEP2023
*/

func LoginValidationforWeb(w http.ResponseWriter, req *http.Request) {

	//
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	reqDtl := apigate.GetRequestorDetail(req)

	log.Println("LoginValidationforWeb (+)", req.Method)

	switch req.Method {
	case "POST":
		var lauthInput LoginMRequest
		var lWallRes LoginRespStruct
		var lLoginResp Response
		lBody, lErr1 := ioutil.ReadAll(req.Body)
		// log.Println(string(lBody))
		if lErr1 != nil {
			log.Println("LLFM01", lErr1)
			lWallRes.Status = common.ErrorCode
			lWallRes.Errmsg = lErr1.Error()
		} else {
			// Unmarshal the request body values in lReqRec variable
			lErr2 := json.Unmarshal(lBody, &lauthInput)
			if lErr2 != nil {
				log.Println("LLFM02", lErr2)
				lWallRes.Status = common.ErrorCode
				lWallRes.Errmsg = lErr2.Error()
			} else {

				// lEncodepwd := common.EncodeToSHA256(lauthInput.Password)
				// lauthInput.Password = lEncodepwd
				lJDataResp, lErr3 := validateLoginDetails(lauthInput)
				if lErr3 != nil {
					log.Println("LLFM03", lErr3)
					lWallRes.Status = common.ErrorCode
					lWallRes.Errmsg = lErr3.Error()
				} else {
					if lJDataResp.Stat == "Ok" {
						db, lErr4 := ftdb.LocalDbConnect(ftdb.SSODB)
						if lErr4 != nil {
							log.Println("LLFM04", lErr4)
							lWallRes.Status = common.ErrorCode
							lWallRes.Errmsg = lErr4.Error()
						} else {
							defer db.Close()
							lLoginResp.Email = lJDataResp.EMail
							lLoginResp.LoginId = lJDataResp.Actid
							lLoginResp.Message = lJDataResp.Emsg
							lLoginResp.Name = lJDataResp.Uname
							lLoginResp.Status = lJDataResp.Stat
							//generate session
							Token := GenerateAbhiSessionToken(db, lLoginResp.LoginId, common.ABHIAppName, reqDtl, 0)
							if Token == "" {
								lWallRes.Status = common.LoginFailure
							} else {
								cookie := http.Cookie{Name: common.ABHICookieName, Value: Token, MaxAge: 21600, HttpOnly: true, Secure: false, Path: "/", Domain: common.ABHIDomain, SameSite: 0}
								http.SetCookie(w, &cookie)
								lWallRes.ClientId = lJDataResp.Actid
								lWallRes.Status = common.SuccessCode

								//added by naveen for insert token with source
								//This method inserting the token in db new code
								//lErr4 := InsertToken(lLoginResp, Token)
								source := common.Mobile //from mobile or web
								lErr4 := InsertToken(lLoginResp, Token, source)
								if lErr4 != nil {
									log.Println("LLFM04", lErr4)
									lWallRes.Status = common.ErrorCode
									lWallRes.Errmsg = lErr4.Error()
								}
							}
						}
					} else {
						lWallRes.Status = common.LoginFailure
						lWallRes.Errmsg = lJDataResp.Emsg
					}
				}
				data, lErr5 := json.Marshal(lWallRes)
				if lErr5 != nil {
					log.Println(lErr5)
					lWallRes.Status = common.ErrorCode
					lWallRes.Errmsg = lErr5.Error()
				}
				fmt.Fprintf(w, string(data))
			}
		}
		log.Println("LoginValidationforWeb (-)", req.Method)
	}
}

func validateLoginDetails(ploginuser LoginMRequest) (JDataRespStruct, error) {
	log.Println("validateLoginDetails (+)")
	var JDataReq JData
	lLoginFile := common.ReadTomlConfig("toml/loginCredentials.toml")
	JDataReq.Apkversion = fmt.Sprintf("%v", lLoginFile.(map[string]interface{})["apkversion"])
	JDataReq.Imei = fmt.Sprintf("%v", lLoginFile.(map[string]interface{})["imei"])
	JDataReq.Source = fmt.Sprintf("%v", lLoginFile.(map[string]interface{})["source"])
	JDataReq.Vc = fmt.Sprintf("%v", lLoginFile.(map[string]interface{})["vc"])
	JDataReq.Appkey = fmt.Sprintf("%v", lLoginFile.(map[string]interface{})["appkey"])
	Url := fmt.Sprintf("%v", lLoginFile.(map[string]interface{})["url"])

	inputString := ploginuser.ClientId + "|" + JDataReq.Appkey

	// convert the string to sha256 format
	lEncodedString := common.EncodeToSHA256(inputString)
	JDataReq.Appkey = lEncodedString
	JDataReq.Pwd = ploginuser.Password
	JDataReq.Factor2 = ploginuser.Otp
	JDataReq.Uid = ploginuser.ClientId
	var lHeaderRec []apiUtil.HeaderDetails
	var lJDataResp JDataRespStruct

	lRequest, lErr1 := json.Marshal(JDataReq)
	if lErr1 != nil {
		log.Println("LVLD01", lErr1)
		return lJDataResp, lErr1
	} else {
		lString := "jData=" + string(lRequest)
		Response, lErr2 := apiUtil.Api_call(Url, "POST", lString, lHeaderRec, "novodev")
		if lErr2 != nil {
			log.Println("LVLD02", lErr2)
			return lJDataResp, lErr2
		} else {
			log.Println("Response", Response)
			lErr3 := json.Unmarshal([]byte(Response), &lJDataResp)
			if lErr3 != nil {
				log.Println("LVLD03", lErr3)
				return lJDataResp, lErr3
			}
		}
	}

	log.Println("validateLoginDetails (-)")
	return lJDataResp, nil
}
