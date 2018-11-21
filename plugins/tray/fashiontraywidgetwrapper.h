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

#ifndef FASHIONTRAYWIDGETWRAPPER_H
#define FASHIONTRAYWIDGETWRAPPER_H

#include "abstracttraywidget.h"

#include <QWidget>
#include <QVBoxLayout>

class FashionTrayWidgetWrapper : public QWidget
{
    Q_OBJECT
public:
    FashionTrayWidgetWrapper(const QString &itemKey, AbstractTrayWidget *absTrayWidget, QWidget *parent = nullptr);

    AbstractTrayWidget *absTrayWidget() const;
    QString itemKey() const;

    bool attention() const;
    void setAttention(bool attention);

Q_SIGNALS:
    void attentionChanged(const bool attention);

protected:
    void paintEvent(QPaintEvent *event) Q_DECL_OVERRIDE;

private:
    void onTrayWidgetNeedAttention();
    void onTrayWidgetClicked();

private:
    AbstractTrayWidget *m_absTrayWidget;
    QVBoxLayout *m_layout;

    bool m_attention;
    QString m_itemKey;
};

#endif //FASHIONTRAYWIDGETWRAPPER_H
