package cbufftcpserver

/*
EventHandler server callback control handler
*/
type EventHandler interface {
    // new connections event
    OnNewConnection(clientID uint64, clientIP string, clientAddr string)
    // disconnected event
    OnDisconnected(clientID uint64, clientIP string, clientAddr string)
    // receive data event
    OnReceiveData(clientID uint64, clientIP string, clientAddr string, pPacks []interface{})
    // data already sended event
    OnSendedData(clientID uint64, clientIP string, clientAddr string, msg interface{}, bysSended []byte, length int)
    // error
    OnError(msg string, err error)
    // error
    OnCliError(clientID uint64, clientIP string, clientAddr string, msg string, err error)
    // error
    OnCliErrorStr(clientID uint64, clientIP string, clientAddr string, msg string)

    // data protocol
    ProtocolIF
}
