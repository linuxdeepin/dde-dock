package iso639

// 文件 auto.go 是自动生成的，生成工具在 https://gitee.com/electricface/codes/517eok9msxijwpnz6hld462

// convert ISO 639-1 to ISO 639-2 T/B
func ConvertA2ToA3(in string) []string {
	for _, lang := range allLanguages {
		if lang.A2 == in {
			if lang.A3T == lang.A3B {
				return []string{lang.A3T}
			}
			return []string{lang.A3T, lang.A3B}
		}
	}
	return nil
}
