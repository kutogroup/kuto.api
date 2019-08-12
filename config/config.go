package config

import (
	"time"
)

//DebugMode 调试模式
const DebugMode = true

//cdn配置
const (
	//CDNHost cdn主机地址
	CDNHost = "rc.kutoapps.com"
	//CDNAccessKey cdn access key
	CDNAccessKey = "1SYOD1M9RIAS97U1X1RZ"
	//CDNSecretKey cdn secret key
	CDNSecretKey = "anBa2vd6jZtIgNZ0QNJaYAMKvhE3Lp7gQRQLFt3H"
	//CDNTimeout cdn超时时间
	CDNTimeout = time.Minute
)

//数据库配置
const (
	//DBHost 数据库主机地址
	DBHost = "localhost"
	//DBTable 数据库表
	DBTable = "kuto"
	//DBUser 数据库用户
	DBUser = "root"
	//DBPwd 数据库密码
	DBPwd = "root"
)

//缓存
const (
	//CacheHost 缓存地址
	CacheHost = ":6379"
	//CachePoolSize 缓存池大小
	CachePoolSize = 6
)
