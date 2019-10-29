package shortcuts

import (
	"strings"

	"pkg.deepin.io/lib/gettext"
)

func pinYinDistinguish(search string, searchIndex int, pinYin []string, wordIndex int, wordStart int) bool {
	if searchIndex == 0 {
		return search[0] == pinYin[0][0] && // 第一个必须匹配
			(len(search) == 1 || pinYinDistinguish(search, 1, pinYin, 0, 1)) //如果仅是1个字符，算匹配，否则从第二个字符开始比较
	}
	if len(pinYin[wordIndex]) > wordStart && //判断不越界
		search[searchIndex] == pinYin[wordIndex][wordStart] { //判断匹配
		if searchIndex == len(search)-1 {
			//如果这是最后一个字符，检查之前的声母是否依次出现
			return pinYinDistinguishAux(search, pinYin, wordIndex)
		} else {
			return pinYinDistinguish(search, searchIndex+1, pinYin, wordIndex, wordStart+1) //同一个字拼音连续匹配
		}
	} else if len(pinYin) > wordIndex+1 && //判断不越界
		search[searchIndex] == pinYin[wordIndex+1][0] { //不能拼音连续匹配的情况下，看看下一个字拼音的首字母是否能匹配
		if searchIndex == len(search)-1 {
			return pinYinDistinguishAux(search, pinYin, wordIndex) //如果这是最后一个字符，检查之前的声母是否依次出现
		} else {
			return pinYinDistinguish(search, searchIndex+1, pinYin, wordIndex+1, 1) //去判断下一个字拼音的第二个字母
		}
	} else if len(pinYin) > wordIndex+1 {
		// //回退试试看  对于zhuang an lan  searchIndex > 1
		for i := 1; i < searchIndex; i++ {
			if pinYinDistinguish(search, searchIndex-i, pinYin, wordIndex+1, 0) {
				return true
			}
		}
	}

	return false
}

// 检查之前的声母是否依次出现
// 辅佐函数，确保pinYin[n][0] (n<=wordIndex)都按顺序依次出现在search里面
//     * 防止zhou ming zhong匹配zz，跳过了ming
func pinYinDistinguishAux(search string, pinYin []string, wordIndex int) bool {
	lastIndex := 0
	for i := 0; i < wordIndex; i++ {
		lastIndex = indexOf(search, pinYin[i][0], lastIndex)
		if lastIndex == -1 {
			return false
		}
		lastIndex++
	}
	return true
}

func indexOf(str string, b byte, fromIndex int) int {
	result := strings.IndexByte(str[fromIndex:], b)
	if result == -1 {
		return -1
	}
	return result + fromIndex
}

func isZH() bool {
	lang := gettext.QueryLang()
	return strings.HasPrefix(lang, "zh")
}
