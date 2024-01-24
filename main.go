package main

import (
	"crypto/tls"
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
	"fcs23pkg/env"
	logfiledownload "fcs23pkg/logFileDownload"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

func autoRestart() {

	for {
		now := time.Now()
		//resart the program everyday at 4am
		//at 3am, the program goes for 1 hour sleep and after that it will restart
		if now.Hour() == 3 {
			//sleep for an hour so that the hour changes to 4 and this condition
			//in the loop does not  continue again in next iteration
			time.Sleep(60 * 61 * time.Second)
			fmt.Println(now.Hour(), now.Minute(), now.Second())
			log.Println(now.Hour(), now.Minute(), now.Second())
			// Restart the program
			fmt.Println("Restarting the program...")
			log.Println("Restarting the program...")
			execPath, err := os.Executable()
			if err != nil {
				fmt.Println("Error getting executable path:", err)
				log.Println("Error getting executable path:", err)

				return
			}
			cmd := exec.Command(execPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Start()
			if err != nil {
				fmt.Println("Error restarting program:", err)
				log.Println("Error restarting program:", err)
				return
			}
			os.Exit(0)

		}
		time.Sleep(60 * 30 * time.Second)
	}
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

	go autoRestart()

	tlsConfig := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		//GetCertificate: m.GetCertificate,
	}

	pemCert, err := ioutil.ReadFile("flattrade.crt")
	log.Println(err)
	pemKey, err := ioutil.ReadFile("flattrade.key")
	log.Println(err)
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], err = tls.X509KeyPair(pemCert, pemKey)
	log.Println(err)

	router := mux.NewRouter()

	//endpoints used to generate application session
	//=============================Login==================================================
	router.HandleFunc("/token", abhilogin.ValidateAbhiPageHSToken).Methods("GET", "OPTIONS")
	router.HandleFunc("/tokenValidation", abhilogin.ValidateAbhiToken).Methods("GET", "OPTIONS")
	//=============================Login==================================================

	//=============================Logout=================================================
	router.HandleFunc("/logout", abhilogin.DeleteTokenFromCookies).Methods("GET", "OPTIONS")

	//=============================Logout=================================================

	//TO know if the user is admin
	router.HandleFunc("/verifyClient", menu.VerifyAccess).Methods("GET", "OPTIONS") //<----- Changed adminaccess to menu
	// To add user in userauth database
	router.HandleFunc("/addUser", adminaccess.Adduser).Methods("POST", "OPTIONS")
	//to change password
	router.HandleFunc("/changePassword", adminaccess.ChangePassword).Methods("POST", "OPTIONS")
	// To get admin details in userauth database
	router.HandleFunc("/getAdminList", adminaccess.GetAdminList).Methods("GET", "OPTIONS")
	//UpiNames
	//commented by pavithra
	// router.HandleFunc("/getUpi", localdetail.GetUpi).Methods("GET", "OPTIONS")
	//History of IPO
	router.HandleFunc("/getIpoHistory", localdetail.GetIpoHistory).Methods("GET", "OPTIONS")
	//History Record
	router.HandleFunc("/getHistoryRecords", localdetail.GetHistoryRecords).Methods("GET", "OPTIONS")
	//Active IPOs Details
	router.HandleFunc("/getIpoMaster", localdetail.GetIpoMaster).Methods("GET", "OPTIONS")
	//Placing Order
	router.HandleFunc("/placeOrder", placeorder.PlaceOrder).Methods("POST", "OPTIONS")
	//GetModify Details
	router.HandleFunc("/getModify", localdetail.GetModifyDetails).Methods("GET", "OPTIONS")
	//Get Reports
	router.HandleFunc("/getReport", localdetail.GetReport).Methods("POST", "OPTIONS")
	// Get Report initially
	router.HandleFunc("/getDefaultReport", localdetail.DefaultReport).Methods("POST", "OPTIONS")
	//Manual Fetch
	router.HandleFunc("/getManual", localdetail.ManualFetch).Methods("GET", "OPTIONS")

	//---------------------------CallBack End Points--------------------------------
	//CallBack For NSE IPO
	router.HandleFunc("/v1/appdpstatus", callback.AppDpStatus).Methods("POST", "OPTIONS")
	//CallBack For Upi
	router.HandleFunc("/v1/apppaystatus", callback.AppPayStatus).Methods("POST", "OPTIONS")
	//---------------------------Schedule End Points--------------------------------
	//Fetch Master Schedule
	router.HandleFunc("/fetchMasterSch", scheduler.FetchIpoMasterSch).Methods("GET", "OPTIONS")
	//Pffline Application Schedule
	router.HandleFunc("/offlineProcessSch", scheduler.OfflineScheduler).Methods("GET", "OPTIONS")
	//----------------------------- July 20 ----------------------------------->>>>>>>>>>>>>>
	router.HandleFunc("/deleteCookie", abhilogin.DeleteTokenFromCookies).Methods("GET", "OPTIONS")
	router.HandleFunc("/validRouter", routeraccess.RouterAccess).Methods("GET", "OPTIONS")
	//<<<<<<<<<<<<------------------July 26-->>>>>>>>>>>>>>>
	//get blog links details
	router.HandleFunc("/getBlogLink", bloglinks.GetBlogLink).Methods("GET", "OPTIONS")
	//add new or existing blog links
	router.HandleFunc("/addBlogLink", bloglinks.AddBlogLinks).Methods("POST", "OPTIONS")

	//CallBack For BSE IPO
	router.HandleFunc("/v1/dpstatus", callback.BseDpStatus).Methods("POST", "OPTIONS")
	//CallBack For Upi
	router.HandleFunc("/v1/upistatus", callback.BsePayStatus).Methods("POST", "OPTIONS")
	// -----------------------Login for mobile----------------------------------
	router.HandleFunc("/webAuthlogin", abhilogin.LoginValidationforWeb).Methods("POST", "OPTIONS")
	// ----------------------- Aug3 ----------------------------------
	// origin := apigate.GetAllowOrigin()
	// log.Println("Origin: ", origin)
	// common.ABHIAllowOrigin = origin
	//<<<<<<<<<<<<------------------Aug 9-->>>>>>>>>>>>>>>
	//get sgb master details
	router.HandleFunc("/getSgbMaster", localdetails.GetSgbMaster).Methods("GET", "OPTIONS")
	// SetMemberDetail
	router.HandleFunc("/setMemberDetail", memberdetail.SetMemberDetails).Methods("POST", "OPTIONS")
	// GetMemberDetail
	router.HandleFunc("/getMemberDetail", memberdetail.GetMemberDetail).Methods("GET", "OPTIONS")
	//---------------------------------- AUG 16 ------------------------
	//get menu items
	router.HandleFunc("/getSubMenu", menu.GetSubMenu).Methods("GET", "OPTIONS")
	//---------------------------------- AUG 17 -------------------------
	//sgb place order
	router.HandleFunc("/sgb/placeOrder", sgbplaceorder.SgbPlaceOrder).Methods("POST", "OPTIONS")
	//-------------------------------- AUG 18 ---------------------------
	//get the client fund detail
	router.HandleFunc("/sgb/fetchFund", sgbplaceorder.FetchClientFund).Methods("GET", "OPTIONS")
	//----------------------------- AUG 22 -------------------------------
	//get the sgb order modify datas
	router.HandleFunc("/sgb/getModify", localdetails.GetSgbModifyDetail).Methods("GET", "OPTIONS")
	//get the sgb order history
	router.HandleFunc("/sgb/getSgbHistory", localdetails.GetSgbOrderHistory).Methods("GET", "OPTIONS")
	//--------------------SGB Schedule End Points-------------------SEP 01-------------
	//SGB master fetching from Exchange
	// commented by pavithra
	// router.HandleFunc("/sgb/fetchMaster", sgbschedule.FetchSgbMasterSch).Methods("GET", "OPTIONS")
	//------------------RoleMaster and TaskMaster End Points-----------SEP 01-------------
	router.HandleFunc("/GetRoleTask", roleTask.GetRoleTask).Methods("GET", "OPTIONS")
	router.HandleFunc("/setRole", roleTask.SetRole).Methods("PUT", "OPTIONS")
	router.HandleFunc("/setTask", roleTask.SetTask).Methods("PUT", "OPTIONS")
	router.HandleFunc("/GetRoleTaskMaster", roleTask.GetRoleTaskMaster).Methods("GET", "OPTIONS")
	router.HandleFunc("/setRoleTaskMaster", roleTask.SetRoleTaskMaster).Methods("PUT", "OPTIONS")
	//<<<<<<<<<<<<------------------ July 27 -->>>>>>>>>>>>>>>
	// commented by nithish to change the valid endpoint name based on standards
	// router.HandleFunc("/getDomainList", brokers.GetBrokerList).Methods("GET", "OPTIONS")
	router.HandleFunc("/getBrokerList", brokers.GetBrokerList).Methods("GET", "OPTIONS")

	router.HandleFunc("/setDomain", brokers.SetBroker).Methods("POST", "OPTIONS")

	// router.HandleFunc("/getBrokerList", adminaccess.GetAdminList).Methods("GET", "OPTIONS")

	router.HandleFunc("/getMemberUser", adminaccess.GetMemberUser).Methods("GET", "OPTIONS")
	//<<<<<<<<<<<<------------------ Sep 12 -->>>>>>>>>>>>>>>
	router.HandleFunc("/getRedirectUrl", abhilogin.GetRedirectURL).Methods("GET", "OPTIONS")
	//<<<<<<<<<<<<----------To Get The ClientName Sep 20 -->>>>>>>>>>>>>>>
	router.HandleFunc("/getClientName", client.GetClientName).Methods("GET", "OPTIONS")
	//================================================================
	router.HandleFunc("/dashboard", dashboard.GetDashboardDetail).Methods("GET", "OPTIONS")
	//============================= 16 OCT 2023 ===================================
	router.HandleFunc("/LookUpHeader", lookup.GetLookUpHeader).Methods("GET", "OPTIONS")
	router.HandleFunc("/LookUpDetails", lookup.GetLookUpDetails).Methods("PUT", "OPTIONS")
	router.HandleFunc("/InsertHeader", lookup.AddLookUpHeader).Methods("PUT", "OPTIONS")
	router.HandleFunc("/InsertDetails", lookup.AddLookUpDetails).Methods("PUT", "OPTIONS")
	router.HandleFunc("/UpdateHeader", lookup.UpdateLookUpHeader).Methods("PUT", "OPTIONS")
	router.HandleFunc("/UpdateDetails", lookup.UpdateLookUpDetails).Methods("PUT", "OPTIONS")
	//================================================================
	// commented by pavithra-changed to POST method
	// router.HandleFunc("/sgbEndDateSch", scheduler.SgbEndDate).Methods("GET", "OPTIONS")
	router.HandleFunc("/sgbEndDateSch", scheduler.SgbEndDate).Methods("POST", "OPTIONS")

	// endpoint for SGB to fetch application status
	// commented by pavithra-changed to POST method
	// router.HandleFunc("/sgbFetchStatus", scheduler.SgbStatusScheduler).Methods("GET", "OPTIONS")
	router.HandleFunc("/sgbFetchStatus", scheduler.SgbStatusScheduler).Methods("POST", "OPTIONS")

	router.HandleFunc("/getCurVersion", versioncontrol.GetCurVersion).Methods("GET", "OPTIONS")
	router.HandleFunc("/setAppVersion", versioncontrol.SetAppVersion).Methods("PUT", "OPTIONS")
	router.HandleFunc("/getAllVersions", versioncontrol.GetAllVersions).Methods("GET", "OPTIONS")

	// Setting Environment
	env.SetNovoEnvironment()

	// Manually run OfflineScheduler from screen
	router.HandleFunc("/manualOfflineSch", localdetail.ManualOfflineFetch).Methods("GET", "OPTIONS")

	//getCategory for Placing IPO order
	router.HandleFunc("/getCategory", localdetail.GetCategory).Methods("GET", "OPTIONS")

	// GET category purchased Falg for IPO order
	router.HandleFunc("/getCategoryPurFlag", localdetail.GetCategoryPurFlag).Methods("GET", "OPTIONS")

	router.HandleFunc("/LogDownload", logfiledownload.LogFileDownload).Methods("GET", "OPTIONS")

	// ===========================prashanth=============================11 Dec 23 ===============================================
	router.HandleFunc("/getRegistrarDetails", registrardetails.GetRegistrarDetails).Methods("GET", "OPTIONS")
	router.HandleFunc("/setRegisterDetails", registrardetails.SetRegisterDetails).Methods("PUT", "OPTIONS")

	// Get the IPO's current market demand scheduler
	router.HandleFunc("/fetchIpoMktData", scheduler.FetchIpoMktDataSch).Methods("GET", "OPTIONS")
	// Get the fetched demand details to the screen
	router.HandleFunc("/getIpoMktData", localdetail.GetIpoMKtData).Methods("GET", "OPTIONS")

	//===================================== NCB ============================================================
	router.HandleFunc("/getNcbMaster", ncblocaldetails.GetNcbMaster).Methods("GET", "OPTIONS")
	router.HandleFunc("/getNcbOrderHistory", ncblocaldetails.GetNcbOrderHistory).Methods("GET", "OPTIONS")
	router.HandleFunc("/ncb/getNcbModify", ncblocaldetails.GetNcbModifyDetail).Methods("GET", "OPTIONS")
	router.HandleFunc("/ncb/ncbPlaceOrder", ncbplaceorder.NcbPlaceOrder).Methods("POST", "OPTIONS")
	router.HandleFunc("/NcbEndDateSch", scheduler.NcbEndDate).Methods("GET", "OPTIONS")
	http.HandleFunc("/ncbFetchStatus", scheduler.NcbStatusScheduler)
	router.HandleFunc("/ncb/fetchFund", ncbplaceorder.NcbFetchClientFund).Methods("GET", "OPTIONS")

	// router.Use(apigate.RequestMiddleware)
	// handler := apigate.ResponseMiddleware(router, apigate.CustomWriteTimeout)
	handler := apigate.LogMiddleware(router)

	srv := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      handler,
		Addr:         ":29091",
	}
	log.Println("Server Started")
	log.Println(srv.ListenAndServeTLS("", ""))

}
