package calendar

import (
	"github.com/mozillazg/go-pinyin"
	"pkg.deepin.io/lib/strv"
	"regexp"
	"strings"
)

/* 合法单拼音表 */
var validPinyinMap map[string]bool = map[string]bool{}

/* 单个合法拼音的最大长度 */
var singlePinyinMaxLength = 0

/* 初始化合法拼音表 */
func init() {
	for i := range validPinyinList {
		var key = validPinyinList[i]
		validPinyinMap[key] = true
		if len(key) > singlePinyinMaxLength {
			singlePinyinMaxLength = len(key)
		}
	}
}

/* 判断字符串是否只可以进行拼音查询 */
func canQueryByPinyin(str string) bool {
	var expr, _ = regexp.Compile("^[a-zA-Z]+$")
	return expr.MatchString(str)
}

/* 创建拼音字符串 */
func createPinyin(zh string) string {
	var args = pinyin.NewArgs()
	args.Heteronym = true
	var pyList = pinyin.Pinyin(zh, args)
	var pyStr string
	for i := range pyList {
		var subList = strv.Strv.Uniq(pyList[i])
		pyStr += "[" + strings.Join(subList, "|") + "]"
	}

	return pyStr
}

/* 构造拼音查询表达式 */
func createPinyinQuery(pinyin string) string {
	var expr string
	for len(pinyin) > 0 {
		var i = singlePinyinMaxLength
		if i > len(pinyin) {
			i = len(pinyin)
		}
		for ; i > 1; i-- {
			var key = pinyin[:i]
			var _, exist = validPinyinMap[key]
			if exist {
				break
			}
		}
		var key = pinyin[:i]
		pinyin = pinyin[i:]
		expr += "[%" + key + "%]"
	}

	return expr
}

/* 构造拼音查询正则表达式 */
func createPinyinRegexp(pinyin string) string {
	var expr string
	for len(pinyin) > 0 {
		var i = singlePinyinMaxLength
		if i > len(pinyin) {
			i = len(pinyin)
		}
		for ; i > 1; i-- {
			var key = pinyin[:i]
			var _, exist = validPinyinMap[key]
			if exist {
				break
			}
		}
		var key = pinyin[:i]
		pinyin = pinyin[i:]
		expr += "\\[[a-z\\|]*" + key + "[a-z\\|]*\\]"
	}

	return expr
}

/* 判断汉字和拼音是否匹配 */
func pinyinMatch(zh string, py string) bool {
	var zhPinyin = createPinyin(zh)
	var expr = createPinyinRegexp(py)
	var pattern, _ = regexp.Compile(expr)
	return pattern.MatchString(zhPinyin)
}

/* 合法拼音列表 */
var validPinyinList []string = []string{
	"a",
	"ai",
	"an",
	"ang",
	"ao",
	"ba",
	"bai",
	"ban",
	"bang",
	"bao",
	"bei",
	"ben",
	"beng",
	"bi",
	"bian",
	"biao",
	"bie",
	"bin",
	"bing",
	"bo",
	"bu",
	"ca",
	"cai",
	"can",
	"cang",
	"cao",
	"ce",
	"cen",
	"ceng",
	"cha",
	"chai",
	"chan",
	"chang",
	"chao",
	"che",
	"chen",
	"cheng",
	"chi",
	"chong",
	"chou",
	"chu",
	"chua",
	"chuai",
	"chuan",
	"chuang",
	"chui",
	"chun",
	"chuo",
	"ci",
	"cong",
	"cou",
	"cu",
	"cuan",
	"cui",
	"cun",
	"cuo",
	"da",
	"dai",
	"dan",
	"dang",
	"dao",
	"de",
	"dei",
	"den",
	"deng",
	"di",
	"dia",
	"dian",
	"diao",
	"die",
	"ding",
	"diu",
	"dong",
	"dou",
	"du",
	"duan",
	"dui",
	"dun",
	"duo",
	"e",
	"en",
	"eng",
	"er",
	"fa",
	"fan",
	"fang",
	"fei",
	"fen",
	"feng",
	"fiao",
	"fo",
	"fou",
	"fu",
	"ga",
	"gai",
	"gan",
	"gang",
	"gao",
	"ge",
	"gei",
	"gen",
	"geng",
	"gong",
	"gou",
	"gu",
	"gua",
	"guai",
	"guan",
	"guang",
	"gui",
	"gun",
	"guo",
	"ha",
	"hai",
	"han",
	"hang",
	"hao",
	"he",
	"hei",
	"hen",
	"heng",
	"hong",
	"hou",
	"hu",
	"hua",
	"huai",
	"huan",
	"huang",
	"hui",
	"hun",
	"huo",
	"ji",
	"jia",
	"jian",
	"jiang",
	"jiao",
	"jie",
	"jin",
	"jing",
	"jiong",
	"jiu",
	"ju",
	"juan",
	"jue",
	"#NAME?",
	"ka",
	"kai",
	"kan",
	"kang",
	"kao",
	"ke",
	"ken",
	"keng",
	"kong",
	"kou",
	"ku",
	"kua",
	"kuai",
	"kuan",
	"kuang",
	"kui",
	"kun",
	"kuo",
	"la",
	"lai",
	"lan",
	"lang",
	"lao",
	"le",
	"lei",
	"leng",
	"li",
	"lia",
	"lian",
	"liang",
	"liao",
	"lie",
	"lin",
	"ling",
	"liu",
	"lo",
	"long",
	"lou",
	"lu",
	"luan",
	"lun",
	"luo",
	"lv",
	"lve",
	"ma",
	"mai",
	"man",
	"mang",
	"mao",
	"me",
	"mei",
	"men",
	"meng",
	"mi",
	"mian",
	"miao",
	"mie",
	"min",
	"ming",
	"miu",
	"mo",
	"mou",
	"mu",
	"na",
	"nai",
	"nan",
	"nang",
	"nao",
	"ne",
	"nei",
	"nen",
	"neng",
	"ni",
	"nian",
	"niang",
	"niao",
	"nie",
	"nin",
	"ning",
	"niu",
	"nong",
	"nou",
	"nu",
	"nuan",
	"nun",
	"nuo",
	"nv",
	"nve",
	"o",
	"ou",
	"pa",
	"pai",
	"pan",
	"pang",
	"pao",
	"pei",
	"pen",
	"peng",
	"pi",
	"pian",
	"piao",
	"pie",
	"pin",
	"ping",
	"po",
	"pou",
	"pu",
	"qi",
	"qia",
	"qian",
	"qiang",
	"qiao",
	"qie",
	"qin",
	"qing",
	"qiong",
	"qiu",
	"qu",
	"quan",
	"que",
	"qun",
	"ran",
	"rang",
	"rao",
	"re",
	"ren",
	"reng",
	"ri",
	"rong",
	"rou",
	"ru",
	"rua",
	"ruan",
	"rui",
	"run",
	"ruo",
	"sa",
	"sai",
	"san",
	"sang",
	"sao",
	"se",
	"sen",
	"seng",
	"sha",
	"shai",
	"shan",
	"shang",
	"shao",
	"she",
	"shei",
	"shen",
	"sheng",
	"shi",
	"shou",
	"shu",
	"shua",
	"shuai",
	"shuan",
	"shuang",
	"shui",
	"shun",
	"shuo",
	"si",
	"song",
	"sou",
	"su",
	"suan",
	"sui",
	"sun",
	"suo",
	"ta",
	"tai",
	"tan",
	"tang",
	"tao",
	"te",
	"tei",
	"teng",
	"ti",
	"tian",
	"tiao",
	"tie",
	"ting",
	"tong",
	"tou",
	"tu",
	"tuan",
	"tui",
	"tun",
	"tuo",
	"wa",
	"wai",
	"wan",
	"wang",
	"wei",
	"wen",
	"weng",
	"wo",
	"wu",
	"xi",
	"xia",
	"xian",
	"xiang",
	"xiao",
	"xie",
	"xin",
	"xing",
	"xiong",
	"xiu",
	"xu",
	"xuan",
	"xue",
	"xun",
	"ya",
	"yan",
	"yang",
	"yao",
	"ye",
	"yi",
	"yin",
	"ying",
	"yo",
	"yong",
	"you",
	"yu",
	"yuan",
	"yue",
	"yun",
	"za",
	"zai",
	"zan",
	"zang",
	"zao",
	"ze",
	"zei",
	"zen",
	"zeng",
	"zha",
	"zhai",
	"zhan",
	"zhang",
	"zhao",
	"zhe",
	"zhei",
	"zhen",
	"zheng",
	"zhi",
	"zhong",
	"zhou",
	"zhu",
	"zhua",
	"zhuai",
	"zhuan",
	"zhuang",
	"zhui",
	"zhun",
	"zhuo",
	"zi",
	"zong",
	"zou",
	"zu",
	"zuan",
	"zui",
	"zun",
	"zuo",
}
