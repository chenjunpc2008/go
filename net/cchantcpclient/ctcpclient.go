package cchantcpclient

import (
    "errors"
    "net"
    "sync"
)

const (
    // DefaultSendBuffSize default send buff size
    DefaultSendBuffSize = 50 * 10000

    // MaxRcvBufferCapSize cap size
    MaxRcvBufferCapSize = 15728640 // 1024*1024*15
)

// Config extra config
type Config struct {
    // tcp send buff size
    SendBuffsize int

    // after recieve a whole package, the receive callback will go sync or async
    AsyncReceive bool
}

// DefaultConfig default Config
func DefaultConfig() Config {
    return Config{
        SendBuffsize: DefaultSendBuffSize,
        AsyncReceive: true,
    }
}

// CtcpCli tcp client
type CtcpCli struct {
    svrIP   string
    svrPort uint16
    cnf     Config

    conn         net.Conn
    bIsConnected bool
    handler      EventHandler
    chExit       chan int // notify all goroutines to shutdown

    lock           sync.Mutex // lock for below values
    chMsgsToBeSend chan interface{}
}

// New new tcp client
func New(eventCb EventHandler, cnf Config) *CtcpCli {

    var cli = &CtcpCli{bIsConnected: false}

    cli.chExit = make(chan int)
    cli.handler = eventCb
    cli.cnf = cnf
    cli.chMsgsToBeSend = make(chan interface{}, cli.cnf.SendBuffsize)

    return cli
}

// ConnectToServer connect to remote server
func (cli *CtcpCli) ConnectToServer(ip string, port uint16) error {
    return connectToServer(ip, port, cli)
}

// Close close connection
func (cli *CtcpCli) Close() {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    if nil != cli.conn {
        close(cli.chExit)
        cli.conn.Close()
        cli.conn = nil
    }

    cli.svrIP = ""
    cli.svrPort = 0
    cli.bIsConnected = false

    if nil != cli.chMsgsToBeSend {
        close(cli.chMsgsToBeSend)
        cli.chMsgsToBeSend = nil
    }
}

// new conn
func (cli *CtcpCli) addNewConnection(conn net.Conn, ip string, port uint16) {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    cli.conn = conn
    cli.svrIP = ip
    cli.svrPort = port
    cli.bIsConnected = true
    cli.chMsgsToBeSend = make(chan interface{}, cli.cnf.SendBuffsize)
}

// disconnected info update
func (cli *CtcpCli) disconnected() {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    cli.conn = nil
    cli.svrIP = ""
    cli.svrPort = 0
    cli.bIsConnected = false

    if nil != cli.chMsgsToBeSend {
        close(cli.chMsgsToBeSend)
        cli.chMsgsToBeSend = nil
    }
}

/*
SendToServer send message to server

@return busy bool : true -- buff is full, you may need to try again
*/
func (cli *CtcpCli) SendToServer(msg interface{}) (busy bool, retErr error) {
    // lock
    cli.lock.Lock()

    if !cli.bIsConnected {
        // unlock
        cli.lock.Unlock()
        return false, errors.New("not connected")
    }

    if nil == cli.chMsgsToBeSend {
        // unlock
        cli.lock.Unlock()
        return false, errors.New("nil chMsgsToBeSend")
    }

    curBuffSize := len(cli.chMsgsToBeSend)

    if curBuffSize >= cli.cnf.SendBuffsize-1 {
        // unlock
        cli.lock.Unlock()
        return true, nil
    }

    // push
    cli.chMsgsToBeSend <- msg

    // unlock
    cli.lock.Unlock()

    return false, nil
}
