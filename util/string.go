package util

import "github.com/axgle/mahonia"

func ConvertToString(src string, srcCode string, tagCode string) string {
	if srcCode == tagCode {
		return src
	}
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
