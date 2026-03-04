package test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

const addr = "http://127.0.0.1:9000"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

type TestRequestData struct {
	txamt       int
	clisn       string
	orig_chnlsn string
	tip_amt     int
	orig_txamt  int
}

func GeneratePaymentData(pos_ext map[string]interface{}, trd TestRequestData) []byte {
	data := map[string]interface{}{
		"code":          "",
		"opuid":         "0",
		"pay_tag":       "",
		"code_verifier": "",
		"appver":        "4.34.16",
		"channel_infos": map[string]interface{}{
			"mackey":  "",
			"partner": "",
			"zpk":     "2c55ae9dfe263c1bc63199f9638a416c",
			"zmk":     "4d51c694763a2a4a3e39a65452b574e7",
		},
		"paysource":       "",
		"limit_pay":       "",
		"txdtm":           time.Now().Format("2006-01-02 15:04:05"),
		"prepay_amt":      0,
		"pos_ext":         pos_ext,
		"syssn":           time.Now().Format("20060102150405") + "000072424",
		"expired_time":    "",
		"sysdtm":          time.Now().Format("2006-01-02 15:04:05"),
		"goods_name":      "dasd-确认码:2424",
		"access_token":    "",
		"userid":          1130081517,
		"alipay_pcode":    "",
		"udid":            "869071055892963",
		"appid":           "",
		"return_url":      "",
		"store_id":        "1130081517",
		"businm":          "",
		"tip_amt":         trd.tip_amt,
		"currency":        "344",
		"chnlid":          0,
		"div_template_id": "",
		"clisn":           trd.clisn, // 使用随机生成的流水号
		"alipay_params":   "",
		"auth_code":       "",
		"out_trade_no":    "",
		"customer_id":     "",
		"groupid":         0,
		"clientip":        "",
		"openid":          "",
		"goods_detail":    "",
		"txcurrcd":        "344",
		"lnglat":          []interface{}{"0", "0"},
		"failed_url":      "",
		"txzone":          "+0800",
		"lang":            "",
		"mch_infos": map[string]interface{}{
			"sub_mchntid": "12345678",
			"key2":        "",
			"mchnt_name":  "dasd",
			"sub_appid":   "",
			"memo":        "",
			"tag":         "",
			"postcode":    "",
			"key":         "",
			"mcc":         "",
			"city":        "",
			"userid":      1130081517,
			"ext":         nil,
			"appid":       "",
			"groupid":     1130000007,
			"mchntid":     "188000344333",
		},
		"swipe_method": "",
		"txamt":        trd.txamt,
		"busicd":       "802808",
		"phonemodel":   "A8S",
		"goods_info":   "",
		"os":           "Android",
		"orig_chnlsn":  trd.orig_chnlsn,
		"orig_txamt":   trd.orig_txamt,
	}
	ret, _ := json.Marshal(data)
	return ret
}

func TestPing(t *testing.T) {
	t.Log("ping")
	client := resty.New()
	client.SetDebug(true)
	path := "/ping"
	resp, err := client.R().Get(addr + path)
	if err != nil {
		t.Errorf("err: %s", err.Error())
		return
	}
	t.Logf("resp: %s", resp.String())
	if resp.StatusCode() != 200 {
		t.Errorf("resp: %s", resp.String())
		return
	}
	assert.Equal(t, resp.String(), "OK")
}

func TestTrade(t *testing.T) {
	// 消费
	fullURL := addr + "/trade/v1/swipe"
	client := resty.New()
	//pos_ext := "0zQnyvKPMydDHW/wTxwvq02RRYYuXt+pGCuSwZsE6P9puViFL2jHbAU/hivEqXz/H4Nidsi8o8cwSIIMGbZYp/wRubOTnYILw+OcBhq5zcQ1eOfsrs1zrSWGE4iKTzri6yDiEsdzSI5xN2Kf+UOhd3ZT21THFSn6IjrwV7wVgFvRMB4yAoTWCG+67eV2p0nLsoey7uS5P8+WFrC216RIhLq36n3GjG2Xf7ZCov3+jsYva0Uvabdvo5rIB6SbxUbxjUdAOlaww302E2bJmXMIL33BfHwUqcjS443qj9HqFGCH2RryoI1iy9+vt8dj5YT58lztSIdpdZMfYJem38hGrQI2obOkFp/Vli5Hc35ew0kdj7WERaNfzaIwNxtH3SxwIc9BDC6dPDCQeX4CMRwB1DdynVnWVgLiAtdbz8fDqgwOQrrGdlBf8z79HKcM/NX7qBe3dyCElBIqPMhySSdbanVW5cTy6kIBU1ps/JTmmHx4jg65qrX7usRaHzU9dqXUH/N8CrJto177x0AikDfIJl/vINeKpvjRxGUJustQaOXTEWdExJbgHgg4JMN6HyhAo5qD7s08Q4GTWGs7ayjqmo84je9QLcl7CpGrshgX51+hnh+BQpclFTgQnzJ4JA9F3yzCQ94dqoq4LuWVCwNIqA=="
	pos_ext := map[string]interface{}{
		"cardseqnum":   "001",
		"entry_mode":   "072",
		"iccdata":      "9F26082BE1761CCDE654E89F2701809F101706011203A000000F00564953414C3354455354434153459F3704FC5E4D3A9F36020001950500000000009A032601229C01009F02060000000110005F2A020344820200009F1A0203449F3303E028C89F3501228407A00000000310109F0902008C9F34031F03029F1E0843415353383334389F03060000000000005F340101",
		"macString":    "5D3DE4BCA70F5050",
		"pinBlock":     "",
		"cardNo":       "4761731000000100",
		"track1":       "",
		"track2":       "4761731000000100D311220113202839",
		"track3":       "",
		"expired_date": "2612",
		"trade_type":   "trade",
	}
	txamt := 10000
	tipAmt := 0
	clisn := generateVerificationCode()
	trd := TestRequestData{
		txamt:       txamt,
		clisn:       clisn,
		orig_chnlsn: "",
		tip_amt:     tipAmt,
		orig_txamt:  0,
	}
	content := GeneratePaymentData(pos_ext, trd)
	client.SetTimeout(15 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("X-Request-ID", "E10000000").
		SetBody(content).
		Post(fullURL)
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}

func TestVoidTrade(t *testing.T) {
	// 消费撤销
	fullURL := addr + "/trade/v1/void_trade"
	client := resty.New()
	pos_ext := map[string]interface{}{
		"cardseqnum":   "001",
		"entry_mode":   "052",
		"iccdata":      "",
		"macString":    "0BD57554F4B099B7",
		"pinBlock":     "",
		"cardNo":       "4336680006896670",
		"track1":       "",
		"track2":       "",
		"track3":       "",
		"expired_date": "2612",
		"trade_type":   "void_trade",
	}
	txamt := 100
	tipAmt := 0
	clisn := generateVerificationCode()
	orig_chnlsn := "531458030873"
	trd := TestRequestData{
		txamt:       txamt,
		clisn:       clisn,
		orig_chnlsn: orig_chnlsn,
		tip_amt:     tipAmt,
		orig_txamt:  0,
	}
	content := GeneratePaymentData(pos_ext, trd)
	client.SetTimeout(5 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(content).
		Post(fullURL)
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}

func TestRefund(t *testing.T) {
	// 消费退款
	fullURL := addr + "/trade/v1/refund"
	client := resty.New()
	//pos_ext := "0zQnyvKPMydDHW/wTxwvq02RRYYuXt+pGCuSwZsE6P9puViFL2jHbAU/hivEqXz/H4Nidsi8o8cwSIIMGbZYp/wRubOTnYILw+OcBhq5zcQ1eOfsrs1zrSWGE4iKTzri6yDiEsdzSI5xN2Kf+UOhd3ZT21THFSn6IjrwV7wVgFvRMB4yAoTWCG+67eV2p0nLsoey7uS5P8+WFrC216RIhLq36n3GjG2Xf7ZCov3+jsYva0Uvabdvo5rIB6SbxUbxjUdAOlaww302E2bJmXMIL33BfHwUqcjS443qj9HqFGCH2RryoI1iy9+vt8dj5YT58lztSIdpdZMfYJem38hGrQI2obOkFp/Vli5Hc35ew0kdj7WERaNfzaIwNxtH3SxwIc9BDC6dPDCQeX4CMRwB1DdynVnWVgLiAtdbz8fDqgwOQrrGdlBf8z79HKcM/NX7qBe3dyCElBIqPMhySSdbanVW5cTy6kIBU1ps/JTmmHx4jg65qrX7usRaHzU9dqXUH/N8CrJto177x0AikDfIJl/vINeKpvjRxGUJustQaOXTEWdExJbgHgg4JMN6HyhAo5qD7s08Q4GTWGs7ayjqmo84je9QLcl7CpGrshgX51+hnh+BQpclFTgQnzJ4JA9FfNLe5ZiN6AAJxKSkhRG/BQ=="
	pos_ext := map[string]interface{}{
		"cardseqnum":   "001",
		"entry_mode":   "052",
		"iccdata":      "9F26089784558B3CD1458B9F2701809F100706010A03A000009F370442B233069F360209FD950580C00818009A032510239C01009F02060000000000015F2A02015682021C009F1A0201569F3303E0F8C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000",
		"macString":    "0BD57554F4B099B7",
		"pinBlock":     "",
		"cardNo":       "4336680006896670",
		"track1":       "",
		"track2":       "4336680006896670D22022011193265100000",
		"track3":       "",
		"expired_date": "2612",
		"trade_type":   "refund",
	}
	txamt := 100
	tipAmt := 0
	clisn := generateVerificationCode()
	orig_chnlsn := "531057742951"
	trd := TestRequestData{
		txamt:       txamt,
		clisn:       clisn,
		orig_chnlsn: orig_chnlsn,
		tip_amt:     tipAmt,
		orig_txamt:  0,
	}
	content := GeneratePaymentData(pos_ext, trd)
	client.SetTimeout(5 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(content).
		Post(fullURL)
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}

func TestVoidRefund(t *testing.T) {
	// 撤销退款
	fullURL := addr + "/trade/v1/void_refund"
	client := resty.New()
	//pos_ext := "0zQnyvKPMydDHW/wTxwvq02RRYYuXt+pGCuSwZsE6P9puViFL2jHbAU/hivEqXz/H4Nidsi8o8cwSIIMGbZYp/wRubOTnYILw+OcBhq5zcQ1eOfsrs1zrSWGE4iKTzri6yDiEsdzSI5xN2Kf+UOhd3ZT21THFSn6IjrwV7wVgFvRMB4yAoTWCG+67eV2p0nLsoey7uS5P8+WFrC216RIhLq36n3GjG2Xf7ZCov3+jsYva0Uvabdvo5rIB6SbxUbxjUdAOlaww302E2bJmXMIL33BfHwUqcjS443qj9HqFGCH2RryoI1iy9+vt8dj5YT58lztSIdpdZMfYJem38hGrQI2obOkFp/Vli5Hc35ew0kdj7WERaNfzaIwNxtH3SxwIc9BDC6dPDCQeX4CMRwB1DdynVnWVgLiAtdbz8fDqgwOQrrGdlBf8z79HKcM/NX7qBe3dyCElBIqPMhySSdbanVW5cTy6kIBU1ps/JTmmHx4jg65qrX7usRaHzU9dqXUH/N8CrJto177x0AikDfIJl/vINeKpvjRxGUJustQaOXTEWdExJbgHgg4JMN6HyhAo5qD7s08Q4GTWGs7ayjqmo84je9QLcl7CpGrshgX51+hnh+BQpclFTgQnzJ4JA9Fy83dzTs5OSryzcfS1rOJ7w=="
	pos_ext := map[string]interface{}{
		"cardseqnum":   "001",
		"entry_mode":   "052",
		"iccdata":      "9F26089784558B3CD1458B9F2701809F100706010A03A000009F370442B233069F360209FD950580C00818009A032510239C01009F02060000000000015F2A02015682021C009F1A0201569F3303E0F8C89F34031E03009F3501228407A00000000310109F0902008C9F1E0843415357383332309F0306000000000000",
		"macString":    "0BD57554F4B099B7",
		"pinBlock":     "",
		"cardNo":       "4336680006896670",
		"track1":       "",
		"track2":       "4336680006896670D22022011193265100000",
		"track3":       "",
		"expired_date": "2612",
		"trade_type":   "void_refund",
	}
	txamt := 100
	clisn := generateVerificationCode()
	orig_chnlsn := "531057743268"
	tipAmt := 0
	trd := TestRequestData{
		txamt:       txamt,
		clisn:       clisn,
		orig_chnlsn: orig_chnlsn,
		tip_amt:     tipAmt,
		orig_txamt:  0,
	}
	content := GeneratePaymentData(pos_ext, trd)
	client.SetTimeout(5 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(content).
		Post(fullURL)
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}

func TestKeyExchange(t *testing.T) {
	fullURL := addr + "/trade/v1/key_exchange"
	client := resty.New()
	content := `{"clisn": "000275", "mch_infos": {"sub_mchntid": "12345678", "key2": "", "mchnt_name": "dasd", "sub_appid": "", "memo": "", "tag": "", "postcode": "", "key": "", "mcc": "", "city": "", "userid": 1130081517, "ext": null, "appid": "", "groupid": 1130000007, "mchntid": "188000344333"}}`
	client.SetTimeout(5 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(content).
		Post(fullURL)

	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}

func Test96Trade(t *testing.T) {
	// 消费
	fullURL := addr + "/trade/v1/swipe"
	client := resty.New()
	//pos_ext := "0zQnyvKPMydDHW/wTxwvq02RRYYuXt+pGCuSwZsE6P9puViFL2jHbAU/hivEqXz/H4Nidsi8o8cwSIIMGbZYp/wRubOTnYILw+OcBhq5zcQ1eOfsrs1zrSWGE4iKTzri6yDiEsdzSI5xN2Kf+UOhd3ZT21THFSn6IjrwV7wVgFvRMB4yAoTWCG+67eV2p0nLsoey7uS5P8+WFrC216RIhLq36n3GjG2Xf7ZCov3+jsYva0Uvabdvo5rIB6SbxUbxjUdAOlaww302E2bJmXMIL33BfHwUqcjS443qj9HqFGCH2RryoI1iy9+vt8dj5YT58lztSIdpdZMfYJem38hGrQI2obOkFp/Vli5Hc35ew0kdj7WERaNfzaIwNxtH3SxwIc9BDC6dPDCQeX4CMRwB1DdynVnWVgLiAtdbz8fDqgwOQrrGdlBf8z79HKcM/NX7qBe3dyCElBIqPMhySSdbanVW5cTy6kIBU1ps/JTmmHx4jg65qrX7usRaHzU9dqXUH/N8CrJto177x0AikDfIJl/vINeKpvjRxGUJustQaOXTEWdExJbgHgg4JMN6HyhAo5qD7s08Q4GTWGs7ayjqmo84je9QLcl7CpGrshgX51+hnh+BQpclFTgQnzJ4JA9F3yzCQ94dqoq4LuWVCwNIqA=="
	pos_ext := map[string]interface{}{
		"cardseqnum":   "001",
		"entry_mode":   "072",
		"iccdata":      "9F260869B5841B37D49D749F2701809F10201F220100A000000000564953414C3354455354434153450000000000000000009F3704C8BE1FAD9F36020002950500000000009A032602109C01009F02060000000102005F2A020344820200209F1A0203449F3303E0B8C89F3501228407A00000000310109F0902008C9F34031F03029F1E0843415353383334389F03060000000000005F340101",
		"macString":    "A6B994899F5C0B97",
		"pinBlock":     "",
		"cardNo":       "4761731000000027",
		"track1":       "",
		"track2":       "4761731000000027D311220119058101",
		"track3":       "",
		"expired_date": "2612",
		"trade_type":   "trade",
	}
	txamt := 10200
	tipAmt := 0
	clisn := generateVerificationCode()
	trd := TestRequestData{
		txamt:       txamt,
		clisn:       clisn,
		orig_chnlsn: "",
		tip_amt:     tipAmt,
		orig_txamt:  0,
	}
	content := GeneratePaymentData(pos_ext, trd)
	client.SetTimeout(15 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("X-Request-ID", "E10000000").
		SetBody(content).
		Post(fullURL)
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}

func TestTip1Trade(t *testing.T) {
	// 消费
	fullURL := addr + "/trade/v1/swipe"
	client := resty.New()
	//pos_ext := "0zQnyvKPMydDHW/wTxwvq02RRYYuXt+pGCuSwZsE6P9puViFL2jHbAU/hivEqXz/H4Nidsi8o8cwSIIMGbZYp/wRubOTnYILw+OcBhq5zcQ1eOfsrs1zrSWGE4iKTzri6yDiEsdzSI5xN2Kf+UOhd3ZT21THFSn6IjrwV7wVgFvRMB4yAoTWCG+67eV2p0nLsoey7uS5P8+WFrC216RIhLq36n3GjG2Xf7ZCov3+jsYva0Uvabdvo5rIB6SbxUbxjUdAOlaww302E2bJmXMIL33BfHwUqcjS443qj9HqFGCH2RryoI1iy9+vt8dj5YT58lztSIdpdZMfYJem38hGrQI2obOkFp/Vli5Hc35ew0kdj7WERaNfzaIwNxtH3SxwIc9BDC6dPDCQeX4CMRwB1DdynVnWVgLiAtdbz8fDqgwOQrrGdlBf8z79HKcM/NX7qBe3dyCElBIqPMhySSdbanVW5cTy6kIBU1ps/JTmmHx4jg65qrX7usRaHzU9dqXUH/N8CrJto177x0AikDfIJl/vINeKpvjRxGUJustQaOXTEWdExJbgHgg4JMN6HyhAo5qD7s08Q4GTWGs7ayjqmo84je9QLcl7CpGrshgX51+hnh+BQpclFTgQnzJ4JA9F3yzCQ94dqoq4LuWVCwNIqA=="
	pos_ext := map[string]interface{}{
		"cardseqnum":   "001",
		"entry_mode":   "072",
		"iccdata":      "9F26082BE1761CCDE654E89F2701809F101706011203A000000F00564953414C3354455354434153459F3704FC5E4D3A9F36020001950500000000009A032601229C01009F02060000000110005F2A020344820200009F1A0203449F3303E028C89F3501228407A00000000310109F0902008C9F34031F03029F1E0843415353383334389F03060000000000005F340101",
		"macString":    "5D3DE4BCA70F5050",
		"pinBlock":     "",
		"cardNo":       "4761731000000100",
		"track1":       "",
		"track2":       "4761731000000100D311220113202839",
		"track3":       "",
		"expired_date": "2612",
		"trade_type":   "trade",
	}
	txamt := 90000
	tipAmt := 0
	clisn := generateVerificationCode()
	trd := TestRequestData{
		txamt:       txamt,
		clisn:       clisn,
		orig_chnlsn: "",
		tip_amt:     tipAmt,
		orig_txamt:  0,
	}
	content := GeneratePaymentData(pos_ext, trd)
	client.SetTimeout(15 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("X-Request-ID", "E10000000").
		SetBody(content).
		Post(fullURL)
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}

func TestTip2Trade(t *testing.T) {
	// 消费
	fullURL := addr + "/trade/v1/tips"
	client := resty.New()
	//pos_ext := "0zQnyvKPMydDHW/wTxwvq02RRYYuXt+pGCuSwZsE6P9puViFL2jHbAU/hivEqXz/H4Nidsi8o8cwSIIMGbZYp/wRubOTnYILw+OcBhq5zcQ1eOfsrs1zrSWGE4iKTzri6yDiEsdzSI5xN2Kf+UOhd3ZT21THFSn6IjrwV7wVgFvRMB4yAoTWCG+67eV2p0nLsoey7uS5P8+WFrC216RIhLq36n3GjG2Xf7ZCov3+jsYva0Uvabdvo5rIB6SbxUbxjUdAOlaww302E2bJmXMIL33BfHwUqcjS443qj9HqFGCH2RryoI1iy9+vt8dj5YT58lztSIdpdZMfYJem38hGrQI2obOkFp/Vli5Hc35ew0kdj7WERaNfzaIwNxtH3SxwIc9BDC6dPDCQeX4CMRwB1DdynVnWVgLiAtdbz8fDqgwOQrrGdlBf8z79HKcM/NX7qBe3dyCElBIqPMhySSdbanVW5cTy6kIBU1ps/JTmmHx4jg65qrX7usRaHzU9dqXUH/N8CrJto177x0AikDfIJl/vINeKpvjRxGUJustQaOXTEWdExJbgHgg4JMN6HyhAo5qD7s08Q4GTWGs7ayjqmo84je9QLcl7CpGrshgX51+hnh+BQpclFTgQnzJ4JA9F3yzCQ94dqoq4LuWVCwNIqA=="

	pos_ext := map[string]interface{}{
		"cardseqnum":   "001",
		"entry_mode":   "072",
		"iccdata":      "9F26082BE1761CCDE654E89F2701809F101706011203A000000F00564953414C3354455354434153459F3704FC5E4D3A9F36020001950500000000009A032601229C01009F02060000000110005F2A020344820200009F1A0203449F3303E028C89F3501228407A00000000310109F0902008C9F34031F03029F1E0843415353383334389F03060000000000005F340101",
		"macString":    "5D3DE4BCA70F5050",
		"pinBlock":     "",
		"cardNo":       "4761731000000100",
		"track1":       "",
		"track2":       "4761731000000100D311220113202839",
		"track3":       "",
		"expired_date": "2612",
		"trade_type":   "tip",
	}
	/*
		pos_ext := map[string]interface{}{
			"cardseqnum": "001",
			"entry_mode": "072",
			"iccdata":    "",
			"macString":  "",
			"pinBlock":   "",
			"cardNo":     "4761731000000100",
			"track1":     "",
			"track2":     "",
			"track3":     "",
			"expired_date": "2612",
			"trade_type": "tip",
		}
	*/

	txamt := 1000
	orig_txamt := 10000
	tipAmt := 0
	orig_chnlsn := "606366694016"
	clisn := generateVerificationCode()
	trd := TestRequestData{
		txamt:       txamt,
		clisn:       clisn,
		orig_chnlsn: orig_chnlsn,
		tip_amt:     tipAmt,
		orig_txamt:  orig_txamt,
	}
	content := GeneratePaymentData(pos_ext, trd)
	client.SetTimeout(15 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("X-Request-ID", "E10000000").
		SetBody(content).
		Post(fullURL)
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}

func TestKeyInTrade(t *testing.T) {
	// 消费
	fullURL := addr + "/trade/v1/key_in"
	client := resty.New()
	//pos_ext := "0zQnyvKPMydDHW/wTxwvq02RRYYuXt+pGCuSwZsE6P9puViFL2jHbAU/hivEqXz/H4Nidsi8o8cwSIIMGbZYp/wRubOTnYILw+OcBhq5zcQ1eOfsrs1zrSWGE4iKTzri6yDiEsdzSI5xN2Kf+UOhd3ZT21THFSn6IjrwV7wVgFvRMB4yAoTWCG+67eV2p0nLsoey7uS5P8+WFrC216RIhLq36n3GjG2Xf7ZCov3+jsYva0Uvabdvo5rIB6SbxUbxjUdAOlaww302E2bJmXMIL33BfHwUqcjS443qj9HqFGCH2RryoI1iy9+vt8dj5YT58lztSIdpdZMfYJem38hGrQI2obOkFp/Vli5Hc35ew0kdj7WERaNfzaIwNxtH3SxwIc9BDC6dPDCQeX4CMRwB1DdynVnWVgLiAtdbz8fDqgwOQrrGdlBf8z79HKcM/NX7qBe3dyCElBIqPMhySSdbanVW5cTy6kIBU1ps/JTmmHx4jg65qrX7usRaHzU9dqXUH/N8CrJto177x0AikDfIJl/vINeKpvjRxGUJustQaOXTEWdExJbgHgg4JMN6HyhAo5qD7s08Q4GTWGs7ayjqmo84je9QLcl7CpGrshgX51+hnh+BQpclFTgQnzJ4JA9F3yzCQ94dqoq4LuWVCwNIqA=="
	pos_ext := map[string]interface{}{
		"cardseqnum":   "001",
		"entry_mode":   "012",
		"iccdata":      "9F26082BE1761CCDE654E89F2701809F101706011203A000000F00564953414C3354455354434153459F3704FC5E4D3A9F36020001950500000000009A032601229C01009F02060000000110005F2A020344820200009F1A0203449F3303E028C89F3501228407A00000000310109F0902008C9F34031F03029F1E0843415353383334389F03060000000000005F340101",
		"macString":    "5D3DE4BCA70F5050",
		"pinBlock":     "",
		"cardNo":       "4761731000000100",
		"track1":       "",
		"track2":       "4761731000000100D311220113202839",
		"track3":       "",
		"expired_date": "2612",
		"trade_type":   "key_in",
	}
	txamt := 20000
	tipAmt := 0
	clisn := generateVerificationCode()
	trd := TestRequestData{
		txamt:       txamt,
		clisn:       clisn,
		orig_chnlsn: "",
		tip_amt:     tipAmt,
		orig_txamt:  0,
	}
	content := GeneratePaymentData(pos_ext, trd)
	client.SetTimeout(15 * time.Second)
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("X-Request-ID", "E10000000").
		SetBody(content).
		Post(fullURL)
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
}
