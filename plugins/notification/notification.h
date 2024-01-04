// SPDX-FileCopyrightText: 2024 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later
#ifndef NOTIFICATION_H
#define NOTIFICATION_H

#include <DGuiApplicationHelper>

#include <QWidget>
#include <QIcon>
#include <QDBusVariant>
#include <QDBusInterface>

class Notification : public QWidget
{
    Q_OBJECT

public:
    explicit Notification(QWidget *parent = nullptr);
    QIcon icon() const;

    bool dndMode() const;
    void setDndMode(bool dnd);

Q_SIGNALS:
    void iconRefreshed();
    void dndModeChanged(bool dnd);

public Q_SLOTS:
    void refreshIcon();

private Q_SLOTS:
    void onSystemInfoChanged(quint32 info, QDBusVariant value);

protected:
    void paintEvent(QPaintEvent *e) override;

private:
    QIcon m_icon;
    QDBusInterface *m_dbus;
    bool m_dndMode;
};

#endif  // NOTIFICATION_H
