package format

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

type RequestData struct {
	AccessToken   string       `json:"access_token"`
	AlipayParams  string       `json:"alipay_params"`
	AlipayPcode   string       `json:"alipay_pcode"`
	Appid         string       `json:"appid"`
	Appver        string       `json:"appver"`
	AuthCode      string       `json:"auth_code"`
	Busicd        string       `json:"busicd"`
	Businm        string       `json:"businm"`
	ChannelInfos  ChannelInfo  `json:"channel_infos"`
	ChnlExt       string       `json:"pos_ext"`
	Chnlid        int          `json:"chnlid"`
	Clientip      string       `json:"clientip"`
	Clisn         string       `json:"clisn"`
	Code          string       `json:"code"`
	CodeVerifier  string       `json:"code_verifier"`
	Currency      string       `json:"currency"`
	CustomerID    string       `json:"customer_id"`
	DivTemplateID string       `json:"div_template_id"`
	ExpiredTime   string       `json:"expired_time"`
	FailedURL     string       `json:"failed_url"`
	GoodsDetail   string       `json:"goods_detail"`
	GoodsInfo     string       `json:"goods_info"`
	GoodsName     string       `json:"goods_name"`
	Groupid       int          `json:"groupid"`
	Lang          string       `json:"lang"`
	LimitPay      string       `json:"limit_pay"`
	Lnglat        []string     `json:"lnglat"`
	MchInfos      MerchantInfo `json:"mch_infos"`
	Openid        string       `json:"openid"`
	Opuid         string       `json:"opuid"`
	Os            string       `json:"os"`
	OutTradeNo    string       `json:"out_trade_no"`
	PayTag        string       `json:"pay_tag"`
	Paysource     string       `json:"paysource"`
	Phonemodel    string       `json:"phonemodel"`
	PrepayAmt     int          `json:"prepay_amt"`
	ReturnURL     string       `json:"return_url"`
	StoreID       string       `json:"store_id"`
	SwipeMethod   string       `json:"swipe_method"`
	Sysdtm        string       `json:"sysdtm"`
	Syssn         string       `json:"syssn"`
	TipAmt        int          `json:"tip_amt"`
	Txamt         int          `json:"txamt"`
	Txcurrcd      string       `json:"txcurrcd"`
	Txdtm         string       `json:"txdtm"`
	Txzone        string       `json:"txzone"`
	Udid          string       `json:"udid"`
	Userid        int          `json:"userid"`
	OrigChnlSn    string       `json:"orig_chnlsn"`
}

const (
	TRADE_TYPE_TRADE                     = "trade"
	TRADE_TYPE_REFUND                    = "refund"
	TRADE_TYPE_VOID_TRADE                = "void_trade"
	TRADE_TYPE_VOID_REFUND               = "void_refund"
	TRADE_TYPE_TRADE_PROSSING_CODE       = "000000"
	TRADE_TYPE_REFUND_PROSSING_CODE      = "200000"
	TRADE_TYPE_VOID_TRADE_PROSSING_CODE  = "020000"
	TRADE_TYPE_VOID_REFUND_PROSSING_CODE = "220000"
)

func (rd *RequestData) Request2TransactionPacket(ctx context.Context) (string, error) {
	packet := ""

	TradeTypeProcessingCode := make(map[string]string)
	TradeTypeProcessingCode[TRADE_TYPE_TRADE] = TRADE_TYPE_TRADE_PROSSING_CODE
	TradeTypeProcessingCode[TRADE_TYPE_REFUND] = TRADE_TYPE_REFUND_PROSSING_CODE
	TradeTypeProcessingCode[TRADE_TYPE_VOID_TRADE] = TRADE_TYPE_VOID_TRADE_PROSSING_CODE
	TradeTypeProcessingCode[TRADE_TYPE_VOID_REFUND] = TRADE_TYPE_VOID_REFUND_PROSSING_CODE

	chnlExtData, err := rd.ToChnlExtData(ctx)
	if err != nil {
		return packet, err
	}
	ph := planet_8583.NewProtoHandler()
	biccdata, err := hex.DecodeString(chnlExtData.Iccdata)
	if err != nil {
		logger.Infof(ctx, "DecodeString err: %s", err.Error())
		return packet, err
	}
	processingCd, flag := TradeTypeProcessingCode[chnlExtData.TradeType]
	if !flag {
		logger.Infof(ctx, "TradeTypeProcessingCode err: %s", chnlExtData.TradeType)
		err = errors.New("trade_type err")
		return packet, err
	}
	pData := &planet_8583.ProtoStruct{
		MsgType: "0200",
		// CardNo:               chnlExtData.CardNo,
		ProcessingCd: processingCd,
		Txamt:        strconv.Itoa(rd.Txamt),
		Syssn:        rd.Clisn,
		NetId:        "226",
		PosCondCd:    "00",
		// TrackData2:           chnlExtData.Track2,
		Tid:                rd.MchInfos.SubMchntid,
		MchntId:            rd.MchInfos.Mchntid,
		CurrencyCd:         rd.Currency,
		Cardsequencenumber: chnlExtData.Cardseqnum,
		PosEntryMode:       chnlExtData.EntryMode,
		// ICCSystemRelatedData: biccdata,
		// Pin:          strings.ToUpper(chnlExtData.PinBlock),
	}
	if chnlExtData.TradeType == TRADE_TYPE_TRADE || chnlExtData.TradeType == TRADE_TYPE_REFUND {
		pData.CardNo = chnlExtData.CardNo
		pData.TrackData2 = chnlExtData.Track2
		pData.ICCSystemRelatedData = biccdata
	}
	if chnlExtData.TradeType == TRADE_TYPE_REFUND || chnlExtData.TradeType == TRADE_TYPE_VOID_REFUND || chnlExtData.TradeType == TRADE_TYPE_VOID_TRADE {
		// TODO 从哪里获取原交易通道返回的流水号
		pData.RetrievalReferenceNumber = rd.OrigChnlSn
	}
	logger.Debugf(ctx, "request pData: %+v", pData)
	pData.Domain63Tags = make(map[string][]byte)

	tag12 := &planet_8583.Tag12{
		Len:       "0003",
		Tag:       "12",
		IndiCator: "X",
	}

	tagIA := &planet_8583.TagIA{
		Len:          "0004",
		Tag:          "IA",
		HostKeyIndex: "220",
	}

	tagIB := &planet_8583.TagIB{
		Len:            "0006",
		Tag:            "IB",
		MacCheckDigits: "F9EA",
	}

	tagIC := &planet_8583.TagIC{
		Len:                  "0003",
		Tag:                  "IC",
		InteracTerminalClass: "03",
	}
	tagID := &planet_8583.TagID{
		Len:                    "0003",
		Tag:                    "ID",
		InteracCustomerPresent: "1",
	}

	tagIE := &planet_8583.TagIE{
		Len:                "0003",
		Tag:                "IE",
		InteracCardPresent: "0",
	}

	tagIF := &planet_8583.TagIF{
		Len:                          "0003",
		Tag:                          "IF",
		InteracCardCaptureCapability: "0",
	}

	tagIG := &planet_8583.TagIG{
		Len:               "0003",
		Tag:               "IG",
		BalanceinResponse: "0",
	}

	tagIH := &planet_8583.TagIH{
		Len:             "0003",
		Tag:             "IH",
		InteracSecurity: "0",
	}

	tagIL := &planet_8583.TagIL{
		Len:             "0010",
		Tag:             "IL",
		InteracSecurity: "0000702940000850",
	}

	ph.RegisterD63Tag(ctx, "12", pData, tag12)
	ph.RegisterD63Tag(ctx, "IA", pData, tagIA)
	ph.RegisterD63Tag(ctx, "IB", pData, tagIB)
	ph.RegisterD63Tag(ctx, "IC", pData, tagIC)
	ph.RegisterD63Tag(ctx, "ID", pData, tagID)
	ph.RegisterD63Tag(ctx, "IE", pData, tagIE)
	ph.RegisterD63Tag(ctx, "IF", pData, tagIF)
	ph.RegisterD63Tag(ctx, "IG", pData, tagIG)
	ph.RegisterD63Tag(ctx, "IH", pData, tagIH)
	ph.RegisterD63Tag(ctx, "IL", pData, tagIL)

	for _, k := range pData.Domain64TagKey {
		logger.Debugf(ctx, "tag: %s", k)
	}

	_, err = ph.PackStru(ctx, pData)
	if err != nil {
		logger.Infof(ctx, "PackStru err: %s", err.Error())
		return packet, err
	}
	err = ph.PackMac(ctx, "BBEFB74400000000")
	if err != nil {
		logger.Debugf(ctx, "PackMac err: %s", err.Error())
		return packet, err
	}
	ph.Pack(ctx)
	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)

	pd := &PacketData{}
	finishFd, err := pd.BuildPacket(ctx, ph.Tbuf)
	packet = strings.ToUpper(string(finishFd))
	logger.Debugf(ctx, "finish packet: %s", packet)
	return packet, err
}

func (rd *RequestData) Request2KeyExchangePacket(ctx context.Context) (string, error) {
	packet := ""
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0800",
		ProcessingCd: "920000",
		Syssn:        rd.Clisn,
		Tid:          rd.MchInfos.SubMchntid,
		MchntId:      rd.MchInfos.Mchntid,
		NetId:        "226",
	}
	tagIL := &planet_8583.TagIL{
		Len:             "0010",
		Tag:             "IL",
		InteracSecurity: "0000702940000850", // TODO通道给
	}

	tagPP := &planet_8583.TagPP{
		Len:                   "0018",
		Tag:                   "PP",
		PlanetPaymentPassword: "24504C414E4554245041594D454E5424", // TODO通道给
	}

	ph.RegisterD63Tag(ctx, "IL", pData, tagIL)
	ph.RegisterD63Tag(ctx, "PP", pData, tagPP)

	_, err := ph.PackStru(ctx, pData)
	if err != nil {
		logger.Infof(ctx, "PackStru err: %s", err.Error())
		return packet, err
	}

	ph.Pack(ctx)

	fs := planet_8583.FormatByte(ctx, ph.Tbuf)
	logger.Debugf(ctx, "bcd: %s", fs)
	pd := &PacketData{}
	finishFd, err := pd.BuildPacket(ctx, ph.Tbuf)
	packet = strings.ToUpper(string(finishFd))
	logger.Debugf(ctx, "finish packet: %s", packet)
	return packet, err
}

func DecryptChnlExt(ciphertext, key, iv string) (string, error) {
	// Base64解码
	ct, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %v", err)
	}

	// 创建AES cipher
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("aes new cipher error: %v", err)
	}

	// 检查IV长度
	if len(iv) != aes.BlockSize {
		return "", errors.New("IV length must equal block size")
	}

	// 创建CBC解密器
	mode := cipher.NewCBCDecrypter(block, []byte(iv))

	// 解密
	plaintext := make([]byte, len(ct))
	mode.CryptBlocks(plaintext, ct)

	// 去除PKCS#7 padding
	plaintext, err = unpadPKCS7(plaintext)
	if err != nil {
		return "", fmt.Errorf("unpad error: %v", err)
	}

	return string(plaintext), nil
}

// PKCS#7 unpadding
func unpadPKCS7(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[len(data)-1])
	if padding < 1 || padding > aes.BlockSize {
		return nil, errors.New("invalid padding")
	}
	return data[:len(data)-padding], nil
}

func (rd *RequestData) ToChnlExtData(ctx context.Context) (*ChnlExtData, error) {
	key := "hV5kK6boR+9FosqRPBKB2XosT4u1c608"
	iv := "qfpay202302_hjsh"
	chnlExtStr, err := DecryptChnlExt(rd.ChnlExt, key, iv)
	logger.Debugf(ctx, "chnlExtStr: %s, err: %v", chnlExtStr, err)
	if err != nil {
		return nil, err
	}
	chnlExtData := &ChnlExtData{}
	if err := json.Unmarshal([]byte(chnlExtStr), chnlExtData); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %v", err)
	}
	return chnlExtData, nil
}

type ChnlExtData struct {
	Cardseqnum string `json:"cardseqnum"`
	Iccdata    string `json:"iccdata"`
	MacString  string `json:"macString"`
	PinBlock   string `json:"pinBlock"`
	TrackData  string `json:"trackData"`
	CardNo     string `json:"cardNo"`
	Track1     string `json:"track1"`
	Track2     string `json:"track2"`
	Track3     string `json:"track3"`
	Terminalid string `json:"terminalid"`
	TradeType  string `json:"trade_type"`
	EntryMode  string `json:"entry_mode"`
}

type ChannelInfo struct {
	Mackey  string `json:"mackey"`
	Partner string `json:"partner"`
	Zmk     string `json:"zmk"`
	Zpk     string `json:"zpk"`
}

type MerchantInfo struct {
	Appid      string      `json:"appid"`
	City       string      `json:"city"`
	Ext        interface{} `json:"ext"`
	Groupid    int         `json:"groupid"`
	Key        string      `json:"key"`
	Key2       string      `json:"key2"`
	Mcc        string      `json:"mcc"`
	MchntName  string      `json:"mchnt_name"`
	Mchntid    string      `json:"mchntid"`
	Memo       string      `json:"memo"`
	Postcode   string      `json:"postcode"`
	SubAppid   string      `json:"sub_appid"`
	SubMchntid string      `json:"sub_mchntid"`
	Tag        string      `json:"tag"`
	Userid     int         `json:"userid"`
}

type PacketData struct {
}

func (pd *PacketData) BuildPacket(ctx context.Context, tbuf []byte) ([]byte, error) {
	tpdu := []byte{0x60, 0x00, 0x00, 0x00, 0x10}

	length := len(tpdu) + len(tbuf)
	bufLength := make([]byte, 2)
	binary.BigEndian.PutUint16(bufLength, uint16(length))

	// 转换为BCD格式
	fb := []byte{}
	bcdLen := hex.EncodeToString(bufLength)
	for i := 0; i < len(bcdLen); i++ {
		fb = append(fb, byte(bcdLen[i]))
	}
	bcdTPDU := hex.EncodeToString(tpdu)
	for i := 0; i < len(bcdTPDU); i++ {
		fb = append(fb, byte(bcdTPDU[i]))
	}
	bcd := hex.EncodeToString(tbuf)
	for i := 0; i < len(bcd); i++ {
		fb = append(fb, byte(bcd[i]))
	}
	return fb, nil
}

func (rd *RequestData) Packet2Request(ctx context.Context, packet string) (*planet_8583.ProtoStruct, error) {
	logger.Infof(ctx, "packet: %s", packet)
	b, _ := hex.DecodeString(packet)
	uph := planet_8583.NewProtoHandler()
	ups := &planet_8583.ProtoStruct{}
	err := uph.Unpack(ctx, b, ups)
	if err != nil {
		logger.Debugf(ctx, "unpack err: %s", err.Error())
		return ups, err
	}
	return ups, nil
}
