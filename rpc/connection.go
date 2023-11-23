package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rahnit-rmm/pki"
	"rahnit-rmm/util"
	"sync"

	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
)

type RpcConnectionState int16

const (
	RpcConnectionOpen RpcConnectionState = iota
	RpcConnectionStopped
)

type RpcConnectionRole int16

const (
	RpcRoleAgent RpcConnectionRole = iota
	RpcRoleServer
	RpcRoleClient
	RpcRoleInit
)

type RpcConnection[T any] struct {
	deps           T
	connection     quic.Connection
	server         *RpcServer[T]
	partner        *pki.Certificate
	uuid           uuid.UUID
	state          RpcConnectionState
	role           RpcConnectionRole
	activeSessions map[quic.StreamID]*RpcSession[T]
	mutex          sync.Mutex
	nonceStorage   *util.NonceStorage
	protocol       TlsConnectionProto
	credentials    pki.Credentials
	verifier       pki.Verifier
}

func newRpcConnection[T any](conn quic.Connection,
	server *RpcServer[T],
	role RpcConnectionRole,
	nonceStorage *util.NonceStorage,
	partner *pki.Certificate,
	protocol TlsConnectionProto,
	credentials pki.Credentials,
	verifier pki.Verifier,
) *RpcConnection[T] {
	return &RpcConnection[T]{
		connection:     conn,
		server:         server,
		partner:        partner,
		uuid:           uuid.New(),
		state:          RpcConnectionOpen,
		role:           role,
		activeSessions: make(map[quic.StreamID]*RpcSession[T]),
		mutex:          sync.Mutex{},
		nonceStorage:   nonceStorage,
		protocol:       protocol,
		credentials:    credentials,
		verifier:       verifier,
	}
}

func (conn *RpcConnection[T]) serveRpc(commands *CommandCollection[T]) error {
	defer conn.Close(500, "RPC connection closed")

	if conn.partner == nil {
		return fmt.Errorf("no partner provided")
	}

	if conn.server != nil {
		conn.server.devices.UpdateDeviceStatus(conn.partner.GetPublicKey().Base64Encode(), func(device *DeviceInfo) *DeviceInfo {
			device.Online = true
			return device
		})
	}

	err := conn.EnsureState(RpcConnectionOpen)
	if err != nil {
		return fmt.Errorf("error ensuring RPC connection is open: %w", err)
	}

	if conn.protocol != ProtoRpc {
		return fmt.Errorf("tried to serve RPC to non-RPC connection")
	}

	log.Printf("Connection accepted, serving RPC")
	for {
		log.Printf("Waiting for incoming QUIC stream...")

		session, err := conn.AcceptSession(context.Background())

		log.Printf("Session requested")
		if err != nil {
			log.Printf("error accepting QUIC stream: %v\n", err)
		}

		stateErr := conn.EnsureState(RpcConnectionOpen)
		if stateErr != nil {
			log.Printf("error ensuring RPC connection is open: %v", stateErr)
			return fmt.Errorf("RPC connection not open anymore")
		}

		if err != nil {
			log.Printf("error accepting QUIC stream: %v", err)
			if errors.Is(err, &quic.ApplicationError{}) {
				return err
			}
			return nil
		}

		log.Printf("RPC session opened, handling incoming commands")
		go func() {
			err := session.handleIncoming(commands)
			if err != nil {
				log.Printf("error handling incoming session: %v", err)
			}
		}()
	}

}

func (conn *RpcConnection[T]) AcceptSession(context.Context) (*RpcSession[T], error) {
	stream, err := conn.connection.AcceptStream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error accepting QUIC stream: %w", err)
	}
	session := newRpcSession(stream, conn)

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	conn.activeSessions[stream.StreamID()] = session

	return session, nil
}

func (conn *RpcConnection[T]) OpenSession(ctx context.Context) (*RpcSession[T], error) {
	err := conn.EnsureState(RpcConnectionOpen)
	if err != nil {
		return nil, fmt.Errorf("error ensuring RPC connection is open: %w", err)
	}

	stream, err := conn.connection.OpenStreamSync(ctx)
	if err != nil {
		return nil, fmt.Errorf("error opening QUIC stream: %w", err)
	}

	return newRpcSession(stream, conn), nil
}

func (conn *RpcConnection[T]) MutateState(from RpcConnectionState, to RpcConnectionState) error {
	conn.mutex.Lock()
	if conn.state != from {
		conn.mutex.Unlock()
		return fmt.Errorf("RPC session not in state %v", from)
	}
	conn.state = to
	conn.mutex.Unlock()
	return nil
}

func (conn *RpcConnection[T]) EnsureState(state RpcConnectionState) error {
	conn.mutex.Lock()
	if conn.state != state {
		conn.mutex.Unlock()
		return fmt.Errorf("RPC session not in state %v", state)
	}
	conn.mutex.Unlock()
	return nil
}

func (conn *RpcConnection[T]) removeSession(id quic.StreamID) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	delete(conn.activeSessions, id)
}

func (conn *RpcConnection[T]) Close(code quic.ApplicationErrorCode, msg string) error {

	if err := conn.MutateState(RpcConnectionOpen, RpcConnectionStopped); err != nil {
		return fmt.Errorf("error closing connection: %w", err)
	}

	conn.mutex.Lock()
	sessionsToClose := conn.activeSessions
	conn.mutex.Unlock()

	// tell all connections to close
	errChan := make(chan error)
	wg := sync.WaitGroup{}

	errorList := make([]error, 0)

	for _, session := range sessionsToClose {
		wg.Add(1)
		go func(session *RpcSession[T]) {
			err := session.Close()
			if err != nil {
				errChan <- err
			}
			wg.Done()
		}(session)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		errorList = append(errorList, err)
	}

	err := conn.connection.CloseWithError(code, msg)
	if err != nil {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		err = fmt.Errorf("error closing sessions: %w", errors.Join(errorList...))
	}

	if conn.server != nil {
		conn.server.removeConnection(conn.uuid)
	}

	return err
}

func (conn *RpcConnection[T]) GetProtocol() TlsConnectionProto {
	return conn.protocol
}
