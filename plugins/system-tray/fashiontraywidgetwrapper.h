#ifndef FASHIONTRAYWIDGETWRAPPER_H
#define FASHIONTRAYWIDGETWRAPPER_H

#include "abstracttraywidget.h"

#include <QWidget>
#include <QVBoxLayout>
#include <QTimer>

class FashionTrayWidgetWrapper : public QWidget
{
    Q_OBJECT
public:
    FashionTrayWidgetWrapper(AbstractTrayWidget *absTrayWidget, QWidget *parent = nullptr);

    AbstractTrayWidget *absTrayWidget() const;

    bool attention() const;

Q_SIGNALS:
    void attentionChanged(const bool attention);

protected:
    void paintEvent(QPaintEvent *event) Q_DECL_OVERRIDE;

private:
    void setAttention(bool attention);
    void onTrayWidgetIconChanged();
    void onTrayWidgetClicked();

private:
    AbstractTrayWidget *m_absTrayWidget;
    QVBoxLayout *m_layout;

    bool m_attention;
};

#endif //FASHIONTRAYWIDGETWRAPPER_H
