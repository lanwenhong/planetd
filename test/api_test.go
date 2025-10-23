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
	fullURL := addr + "/trade/v1/swipe"
	client := resty.New()
	content := `{"code": "", "opuid": "0", "pay_tag": "", "code_verifier": "", "appver":"4.34.16", "channel_infos": {"mackey": "", "partner": "", "zpk": "2c55ae9dfe263c1bc63199f9638a416c", "zmk": "4d51c694763a2a4a3e39a65452b574e7"}, "paysource": "", "limit_pay": "", "txdtm": "2025-10-16 15:21:24", "prepay_amt": 0, "chnl_ext": "qiO1a1+TaVdi9fHGPBHcBYGSgLOVSlI7Q4det8hIl3jRlBbkD0xinrsBz6JnIV/qGU+lLzKEkKKHUsTjkjSxaGMgpyiWAxwfiJz6FtRFVpOcP1UiH3agOqJEggzjY61ulsxJtktaSI6TAk81xsqZZk9EOcekYHFs2vOjMiT2FbVcJ+998SfuxSNFQNufTHCRC7CUG7QxN3afiR8ZHqwm5CFs4Ti5lVLFzKw5+mpzZ9N2BdJh1hKtaCAILiIV8P1aXnbhaTwaiOkRFMwzBNoPOse6dMRxmWcOKNWxZ8eBDSGvBPRdwIC42owF5MbkN4enusvNfrxlp+h4cApGdviyeQd0nxUP8M/ygFex8hYjNI5hUr8DxE9+aD2ZgI90Rk1H3CU6uGkF3JsmXjFdEgtHCyFcMKSINnN7mUwXobsGN7Y0yWvSn+df6YQiYduZ34M7jyFR4khgN6Vo+KsmuQX7L2Ok+UPJBOnXkjqMP+C6GLXe2BQqNUsuT4kooyB0hSmUHXXoHwrFa3/NR0S8WWT9HoWjUbP7mklYJt7HSwJN0h/NpaF0+FsBMOE37pEp9wi2jANOEVzuFGm8pbtfxZ6FSQwoeGUu5OS8ULf2+TGvDl9eCmbpr1df41SxLO/jLVNh", "syssn": "20251021070100020000072424", "expired_time": "", "sysdtm": "2025-10-21 17:10:49", "goods_name": "dasd-\\u786e\\u8ba4\\u7801:2424", "access_token": "", "userid": 1130081517, "alipay_pcode": "", "udid": "869071055892963", "appid": "", "return_url": "", "store_id": "1130081517", "businm": "", "tip_amt": 0, "currency": "344", "chnlid": 0, "div_template_id": "", "clisn": "000276", "alipay_params": "", "auth_code": "", "out_trade_no": "", "customer_id": "", "groupid": 0, "clientip": "", "openid": "", "goods_detail": "", "txcurrcd": "344", "lnglat": ["0", "0"], "failed_url": "", "txzone": "+0800", "lang": "", "mch_infos": {"sub_mchntid": "12345678", "key2": "", "mchnt_name": "dasd", "sub_appid": "", "memo": "", "tag": "", "postcode": "", "key": "", "mcc": "", "city": "", "userid": 1130081517, "ext": null, "appid": "", "groupid": 1130000007, "mchntid": "188000344333"}, "swipe_method": "", "txamt": 100, "busicd": "802808", "phonemodel": "A8S", "goods_info": "", "os": "Android"}`
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
