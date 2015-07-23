#ifndef REFLECTIONEFFECT_H
#define REFLECTIONEFFECT_H

#include <QWidget>

class QPaintEvent;
class ReflectionEffect : public QWidget
{
    Q_OBJECT
public:
    ReflectionEffect(QWidget * source, QWidget *parent = 0);

    qreal opacity() const;
    void setOpacity(const qreal &opacity);
    void updateReflection();

protected:
    void paintEvent(QPaintEvent * event) Q_DECL_OVERRIDE;

private:
    QWidget * m_source;
    qreal m_opacity;
};

#endif // REFLECTIONEFFECT_H
