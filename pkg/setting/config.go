package setting

import (
	"fmt"
	"time"

	"github.com/lanwenhong/lgobase/confparse"
	"github.com/lanwenhong/lgobase/gconfig"
)

type Config struct {
	Addr         string        `confpos:"server:addr" dtype:"base"`
	Port         int           `confpos:"server:port" dtype:"base"`
	ReadTimeout  time.Duration `confpos:"server:read_timeout" dtype:"base"`
	WriteTimeout time.Duration `confpos:"server:write_timeout" dtype:"base"`
	RunMode      string        `confpos:"server:run_mode" dtype:"base"`

	LogFile    string `confpos:"log:logfile" dtype:"base"`
	LogFileErr string `confpos:"log:logfile_err" dtype:"base"`
	LogDir     string `confpos:"log:logdir" dtype:"base"`
	LogLevel   string `confpos:"log:loglevel" dtype:"base"`
	LogStdOut  bool   `confpos:"log:logstdout" dtype:"base"`
	Colorful   bool   `confpos:"log:colorfull" dtype:"base"`
}

var Conf *Config = new(Config)

func ParseConf(filename string) error {
	cfg := gconfig.NewGconf(filename)
	err := cfg.GconfParse()
	if err != nil {
		fmt.Printf("parse %s %s", filename, err.Error())
		return err
	}
	cp := confparse.CpaseNew(filename)
	err = cp.CparseGo(Conf, cfg)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("Addr: %s\n", Conf.Addr)
	fmt.Printf("Port: %d\n", Conf.Port)

	return err
}
