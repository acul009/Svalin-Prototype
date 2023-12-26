package server

import (
	"log"

	"github.com/rahn-it/svalin/pki"
	"github.com/rahn-it/svalin/system"
	"github.com/rahn-it/svalin/util"
)

var _ util.ObservableMap[string, *system.DeviceInfo] = (*DeviceList)(nil)

type DeviceList struct {
	observerHandler *util.MapObserverHandler[string, *system.DeviceInfo]
	deviceStore     *deviceStore
	online          map[string]bool
}

func newDeviceList(deviceStore *deviceStore) *DeviceList {
	return &DeviceList{
		observerHandler: util.NewMapObserverHandler[string, *system.DeviceInfo](),
		deviceStore:     deviceStore,
		online:          make(map[string]bool),
	}
}

func (d *DeviceList) isOnline(key string) bool {
	online, ok := d.online[key]
	if !ok {
		return false
	}

	if ok && !online {
		delete(d.online, key)
	}
	return online
}

func (d *DeviceList) ForEach(fn func(key string, value *system.DeviceInfo) error) error {
	return d.deviceStore.ForEach(func(key string, cert *pki.Certificate) error {
		di := &system.DeviceInfo{
			Certificate: cert,
			LiveInfo: system.LiveDeviceInfo{
				Online: d.isOnline(key),
			},
		}

		return fn(key, di)
	})
}

func (d *DeviceList) Subscribe(onUpdate func(string, *system.DeviceInfo), onRemove func(string, *system.DeviceInfo)) func() {
	return d.observerHandler.Subscribe(onUpdate, onRemove)
}

func (d *DeviceList) setOnlineStatus(key string, online bool) {
	pubKey, err := pki.PublicKeyFromBase64(key)
	if err != nil {
		log.Printf("Error parsing public key: %v", err)
		return
	}

	cert, err := d.deviceStore.GetDevice(pubKey)
	if err != nil {
		log.Printf("Error getting device: %v", err)
		return
	}

	if online {
		d.online[key] = true
	} else {
		delete(d.online, key)
	}

	di := &system.DeviceInfo{
		Certificate: cert,
		LiveInfo: system.LiveDeviceInfo{
			Online: d.isOnline(key),
		},
	}

	d.observerHandler.NotifyUpdate(key, di)
}
