/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
 *
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

#ifndef SWITCHITEM_H
#define SWITCHITEM_H

#include <DSwitchButton>

DWIDGET_USE_NAMESPACE

class QLabel;
class SwitchItem : public QWidget
{
    Q_OBJECT
public:
    explicit SwitchItem(QWidget *parent = nullptr);
    void setChecked(const bool checked = true,bool notify = false);
    void setTitle(const QString &title);
    inline bool checkState() { return m_checkState; }

    inline bool isdefault() { return m_default; }
    inline void setDefault(bool def) { m_default = def; }

signals:
    void checkedChanged(bool checked);
    void justUpdateView(bool checked);

private:
    QLabel *m_title;
    DSwitchButton *m_switchBtn;
    bool m_default;
    bool m_checkState;
};

#endif // SWITCHITEM_H
