/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     fpc_diesel <fanpengcheng@uniontech.com>
 *
 * Maintainer: fpc_diesel <fanpengcheng@uniontech.com>
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
#include "accessible.h"
#define SEPARATOR "_"

QString getAccesibleName(QWidget *w, QAccessible::Role r, QString fallback)
{
    // 避免重复生成
    static QMap< QObject *, QString > objnameMap;
    if (!objnameMap[w].isEmpty())
        return objnameMap[w];

    static QMap< QAccessible::Role, QList< QString > > accessibleMap;
    QString oldAccessName = w->accessibleName();
    oldAccessName.replace(SEPARATOR, "");

    // 按照类型添加固定前缀
    QMetaEnum metaEnum = QMetaEnum::fromType<QAccessible::Role>();
    QByteArray prefix = metaEnum.valueToKeys(r);
    switch (r) {
    case QAccessible::Button:       prefix = "Btn";         break;
    case QAccessible::StaticText:   prefix = "Label";       break;
    default:                        break;
    }

    // 再加上标识
    QString accessibleName = QString::fromLatin1(prefix) + SEPARATOR;
    accessibleName += oldAccessName.isEmpty() ? fallback : oldAccessName;

    // 检查名称是否唯一
    if (accessibleMap[r].contains(accessibleName)) {
        // 获取编号，然后+1
        int pos = accessibleName.indexOf(SEPARATOR);
        int id = accessibleName.mid(pos + 1).toInt();

        QString newAccessibleName;
        do {
            // 一直找到一个不重复的名字
            newAccessibleName = accessibleName + SEPARATOR + QString::number(++id);
        } while (accessibleMap[r].contains(newAccessibleName));

        accessibleMap[r].append(newAccessibleName);
        objnameMap.insert(w, newAccessibleName);

        return newAccessibleName;
    } else {
        accessibleMap[r].append(accessibleName);
        objnameMap.insert(w, accessibleName);

        return accessibleName;
    }
}

QString AccessibleMainWindow::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "mainwindow");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleMainPanelControl::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "mainpanelcontrol");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleTipsWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "Tips");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleDockPopupWindow::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "PopupWindow");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleLauncherItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "launcheritem");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleAppItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "AppItem");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessiblePreviewContainer::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "previewcontainer");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessiblePluginsItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        qDebug() << m_w->pluginName();
        return getAccesibleName(m_w, this->role(), m_w->pluginName());
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleTrayPluginItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), m_w->pluginName());
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessiblePlaceholderItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "placeholderitem");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleAppDragWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "appdragwidget");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleAppSnapshot::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "appsnapshot");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleFloatingPreview::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "floatingpreview");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleSNITrayWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), m_w->itemKeyForConfig());
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleSystemTrayItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), m_w->itemKeyForConfig());
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleFashionTrayItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "fashiontrayitem");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleFashionTrayWidgetWrapper::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "fashiontraywrapper");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleFashionTrayControlWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "fashiontraycontrolwidget");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleAttentionContainer::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "attentioncontainer");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleHoldContainer::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "holdcontainer");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleNormalContainer::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "normalcontainer");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleSpliterAnimated::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "spliteranimated");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleIndicatorTrayWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), m_w->itemKeyForConfig());
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleXEmbedTrayWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), m_w->itemKeyForConfig());
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleSoundItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "plugin-sounditem");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleSoundApplet::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "soundapplet");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleSinkInputWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "sinkinputwidget");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleVolumeSlider::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "volumeslider");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleHorizontalSeparator::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "horizontalseparator");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleDatetimeWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "plugin-datetime");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleOnboardItem::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "plugin-onboard");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleTrashWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "plugin-trash");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessiblePopupControlWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "popupcontrolwidget");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleShutdownWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "plugin-shutdown");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleMultitaskingWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "plugin-multitasking");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleShowDesktopWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), "plugin-showdesktop");
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

//QString AccessibleOverlayWarningWidget::text(QAccessible::Text t) const
//{
//    switch (t) {
//    case QAccessible::Name:
//        return "overlaywarningwidget";
//    case QAccessible::Description:
//        return m_description;
//    default:
//        return QString();
//    }
//}

QString AccessibleQWidget::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return m_w->objectName();
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleDImageButton::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), m_w->objectName());
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

QString AccessibleDSwitchButton::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w, this->role(), m_w->text());
    case QAccessible::Description:
        return m_description;
    default:
        return QString();
    }
}

