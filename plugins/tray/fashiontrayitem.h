/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#ifndef FASHIONTRAYITEM_H
#define FASHIONTRAYITEM_H

#include "constants.h"
#include "fashiontraywidgetwrapper.h"
#include "fashiontraycontrolwidget.h"

#include <QWidget>
#include <QPointer>
#include <QBoxLayout>
#include <QLabel>

#include <abstracttraywidget.h>

class FashionTrayItem : public QWidget
{
    Q_OBJECT

public:
    explicit FashionTrayItem(Dock::Position pos, QWidget *parent = 0);

    void setTrayWidgets(const QMap<QString, AbstractTrayWidget *> &itemTrayMap);
    void trayWidgetAdded(const QString &itemKey, AbstractTrayWidget *trayWidget);
    void trayWidgetRemoved(AbstractTrayWidget *trayWidget);
    void clearTrayWidgets();

    void setDockPostion(Dock::Position pos);

public slots:
    void onTrayListExpandChanged(const bool expand);
    void setSuggestIconSize(QSize size);
    void setRightSplitVisible(const bool visible);

protected:
    void showEvent(QShowEvent *event) Q_DECL_OVERRIDE;
    void hideEvent(QHideEvent *event) Q_DECL_OVERRIDE;
    void resizeEvent(QResizeEvent *event) Q_DECL_OVERRIDE;

private:
    QSize sizeHint() const Q_DECL_OVERRIDE;
    QSize wantedTotalSize() const;

private Q_SLOTS:
    void onTrayAttentionChanged(const bool attention);
    void setCurrentAttentionTray(FashionTrayWidgetWrapper *attentionWrapper);
    void requestResize();
    void moveOutAttionTray();
    void moveInAttionTray();
    void switchAttionTray(FashionTrayWidgetWrapper *attentionWrapper);
    void requestWindowAutoHide(const bool autoHide);
    void requestRefershWindowVisible();
    void expandWithAnimation();
    void foldWithAnimation();
    void refreshTraysVisible();

private:
    QList<QPointer<FashionTrayWidgetWrapper>> m_wrapperList;
    QBoxLayout *m_mainBoxLayout;
    QBoxLayout *m_trayBoxLayout;
    QLabel *m_leftSpliter;
    QLabel *m_rightSpliter;
    FashionTrayControlWidget *m_controlWidget;
    FashionTrayWidgetWrapper *m_currentAttentionTray;
    QTimer *m_attentionDelayTimer;

    Dock::Position m_dockPosistion;

    static int TrayWidgetWidth;
    static int TrayWidgetHeight;
};

#endif // FASHIONTRAYITEM_H
