package e

var MsgFlags = map[string]string{
	SUCCESS: "成功",
	ERROR:   "内部错误",
}

func GetMsg(code string) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
