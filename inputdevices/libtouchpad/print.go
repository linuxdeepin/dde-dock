/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package libtouchpad

func (touchpad *Touchpad) printInfo(format string, v ...interface{}) {
	if touchpad.logger == nil {
		return
	}

	touchpad.logger.Infof(format, v...)
}

func (touchpad *Touchpad) debugInfo(format string, v ...interface{}) {
	if touchpad.logger == nil {
		return
	}

	touchpad.logger.Debugf(format, v...)
}

func (touchpad *Touchpad) warningInfo(format string, v ...interface{}) {
	if touchpad.logger == nil {
		return
	}

	touchpad.logger.Warningf(format, v...)
}

func (touchpad *Touchpad) errorInfo(format string, v ...interface{}) {
	if touchpad.logger == nil {
		return
	}

	touchpad.logger.Errorf(format, v...)
}
