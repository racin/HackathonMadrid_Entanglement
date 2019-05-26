package main

import (
	"fmt"
	"github.com/racin/HackathonMadrid_Entanglement/Code/Entangler"
	"github.com/racin/HackathonMadrid_Entanglement/Code/SwarmConnector"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

	if _, err := os.Create(Entangler.TempDirectory + handler.Filename); err == nil {

	} else {
		fmt.Println("Fatal error ... " + err.Error())
		os.Exit(1)
	}
	ioutil.WriteFile(Entangler.TempDirectory+handler.Filename, fileBytes, os.ModeAppend)

	// Chunker & Entangler
	Entangler.EntangleFile(Entangler.TempDirectory + handler.Filename)

	// Upload
	SwarmConnector.UploadAllChunks()

	allFile, _ := ioutil.ReadFile("../retrives.txt")
	fmt.Fprintf(w, string(allFile))
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Download Endpoint Hit")

	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("Param 'ID' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	fmt.Println("Url Param 'key' is: " + string(key))

	bytes, _ := ioutil.ReadFile(Entangler.ChunkDirectory + "BC_Logo_.png")
	/*if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}*/

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	if _, err := w.Write(bytes); err != nil {
		fmt.Println("unable to write image.")
	}

	//fmt.Fprintf(w, string(bytes))
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	err := http.ListenAndServe(":8081", nil)
	fmt.Println(err.Error())
}

func main() {
	fmt.Println("Hello World")
	//setupRoutes()
	SwarmConnector.DownloadAndReconstruct(1, 2, 3, 5, 6, 7)
}
