package handlers

import (
	"context"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/gen/converter"
	device "github.com/tigorlazuardi/redmage/gen/proto/device/v1"
	deviceConnect "github.com/tigorlazuardi/redmage/gen/proto/device/v1/v1connect"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

var (
	_         deviceConnect.DeviceServiceHandler = (*Device)(nil)
	convert   converter.ConverterImpl
	validator = func() *protovalidate.Validator {
		v, err := protovalidate.New()
		if err != nil {
			panic(err)
		}
		return v
	}()
)

type Device struct {
	deviceConnect.UnimplementedDeviceServiceHandler

	API *api.API
}

func (de *Device) CreateDevice(ctx context.Context, req *connect.Request[device.CreateDeviceRequest]) (*connect.Response[device.CreateDeviceResponse], error) {
	ctx, span := tracer.Start(ctx, "Device.CreateDevice")
	defer span.End()

	if err := validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	requestDao := convert.CreateDeviceRequestToModelsDevice(req.Msg)

	dev, err := de.API.DevicesCreate(ctx, requestDao)
	if err != nil {
		code := errs.ExtractConnectCode(err)
		return nil, connect.NewError(code, err)
	}

	createResponse := convert.ModelsDeviceToCreateDeviceResponse(dev)
	return connect.NewResponse(createResponse), nil
}
