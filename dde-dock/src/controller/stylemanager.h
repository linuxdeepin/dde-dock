/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef STYLEMANAGER_H
#define STYLEMANAGER_H

#include <QObject>
#include <QSettings>

class StyleManager : public QObject
{
    Q_OBJECT
public:
    static StyleManager *instance();
    QStringList styleNameList();
    QString currentStyle();
    void applyStyle(const QString &styleName);
    void initStyleSheet();

private:
    explicit StyleManager(QObject *parent = 0);
    QStringList getStyleFromFilesystem();
    void initSettings();
    void applyDefaultStyle(const QString &name);
    void applyThirdPartyStyle(const QString &name);

private:
    static StyleManager *m_styleManager;
    QSettings *m_settings;
};

#endif // STYLEMANAGER_H
