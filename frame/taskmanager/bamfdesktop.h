// SPDX-FileCopyrightText: 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef BAMFDESKTOP_H
#define BAMFDESKTOP_H

#include <QObject>
#include <QMap>

struct BamfData {
    QString directory;
    QString instanceName;
    QString lineData;
};

class BamfDesktop
{
public:
    static BamfDesktop *instance();
    QString fileName(const QString &instanceName) const;

protected:
    BamfDesktop();
    ~BamfDesktop();

private:
    QStringList applicationDirs() const;

    void loadDesktopFiles();

private:
    QList<BamfData> m_bamfLineData;
};

#endif // BAMFDESKTOP_H
