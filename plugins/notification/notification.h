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

    Q_PROPERTY(bool dndMode READ dndMode WRITE setDndMode NOTIFY dndModeChanged)
    Q_PROPERTY(uint notificationCount READ notificationCount NOTIFY notificationCountChanged)

public:
    explicit Notification(QWidget *parent = nullptr);
    QIcon icon() const;

    bool dndMode() const;
    void setDndMode(bool dnd);
    uint notificationCount() const;

Q_SIGNALS:
    void iconRefreshed();
    void dndModeChanged(bool dnd);
    void notificationCountChanged(uint count);

public Q_SLOTS:
    void refreshIcon();

private Q_SLOTS:
    void onSystemInfoChanged(quint32 info, QDBusVariant value);
    void setNotificationCount(uint count);

protected:
    void paintEvent(QPaintEvent *e) override;

private:
    QIcon m_icon;
    uint m_notificationCount;
    QScopedPointer<QDBusInterface> m_dbus;
    bool m_dndMode;
};

#endif  // NOTIFICATION_H
