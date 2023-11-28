package izips

import (
    "bytes"
    "compress/gzip"
    "io/ioutil"
    "testing"

    "encoding/base64"

    "github.com/stretchr/testify/assert"
)

func TestGzip_orign(t *testing.T) {
    toBeZip := "hello world"
    btTobezip := []byte(toBeZip)

    var res bytes.Buffer
    zwriter, _ := gzip.NewWriterLevel(&res, 7)
    if nil == zwriter {
        assert.Equal(t, 1, 0)
    }

    _, err := zwriter.Write(btTobezip)
    if nil != err {
        assert.Equal(t, 1, 0)
    }
    zwriter.Close()

    byteRes := res.Bytes()

    t.Logf("%v\n", byteRes)

    strZipped := base64.StdEncoding.EncodeToString(byteRes)
    sExpectedZipp := "H4sIAAAAAAAA/8pIzcnJVyjPL8pJAQQAAP//hRFKDQsAAAA="
    assert.Equal(t, sExpectedZipp, strZipped)

    var rder *bytes.Buffer
    rder = bytes.NewBuffer(byteRes)

    zreader, err := gzip.NewReader(rder)
    if nil != err {
        assert.Equal(t, 1, 0)
    }

    if nil == zreader {
        assert.Equal(t, 1, 0)
    }

    byteUnpressRes, err := ioutil.ReadAll(zreader)
    if nil != err {
        assert.Equal(t, 1, 0)
    }

    sUnpressRes := string(byteUnpressRes)
    assert.Equal(t, toBeZip, sUnpressRes)

    strBase64 := base64.StdEncoding.EncodeToString(btTobezip)
    assert.Equal(t, "aGVsbG8gd29ybGQ=", strBase64)
}

func TestGzip_CrossLang1(t *testing.T) {
    toBeZip := "hello world"
    btTobezip := []byte(toBeZip)

    bytsZiped, err := GzipEncode(btTobezip)
    if nil != err {
        t.FailNow()
    }

    t.Logf("%v\n", bytsZiped)

    bytsUnziped, err := GzipDecode(bytsZiped)

    sUnpressRes := string(bytsUnziped)
    assert.Equal(t, toBeZip, sUnpressRes)

    bytsJavaZiped := []byte{byte(31), byte(139), byte(8), byte(0), byte(0),
        byte(0), byte(0), byte(0), byte(0), byte(0),
        byte(203), byte(72), byte(205), byte(201), byte(201),
        byte(87), byte(40), byte(207), byte(47), byte(202),
        byte(73), byte(1), byte(0), byte(133), byte(17),
        byte(74), byte(13), byte(11), byte(0), byte(0),
        byte(0)}

    bytsUnzipedJava, err := GzipDecode(bytsJavaZiped)

    sUnpressJavaRes := string(bytsUnzipedJava)
    assert.Equal(t, toBeZip, sUnpressJavaRes)
}

func TestGzip_base64(t *testing.T) {
    toBeZip := "[{'msgseq':'num_c','key':'matchRule','value':{'rules':" +
        "[{'ruleid':'6109','ruleFlag':'102_1','judge_cond':[{'field':'1180','value':'020','comp_type':'string'}," +
        "{'field':'35','value':'D','comp_type':'string'},{'field':'38','value':'','comp_type':'int_range','low':1600,'high':2000}]," +
        "'timestamp':1561528389,'conf':{'quotVld':0,'mktAmtVld':0}," +
        "'check':'','confirm':{'11':'@11','14':'0.00','17':'$EXEC_ID','22':'@22','35':'8','37':'$EXEC_ID'," +
        "'38':'@38','39':'0','40':'@40','44':'@44','48':'@48','54':'@54','59':'@59','60':'$LOCAL_TIME_STAMP'," +
        "'77':'@77','110':'0.00','150':'0','151':'@38','203':'@203','453':[{'447':'5','448':'$ACCOUNT','452':'5'}," +
        "{'447':'C','448':'$SEAT','452':'1'},{'447':'C','448':'$SEAT','452':'27'},{'447':'F','448':'01','452':'4'}]," +
        "'522':'@522','544':'@544','1090':'0','1180':'@1180'}," +
        "'match':[{'repeat':1,'content':{'11':'@11','14':'@38','17':'$EXEC_ID','22':'@22','31':'@44','32':'@38'," +
        "'35':'8','37':'$EXEC_ID','38':'@38','39':'2','40':'@40','41':'','44':'','48':'@48','54':'@54','58':''," +
        "'59':'','60':'$LOCAL_TIME_STAMP','77':'','99':'','110':'','150':'F','151':'0.00','152':'','203':''," +
        "'453':[{'447':'5','448':'$ACCOUNT','452':'5'},{'447':'C','448':'$SEAT','452':'1'},{'447':'C','448':'$SEAT','452':'27'}," +
        "{'447':'F','448':'01','452':'4'}],'522':'@522','541':'','544':'','664':'','1090':'','1093':'','1180':'@1180','8906':''" +
        ",'8911':'','10189':''}}],'have_settlement':0}]}}]"
    btTobezip := []byte(toBeZip)
    t.Logf("length:%v\n", len(toBeZip))

    bytsZiped, err := GzipEncode(btTobezip)
    if nil != err {
        t.FailNow()
    }

    t.Logf("len:%v, %v\n", len(bytsZiped), bytsZiped)

    strBase64 := base64.StdEncoding.EncodeToString(bytsZiped)
    t.Logf("length:%v, %v\n", len(strBase64), strBase64)

    bytsUnziped, err := GzipDecode(bytsZiped)

    sUnpressRes := string(bytsUnziped)
    assert.Equal(t, toBeZip, sUnpressRes)
}

func TestGzip_Packed1(t *testing.T) {
    toBeZip := "[{'msgseq':'num_c','key':'matchRule','value':{'rules':" +
        "[{'ruleid':'6109','ruleFlag':'102_1','judge_cond':[{'field':'1180','value':'020','comp_type':'string'}," +
        "{'field':'35','value':'D','comp_type':'string'},{'field':'38','value':'','comp_type':'int_range','low':1600,'high':2000}]," +
        "'timestamp':1561528389,'conf':{'quotVld':0,'mktAmtVld':0}," +
        "'check':'','confirm':{'11':'@11','14':'0.00','17':'$EXEC_ID','22':'@22','35':'8','37':'$EXEC_ID'," +
        "'38':'@38','39':'0','40':'@40','44':'@44','48':'@48','54':'@54','59':'@59','60':'$LOCAL_TIME_STAMP'," +
        "'77':'@77','110':'0.00','150':'0','151':'@38','203':'@203','453':[{'447':'5','448':'$ACCOUNT','452':'5'}," +
        "{'447':'C','448':'$SEAT','452':'1'},{'447':'C','448':'$SEAT','452':'27'},{'447':'F','448':'01','452':'4'}]," +
        "'522':'@522','544':'@544','1090':'0','1180':'@1180'}," +
        "'match':[{'repeat':1,'content':{'11':'@11','14':'@38','17':'$EXEC_ID','22':'@22','31':'@44','32':'@38'," +
        "'35':'8','37':'$EXEC_ID','38':'@38','39':'2','40':'@40','41':'','44':'','48':'@48','54':'@54','58':''," +
        "'59':'','60':'$LOCAL_TIME_STAMP','77':'','99':'','110':'','150':'F','151':'0.00','152':'','203':''," +
        "'453':[{'447':'5','448':'$ACCOUNT','452':'5'},{'447':'C','448':'$SEAT','452':'1'},{'447':'C','448':'$SEAT','452':'27'}," +
        "{'447':'F','448':'01','452':'4'}],'522':'@522','541':'','544':'','664':'','1090':'','1093':'','1180':'@1180','8906':''" +
        ",'8911':'','10189':''}}],'have_settlement':0}]}}]"

    strEncoded, err := EncodeStrToGzipBase64Str(toBeZip)
    if nil != err {
        t.FailNow()
    }

    t.Logf("strEncoded:%v\n", strEncoded)

    strDecoded, err := DecodeGzipBase64Str(strEncoded)
    if nil != err {
        t.FailNow()
    }
    assert.Equal(t, toBeZip, strDecoded)

    // before zip len:1389, after zip len:736, compress percent 52.99%
    t.Logf("before zip len:%v, after zip len:%v, compress percent:%v\n", len(toBeZip), len(strEncoded), float64(len(strEncoded))/float64(len(toBeZip)))
}

func TestGzip_CrossLanguage(t *testing.T) {
    toBeZip := "[{'msgseq':'num_c','key':'matchRule','value':{'rules':" +
        "[{'ruleid':'6109','ruleFlag':'102_1','judge_cond':[{'field':'1180','value':'020','comp_type':'string'}," +
        "{'field':'35','value':'D','comp_type':'string'},{'field':'38','value':'','comp_type':'int_range','low':1600,'high':2000}]," +
        "'timestamp':1561528389,'conf':{'quotVld':0,'mktAmtVld':0}," +
        "'check':'','confirm':{'11':'@11','14':'0.00','17':'$EXEC_ID','22':'@22','35':'8','37':'$EXEC_ID'," +
        "'38':'@38','39':'0','40':'@40','44':'@44','48':'@48','54':'@54','59':'@59','60':'$LOCAL_TIME_STAMP'," +
        "'77':'@77','110':'0.00','150':'0','151':'@38','203':'@203','453':[{'447':'5','448':'$ACCOUNT','452':'5'}," +
        "{'447':'C','448':'$SEAT','452':'1'},{'447':'C','448':'$SEAT','452':'27'},{'447':'F','448':'01','452':'4'}]," +
        "'522':'@522','544':'@544','1090':'0','1180':'@1180'}," +
        "'match':[{'repeat':1,'content':{'11':'@11','14':'@38','17':'$EXEC_ID','22':'@22','31':'@44','32':'@38'," +
        "'35':'8','37':'$EXEC_ID','38':'@38','39':'2','40':'@40','41':'','44':'','48':'@48','54':'@54','58':''," +
        "'59':'','60':'$LOCAL_TIME_STAMP','77':'','99':'','110':'','150':'F','151':'0.00','152':'','203':''," +
        "'453':[{'447':'5','448':'$ACCOUNT','452':'5'},{'447':'C','448':'$SEAT','452':'1'},{'447':'C','448':'$SEAT','452':'27'}," +
        "{'447':'F','448':'01','452':'4'}],'522':'@522','541':'','544':'','664':'','1090':'','1093':'','1180':'@1180','8906':''" +
        ",'8911':'','10189':''}}],'have_settlement':0}]}}]"

    // java
    strJavaZiped := "H4sIAAAAAAAAANVU32+bMBD+ZyL5BU0+gwnwVMQSqVK7TmtWTaomhIgDLEDS4HSqqvzvu7MJoGg/" +
        "8roXfL77znf3fZc8v7OmKzr1wiLWHps0Zw7bqje8NZnOyy/" +
        "HWqHnNauPikXv7ID3jkXP1qrWiPOBhwih+7LOCvQAFymg68dxXag037Vrk7GpVE0JAAEf32Rc0C3fNftUv+3J0+lD1Rbs5Iw5rpxkfLwCH0zwF/" +
        "Cq1ekhawsarN79ZBH4nDusrIqSRYJzfvruMF01qtNZs8ew9EGKwA1CeqfdEA8vx51+okqY2Gx13PS3E0JKlW/" +
        "PVdtNdWgoAQBdN0C0gEdTf+A0NszRni2+LZL0lsYSgmD4dWjkiNEY7gUGZ0OMmdAN6Sk0PE4+z5ieMT0yDdIjpDReSV4ZGpNE8yltdveQxHfp6vZ+kT6u4vvPGJhTzRv8YovAJ/" +
        "1K3lcECUMbgrumbTywqHSN3J5Hb0jTEfUxi5Pk4eunlYEIEyLFLCwZYY+LeMTAFRgxn4CWA4jDgPAYSSotudKwK72eEGIEN3gYi5bTKIUnqml+" +
        "BXbh1V5lGtfB6KpVq3+nqyXkb7LCII8rRiGvFVtciA1208w0f1Y8sFGj/" +
        "D9lxzPsgVb7QffloPuwDsLG7Qb8R+r3vMkzcb7fG/0uWMs90zBZCocFIfdtIAgBzmAIDGknKlZmryrtlNa1asyi4H8KBn4BJ10NPW0FAAA="

    t.Logf("len before ziped:%v, len after ziped:%v\n", len(toBeZip), len(strJavaZiped))

    strDecoded, err := DecodeGzipBase64Str(strJavaZiped)
    if nil != err {
        t.FailNow()
    }
    assert.Equal(t, toBeZip, strDecoded)

    // go windows
    strGoWindowsZiped := "H4sIAAAAAAAA/9ST32vbMBDH/5nCvYhxJ1uOraeGLIVCu461G4MSjHEUx4vtpLHcUYr/93GS44SwH3nd21n6nu7u+zk/v0PdFq" +
        "15AQ1NV6c5CNiYN9BQZzZff+kqAwJes6ozoN9h31WmBf3so3IJGiLCBIT7vqmyAjQQypRAwI9uWZg03zZLl7EqTcUJRDEe3wSU/JVv611q33Z80tp92" +
        "RTQi2NOoE4yPl6gj0/0Z/Kysek+awoerNr+BE0RooB1WaxBS0TsFwJsWZvWZvUONKmIlIyDOOF3mhX78NJt7TeuhALqjZ3Ww1cvIF+bfHOo2qzKfc0JR" +
        "KDhmtgWCnnqD8hj0wQ0XM2/z2fpLY8lJcukBMEja+AxgjNNELPGTRgk/BQICJHPQheGLgw5dMqQlcqdKj5ViQsZWsRpV3cPs+ld+nR7P08fn6b3n0HAhG" +
        "teTybcIuFJvwqHiqRobENi4NrGgIuqwOEOQ35DuY64j6vpbPbw9dOTk0h3xcS8bHaUPc6nRw1doJGTE9HNKEIaFSEwUuXNVc5dFQ6GsCOEyTgWL6cjFSP" +
        "0wv8FfuHNzmQWNDmu1jT2d1y9IX/DSiOeQB5BXgpbnsEmv2lumj8Tj/2tI/9P7CAgGYSe/cj9ZuQ+roP0934D/iP6g2/qYFwUDcGwCz4KDjacLIWAOMHI" +
        "X8QJ0UFMsTOt52Lr7NWkrbG2MrVbFOwXfb/4FQAA//8nXQ09bQUAAA=="

    strDecoded, err = DecodeGzipBase64Str(strGoWindowsZiped)
    if nil != err {
        t.FailNow()
    }
    assert.Equal(t, toBeZip, strDecoded)

    // c++ windows
    strCWindowsZiped := "H4sIAAAAAAAAC9VU32+bMBD+ZyL5BU0+gwnwVMQSqVK7TmtWTaomhIgDLEDS4HSqqvzvu7MJoGg/8roXfL77znf3fZc8v7OmKzr1wi" +
        "LWHps0Zw7bqje8NZnOyy/HWqHnNauPikXv7ID3jkXP1qrWiPOBhwih+7LOCvQAFymg68dxXag037Vrk7GpVE0JAAEf32Rc0C3fNftUv+3J0+lD1Rbs5Iw5" +
        "rpxkfLwCH0zwF/Cq1ekhawsarN79ZBH4nDusrIqSRYJzfvruMF01qtNZs8ew9EGKwA1CeqfdEA8vx51+okqY2Gx13PS3E0JKlW/PVdtNdWgoAQBdN0C0gE" +
        "dTf+A0NszRni2+LZL0lsYSgmD4dWjkiNEY7gUGZ0OMmdAN6Sk0PE4+z5ieMT0yDdIjpDReSV4ZGpNE8yltdveQxHfp6vZ+kT6u4vvPGJhTzRv8YovAJ/1K3" +
        "lcECUMbgrumbTywqHSN3J5Hb0jTEfUxi5Pk4eunlYEIEyLFLCwZYY+LeMTAFRgxn4CWA4jDgPAYSSotudKwK72eEGIEN3gYi5bTKIUnqml+BXbh1V5lGtfB" +
        "6KpVq3+nqyXkb7LCII8rRiGvFVtciA1208w0f1Y8sFGj/D9lxzPsgVb7QffloPuwDsLG7Qb8R+r3vMkzcb7fG/0uWMs90zBZCocFIfdtIAgBzmAIDGknKlZ" +
        "mryrtlNa1asyi4H8KBn4BJ10NPW0FAAA="

    strDecoded, err = DecodeGzipBase64Str(strCWindowsZiped)
    if nil != err {
        t.FailNow()
    }
    assert.Equal(t, toBeZip, strDecoded)
}

func TestGzip_CrossLanguage2(t *testing.T) {

    toBeZip := "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijmortenpunnerudengelstadrocksklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv" +
        "abcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuvabcdefghijklmnopqrstuv123"

    // python
    strPythonZiped := "H4sIAFs/d2MC/+3QsQGAIAwEwJnUiSIEVCDBQJzfNSi+vfLoDJFTvu6n1CbaXxvTP4JCoVDoktrUJkt3ETaPLJnrmBRNQxk4gkKhUOgquu3HDzQpmwFpDAAA"

    t.Logf("len before ziped:%v, len after ziped:%v\n", len(toBeZip), len(strPythonZiped))

    strDecoded, err := DecodeGzipBase64Str(strPythonZiped)
    if nil != err {
        t.FailNow()
    }
    assert.Equal(t, toBeZip, strDecoded)
    assert.Equal(t, len(toBeZip), len(strDecoded))
}
