package server

import (
	"fmt"

	"github.com/rahn-it/svalin/config"
	"github.com/rahn-it/svalin/pki"
	"github.com/rahn-it/svalin/rmm"
	"github.com/rahn-it/svalin/rpc"
	"github.com/rahn-it/svalin/util"

	"github.com/google/uuid"
)

type Server struct {
	*rpc.RpcServer
	profile       *config.Profile
	devices       util.UpdateableMap[string, *DeviceInfo]
	configManager *ConfigManager
}

func Open(profile *config.Profile) (*Server, error) {

	verifier, err := pki.NewLocalVerify()
	if err != nil {
		return nil, fmt.Errorf("error creating local verify: %w", err)
	}

	ConfigManager := NewConfigManager(verifier, nil)

	devices := NewDeviceList(profile.Scope().Scope([]byte("devices")))

	cmds := rpc.NewCommandCollection(
		rpc.PingHandler,
		rpc.RegisterUserHandler,
		// rpc.GetPendingEnrollmentsHandler,
		rpc.EnrollAgentHandler,
		rmm.CreateGetDevicesCommandHandler(devices),
		rpc.ForwardCommandHandler,
		rpc.VerifyCertificateChainHandler,
		// CreateHostConfigCommandHandler[*TunnelConfig],
	)

	rpcS, err := rpc.NewRpcServer(listenAddr, cmds, verifier, credentials)
	if err != nil {
		return nil, fmt.Errorf("error creating rpc server: %w", err)
	}

	rpcS.Connections().Subscribe(
		func(_ uuid.UUID, rc *rpc.RpcConnection) {
			key := rc.Partner().PublicKey().Base64Encode()
			devices.Update(key, func(d *DeviceInfo, found bool) (*DeviceInfo, bool) {
				if !found {
					return nil, false
				}

				d.Online = true
				return d, true
			})
		},
		func(_ uuid.UUID, rc *rpc.RpcConnection) {
			key := rc.Partner().PublicKey().Base64Encode()
			devices.Update(key, func(d *DeviceInfo, found bool) (*DeviceInfo, bool) {
				if !found {
					return nil, false
				}

				d.Online = false
				return d, true
			})
		},
	)

	s := &Server{
		RpcServer:     rpcS,
		devices:       devices,
		configManager: ConfigManager,
	}

	return s, nil
}

func (s *Server) Run() error {
	return s.RpcServer.Run()
}
