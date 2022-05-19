package cbufftcpserver

// ProtocolIF interface to self define pack and depack
type ProtocolIF interface {
    // pack message into the []byte to be written
    Pack(clientID uint64, msg interface{}) ([]byte, error)

    // depack the message packages from read []byte
    Depack(rawData []byte) ([]byte, []interface{})
}
