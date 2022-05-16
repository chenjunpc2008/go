package ctcpclient

import (
    "container/list"
    "errors"
    "net"
    "sync"
)

const (
    // MaxRcvBufferCapSize cap size
    MaxRcvBufferCapSize = 15728640 // 1024*1024*15
)

// Config extra config
type Config struct {
    // after recieve a whole package, the receive callback will go sync or async
    AsyncReceive bool
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

    lock         sync.Mutex // lock for below values
    msgsToBeSend *list.List
}

// New new tcp client
func New(eventCb EventHandler, cnf Config) *CtcpCli {

    var cli = &CtcpCli{bIsConnected: false}

    cli.chExit = make(chan int)
    cli.handler = eventCb
    cli.cnf = cnf

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
    cli.msgsToBeSend = nil
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
    cli.msgsToBeSend = list.New()
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
    cli.msgsToBeSend = nil
}

// SendToServer send message to server
func (cli *CtcpCli) SendToServer(msg interface{}) error {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    if !cli.bIsConnected {
        return errors.New("not connected")
    }

    cli.msgsToBeSend.PushBack(msg)
    return nil
}

// SendPriorToServer send prior message to server
func (cli *CtcpCli) SendPriorToServer(msg interface{}) error {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    if !cli.bIsConnected {
        return errors.New("not connected")
    }

    cli.msgsToBeSend.PushFront(msg)
    return nil
}

// SendPrioresToServer send prior message to server
func (cli *CtcpCli) SendPrioresToServer(msgs []interface{}) error {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    if !cli.bIsConnected {
        return errors.New("not connected")
    }

    for _, v := range msgs {
        cli.msgsToBeSend.PushFront(v)
    }

    return nil
}

// SendBuffToServer send buff message to server
func (cli *CtcpCli) SendBuffToServer(buff *list.List) error {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    if !cli.bIsConnected {
        return errors.New("not connected")
    }

    cli.msgsToBeSend.PushBackList(buff)
    return nil
}

// get message to send
func (cli *CtcpCli) getSendMsg() (msg interface{}, bHaveMsg bool, bConnected bool, err error) {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    if !cli.bIsConnected {
        return nil, false, false, nil
    }

    if nil == cli.msgsToBeSend {
        return nil, false, false, errors.New("nil msgsToBeSend")
    }

    // get msg
    elemFront := cli.msgsToBeSend.Front()
    if nil == elemFront {
        return nil, false, true, nil
    }

    var msgToSend interface{} = elemFront.Value
    cli.msgsToBeSend.Remove(elemFront)
    return msgToSend, true, true, nil
}

// DumpSendBuffer dump send buffer out, session's send buffer will replace with empty one
func (cli *CtcpCli) DumpSendBuffer() *list.List {
    // lock
    cli.lock.Lock()
    defer cli.lock.Unlock()

    if !cli.bIsConnected {
        return nil
    }

    if nil == cli.msgsToBeSend {
        return nil
    }

    var buffOut = cli.msgsToBeSend
    cli.msgsToBeSend = list.New()

    return buffOut
}
