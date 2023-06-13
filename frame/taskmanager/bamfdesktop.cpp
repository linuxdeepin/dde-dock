// SPDX-FileCopyrightText: 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "bamfdesktop.h"

#include <QDir>
#include <qstandardpaths.h>

#define BAMF_INDEX_NAME "bamf-2.index"

BamfDesktop *BamfDesktop::instance()
{
    static BamfDesktop instance;
    return &instance;
}

QString BamfDesktop::fileName(const QString &instanceName) const
{
    for (const BamfData &lineData: m_bamfLineData) {
        if (lineData.instanceName.toLower() == instanceName.toLower()) {
            QString name = lineData.lineData.split("\t").first();
            return QString("%1%2").arg(lineData.directory).arg(name);
        }
    }

    // 如果根据instanceName没有找到，则根据空格来进行分隔
    for (const BamfData &lineData: m_bamfLineData) {
        QStringList lines = lineData.lineData.split("\t");
        if (lines.size() < 2)
            continue;

        QStringList cmds = lines[2].split(" ");
        if (cmds.size() > 1 && cmds[1].toLower() == instanceName.toLower())
            return QString("%1%2").arg(lineData.directory).arg(lines.first());
    }

    return instanceName;
}

BamfDesktop::BamfDesktop()
{
    loadDesktopFiles();
}

BamfDesktop::~BamfDesktop()
{
}

QStringList BamfDesktop::applicationDirs() const
{
    QStringList appDirs = QStandardPaths::standardLocations(QStandardPaths::ApplicationsLocation);
    QStringList directions;
    for (auto appDir : appDirs)
        directions << appDir;

    return directions;
}

void BamfDesktop::loadDesktopFiles()
{
    QStringList directions = applicationDirs();
    for (const QString &direction : directions) {
        // 读取后缀名为
        QDir dir(direction);
        dir.setNameFilters(QStringList() << BAMF_INDEX_NAME);
        QFileInfoList fileList = dir.entryInfoList();
        if (fileList.size() == 0)
            continue;

        QFileInfo fileInfo = fileList.at(0);
        QFile file(fileInfo.absoluteFilePath());
        if (!file.open(QIODevice::ReadOnly | QIODevice::Text))
            continue;

        QList<QPair<QString, QString>> bamfLine;
        while (!file.atEnd()) {
            QString line = file.readLine();
            QStringList part = line.split("\t");
            BamfData bamf;
            bamf.directory = direction;
            if (part.size() > 2)
                bamf.instanceName = part[2].trimmed();
            bamf.lineData = line;
            m_bamfLineData << bamf;
        }
    }
}
