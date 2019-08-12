package main

import (
	"errors"
	"net/http"
	"time"

	"kuto/pkg"
	"kuto/utils"

	"github.com/gorilla/mux"
)

//HTTP http控制器
type HTTP struct {
	CDN    *pkg.WahaCDN
	DB     *pkg.WahaDB
	Log    *pkg.WahaLogger
	Cache  *pkg.WahaCache
	mux    *mux.Router
	custom map[string]interface{}
}

//NewHTTP 新建http对象
func NewHTTP(cdn *pkg.WahaCDN, db *pkg.WahaDB, log *pkg.WahaLogger, cache *pkg.WahaCache) *HTTP {
	waha := &HTTP{
		CDN:   cdn,
		DB:    db,
		Log:   log,
		Cache: cache,
	}

	waha.mux = mux.NewRouter()
	waha.custom = make(map[string]interface{})
	return waha
}

//AddExtras 添加自定义参数
func (c *HTTP) AddExtras(key string, value interface{}) *HTTP {
	c.custom[key] = value
	return c
}

//GetExtras 获取自定义参数
func (c *HTTP) GetExtras(key string) (interface{}, error) {
	if _, ok := c.custom[key]; !ok {
		return nil, errors.New("no key exist, key=" + key)
	}

	return c.custom[key], nil
}

//Handle 自定义path和函数的对应关系，注意这里的handle尽量使用unexported函数，否则会重复注册
func (c *HTTP) Handle(path string, handle func(http.ResponseWriter, *http.Request)) {
	c.mux.HandleFunc(path, handle)
}

//HandleFile 自定义path路径
func (c *HTTP) HandleFile(path string) {
	c.mux.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir("."+path))))
}

//Serve 监听http
func (c *HTTP) Serve(addr string, timeout time.Duration) {
	k, _ := utils.ReflectGetVTByOrigin(c)
	for i := 0; i < k.NumMethod(); i++ {
		method := k.Method(i)
		methodName := method.Name

		if f, ok := method.Func.Interface().(func(*HTTP, http.ResponseWriter, *http.Request)); ok {
			path := utils.ConvertCamel2Line(methodName)
			c.mux.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
				f(c, w, r)
			})
		}
	}

	srv := &http.Server{
		Handler: c.mux,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}

	c.Log.E("http error=", srv.ListenAndServe())
}
