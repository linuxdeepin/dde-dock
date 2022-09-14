// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef VOLUMESLIDER_H
#define VOLUMESLIDER_H

#include <DSlider>
#include <QTimer>

class VolumeSlider : public DTK_WIDGET_NAMESPACE::DSlider
{
    Q_OBJECT

public:
    explicit VolumeSlider(QWidget *parent = 0);

    void setValue(const int value);

signals:
    void requestPlaySoundEffect() const;

protected:
    void mousePressEvent(QMouseEvent *e) override;
    void mouseMoveEvent(QMouseEvent *e) override;
    void mouseReleaseEvent(QMouseEvent *e) override;
    void wheelEvent(QWheelEvent *e) override;

private slots:
    void onTimeout();

private:
    bool m_pressed;
    QTimer *m_timer;
};

#endif // VOLUMESLIDER_H
