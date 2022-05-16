package onlinepprof

import (
	"fmt"
	"net/http"

	// pprof
	_ "net/http/pprof"
)

/*
StartOnlinePprof net pprof
*/
func StartOnlinePprof(bEnable bool, httpPort uint16, bPanicIfFailed bool) error {
	if bEnable {
		sAddr := fmt.Sprintf(":%d", httpPort)
		go procPprofListen(sAddr, bPanicIfFailed)
	}

	return nil
}

func procPprofListen(addr string, bPanicIfFailed bool) {
	const ftag = "online_pprof.procPprofListen()"

	/*
		http://127.0.0.1:10108/debug/pprof/
	*/

	// For load balance keep alive and pprof debug
	// 提供给负载均衡探活以及pprof调试
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	var err error = http.ListenAndServe(addr, nil)
	if nil != err {
		sErrMsg := fmt.Sprintf("%v http.ListenAndServe failed, err:%v", ftag, err)
		fmt.Printf(sErrMsg)
		if bPanicIfFailed {
			panic(sErrMsg)
		}
	}
}
