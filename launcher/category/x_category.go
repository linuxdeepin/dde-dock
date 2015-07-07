package category

import (
	. "pkg.deepin.io/dde-daemon/launcher/interfaces"
)

// AudioVideo/Audio/Video -> multimedia
// Development -> Development
// Education -> Education
// Game -> Games
// Graphics -> Graphics
// Network -> Network
// System/settings -> System
// Utility -> Utilities
// office -> Productivity
// -> Industry
var xCategoryNameIdMap map[string]CategoryId = map[string]CategoryId{
	"network":           NetworkID,
	"webbrowser":        NetworkID,
	"email":             NetworkID,
	"contactmanagement": NetworkID, // productivity
	"filetransfer":      NetworkID,
	"p2p":               NetworkID,
	"instantmessaging":  NetworkID,
	"chat":              NetworkID,
	"ircclient":         NetworkID,
	"news":              NetworkID,
	"remoteaccess":      NetworkID,

	"tv":                MultimediaID,
	"multimedia":        MultimediaID,
	"audio":             MultimediaID,
	"video":             MultimediaID,
	"audiovideo":        MultimediaID,
	"audiovideoediting": MultimediaID,
	"discburning":       MultimediaID,
	"midi":              MultimediaID,
	"mixer":             MultimediaID,
	"player":            MultimediaID,
	"music":             MultimediaID,
	"recorder":          MultimediaID,
	"sequencer":         MultimediaID,
	"tuner":             MultimediaID,

	"game":          GamesID,
	"amusement":     GamesID,
	"actiongame":    GamesID,
	"adventuregame": GamesID,
	"arcadegame":    GamesID,
	"emulator":      GamesID, // system or game
	"simulation":    GamesID,
	"kidsgame":      GamesID,
	"logicgame":     GamesID,
	"roleplaying":   GamesID,
	"sportsgame":    GamesID,
	"strategygame":  GamesID,

	"graphics":        GraphicsID,
	"2dgraphics":      GraphicsID,
	"3dgraphics":      GraphicsID,
	"imageprocessing": GraphicsID, // education
	"ocr":             GraphicsID,
	"photography":     GraphicsID,
	"rastergraphics":  GraphicsID,
	"vectorgraphics":  GraphicsID,
	"viewer":          GraphicsID,

	"office":            ProductivityID,
	"spreadsheet":       ProductivityID,
	"wordprocessor":     ProductivityID,
	"projectmanagement": ProductivityID,
	"chart":             ProductivityID,
	"numericalanalysis": ProductivityID, // education
	"presentation":      ProductivityID,
	"scanning":          ProductivityID, // graphics
	"printing":          ProductivityID, // system

	"engineering":     IndustryID,
	"telephonytools":  IndustryID, // utilities
	"telephony":       IndustryID, // network
	"finance":         IndustryID, // productivity
	"hamradio":        IndustryID,
	"medicalsoftware": IndustryID, // education
	"publishing":      IndustryID,

	"education":              EducationID,
	"art":                    EducationID,
	"literature":             EducationID,
	"dictionary":             EducationID, // productivity
	"artificialintelligence": EducationID,
	"electricity":            EducationID,
	"robotics":               EducationID,
	"geography":              EducationID,
	"computerscience":        EducationID,
	"math":                   EducationID,
	"biology":                EducationID,
	"physics":                EducationID,
	"chemistry":              EducationID,
	"electronics":            EducationID,
	"geology":                EducationID,
	"astronomy":              EducationID,
	"science":                EducationID,

	"development":    DevelopmentID,
	"debugger":       DevelopmentID,
	"ide":            DevelopmentID,
	"building":       DevelopmentID,
	"guidesigner":    DevelopmentID,
	"webdevelopment": DevelopmentID,
	"profiling":      DevelopmentID,
	"transiation":    DevelopmentID,

	"system":         SystemID,
	"settings":       SystemID,
	"monitor":        SystemID,
	"dialup":         SystemID, // network
	"packagemanager": SystemID,
	"filesystem":     SystemID,

	"utility":          UtilitiesID,
	"pda":              UtilitiesID, // productivity
	"accessibility":    UtilitiesID, // system/utilities
	"clock":            UtilitiesID,
	"calendar":         UtilitiesID,
	"calculator":       UtilitiesID,
	"documentation":    UtilitiesID,
	"archiving":        UtilitiesID,
	"compression":      UtilitiesID,
	"filemanager":      UtilitiesID, // system/ utilities
	"filetools":        UtilitiesID, // system/utilities
	"terminalemulator": UtilitiesID, // system
	"texteditor":       UtilitiesID,
	"texttools":        UtilitiesID,
}

var extraXCategoryNameIdMap map[string]CategoryId = map[string]CategoryId{
	"internet":        NetworkID,
	"videoconference": NetworkID,

	"x-jack":           MultimediaID,
	"x-alsa":           MultimediaID,
	"x-multitrack":     MultimediaID,
	"x-sound":          MultimediaID,
	"cd":               MultimediaID,
	"x-midi":           MultimediaID,
	"x-sequencers":     MultimediaID,
	"x-suse-sequencer": MultimediaID,

	"boardgame":                       GamesID,
	"cardgame":                        GamesID,
	"x-debian-applications-emulators": GamesID,
	"puzzlegame":                      GamesID,
	"blocksgame":                      GamesID,
	"x-suse-core-game":                GamesID,

	"x-geeqie": GraphicsID,

	"x-suse-core-office":           ProductivityID,
	"x-mandrivalinux-office-other": ProductivityID,
	"x-turbolinux-office":          ProductivityID,

	"technical":                    IndustryID,
	"x-mandriva-office-publishing": IndustryID,

	"x-kde-edu-misc":     EducationID,
	"translation":        EducationID,
	"x-religion":         EducationID,
	"x-bible":            EducationID,
	"x-islamic-software": EducationID,
	"x-quran":            EducationID,
	"geoscience":         EducationID,
	"meteorology":        EducationID,

	"revisioncontrol": DevelopmentID,

	"trayicon":                    SystemID,
	"x-lxde-settings":             SystemID,
	"x-xfce-toplevel":             SystemID,
	"x-xfcesettingsdialog":        SystemID,
	"x-xfce":                      SystemID,
	"x-kde-utilities-pim":         SystemID,
	"x-kde-internet":              SystemID,
	"x-kde-more":                  SystemID,
	"x-kde-utilities-peripherals": SystemID,
	"kde": SystemID,
	"x-kde-utilities-file":                    SystemID,
	"x-kde-utilities-desktop":                 SystemID,
	"x-gnome-networksettings":                 SystemID,
	"gnome":                                   SystemID,
	"x-gnome-settings-panel":                  SystemID,
	"x-gnome-personalsettings":                SystemID,
	"x-gnome-systemsettings":                  SystemID,
	"desktoputility":                          SystemID,
	"x-misc":                                  SystemID,
	"x-suse-core":                             SystemID,
	"x-red-hat-base-only":                     SystemID,
	"x-novell-main":                           SystemID,
	"x-red-hat-extra":                         SystemID,
	"x-suse-yast":                             SystemID,
	"x-sun-supported":                         SystemID,
	"x-suse-yast-high_availability":           SystemID,
	"x-suse-controlcenter-lookandfeel":        SystemID,
	"x-suse-controlcenter-system":             SystemID,
	"x-red-hat-serverconfig":                  SystemID,
	"x-mandrivalinux-system-archiving-backup": SystemID,
	"x-suse-backup":                           SystemID,
	"x-red-hat-base":                          SystemID,
	"panel":                                   SystemID,
	"x-gnustep":                               SystemID,
	"x-bluetooth":                             SystemID,
	"x-ximian-main":                           SystemID,
	"x-synthesis":                             SystemID,
	"x-digital_processing":                    SystemID,
	"desktopsettings":                         SystemID,
	"x-mandrivalinux-internet-other":          SystemID,
	"systemsettings":                          SystemID,
	"hardwaresettings":                        SystemID,
	"advancedsettings":                        SystemID,
	"x-enlightenment":                         SystemID,
	"compiz":                                  SystemID,

	"consoleonly": UtilitiesID,
	"core":        UtilitiesID,
	"favorites":   UtilitiesID,
	"pim":         UtilitiesID,
	"gpe":         UtilitiesID,
	"motif":       UtilitiesID,
	"applet":      UtilitiesID,
	"accessories": UtilitiesID,
	"wine":        UtilitiesID,
	"wine-programs-accessories": UtilitiesID,
	"playonlinux":               UtilitiesID,
	"screensaver":               UtilitiesID,
	"editors":                   UtilitiesID,
}

// TODO:
// Database:productivity:development:multimedia
// security:system
// flowchart:productivity
// construction:education
// languages: education
// datavisualization: education
// economy: education
// history: education
// sports: education
// parallelcomputing: education
