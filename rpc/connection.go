package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rahnit-rmm/pki"
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

type rpcConnection struct {
	connection     quic.Connection
	server         *RpcServer
	partner        *pki.Certificate
	Uuid           uuid.UUID
	state          RpcConnectionState
	role           RpcConnectionRole
	activeSessions map[uuid.UUID]*RpcSession
	mutex          sync.Mutex
	nonceStorage   *nonceStorage
	protocol       TlsConnectionProto
}

func newRpcConnection(conn quic.Connection,
	server *RpcServer,
	role RpcConnectionRole,
	nonceStorage *nonceStorage,
	partner *pki.Certificate,
	protocol TlsConnectionProto,
) *rpcConnection {
	return &rpcConnection{
		connection:     conn,
		server:         server,
		partner:        partner,
		Uuid:           uuid.New(),
		state:          RpcConnectionOpen,
		role:           role,
		activeSessions: make(map[uuid.UUID]*RpcSession),
		mutex:          sync.Mutex{},
		nonceStorage:   nonceStorage,
		protocol:       protocol,
	}
}

func (conn *rpcConnection) serve(commands *CommandCollection) error {
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
		go session.handleIncoming(commands)
	}

}

func (conn *rpcConnection) AcceptSession(context.Context) (*RpcSession, error) {
	stream, err := conn.connection.AcceptStream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error accepting QUIC stream: %w", err)
	}
	var session *RpcSession = nil

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	for i := 0; i < 10; i++ {
		newSession := newRpcSession(stream, conn)
		if _, ok := conn.activeSessions[newSession.Uuid]; !ok {
			session = newSession
			break
		}
	}

	if session == nil {
		return nil, fmt.Errorf("multiple uuid collisions, this should mathematically be impossible")
	}

	conn.activeSessions[session.Uuid] = session

	return session, nil
}

func (conn *rpcConnection) OpenSession(ctx context.Context) (*RpcSession, error) {
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

func (conn *rpcConnection) MutateState(from RpcConnectionState, to RpcConnectionState) error {
	conn.mutex.Lock()
	if conn.state != from {
		conn.mutex.Unlock()
		return fmt.Errorf("RPC session not in state %v", from)
	}
	conn.state = to
	conn.mutex.Unlock()
	return nil
}

func (conn *rpcConnection) EnsureState(state RpcConnectionState) error {
	conn.mutex.Lock()
	if conn.state != state {
		conn.mutex.Unlock()
		return fmt.Errorf("RPC session not in state %v", state)
	}
	conn.mutex.Unlock()
	return nil
}

func (conn *rpcConnection) removeSession(uuid uuid.UUID) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	delete(conn.activeSessions, uuid)
}

func (conn *rpcConnection) Close(code quic.ApplicationErrorCode, msg string) error {

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
		go func(session *RpcSession) {
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
		conn.server.removeConnection(conn.Uuid)
	}

	return err
}
