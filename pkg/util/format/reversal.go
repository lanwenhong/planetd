package format

import (
	"context"
	"strings"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func (rd *RequestData) PacketReversal(ctx context.Context) (string, error) {
	packet := ""
	var err error

	chnlExtData := rd.ChnlExt
	ph := planet_8583.NewProtoHandler()
	pData := &planet_8583.ProtoStruct{
		MsgType:      "0400",
		ProcessingCd: chnlExtData.ProcessingCd,
		Syssn:        chnlExtData.OrigSTAN,
		PosEntryMode: chnlExtData.EntryMode,
	}
	logger.Debugf(ctx, "request pData: %+v", pData)
	pData.Domain63Tags = make(map[string][]byte)
	tagFA := &planet_8583.TagFA{
		Len:                "0003",
		Tag:                "FA",
		FinalAuthIndicator: chnlExtData.TagFA,
	}

	tagTC := &planet_8583.TagTC{
		Len:                       "0003",
		Tag:                       "TC",
		TerminalEntryCapabilities: chnlExtData.TagTC,
	}

	ph.RegisterD63Tag(ctx, "FA", pData, tagFA)
	ph.RegisterD63Tag(ctx, "TC", pData, tagTC)

	for _, k := range pData.Domain63TagKey {
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
