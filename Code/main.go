package main

import (
	"fmt"
	"github.com/racin/HackathonMadrid_Entanglement/Code/Entangler"
	"github.com/racin/HackathonMadrid_Entanglement/Code/SwarmConnector"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	var boolArr []bool
	var length int = 36 // Hardcoded to 36 for the demo ..
	strSplit := strings.Split(key, ",")
	compare := 0
	for i := 1; i <= length; i++ {
		if compare >= len(strSplit) {
			boolArr = append(boolArr, false)
			continue
		}
		str, _ := strconv.Atoi(strSplit[compare]) // 1,4,5
		if i != str {
			boolArr = append(boolArr, false)
		} else {
			boolArr = append(boolArr, true)
			compare++
		}
	}

	SwarmConnector.DownloadAndReconstruct(Entangler.ChunkDirectory+"reconstruct_swarm_logo.jpeg", boolArr...)

	bytes, _ := ioutil.ReadFile(Entangler.ChunkDirectory + "reconstruct_swarm_logo.jpeg")
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
	//http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	err := http.ListenAndServe(":8081", nil)
	fmt.Println(err.Error())
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}
