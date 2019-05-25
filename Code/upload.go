package main

//import requiered libraries
import (
//    "fmt"
    "log"
    "os"
    "io/ioutil"

    //bzzclient "https://github.com/ethereum/go-ethereum/tree/master/swarm/api/client/client.go"
    bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

var (
    newFile *os.File
    err     error
)

func main() {

var path string
//define Swarm cliente
    client := bzzclient.NewClient("http://127.0.0.1:8500")
//Create file retrieval log info
//    f, retrieves := os.Create("./$home/swarm/files/retrives.txt)


newFile, err = os.Create("files/retrives.txt")
//Read directory

files, err := ioutil.ReadDir("files/")
if err != nil {
    log.Fatal(err)
}

for _, file := range files {
    //fmt.Println(file.Name())
    path = "files/"
    path += file.Name()
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    fileSw, err := bzzclient.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    //Upload file[n]
    manifestHash, err := client.Upload(fileSw, "", false)
    if err != nil {
        log.Fatal(err)
    }
    //Log retrieval information
            f, err := os.OpenFile("files/retrives.txt", os.O_APPEND|os.O_WRONLY, 0644)
            //fmt.Fprintln(f, file.Name())
            f.WriteString (file.Name())
            f.WriteString ("=")
            f.WriteString (manifestHash)
            f.WriteString ("\n")
        //    fmt.Fprintln(f, manifestHash)
          //  defer f.close()
            //fmt.Println(manifestHash) // 2e0849490b62e706a5f1cb8e7219db7b01677f2a859bac4b5f522afd2a5f02c0
}

}
