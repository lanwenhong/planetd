package format

// TradeType 交易类型常量
const (
	TRADE_TYPE_TRADE       = "trade"
	TRADE_TYPE_REFUND      = "refund"
	TRADE_TYPE_TIP         = "tip"
	TRADE_TYPE_VOID_TRADE  = "void_trade"
	TRADE_TYPE_VOID_REFUND = "void_refund"
	TRADE_TYPE_KEY_IN      = "key_in"

	TRADE_TYPE_REVERSAL_TRADE       = "reversal_trade"
	TRADE_TYPE_REVERSAL_REFUND      = "reversal_refund"
	TRADE_TYPE_REVERSAL_VOID_TRADE  = "reversal_void_trade"
	TRADE_TYPE_REVERSAL_VOID_REFUND = "reversal_void_refund"
)

// ProcessingCode 处理码常量
const (
	TRADE_TYPE_TRADE_PROSSING_CODE       = "000000"
	TRADE_TYPE_REFUND_PROSSING_CODE      = "200000"
	TRADE_TYPE_VOID_TRADE_PROSSING_CODE  = "020000"
	TRADE_TYPE_VOID_REFUND_PROSSING_CODE = "220000"
	TRADE_TYPE_TRADE_TIPS_PROSSING_CODE  = "020000"
)

// TradeTypeProcessingCodeMap 交易类型与处理码的映射关系
var TradeTypeProcessingCodeMap = map[string]string{
	TRADE_TYPE_TRADE:       TRADE_TYPE_TRADE_PROSSING_CODE,
	TRADE_TYPE_REFUND:      TRADE_TYPE_REFUND_PROSSING_CODE,
	TRADE_TYPE_TIP:         TRADE_TYPE_TRADE_TIPS_PROSSING_CODE,
	TRADE_TYPE_VOID_TRADE:  TRADE_TYPE_VOID_TRADE_PROSSING_CODE,
	TRADE_TYPE_VOID_REFUND: TRADE_TYPE_VOID_REFUND_PROSSING_CODE,
	TRADE_TYPE_KEY_IN:      TRADE_TYPE_TRADE_PROSSING_CODE,
}

// GetProcessingCode 根据交易类型获取处理码
func GetProcessingCode(tradeType string) (string, bool) {
	processingCd, ok := TradeTypeProcessingCodeMap[tradeType]
	return processingCd, ok
}
