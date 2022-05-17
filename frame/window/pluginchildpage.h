/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#ifndef PLUGINCHILDPAGE_H
#define PLUGINCHILDPAGE_H

#include <QWidget>

class QLabel;
class QVBoxLayout;

class PluginChildPage : public QWidget
{
    Q_OBJECT

public:
    explicit PluginChildPage(QWidget *parent);
    ~PluginChildPage() override;
    void pushWidget(QWidget *widget);
    void setTitle(const QString &text);
    bool isBack();

Q_SIGNALS:
    void back();
    void closeSelf();

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    void initUi();
    void resetHeight();

private:
    QWidget *m_headerWidget;
    QLabel *m_back;
    QLabel *m_title;
    QWidget *m_container;
    QWidget *m_topWidget;
    QVBoxLayout *m_containerLayout;
    bool m_isBack;
};

#endif // PLUGINCHILDPAGE_H
