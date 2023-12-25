// *****************************************************************************
// 作者: lgdz
// 创建时间: 2023/12/25
// 描述：汉字转拼音
// *****************************************************************************

package vingo

import (
	"github.com/mozillazg/go-pinyin"
)

// 获取拼音首字母
func PinyinOfChineseInitial(s string) (res string) {
	words := pinyin.LazyConvert(s, nil)
	for _, word := range words {
		res += string(word[0])
	}
	return
}
