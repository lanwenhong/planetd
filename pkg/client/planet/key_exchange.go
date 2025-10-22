package planet

import (
	"context"
	"planetd/pkg/util/format"

	"github.com/lanwenhong/lgobase/logger"
	"github.com/lanwenhong/planet_8583/planet_8583"
)

func DoExchange(ctx context.Context, rd *format.RequestData) (*planet_8583.ProtoStruct, error) {
	bcd, err := rd.Request2KeyExchangePacket(ctx)
	if err != nil {
		return nil, err
	}
	ret, err := SendPacket(ctx, bcd)
	if err != nil {
		return nil, err
	}
	logger.Debugf(ctx, "ret: %s", ret)
	ps, err := rd.Packet2Request(ctx, ret[10:])
	if err != nil {
		return nil, err
	}
	logger.Debugf(ctx, "ps: %+v", ps)
	return ps, nil
}
