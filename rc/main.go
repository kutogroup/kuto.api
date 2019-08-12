package main

import (
	"encoding/json"
	"fmt"
	"kuto/config"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"waha.api/pkg"
)

var url string
var logger *pkg.WahaLogger
var cdn *pkg.WahaCDN

func main() {
	logger = pkg.NewLogger(os.Stdout, true)
	cdn = pkg.NewCDN(config.CDNHost, config.CDNAccessKey, config.CDNSecretKey, config.CDNSsl, config.CDNTimeout)

	r := mux.NewRouter()
	r.HandleFunc("/upload", upload)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	f, header, err := r.FormFile("file")
	if err != nil {
		logger.Err(err)
		w.WriteHeader(403)
		return
	}

	defer f.Close()

	bucket := r.Form["bucket"]
	if len(bucket) > 0 {
		pos := strings.LastIndex(header.Filename, ".")
		if pos < 0 {
			pos = len(header.Filename)
		}

		//取微秒
		stamp := time.Now().UnixNano() / 1000 / 1000
		newName := fmt.Sprintf("%s/%s/%s%s",
			strconv.FormatInt((stamp/10000/10000)%10000, 16),
			strconv.FormatInt((stamp/10000)%10000, 16),
			strconv.FormatInt(stamp%10000, 16),
			header.Filename[pos:])

		err = cdn.PutObject(bucket[0], newName, f, header.Size, nil)
		if err != nil {
			logger.Err(err)
			w.WriteHeader(403)
			return
		}

		logger.I("bucket=%s, name=%s", bucket[0], newName)
		// url := cdn.GetPresignedURL(bucket[0], newName, 7*24*60*60*time.Second)
		scheme := "http"
		if config.CDNSsl == true {
			scheme = "https"
		}
		url := fmt.Sprintf("%s://%s/%s/%s", scheme, config.CDNHost, bucket[0], newName)

		b, _ := json.Marshal(map[string]interface{}{
			"code": 200,
			"msg":  "Success",
			"data": map[string]string{
				"bucket": bucket[0],
				"name":   newName,
				"url":    url,
			},
		})

		w.Write(b)
	} else {
		logger.E("no bucket field")
		w.WriteHeader(403)
	}
}
