package main

import (
	device "github.com/tigorlazuardi/redmage/gen/proto/device/v1"
	"github.com/tigorlazuardi/redmage/models"
)

// goverter:converter
// goverter:extend BoolToInt32
type Converter interface {
	// goverter:map Nsfw NSFW
	// goverter:ignore CreatedAt
	// goverter:ignore UpdatedAt
	// goverter:ignore R
	ConvertCreateDeviceRequestToModelsDevice(source device.CreateDeviceRequest) *models.Device
}

func BoolToInt32(input bool) int32 {
	if input {
		return 1
	}
	return 0
}
