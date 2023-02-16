// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SOUNDDEVICEPORT_H
#define SOUNDDEVICEPORT_H

#include <QObject>

class SoundDevicePort : public QObject
{
    Q_OBJECT

public:
    enum Direction {
        Out = 1,
        In = 2
    };

    explicit SoundDevicePort(QObject *parent = nullptr);
    virtual ~SoundDevicePort();

    inline QString id() const { return m_id; }
    void setId(const QString &id);

    inline QString name() const { return m_name; }
    void setName(const QString &name);

    inline QString cardName() const { return m_cardName; }
    void setCardName(const QString &cardName);

    inline bool isActive() const { return m_isActive; }
    void setIsActive(bool isActive);

    inline Direction direction() const { return m_direction; }
    void setDirection(const Direction &direction);

    inline uint cardId() const { return m_cardId; }
    void setCardId(const uint &cardId);

Q_SIGNALS:
    void idChanged(const QString &id) const;
    void nameChanged(const QString &name) const;
    void cardNameChanged(const QString &name) const;
    void isActiveChanged(bool ative) const;
    void directionChanged(Direction direction) const;
    void cardIdChanged(uint cardId) const;

private:
    QString m_id;
    QString m_name;
    uint m_cardId;
    QString m_cardName;
    bool m_isActive;
    Direction m_direction;
};

Q_DECLARE_METATYPE(const SoundDevicePort *)

#endif // SOUNDDEVICEPORT_H
