#ifndef HIGHLIGHTEFFECT_H
#define HIGHLIGHTEFFECT_H

#include <QWidget>

class HighlightEffect : public QWidget
{
    Q_OBJECT
public:
    HighlightEffect(QWidget * source, QWidget *parent = 0);

    int lighter() const;
    void setLighter(int lighter);

protected:
    void paintEvent(QPaintEvent * event) Q_DECL_OVERRIDE;

private:
    QWidget * m_source;
    int m_lighter;

    void pixmapLigher(QPixmap * pixmap, int lighter);
};

#endif // HIGHLIGHTEFFECT_H
