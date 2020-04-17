package main

import (
	"flag"
	"fmt"
	. "genids/config"
	"genids/utils/logs"
	"genids/web/mdw"
	"github.com/getsentry/sentry-go"
	_ "github.com/icattlecoder/godaemon" //包的init函数实现，启动时加入 -d=true表示为daemon运行
	"os"
)

var (
	version     string
	buildTime   string
	goVersion   string
	gitCommitid string
)

func main() {
	// 初始化配置文件
	flag.Parse()
	fmt.Print("InitConfig...\r")
	checkErr("InitConfig", InitConfig())
	fmt.Print("InitConfig Success!!!\n")

	baseConf := Config.BaseConf

	initLogger(baseConf)
	initSentry(baseConf)
	startService(baseConf)
}

func initSentry(baseConf BaseConf) {
	// 启动sentry
	if baseConf.SentryDSN != "" {
		checkErr("Sentry initialization", sentry.Init(sentry.ClientOptions{Dsn: baseConf.SentryDSN}))
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("ServiceName", SERVER_NAME)
		})
		fmt.Print("Sentry initialization Success!!!\n")
	}
}

func initLogger(baseConf BaseConf) {
	// 创建文件日志，按天分割，默认日志文件仅保留一周, 此处更新为根据配置文件决定,不支持动态更改
	err := logs.InitLog(baseConf.SystemLogPath, baseConf.LogMaxAge)
	checkErr("InitLog", err)
}

func startService(conf BaseConf) {
	// 开始运行iris框架
	fmt.Printf(
		"genids running:\nversion:	%s\t\nbuildTime:	%s\t\ngoVersion:	%s\t\ngitCommitid:	%s\t\n",
		version, buildTime, goVersion, gitCommitid)

	// 启动API服务
	mdw.RunIris(conf.ServerPort, conf.SentryDSN)
}

// 检查错误
func checkErr(errMsg string, err error) {
	if err != nil {
		fmt.Printf("%s Error: %v\n", errMsg, err)
		os.Exit(1)
	}
}
