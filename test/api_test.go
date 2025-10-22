package test

import (
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

const addr = "http://127.0.0.1:9000"

func TestPing(t *testing.T) {
	t.Log("ping")
	client := resty.New()
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
	// 1. 打印完整请求URL
	fullURL := addr + "/trade/v1/swipe"
	t.Logf("请求URL: %s", fullURL)
	client := resty.New()
	content := `{"code": "", "opuid": "0", "pay_tag": "", "code_verifier": "", "appver":"4.34.16", "channel_infos": {"mackey": "", "partner": "", "zpk": "2c55ae9dfe263c1bc63199f9638a416c", "zmk": "4d51c694763a2a4a3e39a65452b574e7"}, "paysource": "", "limit_pay": "", "txdtm": "2025-10-16 15:21:24", "prepay_amt": 0, "chnl_ext": "qiO1a1+TaVdi9fHGPBHcBV0O8Tf86FOHmR67B7JwNJ18wZBN+AP2L4mjkxne3hqDrtTLWploxpYfThI6UwkP5NcaTLnkgDYdAb0vmUuNbfhAtDRV4mjdnwavmlvZw8oWgj85GWlr58TIYz0TUWd570tFmae700xZEsmbQ+SHK+q1eRWxXIxKN5zuznNpq62ckQtsDGvb3kxY9W2RVCWu+uH3o3hQB4R3TbYX26s/T7xyI5APZfhuESiIEYT0pMofHJmNStkgDu0JszPhjJd98ZsQXvLUR0TJxZQ40S7+26p8j8imGEslyR4OdSCeoocDwlo8+pqWky0kq9yME2SgFiiYbE8DCSYobn+yaQd8VfczcAN+4yzYKj4PjNXgac3zOuTFP7gK9plL7QAJ6vn3fdox8EFaABqbVrsbFzFwoAlkMul67iIQl6987xM3MPJyVKJF3B/QrU4Zi/tdBm5uIDO8qQnJIwA3SjxzokhKqemDLb8EpfpDG99kpsMJTHNLnTpP7BrF6Pv7NqvmE/mcQxd2dkPmNkDmjxvRecVzk2sDSO09t57vpNl6cvhjaOMMS+uj4Gbv3u1ettnPxv9qz0DyBQMR6Nviw+IRuWfEHpCqQ4317Btc43+OXzxuokRV+oOjMttI8YLPqKGd5yleCmE+CEdRoaj4ScKTLc0I3QHldn+Pse4KS5gKQXzJ4xqBwq917mwLx+Vk0z2A1YEkHAsAUb3E38D+XwSfqWbBsvA=", "syssn": "20251021070100020000072424", "expired_time": "", "sysdtm": "2025-10-21 17:10:49", "goods_name": "dasd-\\u786e\\u8ba4\\u7801:2424", "access_token": "", "userid": 1130081517, "alipay_pcode": "", "udid": "869071055892963", "appid": "", "return_url": "", "store_id": "1130081517", "businm": "", "tip_amt": 0, "currency": "344", "chnlid": 0, "div_template_id": "", "clisn": "000275", "alipay_params": "", "auth_code": "", "out_trade_no": "", "customer_id": "", "groupid": 0, "clientip": "", "openid": "", "goods_detail": "", "txcurrcd": "344", "lnglat": ["0", "0"], "failed_url": "", "txzone": "+0800", "lang": "", "mch_infos": {"sub_mchntid": "12345678", "key2": "", "mchnt_name": "dasd", "sub_appid": "", "memo": "", "tag": "", "postcode": "", "key": "", "mcc": "", "city": "", "userid": 1130081517, "ext": null, "appid": "", "groupid": 1130000007, "mchntid": "188000344333"}, "swipe_method": "", "txamt": 100, "busicd": "802808", "phonemodel": "A8S", "goods_info": "", "os": "Android"}`
	// 2. 打印请求体
	t.Logf("请求体: %s", content)

	// 3. 添加超时设置
	client.SetTimeout(5 * time.Second)

	// 4. 启用resty调试
	client.SetDebug(true)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(content).
		Post(fullURL)
	// 5. 打印完整响应详情
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
	t.Logf("完整响应: %+v", resp)
}

func TestKeyExchange(t *testing.T) {
	// 1. 打印完整请求URL
	fullURL := addr + "/trade/v1/key_exchange"
	t.Logf("请求URL: %s", fullURL)
	client := resty.New()
	content := `{"clisn": "000275", "mch_infos": {"sub_mchntid": "12345678", "key2": "", "mchnt_name": "dasd", "sub_appid": "", "memo": "", "tag": "", "postcode": "", "key": "", "mcc": "", "city": "", "userid": 1130081517, "ext": null, "appid": "", "groupid": 1130000007, "mchntid": "188000344333"}}`
	// 2. 打印请求体
	t.Logf("请求体: %s", content)

	// 3. 添加超时设置
	client.SetTimeout(5 * time.Second)

	// 4. 启用resty调试
	client.SetDebug(true)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(content).
		Post(fullURL)
	// 5. 打印完整响应详情
	if err != nil {
		t.Errorf("请求错误: %v", err)
		if resp != nil {
			t.Errorf("响应状态: %d", resp.StatusCode())
			t.Errorf("响应头: %v", resp.Header())
		}
		return
	}
	t.Logf("完整响应: %+v", resp)

}
