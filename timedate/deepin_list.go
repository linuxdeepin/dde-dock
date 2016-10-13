/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package timedate

import (
	. "pkg.deepin.io/lib/gettext"
)

type zoneDesc struct {
	zone string
	desc string
}

// if zoneWhiteList changed, please update dst_data
var zoneWhiteList = []zoneDesc{
	zoneDesc{
		zone: "Pacific/Niue",
		desc: Tr("Niue"),
	},
	zoneDesc{
		zone: "US/Hawaii",
		desc: Tr("Hawaii"),
	},
	zoneDesc{
		zone: "Pacific/Tahiti",
		desc: Tr("Tahiti"),
	},
	zoneDesc{
		zone: "Pacific/Honolulu",
		desc: Tr("Honolulu"),
	},
	zoneDesc{
		zone: "Pacific/Marquesas",
		desc: Tr("Marquesas"),
	},
	zoneDesc{
		zone: "US/Alaska",
		desc: Tr("Alaska"),
	},
	zoneDesc{
		zone: "America/Juneau",
		desc: Tr("Juneau"),
	},
	zoneDesc{
		zone: "Pacific/Gambier",
		desc: Tr("Gambier"),
	},
	zoneDesc{
		zone: "Mexico/BajaNorte",
		desc: Tr("BajaNorte"),
	},
	zoneDesc{
		zone: "America/Vancouver",
		desc: Tr("Vancouver"),
	},
	zoneDesc{
		zone: "America/Metlakatla",
		desc: Tr("Metlakatla"),
	},
	zoneDesc{
		zone: "America/Chihuahua",
		desc: Tr("Chihuahua"),
	},
	zoneDesc{
		zone: "US/Arizona",
		desc: Tr("Arizona"),
	},
	zoneDesc{
		zone: "Mexico/BajaSur",
		desc: Tr("BajaSur"),
	},
	zoneDesc{
		zone: "America/Mexico_City",
		desc: Tr("Mexico City"),
	},
	zoneDesc{
		zone: "America/Chicago",
		desc: Tr("Chicago"),
	},
	zoneDesc{
		zone: "America/Managua",
		desc: Tr("Managua"),
	},
	zoneDesc{
		zone: "America/Monterrey",
		desc: Tr("Monterrey"),
	},
	zoneDesc{
		zone: "America/New_York",
		desc: Tr("New York"),
	},
	zoneDesc{
		zone: "America/Lima",
		desc: Tr("Lima"),
	},
	zoneDesc{
		zone: "America/Bogota",
		desc: Tr("Bogota"),
	},
	zoneDesc{
		zone: "America/Caracas",
		desc: Tr("Caracas"),
	},
	zoneDesc{
		zone: "America/Cuiaba",
		desc: Tr("Cuiaba"),
	},
	zoneDesc{
		zone: "America/Santiago",
		desc: Tr("Santiago"),
	},
	zoneDesc{
		zone: "America/La_Paz",
		desc: Tr("La Paz"),
	},
	zoneDesc{
		zone: "America/Asuncion",
		desc: Tr("Asuncion"),
	},
	zoneDesc{
		zone: "Canada/Newfoundland",
		desc: Tr("Newfoundland"),
	},
	zoneDesc{
		zone: "America/St_Johns",
		desc: Tr("St Johns"),
	},
	zoneDesc{
		zone: "America/Buenos_Aires",
		desc: Tr("Buenos Aires"),
	},
	zoneDesc{
		zone: "America/Cayenne",
		desc: Tr("Cayenne"),
	},
	zoneDesc{
		zone: "Brazil/DeNoronha",
		desc: Tr("DeNoronha"),
	},
	zoneDesc{
		zone: "Atlantic/Azores",
		desc: Tr("Azores"),
	},
	zoneDesc{
		zone: "Atlantic/Cape_Verde",
		desc: Tr("Cape Verde"),
	},
	zoneDesc{
		zone: "Europe/London",
		desc: Tr("London"),
	},
	zoneDesc{
		zone: "Europe/Dublin",
		desc: Tr("Dublin"),
	},
	zoneDesc{
		zone: "Africa/Casablanca",
		desc: Tr("Casablanca"),
	},
	zoneDesc{
		zone: "Africa/Monrovia",
		desc: Tr("Monrovia"),
	},
	zoneDesc{
		zone: "Europe/Madrid",
		desc: Tr("Madrid"),
	},
	zoneDesc{
		zone: "Europe/Paris",
		desc: Tr("Paris"),
	},
	zoneDesc{
		zone: "Europe/Rome",
		desc: Tr("Rome"),
	},
	zoneDesc{
		zone: "Europe/Vienna",
		desc: Tr("Vienna"),
	},
	zoneDesc{
		zone: "Africa/Algiers",
		desc: Tr("Algiers"),
	},
	zoneDesc{
		zone: "Africa/Cairo",
		desc: Tr("Cairo"),
	},
	zoneDesc{
		zone: "Europe/Athens",
		desc: Tr("Athens"),
	},
	zoneDesc{
		zone: "Europe/Istanbul",
		desc: Tr("Istanbul"),
	},
	zoneDesc{
		zone: "Africa/Blantyre",
		desc: Tr("Blantyre"),
	},
	zoneDesc{
		zone: "Africa/Nairobi",
		desc: Tr("Nairobi"),
	},
	zoneDesc{
		zone: "Asia/Tehran",
		desc: Tr("Tehran"),
	},
	zoneDesc{
		zone: "Asia/Muscat",
		desc: Tr("Muscat"),
	},
	zoneDesc{
		zone: "Asia/Baku",
		desc: Tr("Baku"),
	},
	zoneDesc{
		zone: "Europe/Moscow",
		desc: Tr("Moscow"),
	},
	zoneDesc{
		zone: "Asia/Kabul",
		desc: Tr("Kabul"),
	},
	zoneDesc{
		zone: "Asia/Karachi",
		desc: Tr("Karachi"),
	},
	zoneDesc{
		zone: "Asia/Calcutta",
		desc: Tr("Calcutta"),
	},
	zoneDesc{
		zone: "Asia/Kathmandu",
		desc: Tr("Kathmandu"),
	},
	zoneDesc{
		zone: "Asia/Dhaka",
		desc: Tr("Dhaka"),
	},
	zoneDesc{
		zone: "Asia/Yekaterinburg",
		desc: Tr("Yekaterinburg"),
	},
	zoneDesc{
		zone: "Asia/Rangoon",
		desc: Tr("Rangoon"),
	},
	zoneDesc{
		zone: "Asia/Bangkok",
		desc: Tr("Bangkok"),
	},
	zoneDesc{
		zone: "Asia/Jakarta",
		desc: Tr("Jakarta"),
	},
	zoneDesc{
		zone: "Asia/Beijing",
		desc: Tr("Beijing"),
	},
	zoneDesc{
		zone: "Asia/Hong_Kong",
		desc: Tr("Hong Kong"),
	},
	zoneDesc{
		zone: "Asia/Taipei",
		desc: Tr("Taipei"),
	},
	zoneDesc{
		zone: "Asia/Kuala_Lumpur",
		desc: Tr("Kuala Lumpur"),
	},
	zoneDesc{
		zone: "Australia/Perth",
		desc: Tr("Perth"),
	},
	zoneDesc{
		zone: "Australia/Eucla",
		desc: Tr("Eucla"),
	},
	zoneDesc{
		zone: "Asia/Tokyo",
		desc: Tr("Tokyo"),
	},
	zoneDesc{
		zone: "Asia/Seoul",
		desc: Tr("Seoul"),
	},
	zoneDesc{
		zone: "Australia/Darwin",
		desc: Tr("Darwin"),
	},
	zoneDesc{
		zone: "Australia/Sydney",
		desc: Tr("Sydney"),
	},
	zoneDesc{
		zone: "Pacific/Guam",
		desc: Tr("Guam"),
	},
	zoneDesc{
		zone: "Australia/Melbourne",
		desc: Tr("Melbourne"),
	},
	zoneDesc{
		zone: "Australia/Hobart",
		desc: Tr("Hobart"),
	},
	zoneDesc{
		zone: "Australia/Lord_Howe",
		desc: Tr("Lord Howe"),
	},
	zoneDesc{
		zone: "Pacific/Pohnpei",
		desc: Tr("Pohnpei"),
	},
	zoneDesc{
		zone: "Pacific/Norfolk",
		desc: Tr("Norfolk"),
	},
	zoneDesc{
		zone: "Pacific/Auckland",
		desc: Tr("Auckland"),
	},
	zoneDesc{
		zone: "Asia/Anadyr",
		desc: Tr("Anadyr"),
	},
	zoneDesc{
		zone: "Pacific/Chatham",
		desc: Tr("Chatham"),
	},
	zoneDesc{
		zone: "Pacific/Apia",
		desc: Tr("Apia"),
	},
	zoneDesc{
		zone: "Pacific/Fakaofo",
		desc: Tr("Fakaofo"),
	},
	zoneDesc{
		zone: "Asia/Kolkata",
		desc: Tr("Kolkata"),
	},
	zoneDesc{
		zone: "Asia/Colombo",
		desc: Tr("Colombo"),
	},
	zoneDesc{
		zone: "Asia/Pyongyang",
		desc: Tr("Pyongyang"),
	},
}

func getZoneDesc(zone string) string {
	for _, zdesc := range zoneWhiteList {
		if zdesc.zone == zone {
			return zdesc.desc
		}
	}

	return zone
}
