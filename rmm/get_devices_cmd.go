package rmm

import (
	"rahnit-rmm/rpc"
	"rahnit-rmm/util"
)

type getDevicesCommand struct {
	*SyncDownCommand[string, *DeviceInfo]
}

func CreateGetDevicesCommandHandler(m util.UpdateableMap[string, *DeviceInfo]) rpc.RpcCommandHandler {
	return func() rpc.RpcCommand {
		return &getDevicesCommand{
			SyncDownCommand: NewSyncDownCommand[string, *DeviceInfo](m),
		}
	}
}

func NewGetDevicesCommand(targetMap util.UpdateableMap[string, *DeviceInfo]) *getDevicesCommand {
	return &getDevicesCommand{
		SyncDownCommand: NewSyncDownCommand[string, *DeviceInfo](targetMap),
	}
}

func (c *getDevicesCommand) GetKey() string {
	return "get-devices"
}
