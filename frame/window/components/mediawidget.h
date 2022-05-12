#ifndef MEDIAWIDGET_H
#define MEDIAWIDGET_H

#include "mediaplayermodel.h"

#include <DBlurEffectWidget>

class QLabel;
class MusicButton;

DWIDGET_USE_NAMESPACE

class MediaWidget : public DBlurEffectWidget
{
    Q_OBJECT

public:
    explicit MediaWidget(QWidget *parent = nullptr);
    ~MediaWidget() override;

Q_SIGNALS:
    void visibleChanged(bool);

protected:
    void showEvent(QShowEvent *event) override;
    void hideEvent(QHideEvent *event) override;

private Q_SLOTS:
    void statusChanged(const MediaPlayerModel::PlayStatus &newStatus);
    void onPlayClicked();
    void onNext();
    void onUpdateMediaInfo();

private:
    void initUi();
    void initConnection();

private:
    QLabel *m_musicIcon;
    QLabel *m_musicName;
    QLabel *m_musicSinger;
    MusicButton *m_pausePlayButton;
    MusicButton *m_nextButton;
};

// 音乐播放按钮
class MusicButton : public QWidget
{
    Q_OBJECT

Q_SIGNALS:
    void clicked();

public:
    enum ButtonType { Playing = 0, Pause, Next };

public:
    MusicButton(QWidget *parent = Q_NULLPTR);
    ~MusicButton() override;

    inline void setButtonType(const ButtonType &bt) {
        m_buttonType = bt;
        update();
    }

protected:
    void paintEvent(QPaintEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;

private:
    ButtonType m_buttonType;
};

#endif // MEDIAWIDGER_H
