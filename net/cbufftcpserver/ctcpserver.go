/*
Package cbufftcpserver tcp server
*/
package cbufftcpserver

import (
    "errors"
    "fmt"
    "net"
    "sync"
)

func init() {

}

// ServerStatus server status
type ServerStatus int

const (
    // ServerStatusClosed closed
    ServerStatusClosed ServerStatus = 0
    // ServerStatusStarting starting
    ServerStatusStarting ServerStatus = 1
    // ServerStatusRunning running
    ServerStatusRunning ServerStatus = 2
    // ServerStatusStopping stopping
    ServerStatusStopping ServerStatus = 3
)

// Config extra config
type Config struct {
    // AsyncReceive after recieve a whole package, the receive callback will go sync or async
    AsyncReceive bool
}

/*
CtcpsvrSt tcp server instance
*/
type CtcpsvrSt struct {
    port uint16
    cnf  Config

    listener  *net.TCPListener
    clientID  uint64
    cliSns    clientSnsSt
    handler   EventHandler
    chExit    chan int        // notify all goroutines to shutdown
    wg        *sync.WaitGroup // wait for all goroutines to exit
    svrStatus ServerStatus    // server status
}

/*
NewTCPSvr new ctcpserver object
*/
func NewTCPSvr(eventCb EventHandler, cnf Config) *CtcpsvrSt {
    var svr = &CtcpsvrSt{clientID: 0}
    svr.cnf = cnf

    svr.cliSns.init()
    svr.chExit = make(chan int)
    svr.wg = &sync.WaitGroup{}

    svr.handler = eventCb

    svr.svrStatus = ServerStatusClosed

    return svr
}

// StartServer start server
func (t *CtcpsvrSt) StartServer(iPort uint16) error {

    if ServerStatusClosed != t.svrStatus {
        errMsg := fmt.Sprintf("couldn't start server in status:%d", t.svrStatus)
        return errors.New(errMsg)
    }

    t.svrStatus = ServerStatusStarting

    listener, err := createListen(iPort, t)
    if nil != err {
        return err
    }

    go startServer(iPort, listener, t)

    return nil
}

// StopServer top server
func (t *CtcpsvrSt) StopServer() error {

    if ServerStatusRunning != t.svrStatus {
        errMsg := fmt.Sprintf("couldn't stop server in status:%d", t.svrStatus)
        return errors.New(errMsg)
    }

    stopServer(t)

    return nil
}

// SendToClient sent to client
func (t *CtcpsvrSt) SendToClient(clientID uint64, msg interface{}) error {
    if ServerStatusRunning != t.svrStatus {
        errMsg := fmt.Sprintf("couldn't SendToClient, server in status:%d", t.svrStatus)
        return errors.New(errMsg)
    }

    cli, ok := t.cliSns.getClientSession(clientID)
    if !ok {
        return errors.New("couldn't getClientSession")
    }

    // nil
    if nil == cli {
        return errors.New("nil ClientSession")
    }

    err := cli.putSendMsg(msg)
    if nil != err {
        return err
    }

    return nil
}

// SendToAllClients send to all clients
func (t *CtcpsvrSt) SendToAllClients(msg interface{}) error {
    if ServerStatusRunning != t.svrStatus {
        errMsg := fmt.Sprintf("couldn't SendToClient, server in status:%d", t.svrStatus)
        return errors.New(errMsg)
    }

    cliIDs := t.cliSns.getAllClientIDs()
    // nil
    if nil == cliIDs {
        return errors.New("nil []cli-ids")
    }

    var (
        err error
        cli *clientSessnSt
        ok  bool
    )

    for _, v := range cliIDs {
        cli, ok = t.cliSns.getClientSession(v)
        if !ok {
            t.handler.OnCliErrorStr(v, "", "", "couldn't getClientSession")
            continue
        }

        // nil
        if nil == cli {
            t.handler.OnCliErrorStr(v, "", "", "nil ClientSession")
            continue
        }

        err = cli.putSendMsg(msg)
        if nil != err {
            t.handler.OnCliError(v, "", "", "putSendMsg", err)
            continue
        }
    }

    return nil
}

// CloseClient close one client
func (t *CtcpsvrSt) CloseClient(clientID uint64, reason string) {
    closeCli(clientID, reason, t)
}

// CloseClients close clients
func (t *CtcpsvrSt) CloseClients(clientIDs []uint64, reason string) {
    closeClients(clientIDs, reason, t)
}
