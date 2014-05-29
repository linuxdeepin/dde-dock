package main

import . "dlib/gettext"

type profile struct {
	uuid, name string
}

var profiles = []profile{
	profile{SPP_UUID, Tr("Serial Port")},
	profile{DUN_GW_UUID, Tr("Dial-Up Networking")},
	profile{HFP_HS_UUID, Tr("Hands-Free unit")},
	profile{HFP_AG_UUID, Tr("Hands-Free Voice Gateway")},
	profile{HSP_AG_UUID, Tr("Headset Voice Gateway")},
	profile{OBEX_OPP_UUID, Tr("Object Push")},
	profile{OBEX_FTP_UUID, Tr("File Transfer")},
	profile{OBEX_SYNC_UUID, Tr("Synchronization")},
	profile{OBEX_PSE_UUID, Tr("Phone Book Access")},
	profile{OBEX_PCE_UUID, Tr("Phone Book Access Client")},
	profile{OBEX_MAS_UUID, Tr("Message Access")},
	profile{OBEX_MNS_UUID, Tr("Message Notification")},
}
