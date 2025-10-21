package e

var MsgFlags = map[string]string{
	SUCCESS:             "成功",
	ERROR:               "内部错误",
	PARAMERR:            "参数错误",
	DATAERR:             "数据错误",
	REFERRAL:            "Referral",
	INVALID_TRANSACTION: "Invalid Transaction",
}

func GetMsg(code string) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
