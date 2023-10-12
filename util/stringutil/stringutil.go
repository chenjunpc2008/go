/*
Package stringutil string utils
*/
package stringutil

/*
CutByMaxLen 按最大长度截取字符串，如超过最大长度返回原字符串，如未超过则截断尾部

@return
*/
func CutByMaxLen(src string, iMaxLen int) string {
    if 0 > iMaxLen {
        return src
    }

    if 0 == iMaxLen {
        return ""
    }

    if len(src) <= iMaxLen {
        return src
    }

    return src[0:iMaxLen]
}

/*
CutLastByMaxLen 按最大长度从后向前截取字符串，如超过最大长度返回原字符串，如未超过则截断首部

@return
*/
func CutLastByMaxLen(src string, iMaxLen int) string {
    if 0 > iMaxLen {
        return src
    }

    if 0 == iMaxLen {
        return ""
    }

    if len(src) <= iMaxLen {
        return src
    }

    return src[len(src)-iMaxLen:]
}
