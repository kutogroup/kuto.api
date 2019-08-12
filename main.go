package main

import (
	"kuto/config"
	"kuto/pkg"
	"net/http"
	"os"
	"time"
)

func main() {
	h := NewHTTP(
		pkg.NewCDN(config.CDNHost, config.CDNAccessKey, config.CDNSecretKey, config.CDNTimeout),
		pkg.NewDatabase(config.DBHost, config.DBTable, config.DBUser, config.DBPwd),
		pkg.NewLogger(os.Stdout, true),
		pkg.NewCache(config.CacheHost, config.CachePoolSize, time.Minute))
	h.Serve(":1234", time.Second)
}

func (c *HTTP) Upload(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
	{
    "name": "BeJson",
    "url": "http://www.bejson.com",
    "page": 88,
    "isNonProfit": true,
    "address": {
        "street": "科技园路.",
        "city": "江苏苏州",
        "country": "中国"
    },
    "links": [
        {
            "name": "Google",
            "url": "http://www.google.com"
        },
        {
            "name": "Baidu",
            "url": "http://www.baidu.com"
        },
        {
            "name": "SoSo",
            "url": "http://www.SoSo.com"
        }
    ]
}
	`))
}

func (c *HTTP) testUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
