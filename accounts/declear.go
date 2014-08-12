/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package accounts

type UserInfo struct {
	Name   string
	Uid    string
	Gid    string
	Home   string
	Shell  string
	Path   string
	Locked bool
}

const (
	ACCOUNT_DEST         = "com.deepin.daemon.Accounts"
	ACCOUNT_MANAGER_PATH = "/com/deepin/daemon/Accounts"
	ACCOUNT_MANAGER_IFC  = "com.deepin.daemon.Accounts"
	USER_MANAGER_PATH    = "/com/deepin/daemon/Accounts/User"
	USER_MANAGER_IFC     = "com.deepin.daemon.Accounts.User"
)

const (
	ETC_PASSWD          = "/etc/passwd"
	ETC_SHADOW          = "/etc/shadow"
	ETC_SHADOW_BAK      = "/etc/shadow.bak"
	ETC_GROUP           = "/etc/group"
	ETC_DISPLAY_MANAGER = "/etc/X11/default-display-manager"
	ETC_LIGHTDM_CONFIG  = "/etc/lightdm/lightdm.conf"
	ETC_GDM_CONFIG      = "/etc/gdm/custom.conf"
	ETC_KDM_CONFIG      = "/etc/kde4/kdm/kdmrc"
	USER_KDM_CONFIG     = "/usr/share/config/kdm/kdmrc"

	LIGHTDM_AUTOLOGIN_GROUP = "SeatDefaults"
	LIGHTDM_AUTOLOGIN_USER  = "autologin-user"
	GDM_AUTOLOGIN_GROUP     = "daemon"
	GDM_AUTOLOGIN_USER      = "AutomaticLogin"
	KDM_AUTOLOGIN_GROUP     = "X-:0-Core"
	KDM_AUTOLOGIN_ENABLE    = "AutoLoginEnable"
	KDM_AUTOLOGIN_USER      = "AutoLoginUser"

	ETC_PERM         = 0644
	PASSWD_SPLIT_LEN = 7
	SHADOW_SPLIT_LEN = 9
	GROUP_SPLIT_LEN  = 4
)

const (
	SHELL_END_FALSE   = "false"
	SHELL_END_NOLOGIN = "nologin"

	KEY_TYPE_BOOL        = 0
	KEY_TYPE_INT         = 1
	KEY_TYPE_STRING      = 2
	KEY_TYPE_STRING_LIST = 3

	ACCOUNT_TYPE_STANDARD      = 0
	ACCOUNT_TYPE_ADMINISTACTOR = 1

	GUEST_USER_ICON     = "/var/lib/AccountsService/icons/guest.png"
	ACCOUNT_CONFIG_FILE = "/var/lib/AccountsService/accounts.ini"
	ACCOUNT_GROUP_KEY   = "Accounts"
	ACCOUNT_KEY_GUEST   = "AllowGuest"

	USER_ICON_DIR     = "/var/lib/AccountsService/icons/"
	USER_DEFAULT_ICON = USER_ICON_DIR + "1.png"
	USER_CONFIG_DIR   = "/var/lib/AccountsService/users"
	ICON_SYSTEM_DIR   = "/var/lib/AccountsService/icons"
	ICON_LOCAL_DIR    = "/var/lib/AccountsService/icons/local"
	USER_DEFAULT_BG   = "file:///usr/share/backgrounds/default_background.jpg"
)

const (
	CMD_USERADD = "/usr/sbin/useradd"
	CMD_USERDEL = "/usr/sbin/userdel"
	CMD_CHOWN   = "/bin/chown"
	CMD_USERMOD = "/usr/sbin/usermod"
	CMD_GPASSWD = "/usr/bin/gpasswd"

	POLKIT_CHANGED_OWN_DATA = "com.deepin.daemon.accounts.change-own-user-data"
	POLKIT_MANAGER_USER     = "com.deepin.daemon.accounts.user-administration"
	POLKIT_SET_LOGIN_OPTION = "com.deepin.daemon.accounts.set-login-option"
)
