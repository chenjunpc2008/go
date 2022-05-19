package cbufftcpserver

import (
    "container/list"
    "errors"
    "net"
    "sync"
)

/*
single client session
*/
type clientSessnSt struct {
    conn net.Conn
    id   uint64
    ip   string
    addr string

    lock         sync.Mutex // lock for below values
    msgsToBeSend *list.List
    closeOnce    sync.Once // close the conn, once/per instance
    closed       bool
}

func newClientSessnSt(c net.Conn, clientID uint64, cliIP string, cliAddr string) *clientSessnSt {
    return &clientSessnSt{conn: c, id: clientID,
        ip: cliIP, addr: cliAddr, msgsToBeSend: list.New(),
        closed: false}
}

// close client session
func (sn *clientSessnSt) close() {
    // lock
    sn.lock.Lock()
    defer sn.lock.Unlock()

    sn.closeOnce.Do(func() {
        sn.conn.Close()
        sn.closed = true
    })

    sn.msgsToBeSend = nil
}

//
func (sn *clientSessnSt) putSendMsg(msg interface{}) error {
    // lock
    sn.lock.Lock()
    defer sn.lock.Unlock()

    if sn.closed {
        return errors.New("client closed")
    }

    if nil == sn.msgsToBeSend {
        return errors.New("nil msgsToBeSend")
    }

    sn.msgsToBeSend.PushBack(msg)
    return nil
}

// get send msg from session local buffer
func (sn *clientSessnSt) getSendMsg() (msg interface{}, bHaveMsg bool, bConnected bool, err error) {
    // lock
    sn.lock.Lock()
    defer sn.lock.Unlock()

    if sn.closed {
        return nil, false, false, nil
    }

    if nil == sn.msgsToBeSend {
        return nil, false, false, errors.New("nil msgsToBeSend")
    }

    // get msg
    elemFront := sn.msgsToBeSend.Front()
    if nil == elemFront {
        return nil, false, true, nil
    }

    var msgToSend interface{} = elemFront.Value

    sn.msgsToBeSend.Remove(elemFront)
    return msgToSend, true, true, nil
}

/*
client sessions
*/
type clientSnsSt struct {
    initOnce sync.Once // init once

    // lock for values below
    lock sync.Mutex
    // key-clientID
    mapCliSess map[uint64]*clientSessnSt
}

func (sns *clientSnsSt) init() {
    sns.initOnce.Do(func() {
        sns.mapCliSess = make(map[uint64]*clientSessnSt, 0)
    })
}

// add a new client session
func (sns *clientSnsSt) addNewConnection(c net.Conn, clientID uint64, cliIP string, cliAddr string) *clientSessnSt {
    // lock
    sns.lock.Lock()
    defer sns.lock.Unlock()

    var cliSessn = newClientSessnSt(c, clientID, cliIP, cliAddr)

    sns.mapCliSess[clientID] = cliSessn

    // fmt.Printf("%v on new connection, client-id:=%v, ip:%v, addr:%v\n", ftag, clientID, cliIP, cliAddr)
    return cliSessn
}

// delete connetion
func (sns *clientSnsSt) delConnect(clientID uint64, cliIP string, cliAddr string) {
    // lock
    sns.lock.Lock()
    defer sns.lock.Unlock()

    delete(sns.mapCliSess, clientID)

    // fmt.Printf("%v delConnect, client-id:=%v, ip:%v, addr:%v\n", ftag, clientID, cliIP, cliAddr)
}

// get all client ids
func (sns *clientSnsSt) getAllClientIDs() []uint64 {
    // lock
    sns.lock.Lock()
    defer sns.lock.Unlock()

    var clientids = make([]uint64, 0)
    for k := range sns.mapCliSess {
        clientids = append(clientids, k)
    }

    return clientids
}

// get client session object
func (sns *clientSnsSt) getClientSession(clientID uint64) (*clientSessnSt, bool) {
    // lock
    sns.lock.Lock()
    defer sns.lock.Unlock()

    cli, ok := sns.mapCliSess[clientID]
    if !ok {
        return nil, false
    }

    return cli, true
}
