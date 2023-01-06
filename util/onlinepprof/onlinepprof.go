/*
Package onlinepprof net pprof
*/
package onlinepprof

import (
    "fmt"
    "log"
    "net/http"

    // pprof
    _ "net/http/pprof"
)

/*
StartOnlinePprof net pprof
*/
func StartOnlinePprof(bEnable bool, httpPort uint16, bPanicIfFailed bool) (*http.Server, error) {
    if !bEnable {
        return nil, nil
    }

    sAddr := fmt.Sprintf(":%d", httpPort)
    var svr = &http.Server{
        Addr: sAddr,
    }

    /*
       http://127.0.0.1:10108/debug/pprof/
    */

    // For load balance keep alive and pprof debug
    // 提供给负载均衡探活以及pprof调试
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    go procPprofListen(sAddr, bPanicIfFailed, svr)

    return svr, nil
}

func procPprofListen(addr string, bPanicIfFailed bool, svr *http.Server) {
    const ftag = "online_pprof.procPprofListen()"

    var err = svr.ListenAndServe()
    if nil != err {
        sErrMsg := fmt.Sprintf("%v http.ListenAndServe failed, addr:%v, err:%v", ftag, addr, err)
        fmt.Println(sErrMsg)
        log.Println(sErrMsg)

        if http.ErrServerClosed != err {
            if bPanicIfFailed {
                panic(sErrMsg)
            }
        } else {
            log.Println("http server Shutdown or Close")
        }
    }
}
