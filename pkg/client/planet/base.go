package planet

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"planetd/pkg/setting"
	"strconv"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/lgobase/network"
)

func SendPacket(ctx context.Context, bcd string) (string, error) {
	ret := ""
	rootCAData, err := os.ReadFile(setting.Conf.PlanetCertPath)
	if err != nil {
		logger.Debugf(ctx, "read root.crt error: %s", err.Error())
		return ret, err
	}
	rootCAs := x509.NewCertPool()
	// 将 PEM 格式的根证书添加到证书池
	succ := rootCAs.AppendCertsFromPEM(rootCAData)
	if !succ {
		err = fmt.Errorf("append certs from pem error")
		logger.Debugf(ctx, "append certs from pem error: %s", err.Error())
		return ret, err
	}

	tlsConfig := &tls.Config{
		RootCAs:    rootCAs,                          // 信任的根 CA 池（用于验证服务器证书）
		ServerName: "terminal.uat.planetpayment.com", // 服务器证书的域名（必须与证书的 CN 或 SAN 匹配）
		MinVersion: tls.VersionTLS12,                 // 禁用不安全的 TLS 版本
		// 禁用 InsecureSkipVerify（生产环境绝对不能开启！）
		//InsecureSkipVerify: true, // 开启后会跳过证书验证（不安全）
	}

	cTimeout := setting.Conf.PlanetConnTimetout
	rTimeout := setting.Conf.ReadTimeout
	wTimeout := setting.Conf.WriteTimeout
	addr := fmt.Sprintf("%s:%d", setting.Conf.PlanetAddr, setting.Conf.PlanetPort)
	logger.Debugf(ctx, "addr: %s", addr)
	c := network.NewTcpSsslConn(addr, cTimeout, rTimeout, wTimeout, tlsConfig)
	err = c.Open(ctx)
	if err != nil {
		err = fmt.Errorf("open tcp ssl conn error: %s", err.Error())
		logger.Debugf(ctx, "open tcp ssl conn error: %s", err.Error())
		return ret, err
	}

	b, _ := hex.DecodeString(bcd)
	c.Writen(ctx, b)

	head, err := c.Readn(ctx, 2)
	if err != nil {
		logger.Infof(ctx, "Readn err: %s", err.Error())
		return ret, err
	}

	slen := hex.EncodeToString(head)
	logger.Debugf(ctx, "slen: %s", slen)
	blen, err := strconv.ParseUint(slen, 16, 32)
	if err != nil {
		logger.Warnf(ctx, "err: %s", err.Error())
		return ret, err
	}
	logger.Debugf(ctx, "blen: %d", blen)
	data, err := c.Readn(ctx, int(blen))

	if err != nil {
		logger.Infof(ctx, "Readn err: %s", err.Error())
		return ret, err
	}
	bcdData := hex.EncodeToString(data)
	return bcdData, nil
}
