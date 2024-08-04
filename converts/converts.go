package main

import (
	device "github.com/tigorlazuardi/redmage/gen/proto/device/v1"
	"github.com/tigorlazuardi/redmage/models"
)

// goverter:converter
// goverter:extend BoolToInt32
type Converter interface {
	// goverter:map Nsfw NSFW
	// goverter:ignore CreatedAt UpdatedAt R
	CreateDeviceRequestToModelsDevice(source *device.CreateDeviceRequest) *models.Device

	// goverter:ignore state sizeCache unknownFields
	ModelsDeviceToCreateDeviceResponse(source *models.Device) *device.CreateDeviceResponse
}

func BoolToInt32(input bool) int32 {
	if input {
		return 1
	}
	return 0
}
