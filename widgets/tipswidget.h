/*
 * Copyright (C) 2011 ~ 2018 uniontech Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
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
#ifndef TIPSWIDGET_H
#define TIPSWIDGET_H

#include <QFrame>
namespace Dock {
class TipsWidget : public QFrame
{
    Q_OBJECT
    enum ShowType
    {
        SingleLine,
        MultiLine
    };
public:
    explicit TipsWidget(QWidget *parent = nullptr);

    const QString& text(){return m_text;}
    const QStringList &textList() { return  m_textList; }
    void setText(const QString &text);
    void setTextList(const QStringList &textList);
    
protected:
    void paintEvent(QPaintEvent *event) override;
    bool event(QEvent *event) override;

private:
    QString m_text;
    QStringList m_textList;
    int m_width;
    ShowType m_type;
};
}

#endif // TIPSWIDGET_H
