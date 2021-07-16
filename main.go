package main

import (
	"aws_service/aws"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SendMailInAWS(w http.ResponseWriter, r *http.Request) {
	
	//using SES Service
	subject := "메일 제목/ mail subject"
	context := "메일 내용/ mail description"

	aws.AmazonSES(subject, context)

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result" :"success"})
	return
}

func SingleUploadToS3(w http.ResponseWriter, r *http.Request) {

	get_file_name := aws.SingleUploadFile(r)

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result" :get_file_name})
	return
}

func SingleDeleteToS3(w http.ResponseWriter, r *http.Request) {

	err := aws.SingleDeleteFile("delete filename")

	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result" :"success"})
	return
}

func MultiFileUploadToS3(w http.ResponseWriter, r *http.Request) {
	
	get_multi_array_filename := aws.MultiFileUpload(r)

	
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result" : get_multi_array_filename})
	return
}

func MultiFileDeleteToS3(w http.ResponseWriter, r *http.Request) {

	err := aws.MultiFileDeleter("delete image array set")

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result" :"success"})
	return
}

func main() {


	
	r := mux.NewRouter()

	r.HandleFunc("/ses", SendMailInAWS)
	r.HandleFunc("/singleupload", SingleUploadToS3)
	r.HandleFunc("/singledelete", SingleDeleteToS3)
	r.HandleFunc("/multiupload", MultiFileUploadToS3)
	r.HandleFunc("/multidelete", MultiFileDeleteToS3)
	http.ListenAndServe(":8081", r)

}