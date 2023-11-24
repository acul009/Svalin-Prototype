package rmm

import (
	"context"
	"log"
	"rahnit-rmm/pki"
	"rahnit-rmm/rpc"
	"rahnit-rmm/util"
)

type Client struct {
	ep            *rpc.RpcEndpoint
	tunnelHandler *tunnelHandler
	devices       util.ObservableMap[string, *Device]
}

func ClientConnect(ctx context.Context, credentials *pki.PermanentCredentials) (*Client, error) {
	ep, err := rpc.ConnectToUpstream(ctx, credentials)
	if err != nil {
		return nil, err
	}

	c := &Client{
		ep: ep,
	}

	c.initSyncedDeviceList()

	th := newTunnelHandler(ep)
	c.tunnelHandler = th

	return c, nil
}

func (c *Client) Tunnels() *tunnelHandler {
	return c.tunnelHandler
}

func (c *Client) Devices() util.ObservableMap[string, *Device] {
	return c.devices
}

func (c *Client) dispatch() rpc.Dispatcher {
	return c.ep
}

func (c *Client) Close() error {
	return c.ep.Close(200, "Shutdown")
}

func (c *Client) initSyncedDeviceList() {

	var dRunning util.AsyncAction

	devicesInfo := util.NewSyncedMap[string, *DeviceInfo](
		func(m util.ObservableMap[string, *DeviceInfo]) {
			cmd := NewGetDevicesCommand(m)
			running, err := c.dispatch().SendCommand(context.Background(), cmd)
			if err != nil {
				log.Printf("Error subscribing to devices: %v", err)
				return
			}
			dRunning = running
		},
		func(_ util.ObservableMap[string, *DeviceInfo]) {
			err := dRunning.Close()
			if err != nil {
				log.Printf("Error unsubscribing from devices: %v", err)
			}
		},
	)

	var unsub func()

	devices := util.NewSyncedMap[string, *Device](
		func(m util.ObservableMap[string, *Device]) {
			unsub = devicesInfo.Subscribe(
				func(s string, di *DeviceInfo) {
					m.Update(s, func(d *Device, found bool) (*Device, bool) {
						if !found {
							d = &Device{
								c: c,
							}
						}

						d.DeviceInfo = di
						return d, true
					})
				},
				func(s string, _ *DeviceInfo) {
					m.Delete(s)
				},
			)
		},
		func(m util.ObservableMap[string, *Device]) {
			unsub()
		},
	)

	c.devices = devices
}
