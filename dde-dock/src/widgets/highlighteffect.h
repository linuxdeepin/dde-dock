#ifndef HIGHLIGHTEFFECT_H
#define HIGHLIGHTEFFECT_H

#include <QColor>
#include <QWidget>
#include <QBitmap>
#include <QPainter>

class HighlightEffect : public QWidget
{
    Q_OBJECT
public:
    HighlightEffect(QWidget * source, QWidget *parent = 0);

    enum EffectState {
        ESNormal,
        ESLighter,
        ESDarker
    };

    int lighter() const;
    void setLighter(int lighter);
    int darker() const;
    void setDarker(int darker);

    void showDarker();
    void showLighter();
    void showNormal();

protected:
    void resizeEvent(QResizeEvent *) Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *) Q_DECL_OVERRIDE;

private:
    QWidget * m_source;
    int m_lighter = 110;
    int m_darker = 150;
    EffectState m_effectState = ESNormal;


    void pixmapLigher(QPixmap * pixmap);
    void pixmapDarker(QPixmap * pixmap);
};

#endif // HIGHLIGHTEFFECT_H
