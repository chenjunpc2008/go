/*
Package izips i-zips
*/
package izips

import (
    "bytes"
    "compress/gzip"
    "encoding/base64"
    "errors"
    "io"
    "io/ioutil"
)

/*
GzipEncode gzip encode

@return
*/
func GzipEncode(src []byte) (dest []byte, err error) {
    var (
        btsBuffer bytes.Buffer
    )

    writer := gzip.NewWriter(&btsBuffer)
    _, err = writer.Write(src)
    if nil != err {
        writer.Close()
        return
    }

    err = writer.Close()
    if nil != err {
        return
    }

    dest = btsBuffer.Bytes()
    err = nil
    return
}

/*
GzipDecode gzip decode

@return
*/
func GzipDecode(src []byte) (dest []byte, err error) {
    bysReader, err := gzip.NewReader(bytes.NewReader(src))
    if nil != err {
        return
    }

    dest, err = ioutil.ReadAll(bysReader)

    bysReader.Close()
    return
}

/*
EncodeBytesToGzipBase64Str gzip encode and convert to base64 string

@return
*/
func EncodeBytesToGzipBase64Str(bysSrc []byte) (string, error) {
    var (
        btsBuffer bytes.Buffer
        bytsZiped []byte
        strBase64 string
        err       error
    )

    if nil == bysSrc {
        return "", errors.New("nil []byte")
    }

    // gzip encode
    writer := gzip.NewWriter(&btsBuffer)
    _, err = writer.Write(bysSrc)
    if nil != err {
        writer.Close()
        return "", err
    }

    err = writer.Close()
    if nil != err {
        return "", err
    }

    bytsZiped = btsBuffer.Bytes()

    // base64 encode
    strBase64 = base64.StdEncoding.EncodeToString(bytsZiped)

    return strBase64, nil
}

/*
EncodeStrToGzipBase64Str gzip encode and convert to base64 string

@return
*/
func EncodeStrToGzipBase64Str(src string) (string, error) {

    var bysSrc = []byte(src)
    return EncodeBytesToGzipBase64Str(bysSrc)
}

/*
DecodeGzipBase64Str decode base64 gzip string

@return
*/
func DecodeGzipBase64Str(src string) (string, error) {
    var (
        bysSrc           []byte
        err              error
        bysBase64Decoded []byte
        bysGzDecode      []byte
        decodeRes        string
    )

    // base64 decode
    bysSrc = []byte(src)
    base64BuffLen := base64.StdEncoding.DecodedLen(len(bysSrc))
    bysBase64Buff := make([]byte, base64BuffLen)
    nUsed, err := base64.StdEncoding.Decode(bysBase64Buff, bysSrc)
    if nil != err {
        return "", err
    }

    bysBase64Decoded = bysBase64Buff[:nUsed]

    // gzip decode
    bysReader, err := gzip.NewReader(bytes.NewReader(bysBase64Decoded))
    if nil != err {
        return "", err
    }

    bysGzDecode, err = io.ReadAll(bysReader)
    if nil != err {
        return "", err
    }

    bysReader.Close()

    decodeRes = string(bysGzDecode)

    return decodeRes, nil
}
