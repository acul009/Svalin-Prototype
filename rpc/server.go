package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rahn-it/svalin/pki"
	"github.com/rahn-it/svalin/util"

	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
)

type rpcNotRunningError struct {
}

func (e rpcNotRunningError) Error() string {
	return fmt.Errorf("rpc not running anymore").Error()
}

var ErrRpcNotRunning = rpcNotRunningError{}

func (e rpcNotRunningError) Is(target error) bool {
	_, ok := target.(rpcNotRunningError)
	return ok
}

type RpcServer struct {
	listener          *quic.Listener
	rpcCommands       *CommandCollection
	state             RpcServerState
	activeConnections util.UpdateableMap[uuid.UUID, *RpcConnection]
	mutex             sync.Mutex
	nonceStorage      *util.NonceStorage
	credentials       *pki.PermanentCredentials
	enrollment        *enrollmentManager
	verifier          pki.Verifier
	loginHandler      func(session *RpcSession) error
}

type RpcServerState int16

const (
	RpcServerCreated RpcServerState = iota
	RpcServerRunning
	RpcServerStopped
)

func NewRpcServer(listenAddr string, rpcCommands *CommandCollection, verifier pki.Verifier, credentials *pki.PermanentCredentials, root *pki.Certificate) (*RpcServer, error) {
	tlsConf, err := getTlsServerConfig([]TlsConnectionProto{ProtoRpc, ProtoClientLogin, ProtoAgentEnroll})
	if err != nil {
		return nil, fmt.Errorf("error getting server tls config: %w", err)
	}

	quicConf := &quic.Config{
		KeepAlivePeriod: 30 * time.Second,
	}
	listener, err := quic.ListenAddr(listenAddr, tlsConf, quicConf)
	if err != nil {
		return nil, fmt.Errorf("error creating QUIC server: %w", err)
	}

	return &RpcServer{
		listener:          listener,
		rpcCommands:       rpcCommands,
		state:             RpcServerCreated,
		activeConnections: util.NewObservableMap[uuid.UUID, *RpcConnection](),
		mutex:             sync.Mutex{},
		nonceStorage:      util.NewNonceStorage(),
		credentials:       credentials,
		enrollment:        newEnrollmentManager(credentials.Certificate(), root),
		verifier:          verifier,
	}, nil
}

func (s *RpcServer) accept() (*RpcConnection, error) {
	conn, err := s.listener.Accept(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error accepting QUIC connection: %w", err)
	}

	state := conn.ConnectionState()

	peerCertList := state.TLS.PeerCertificates
	if len(peerCertList) > 1 {
		conn.CloseWithError(400, "")
		return nil, fmt.Errorf("peer provided multiple certificates")
	}

	protocol := TlsConnectionProto(state.TLS.NegotiatedProtocol)
	var peerCert *pki.Certificate

	if len(peerCertList) == 1 {
		peerCert, err = pki.ImportCertificate(peerCertList[0])
		if err != nil {
			conn.CloseWithError(400, "")
			return nil, fmt.Errorf("error parsing peer certificate: %w", err)
		}

		log.Printf("Peer certificate:\n%s", peerCert.PemEncode())
	} else {
		peerCert = nil
	}

	switch protocol {
	case ProtoRpc:

		if peerCert == nil {
			conn.CloseWithError(400, "")
			return nil, fmt.Errorf("peer did not provide a any certificate")
		}

		if _, err := s.verifier.Verify(peerCert); err != nil {
			conn.CloseWithError(400, "")
			return nil, fmt.Errorf("peer did not provide a valid certificate: %w", err)
		}

		certType := peerCert.Type()

		switch certType {
		case pki.CertTypeUser, pki.CertTypeRoot, pki.CertTypeAgent:
			log.Println("Valid certificate type")
		default:
			conn.CloseWithError(400, "")
			return nil, fmt.Errorf("peer presented wrong certificate type: %s", certType)
		}

	case ProtoClientLogin, ProtoAgentEnroll:

		if peerCert != nil {
			conn.CloseWithError(400, "")
			return nil, fmt.Errorf("peer provided certificate when it should not have")
		}

	default:
		conn.CloseWithError(400, "")
		return nil, fmt.Errorf("unwanted connection type: wrong tls-next protocol")

	}

	var connection *RpcConnection

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := 0; i < 10; i++ {
		newConnection := newRpcConnection(conn, s, RpcRoleServer, s.nonceStorage, peerCert, protocol, s.credentials, s.verifier)
		if _, ok := s.activeConnections.Get(newConnection.uuid); !ok {
			connection = newConnection
			break
		}
	}

	if connection == nil {
		conn.CloseWithError(400, "")
		return nil, fmt.Errorf("multiple uuid collisions, this should mathematically be impossible")
	}

	s.activeConnections.Set(connection.uuid, connection)

	return connection, nil
}

func (s *RpcServer) Connections() util.ObservableMap[uuid.UUID, *RpcConnection] {
	return s.activeConnections
}

func (s *RpcServer) Enrollments() *enrollmentManager {
	return s.enrollment
}

func (s *RpcServer) removeConnection(uuid uuid.UUID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.activeConnections.Delete(uuid)
}

func (s *RpcServer) Run() error {
	fmt.Println("Starting RPC server")
	s.mutex.Lock()
	if s.state != RpcServerCreated {
		s.mutex.Unlock()
		return fmt.Errorf("RPC server already running")
	}
	s.state = RpcServerRunning
	s.mutex.Unlock()

	go func() {
		for s.state == RpcServerRunning {
			s.cleanup()
			time.Sleep(30 * time.Second)
		}
	}()

	for {
		conn, err := s.accept()

		s.mutex.Lock()
		if s.state != RpcServerRunning {
			s.mutex.Unlock()
			return rpcNotRunningError{}
		}
		s.mutex.Unlock()

		if err != nil {
			log.Printf("error accepting QUIC connection: %v", err)
			continue
		}

		switch conn.protocol {
		case ProtoRpc:
			go conn.serveRpc(s.rpcCommands)

		case ProtoClientLogin:
			if s.loginHandler == nil {
				log.Printf("login request received but login handler is not set")
				conn.Close(400, "login is disabled")
				continue
			}
			go func() {
				err := s.acceptLoginRequest(conn)
				if err != nil {
					log.Printf("error accepting login request: %v", err)
				}
			}()

		case ProtoAgentEnroll:
			go func() {
				err := s.enrollment.startEnrollment(conn)
				if err != nil {
					log.Printf("error accepting agent enroll request: %v", err)
				}
			}()

		default:
			log.Printf("error accepted connection has wrong protocol: %s", conn.protocol)
			conn.Close(400, "wrong protocol")
			continue
		}

		// TODO: check certificate

	}
}

func (s *RpcServer) Close(code quic.ApplicationErrorCode, msg string) error {

	// lock server before closing
	s.mutex.Lock()
	if s.state != RpcServerRunning {
		s.mutex.Unlock()
		return fmt.Errorf("RPC server not running")
	}
	s.state = RpcServerStopped
	s.mutex.Unlock()

	// tell all connections to close
	errChan := make(chan error)
	wg := sync.WaitGroup{}

	errorList := make([]error, 0)

	s.activeConnections.ForEach(func(_ uuid.UUID, connection *RpcConnection) error {
		wg.Add(1)
		go func(connection *RpcConnection) {
			err := connection.Close(code, msg)
			if err != nil {
				errChan <- err
			}
			wg.Done()
		}(connection)
		return nil
	})

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		errorList = append(errorList, err)
	}

	var err error = nil
	if len(errorList) > 0 {
		err = fmt.Errorf("error closing connections: %w", errors.Join(errorList...))
	}

	s.listener.Close()

	return err
}

func (s *RpcServer) cleanup() {
	s.enrollment.cleanup()
	s.nonceStorage.Cleanup(messageExpiration * 2)
}

func (s *RpcServer) getConnectionWith(partner *pki.Certificate) (*RpcConnection, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var conn *RpcConnection
	s.activeConnections.ForEach(func(uuid uuid.UUID, connection *RpcConnection) error {
		if connection.partner == nil {
			return nil
		}
		if connection.partner.Equal(partner) {
			conn = connection
			return errors.New("found connection")
		}
		return nil
	})

	if conn != nil {
		return conn, nil
	}

	return nil, fmt.Errorf("no connection found")

}

func (s *RpcServer) LoginHandler(handler func(*RpcSession) error) {
	s.loginHandler = handler
}
