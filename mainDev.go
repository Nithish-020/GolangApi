package main

import (
	"fcs23pkg/apigate"
	dashboard "fcs23pkg/apps/DashBoard"
	registrardetails "fcs23pkg/apps/Ipo/RegistrarDetails"
	"fcs23pkg/apps/Ipo/bloglinks"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/Ipo/brokers/memberdetail"
	"fcs23pkg/apps/Ipo/callback"
	"fcs23pkg/apps/Ipo/localdetail"
	"fcs23pkg/apps/Ipo/placeorder"
	lookup "fcs23pkg/apps/Lookup"
	ncblocaldetails "fcs23pkg/apps/Ncb/ncbLocaldetails"
	"fcs23pkg/apps/Ncb/ncbplaceorder"
	"fcs23pkg/apps/SGB/localdetails"
	"fcs23pkg/apps/SGB/sgbplaceorder"
	"fcs23pkg/apps/abhilogin"
	"fcs23pkg/apps/roleTask"
	"fcs23pkg/apps/scheduler"
	client "fcs23pkg/apps/validation/Client"
	"fcs23pkg/apps/validation/adminaccess"
	"fcs23pkg/apps/validation/menu"
	"fcs23pkg/apps/validation/routeraccess"
	"fcs23pkg/apps/versioncontrol"
	"fcs23pkg/common"
	"fcs23pkg/env"
	logfiledownload "fcs23pkg/logFileDownload"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		requestorDetail := apigate.GetRequestorDetail(r)
		requestorDetail.Body = string(body)
		apigate.LogRequest("", requestorDetail, "")
		next.ServeHTTP(w, r)
	})
}

// --------------------------------------------------------------------
// main function executed from command
// --------------------------------------------------------------------
func main() {
	log.Println("Server Started")
	f, err := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	//=============================Login==================================================
	http.HandleFunc("/token", abhilogin.ValidateAbhiPageHSToken)

	http.HandleFunc("/tokenValidation", abhilogin.ValidateAbhiToken)

	//=============================Login==================================================

	//=============================Logout=================================================
	http.HandleFunc("/logout", abhilogin.DeleteTokenFromCookies)
	//=============================Logout=================================================

	time := time.Now()
	log.Println("Today", time.Format("15:04:05"))

	//=============================Admin Process=================================================
	//TO know if the user is admin
	// http.HandleFunc("/verifyClient", adminaccess.VerifyAccess)
	http.HandleFunc("/verifyClient", menu.VerifyAccess) // Changed into menu.VerifyAccess

	// To add user in userauth database
	http.HandleFunc("/addUser", adminaccess.Adduser)

	//to change password
	http.HandleFunc("/changePassword", adminaccess.ChangePassword)

	// To get admin details in userauth database
	http.HandleFunc("/getAdminList", adminaccess.GetAdminList)

	//=============================Admin Process=================================================

	//UpiNames
	// http.HandleFunc("/getUpi", localdetail.GetUpi)

	//History of IPOs
	http.HandleFunc("/getIpoHistory", localdetail.GetIpoHistory)

	//History Record
	http.HandleFunc("/getHistoryRecords", localdetail.GetHistoryRecords)

	//Active IPOs Details
	// http.HandleFunc("/getIpoMaster", localdetail.GetIpoMaster)
	http.HandleFunc("/getIpoMaster", localdetail.GetIpoMaster2)

	//Placing Order
	http.HandleFunc("/placeOrder", placeorder.PlaceOrder)

	//GetModify Details
	http.HandleFunc("/getModify", localdetail.GetModifyDetails)

	//Get Reports
	http.HandleFunc("/getReport", localdetail.GetReport)

	// Get Report initially
	http.HandleFunc("/getDefaultReport", localdetail.DefaultReport)

	//Manual Fetch
	http.HandleFunc("/getManual", localdetail.ManualFetch)

	//---------------------------CallBack End Points--------------------------------

	//CallBack For NSE IPO
	http.HandleFunc("/v1/appdpstatus", callback.AppDpStatus)

	//CallBack For Upi
	http.HandleFunc("/v1/apppaystatus", callback.AppPayStatus)

	//---------------------------Schedule End Points--------------------------------

	//Fetch Master Schedule
	http.HandleFunc("/fetchMasterSch", scheduler.FetchIpoMasterSch)

	//offline Application Schedule
	http.HandleFunc("/offlineProcessSch", scheduler.OfflineScheduler)
	//----------------------------- July 20 -----------------------------------------

	//Delete a cookie
	http.HandleFunc("/deleteCookie", abhilogin.DeleteTokenFromCookies)

	//router access for api
	http.HandleFunc("/validRouter", routeraccess.RouterAccess)

	//-------------------------------------July 26---------------------

	//get blog links details
	http.HandleFunc("/getBlogLink", bloglinks.GetBlogLink)

	//add new or existing blog links
	http.HandleFunc("/addBlogLink", bloglinks.AddBlogLinks)
	//CallBack For BSE IPO
	http.HandleFunc("/v1/dpstatus", callback.BseDpStatus)
	//CallBack For Upi
	http.HandleFunc("/v1/upistatus", callback.BsePayStatus)
	// -----------------------Login for mobile----------------------------------
	http.HandleFunc("/webAuthlogin", abhilogin.LoginValidationforWeb)
	// ----------------------- Aug3 ----------------------------------

	origin := apigate.GetAllowOrigin()
	log.Println("Origin: ", origin)
	common.ABHIAllowOrigin = origin

	// Setting Environment
	env.SetNovoEnvironment()
	//<<<<<<<<<<<<------------------Aug 9-->>>>>>>>>>>>>>>

	//get sgb master details
	http.HandleFunc("/getSgbMaster", localdetails.GetSgbMaster)

	// SetMemberDetail
	http.HandleFunc("/setMemberDetail", memberdetail.SetMemberDetails)

	// GetMemberDetail
	http.HandleFunc("/getMemberDetail", memberdetail.GetMemberDetail)

	//---------------------------------- AUG 16 ------------------------
	//get menu items
	http.HandleFunc("/getSubMenu", menu.GetSubMenu)

	//---------------------------------- AUG 17 -------------------------
	//sgb place order
	http.HandleFunc("/sgb/placeOrder", sgbplaceorder.SgbPlaceOrder)

	//-------------------------------- AUG 18 ---------------------------
	//get the client fund detail
	http.HandleFunc("/sgb/fetchFund", sgbplaceorder.FetchClientFund)

	//----------------------------- AUG 22 -------------------------------
	//get the sgb order modify datas
	http.HandleFunc("/sgb/getModify", localdetails.GetSgbModifyDetail)

	//get the sgb order history
	http.HandleFunc("/sgb/getSgbHistory", localdetails.GetSgbOrderHistory)

	//--------------------SGB Schedule End Points-------------------SEP 01-------------
	//SGB master fetching from Exchange
	//commented by pavithra
	// http.HandleFunc("/sgb/fetchMaster", sgbschedule.FetchSgbMasterSch)
	//SGB download Orders from Exchange naveen
	// http.HandleFunc("/downloadStatus", sgbschedule.SgbDownStatusSch)

	// http.HandleFunc("/pendingSgbStatus", sgbschedule.SgbOfflineSch)
	//------------------RoleMaster and TaskMaster End Points-----------SEP 01-------------
	http.HandleFunc("/GetRoleTask", roleTask.GetRoleTask)
	http.HandleFunc("/setRole", roleTask.SetRole)
	http.HandleFunc("/setTask", roleTask.SetTask)
	http.HandleFunc("/GetRoleTaskMaster", roleTask.GetRoleTaskMaster)
	http.HandleFunc("/setRoleTaskMaster", roleTask.SetRoleTaskMaster)
	// http.HandleFunc("/bseoffline", iposchedule.BseOfflineSch)

	//<<<<<<<<<<<<------------------ July 27 -->>>>>>>>>>>>>>>

	// commented by nithish for the incorrect endpoint name
	// http.HandleFunc("/getDomainList", brokers.GetBrokerList)
	http.HandleFunc("/getBrokerList", brokers.GetBrokerList)

	http.HandleFunc("/setDomain", brokers.SetBroker)

	// commented by nithish for the incorrect endpoint name
	// http.HandleFunc("/getBrokerList", adminaccess.GetAdminList)
	// http.HandleFunc("/getAdminList", adminaccess.GetAdminList)

	http.HandleFunc("/getMemberUser", adminaccess.GetMemberUser)

	//<<<<<<<<<<<<------------------ Sep 12 -->>>>>>>>>>>>>>>

	http.HandleFunc("/getRedirectUrl", abhilogin.GetRedirectURL)

	//<<<<<<<<<<<<----------To Get The ClientName Sep 20 -->>>>>>>>>>>>>>>

	http.HandleFunc("/getClientName", client.GetClientName)

	// http.HandleFunc("/sgbFetchstatus", sgbschedule.NseSgbFetchStatus)

	//! Commented temprarly
	// http.HandleFunc("/offlineSch", iposchedule.IpoOfflineProcess)

	//--------------------------------------------------------------------

	// http.HandleFunc("/placeSgb", sgbschedule.PlacingSgbOrder)

	//================================================================
	http.HandleFunc("/dashboard", dashboard.GetDashboardDetail)

	//============================= 16 OCT 2023 ===================================

	http.HandleFunc("/LookUpHeader", lookup.GetLookUpHeader)
	http.HandleFunc("/LookUpDetails", lookup.GetLookUpDetails)
	http.HandleFunc("/InsertHeader", lookup.AddLookUpHeader)
	http.HandleFunc("/InsertDetails", lookup.AddLookUpDetails)
	http.HandleFunc("/UpdateHeader", lookup.UpdateLookUpHeader)
	http.HandleFunc("/UpdateDetails", lookup.UpdateLookUpDetails)

	//================================================================
	http.HandleFunc("/sgbEndDateSch", scheduler.SgbEndDate)

	// need to add in main
	http.HandleFunc("/sgbFetchStatus", scheduler.SgbStatusScheduler)

	//========================= 10 NOV 2023 =============================
	http.HandleFunc("/setDashboard", dashboard.SetDashboardDetail)

	// Manually run OfflineScheduler from screen
	http.HandleFunc("/manualOfflineSch", localdetail.ManualOfflineFetch)

	//getCategory for Placing IPO order
	http.HandleFunc("/getCategory", localdetail.GetCategory)

	// GET category purchased Falg for IPO order
	http.HandleFunc("/getCategoryPurFlag", localdetail.GetCategoryPurFlag)

	http.HandleFunc("/LogDownload", logfiledownload.LogFileDownload)

	// ===========================prashanth=============================11 Dec 23 ======================
	http.HandleFunc("/getRegistrarDetails", registrardetails.GetRegistrarDetails)
	http.HandleFunc("/setRegisterDetails", registrardetails.SetRegisterDetails)
	// =================================================================================================

	// =================================== 19 DEC 2023 =================================================
	http.HandleFunc("/fetchIpoMktData", scheduler.FetchIpoMktDataSch)
	http.HandleFunc("/getIpoMktData", localdetail.GetIpoMKtData)
	// =================================================================================================

	http.HandleFunc("/getCurVersion", versioncontrol.GetCurVersion)
	http.HandleFunc("/setAppVersion", versioncontrol.SetAppVersion)
	http.HandleFunc("/getAllVersions", versioncontrol.GetAllVersions)

	//============================ NCB ==========================================

	// http.HandleFunc("/ncb/ncbSchedule", ncbschedule.FetchNcbMasterSch)
	//get Ncb master details
	http.HandleFunc("/getNcbMaster", ncblocaldetails.GetNcbMaster)
	//get the Ncb order history
	http.HandleFunc("/getNcbOrderHistory", ncblocaldetails.GetNcbOrderHistory)
	//Ncb place order
	http.HandleFunc("/ncb/ncbPlaceOrder", ncbplaceorder.NcbPlaceOrder)
	//get the Ncb order modify datas
	http.HandleFunc("/ncb/getNcbModify", ncblocaldetails.GetNcbModifyDetail)
	// Scheduler
	http.HandleFunc("/ncbEndDateSch", scheduler.NcbEndDate)

	http.HandleFunc("/ncbFetchStatus", scheduler.NcbStatusScheduler)
	//get the client fund detail
	http.HandleFunc("/ncb/fetchFund", ncbplaceorder.NcbFetchClientFund)

	// http.Handle("/test", middleware(http.HandlerFunc(helloworld.Start)))
	// http.HandleFunc("/getToken", ipo.GetToken) //GetToken.go
	// http.HandleFunc("/fetchipomaster", ipo.FetchIPOmaster)
	// http.ListenAndServe(":29091", nil)

	// Using DefaultServeMux, you can also create your own mux if needed
	// handler := apigate.ResponseMiddleware(http.DefaultServeMux)
	handler := apigate.LogMiddleware(http.DefaultServeMux)

	srv := http.Server{
		Addr:    ":29091",
		Handler: handler,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}

}
