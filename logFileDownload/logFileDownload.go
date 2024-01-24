package logfiledownload

import (
	"encoding/json"
	"fcs23pkg/common"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Response struct {
	Status   string `json:"status"`
	ErrMsg   string `json:"errMsg"`
	FileName string `json:"fileName"`
	Content  string `json:"content"`
}

// Purpose:This Api is used to Download the Current Log File
// Response:

// On Sucess
//{
//  Status:"S"
// ErrMsg:""
//}

// On Error
//{
//  Status:"E"
// ErrMsg:Error
//}
//Author: Prabhaharan S
//Date: 13JUN2023
func LogFileDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("LogFileDownload+")
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "GET" {

		var resp Response
		resp.Status = "S"

		// config := genpkg.ReadTomlConfig("./toml/fcsfourtysevenconfig.toml")
		// dirPath := fmt.Sprintf("%v", config.(map[string]interface{})["logfilepath"])
		// dirPath := "/home/dev-51/Documents/CODE/Flow_code/FCS_94_FlowAPI/log" // Replace with the actual directory path
		dirPath := "./log" // Replace with the actual directory path

		// Get all file paths in the directory
		log.Println("dirPath", dirPath)
		filePaths, err := filepath.Glob(filepath.Join(dirPath, "*"))
		if err != nil {
			log.Println(err)
			resp.ErrMsg = "Error" + err.Error()
			resp.Status = "E"

		} else {

			var lastModifiedFile string
			var lastModifiedTime time.Time

			// Iterate through each file path
			for _, filePath := range filePaths {
				// Get file information
				fileInfo, err := os.Stat(filePath)
				if err != nil {
					log.Println(err)
					resp.ErrMsg = "Error" + err.Error()
					resp.Status = "E"
					return
				}

				// Check if the file's modification time is later than the last modified time found
				if fileInfo.ModTime().After(lastModifiedTime) {
					lastModifiedFile = filePath
					lastModifiedTime = fileInfo.ModTime()
				}
			}

			if lastModifiedFile != "" {
				// log.Println("Last Modified File:", lastModifiedFile)
				var split []string
				CheckThisDomain := strings.ToLower(common.ABHIDomain)
				if strings.Contains(CheckThisDomain, "localhost") {
					split = strings.Split(lastModifiedFile, `log/`) // local

				} else {
					// fmt.Println("String does not contain the value.")
					split = strings.Split(lastModifiedFile, `log\`) // prod
				}

				// log.Println("split", split)
				resp.FileName = split[1]
				// log.Println("resp.FileName ", resp.FileName)
				ContentArr, err := ReadRawFileFromPath(lastModifiedFile)

				if err != nil {
					log.Println(err)
					resp.ErrMsg = "Error" + err.Error()
					resp.Status = "E"

				} else {
					// w.write(string(ContentArr))
					resp.Content = string(ContentArr)
					// log.Println("Response", resp)
					// err := DownloadFile(ContentArr, FileName)
					// if err != nil {
					// 	resp.Status = "E"
					// 	resp.ErrMsg = "Error" + err.Error()
					// 	return
					// }
				}

			} else {

				resp.Status = "E"
				resp.ErrMsg = "No files found in the directory."
			}
		}

		//log.Println(resp)
		data, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "Error taking data"+err.Error())
		} else {
			fmt.Fprintf(w, string(data))
		}
		log.Println("LogFileDownload-")
	}
}

// Purpose:This method is used to Read a file
// Parameter:filePath
// Response:[]byte,error
//Author: Prabhaharan S
//Date: 13JUN2023
func ReadRawFileFromPath(filePath string) ([]byte, error) {
	log.Println("ReadRawFileFromPath+")
	log.Println("filePath:", filePath)
	//Read file from the physical path
	dat, err := os.ReadFile(filePath)
	if err != nil {

		return dat, err
	}
	log.Println("ReadRawFileFromPath-")
	return dat, nil
}
