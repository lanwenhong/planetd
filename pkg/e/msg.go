package e

var MsgFlags = map[string]string{
	SUCCESS:               "成功",
	ERROR:                 "内部错误",
	PARAMERR:              "参数错误",
	DATAERR:               "数据错误",
	REFERRAL:              "Referral",
	ReferSpecialCondition: "Refer Special Condition",
	InvalidAcquirer:       "Invalid Acquirer",
	AlreadyReversed:       "Already Reversed",
	PartialAuthorization:  "Partial Authorization",
	InvalidAmount:         "Invalid Amount",
	InvalidAccountNumber:  "Invalid Account Number",
	InvalidIssuer:         "Invalid Issuer",
	ReEnterTransaction:    "Re-Enter Transaction",
	INVALID_TRANSACTION:   "Invalid Transaction",
	CANNOT_VERIFY_PIN:     "Cannot Verify PIN",
}

func GetMsg(code string) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
