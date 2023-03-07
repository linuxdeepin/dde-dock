// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLUGINCHILDPAGE_H
#define PLUGINCHILDPAGE_H

#include <QWidget>
#include <DIconButton>

class QPushButton;
class QLabel;
class QVBoxLayout;

DWIDGET_USE_NAMESPACE

class PluginChildPage : public QWidget
{
    Q_OBJECT

public:
    explicit PluginChildPage(QWidget *parent);
    ~PluginChildPage() override;
    void pushWidget(QWidget *widget);
    void setTitle(const QString &text);

Q_SIGNALS:
    void back();

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    void initUi();
    void initConnection();
    void resetHeight();

private:
    QWidget *m_headerWidget;
    DIconButton *m_back;
    QLabel *m_title;
    QWidget *m_container;
    QWidget *m_topWidget;
    QVBoxLayout *m_containerLayout;
};

#endif // PLUGINCHILDPAGE_H
