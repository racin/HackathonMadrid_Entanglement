package main

import (
	"fmt"
	"github.com/racin/HackathonMadrid_Entanglement/Code/Entangler"
	"github.com/racin/HackathonMadrid_Entanglement/Code/SwarmConnector"
	"io/ioutil"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(10 << 40)

	file, handler, err := r.FormFile("myFile")
	defer file.Close()
	if err != nil {
		fmt.Println("FATAL")
		return
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "OK")

	if _, err := os.Create(Entangler.TempDirectory + handler.Filename); err == nil {

	} else {
		fmt.Println("Fatal error ... " + err.Error())
		os.Exit(1)
	}
	ioutil.WriteFile(Entangler.TempDirectory+handler.Filename, fileBytes, os.ModeAppend)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	err := http.ListenAndServe(":8081", nil)
	fmt.Println(err.Error())
}

func main() {
	fmt.Println("Hello World")
	//setupRoutes()
	SwarmConnector.UploadAllChunks()
}
