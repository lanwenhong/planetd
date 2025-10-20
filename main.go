package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"planetd/pkg/setting"
	"planetd/routers"

	"github.com/lanwenhong/lgobase/logger"
)

func main() {

	flag.Parse()
	arg_num := len(os.Args)
	if arg_num != 2 {
		fmt.Printf("input param error")
		return
	}

	var filename = os.Args[1]
	err := setting.ParseConf(filename)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	loglevel, _ := logger.LoggerLevelIndex(setting.Conf.LogLevel)
	lconf := &logger.Glogconf{
		RotateMethod: logger.ROTATE_FILE_DAILY,
		Stdout:       setting.Conf.LogStdOut,
		Loglevel:     loglevel,
		Colorful:     setting.Conf.Colorful,
	}
	logger.Newglog(setting.Conf.LogDir, setting.Conf.LogFile, setting.Conf.LogFileErr, lconf)
	// ctx, _ := context.WithCancel(context.Background())
	// ctx = context.WithValue(ctx, "trace_id", util.GenXid())
	router := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.Conf.Port),
		Handler:        router,
		ReadTimeout:    setting.Conf.ReadTimeout,
		WriteTimeout:   setting.Conf.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
