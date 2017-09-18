/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

import (
	"github.com/axgle/mahonia"
	"os"
	"strings"
	"unicode/utf8"
)

var localeCharsetMap = map[string]string{
	"es_NI":             "ISO-8859-1",
	"es_BO":             "ISO-8859-1",
	"es_MX":             "ISO-8859-1",
	"he_IL":             "ISO-8859-8",
	"sv_SE":             "ISO-8859-1",
	"it_CH":             "ISO-8859-1",
	"eu_ES":             "ISO-8859-1",
	"id_ID":             "ISO-8859-1",
	"gl_ES@euro":        "ISO-8859-15",
	"sl_SI":             "ISO-8859-2",
	"ar_KW":             "ISO-8859-6",
	"et_EE.ISO-8859-15": "ISO-8859-15",
	"en_AU":             "ISO-8859-1",
	"en_IE@euro":        "ISO-8859-15",
	"ca_IT":             "ISO-8859-15",
	"cs_CZ":             "ISO-8859-2",
	"ca_ES@euro":        "ISO-8859-15",
	"de_BE@euro":        "ISO-8859-15",
	"en_IE":             "ISO-8859-1",
	"es_US":             "ISO-8859-1",
	"es_CL":             "ISO-8859-1",
	"fr_BE":             "ISO-8859-1",
	"fr_FR":             "ISO-8859-1",
	"so_DJ":             "ISO-8859-1",
	"aa_DJ":             "ISO-8859-1",
	"es_PA":             "ISO-8859-1",
	// No decoder found
	//"yi_US": "CP1255",
	"ar_LY":      "ISO-8859-6",
	"ku_TR":      "ISO-8859-9",
	"en_US":      "ISO-8859-1",
	"es_CO":      "ISO-8859-1",
	"es_PY":      "ISO-8859-1",
	"ca_ES":      "ISO-8859-1",
	"sv_FI@euro": "ISO-8859-15",
	"br_FR@euro": "ISO-8859-15",
	"pt_BR":      "ISO-8859-1",
	"ast_ES":     "ISO-8859-15",
	"fi_FI":      "ISO-8859-1",
	"fr_LU":      "ISO-8859-1",
	"an_ES":      "ISO-8859-15",
	"ar_IQ":      "ISO-8859-6",
	"cy_GB":      "ISO-8859-14",
	"de_AT@euro": "ISO-8859-15",
	"de_IT":      "ISO-8859-1",
	"gd_GB":      "ISO-8859-15",
	"so_KE":      "ISO-8859-1",
	"ga_IE":      "ISO-8859-1",
	"fr_BE@euro": "ISO-8859-15",
	"ar_MA":      "ISO-8859-6",
	"ar_YE":      "ISO-8859-6",
	"fi_FI@euro": "ISO-8859-15",
	"tl_PH":      "ISO-8859-1",
	"de_AT":      "ISO-8859-1",
	"oc_FR":      "ISO-8859-1",
	"mg_MG":      "ISO-8859-15",
	"nn_NO":      "ISO-8859-1",
	"xh_ZA":      "ISO-8859-1",
	"ar_AE":      "ISO-8859-6",
	"es_UY":      "ISO-8859-1",
	"en_CA":      "ISO-8859-1",
	"pt_PT":      "ISO-8859-1",
	"af_ZA":      "ISO-8859-1",
	"ca_AD":      "ISO-8859-15",
	"gl_ES":      "ISO-8859-1",
	"lt_LT":      "ISO-8859-13",
	"wa_BE":      "ISO-8859-1",
	"de_LU@euro": "ISO-8859-15",
	"en_ZA":      "ISO-8859-1",
	"fr_FR@euro": "ISO-8859-15",
	"nl_BE":      "ISO-8859-1",
	"en_ZW":      "ISO-8859-1",
	"es_ES@euro": "ISO-8859-15",
	"st_ZA":      "ISO-8859-1",
	"is_IS":      "ISO-8859-1",
	// No decoder found
	//"ka_GE": "GEORGIAN-PS",
	"ro_RO":      "ISO-8859-2",
	"ga_IE@euro": "ISO-8859-15",
	// No decoder found
	//"ko_KR.EUC-KR": "EUC-KR",
	"mt_MT": "ISO-8859-3",
	"kl_GL": "ISO-8859-1",
	// No decoder found
	//"bg_BG": "CP1251",
	"es_DO":         "ISO-8859-1",
	"ru_RU":         "ISO-8859-5",
	"ca_FR":         "ISO-8859-15",
	"ru_UA":         "KOI8-U",
	"el_CY":         "ISO-8859-7",
	"hr_HR":         "ISO-8859-2",
	"ar_LB":         "ISO-8859-6",
	"fr_CA":         "ISO-8859-1",
	"bs_BA":         "ISO-8859-2",
	"es_GT":         "ISO-8859-1",
	"fo_FO":         "ISO-8859-1",
	"zh_CN.GB18030": "GB18030",
	"ar_SY":         "ISO-8859-6",
	"sq_AL":         "ISO-8859-1",
	"fr_CH":         "ISO-8859-1",
	// No decoder found
	//"kk_KZ": "PT154",
	"nl_BE@euro": "ISO-8859-15",
	"zh_TW":      "BIG5",
	// No decoder found
	//"zh_HK": "BIG5-HKSCS",
	"zh_HK":     "BIG5",
	"sv_FI":     "ISO-8859-1",
	"ar_QA":     "ISO-8859-6",
	"th_TH":     "TIS-620",
	"zh_CN.GBK": "GBK",
	"da_DK":     "ISO-8859-1",
	"en_GB":     "ISO-8859-1",
	"en_HK":     "ISO-8859-1",
	"es_AR":     "ISO-8859-1",
	"ms_MY":     "ISO-8859-1",
	// No decoder found
	//"zh_TW.EUC-TW": "EUC-TW",
	"es_ES": "ISO-8859-1",
	// No decoder found
	//"tg_TJ": "KOI8-T",
	"el_GR":        "ISO-8859-7",
	"hu_HU":        "ISO-8859-2",
	"ja_JP.EUC-JP": "EUC-JP",
	"lg_UG":        "ISO-8859-10",
	"hsb_DE":       "ISO-8859-2",
	"mi_NZ":        "ISO-8859-13",
	"es_CR":        "ISO-8859-1",
	// No decoder found
	//"zh_CN": "GB2312",
	"zh_CN":      "GBK",
	"nl_NL@euro": "ISO-8859-15",
	"zu_ZA":      "ISO-8859-1",
	"uk_UA":      "KOI8-U",
	"es_PE":      "ISO-8859-1",
	// No decoder found
	//"zh_SG": "GB2312",
	"zh_SG":        "GBK",
	"kw_GB":        "ISO-8859-1",
	"ru_RU.KOI8-R": "KOI8-R",
	"so_SO":        "ISO-8859-1",
	"ar_BH":        "ISO-8859-6",
	"tr_TR":        "ISO-8859-9",
	"de_LU":        "ISO-8859-1",
	"en_SG":        "ISO-8859-1",
	"es_HN":        "ISO-8859-1",
	"lv_LV":        "ISO-8859-13",
	"fr_LU@euro":   "ISO-8859-15",
	"wa_BE@euro":   "ISO-8859-15",
	"it_IT@euro":   "ISO-8859-15",
	"pl_PL":        "ISO-8859-2",
	"es_PR":        "ISO-8859-1",
	"tr_CY":        "ISO-8859-9",
	"nb_NO":        "ISO-8859-1",
	"de_CH":        "ISO-8859-1",
	"nl_NL":        "ISO-8859-1",
	"pt_PT@euro":   "ISO-8859-15",
	// No decoder found
	//"be_BY": "CP1251",
	"ar_DZ":      "ISO-8859-6",
	"ar_TN":      "ISO-8859-6",
	"sk_SK":      "ISO-8859-2",
	"et_EE":      "ISO-8859-1",
	"en_PH":      "ISO-8859-1",
	"ar_JO":      "ISO-8859-6",
	"de_DE":      "ISO-8859-1",
	"en_DK":      "ISO-8859-1",
	"gv_GB":      "ISO-8859-1",
	"de_DE@euro": "ISO-8859-15",
	"it_IT":      "ISO-8859-1",
	"ar_SD":      "ISO-8859-6",
	"es_VE":      "ISO-8859-1",
	"ar_OM":      "ISO-8859-6",
	"om_KE":      "ISO-8859-1",
	"zh_SG.GBK":  "GBK",
	"uz_UZ":      "ISO-8859-1",
	"ar_EG":      "ISO-8859-6",
	"br_FR":      "ISO-8859-1",
	"en_NZ":      "ISO-8859-1",
	"es_SV":      "ISO-8859-1",
	// No decoder found
	//"hy_AM.ARMSCII-8": "ARMSCII-8",
	"de_BE":      "ISO-8859-1",
	"eu_ES@euro": "ISO-8859-15",
	"en_BW":      "ISO-8859-1",
	"es_EC":      "ISO-8859-1",
	"mk_MK":      "ISO-8859-5",
	"ar_SA":      "ISO-8859-6",
}

func decodeSsid(ssid []byte) string {
	if utf8.Valid(ssid) {
		return string(ssid)
	}

	locale := os.Getenv("LANG")
	if locale == "" {
		logger.Debug("No 'LANG' found")
		return string(ssid)
	}

	charset, ok := localeCharsetMap[locale]
	if !ok {
		charset, ok = localeCharsetMap[strings.Split(locale, ".")[0]]
		if !ok {
			logger.Debug("No charset found for:", locale)
			return string(ssid)
		}
	}

	_, data, err := mahonia.NewDecoder(charset).Translate(ssid, true)
	if err != nil {
		logger.Error("Failed to decode charset:", charset, err)
		return string(ssid)
	}
	return string(data)
}
