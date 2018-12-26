#ifndef FASHIONTRAYHOLDCONTAINER_H
#define FASHIONTRAYHOLDCONTAINER_H

#include "constants.h"
#include "fashiontraywidgetwrapper.h"

#include <QWidget>
#include <QBoxLayout>
#include <QLabel>

class FashionTrayHoldContainer : public QWidget
{
    Q_OBJECT
public:
    explicit FashionTrayHoldContainer(Dock::Position dockPosistion, QWidget *parent = nullptr);

    void setDockPostion(Dock::Position pos);
    void setTrayExpand(const bool expand);

public:
    QSize sizeHint() const Q_DECL_OVERRIDE;
    void resizeEvent(QResizeEvent *event) Q_DECL_OVERRIDE;

private:
    QBoxLayout *m_mainBoxLayout;
    QLabel *m_holdSpliter;

    bool m_expand;

    Dock::Position m_dockPosistion;

    QList<QPointer<FashionTrayWidgetWrapper>> m_holdWrapperList;
};

#endif // FASHIONTRAYHOLDCONTAINER_H
