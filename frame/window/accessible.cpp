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


QString getAccesibleName(QWidget *w,QString fallback)
{
    return w->accessibleName().isEmpty()?fallback:w->accessibleName();
}

QString AccessibleMainPanelControl::text(QAccessible::Text t) const
{
    switch (t) {
    case QAccessible::Name:
        return getAccesibleName(m_w,"mainpanelcontrol");
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
        return getAccesibleName(m_w,"launcheritem");
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
        return getAccesibleName(m_w,m_w->accessibleName());
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
        return getAccesibleName(m_w,"previewcontainer");
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
        return getAccesibleName(m_w,m_w->pluginName());
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
        return getAccesibleName(m_w,m_w->pluginName());
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
        return getAccesibleName(m_w,"placeholderitem");
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
        return getAccesibleName(m_w,"appdragwidget");
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
        return getAccesibleName(m_w,"appsnapshot");
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
        return getAccesibleName(m_w,"floatingpreview");
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
        return getAccesibleName(m_w,m_w->itemKeyForConfig());
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
        return getAccesibleName(m_w,m_w->itemKeyForConfig());
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
        return getAccesibleName(m_w,"fashiontrayitem");
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
        return getAccesibleName(m_w,"fashiontraywrapper");
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
        return getAccesibleName(m_w,"fashiontraycontrolwidget");
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
        return getAccesibleName(m_w,"attentioncontainer");
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
        return getAccesibleName(m_w,"holdcontainer");
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
        return getAccesibleName(m_w,"normalcontainer");
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
        return getAccesibleName(m_w,"spliteranimated");
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
        return getAccesibleName(m_w,m_w->itemKeyForConfig());
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
        return getAccesibleName(m_w,m_w->itemKeyForConfig());
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
        return getAccesibleName(m_w,"showdesktop");
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
        return getAccesibleName(m_w,"sounditem");
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
        return getAccesibleName(m_w,"soundapplet");
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
        return getAccesibleName(m_w,"sinkinputwidget");
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
        return getAccesibleName(m_w,"volumeslider");
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
        return getAccesibleName(m_w,"horizontalseparator");
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
        return getAccesibleName(m_w,m_w->text());
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
        return getAccesibleName(m_w,"datetimewidget");
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
        return getAccesibleName(m_w,"onboarditem");
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
        return getAccesibleName(m_w,"trashwidget");
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
        return getAccesibleName(m_w,"popupcontrolwidget");
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
        return getAccesibleName(m_w,"shutdownwidget");
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
        return getAccesibleName(m_w,"multitaskingwidget");
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
