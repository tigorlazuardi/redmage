package handlers

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tigorlazuardi/redmage/api"
	device "github.com/tigorlazuardi/redmage/gen/proto/device/v1"
	deviceConnect "github.com/tigorlazuardi/redmage/gen/proto/device/v1/v1connect"
)

var _ deviceConnect.DeviceServiceHandler = (*Device)(nil)

type Device struct {
	deviceConnect.UnimplementedDeviceServiceHandler

	API *api.API
}

func (de *Device) CreateDevice(ctx context.Context, req *connect.Request[device.CreateDeviceRequest]) (*connect.Response[device.CreateDeviceResponse], error) {
	ctx, span := tracer.Start(ctx, "Device.CreateDevice")
	defer span.End()

	panic("not implemented") // TODO: Implement
}
