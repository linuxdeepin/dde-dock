/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package bluetooth

import . "pkg.deepin.io/lib/gettext"

type profile struct {
	uuid, name string
}

var profiles = []profile{
	profile{SPP_UUID, Tr("Serial port")},
	profile{DUN_GW_UUID, Tr("Dial-Up networking")},
	profile{HFP_HS_UUID, Tr("Hands-Free device")},
	profile{HFP_AG_UUID, Tr("Hands-Free voice gateway")},
	profile{HSP_AG_UUID, Tr("Headset voice gateway")},
	profile{OBEX_OPP_UUID, Tr("Object push")},
	profile{OBEX_FTP_UUID, Tr("File transfer")},
	profile{OBEX_SYNC_UUID, Tr("Synchronization")},
	profile{OBEX_PSE_UUID, Tr("Phone book access")},
	profile{OBEX_PCE_UUID, Tr("Phone book access client")},
	profile{OBEX_MAS_UUID, Tr("Message access")},
	profile{OBEX_MNS_UUID, Tr("Message notification")},
}
