package main

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	_ "embed"
	"bytes"
	"fmt"
	"log"
	"os"
	"errors"
	"os/exec"
)
//go:embed webpages/portal.html
var indexPage []byte

//go:embed webpages/VenomAntidotum.png
var portal []byte

func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	fileHash := md5.New()
	fileHash.Write(fileBytes)
	fmt.Fprintf(w, "{\"report\":\"./reports/" + handler.Filename + "." + hex.EncodeToString(fileHash.Sum(nil)) + "\"}" )
	go runCapa(fileBytes, handler.Filename, hex.EncodeToString(fileHash.Sum(nil)))
}

func runCapa(fileBytes []byte, fileName string, fileHash string){
	reportName := fileName + "." + fileHash
	cmd :=exec.Command("./capa",fileHash, "-v", "-q")
	fmt.Println("Running capa on " +fileName + ":" + fileHash + "...")
	ioutil.WriteFile(fileHash,fileBytes, 0666)
	outfile, _ := os.Create("./reports/" + reportName)
	output, _ := cmd.CombinedOutput()

	outfile.WriteString(string(output))
	os.Remove(fileHash)
}

func setupRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(indexPage)
	})
	mux.HandleFunc("/portal.png",func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(portal)
	})
	mux.Handle("/reports/",http.FileServer(http.Dir(".")))
	mux.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", mux)
}

var (
	capaUrl = "https://github.com/mandiant/capa/releases/download/v3.0.3/capa-v3.0.3-linux.zip"
	capaVer = "3.0.3"
)

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func DownloadFileUnzip(url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))

	for _, zipFile := range zipReader.File {
		fmt.Println("Reading file:", zipFile.Name)
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			log.Println(err)
			continue
		}

		 ioutil.WriteFile(zipFile.Name, unzippedFileBytes, 0755 )// this is unzipped file bytes
	}
	return err
}

func main() {
	if _, err := os.Stat("./capa"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Capa not found, downloading capa version " + capaVer)
		DownloadFileUnzip(capaUrl)
	}
	if _, err := os.Stat("./reports"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating reports directory...")
		os.Mkdir("./reports",0755)
	}
	setupRoutes()
}
