/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef PLUGINSSETTINGFRAME_H
#define PLUGINSSETTINGFRAME_H

#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QMouseEvent>
#include <QPushButton>
#include <QFrame>
#include <QLabel>
#include <QDebug>

#include <libdui/dswitchbutton.h>
#include <libdui/dseparatorhorizontal.h>

DUI_USE_NAMESPACE

class PluginsSettingLine : public QLabel
{
    Q_OBJECT
public:
    explicit PluginsSettingLine(bool checked = false,
                                const QString &id = "",
                                const QString &title = "",
                                const QPixmap &icon = QPixmap(),
                                QWidget *parent = 0);

    void setIcon(const QPixmap &icon);
    void setTitle(const QString &title);
    void setPluginId(const QString &pluginId);
    void setChecked(const bool checked);

    QString pluginId() const;
    bool checked();

signals:
    void checkedChanged(QString id, bool check);

private:
    QLabel *m_iconLabel = NULL;
    QLabel *m_titleLabel = NULL;
    DSwitchButton *m_switchButton = NULL;

    QString m_pluginId = "";

    const int ICON_SIZE = 16;
    const int ICON_SPACING = 6;
    const int CONTENT_MARGIN = 6;
    const int MAX_TEXT_WIDTH = 125;
};





class PluginsSettingFrame : public QFrame
{
    Q_OBJECT
public:
    explicit PluginsSettingFrame(QWidget *parent = 0);

public slots:
    void onPluginAdd(bool checked = false,
                     const QString &id = "",
                     const QString &title = "",
                     const QPixmap &icon = QPixmap());
    void onPluginRemove(const QString &id);
    void onPluginEnabledChanged(const QString &id, bool enabled);
    void onPluginTitleChanged(const QString &id, const QString &title);

signals:
    void checkedChanged(const QString &id, bool checked);

protected:
    void mouseMoveEvent(QMouseEvent *event);
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);
    bool eventFilter(QObject *obj, QEvent *event);

private:
    void resizeWithLineCount();
    void initCloseTitle();

    bool m_mousePressed;
    QPoint m_pressPosition;
    QVBoxLayout *m_mainLayout;
    QMap<QString, PluginsSettingLine *> m_lineMap;

    const int ICON_SIZE = 24;
    const int CONTENT_MARGIN = 6;
    const int LINE_SPACING = 5;
    const int LINE_HEIGHT = 30;
    const int WIN_WIDTH = 230;
};

#endif // PLUGINSSETTINGFRAME_H
