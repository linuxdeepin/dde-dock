/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
