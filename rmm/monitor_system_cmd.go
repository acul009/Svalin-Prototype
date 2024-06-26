package rmm

import (
	"fmt"
	"log"
	"time"

	"github.com/rahn-it/svalin/rpc"
	"github.com/rahn-it/svalin/util"
)

const reportingInterval = 1 * time.Second

func MonitorSystemCommandHandler() rpc.RpcCommand {
	return &monitorSystemCommand{}
}

type monitorSystemCommand struct {
	static *StaticStats
	active util.UpdateableObservable[*ActiveStats]
}

func NewMonitorSystemCommand(static *StaticStats, activeOb util.UpdateableObservable[*ActiveStats]) *monitorSystemCommand {
	return &monitorSystemCommand{
		static: static,
		active: activeOb,
	}
}

func (cmd *monitorSystemCommand) GetKey() string {
	return "monitor-system"
}

func (cmd *monitorSystemCommand) ExecuteServer(session *rpc.RpcSession) error {

	static, err := GetStaticStats()
	if err != nil {
		session.WriteResponseHeader(rpc.SessionResponseHeader{
			Code: 500,
			Msg:  "Unable to read host info",
		})
		return fmt.Errorf("error reading host info: %w", err)
	}

	err = session.WriteResponseHeader(rpc.SessionResponseHeader{
		Code: 200,
		Msg:  "OK",
	})
	if err != nil {
		return fmt.Errorf("error writing response header: %w", err)
	}

	err = rpc.WriteMessage[*StaticStats](session, static)
	if err != nil {
		return fmt.Errorf("error writing static stats: %w", err)
	}

	for {
		active, err := GetActiveStats()
		if err != nil {
			return fmt.Errorf("error getting active stats: %w", err)
		}

		err = rpc.WriteMessage[*ActiveStats](session, active)
		if err != nil {
			return fmt.Errorf("error writing active stats: %w", err)
		}

		time.Sleep(reportingInterval)
	}
}

func (cmd *monitorSystemCommand) ExecuteClient(session *rpc.RpcSession) error {
	log.Printf("Monitoring remote system...")

	err := rpc.ReadMessage[*StaticStats](session, cmd.static)
	if err != nil {
		return fmt.Errorf("error reading static stats: %w", err)
	}

	active := &ActiveStats{}

	for {
		err = rpc.ReadMessage[*ActiveStats](session, active)
		if err != nil {
			return fmt.Errorf("error reading active stats: %w", err)
		}

		// log.Printf("Received active stats: %+v", active)

		cmd.active.Update(func(_ *ActiveStats) *ActiveStats {
			return active
		})
	}
}
