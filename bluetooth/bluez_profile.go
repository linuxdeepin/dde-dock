package main

import (
	"dlib"
)

type profile struct {
	uuid, name string
}

var profiles = []profile{
	profile{SPP_UUID, dlib.Tr("Serial Port")},
	profile{DUN_GW_UUID, dlib.Tr("Dial-Up Networking")},
	profile{HFP_HS_UUID, dlib.Tr("Hands-Free unit")},
	profile{HFP_AG_UUID, dlib.Tr("Hands-Free Voice Gateway")},
	profile{HSP_AG_UUID, dlib.Tr("Headset Voice Gateway")},
	profile{OBEX_OPP_UUID, dlib.Tr("Object Push")},
	profile{OBEX_FTP_UUID, dlib.Tr("File Transfer")},
	profile{OBEX_SYNC_UUID, dlib.Tr("Synchronization")},
	profile{OBEX_PSE_UUID, dlib.Tr("Phone Book Access")},
	profile{OBEX_PCE_UUID, dlib.Tr("Phone Book Access Client")},
	profile{OBEX_MAS_UUID, dlib.Tr("Message Access")},
	profile{OBEX_MNS_UUID, dlib.Tr("Message Notification")},
}
