#ifndef FASHIONTRAYHOLDCONTAINER_H
#define FASHIONTRAYHOLDCONTAINER_H

#include "constants.h"
#include "trayplugin.h"
#include "fashiontraywidgetwrapper.h"

#include <QWidget>
#include <QBoxLayout>
#include <QLabel>

class FashionTrayHoldContainer : public QWidget
{
    Q_OBJECT
public:
    explicit FashionTrayHoldContainer(TrayPlugin *trayPlugin, QWidget *parent = nullptr);

    void setDockPostion(Dock::Position pos);
    void setTrayExpand(const bool expand);

    bool exists(FashionTrayWidgetWrapper *wrapper) const;
    bool isHoldTrayWrapper(FashionTrayWidgetWrapper *wrapper) const;

    void addTrayWrapper(FashionTrayWidgetWrapper *wrapper);
    bool removeTrayWrapper(FashionTrayWidgetWrapper *wrapper);

    int whereToInsert(FashionTrayWidgetWrapper *wrapper) const;

public:
    QSize sizeHint() const Q_DECL_OVERRIDE;
    void resizeEvent(QResizeEvent *event) Q_DECL_OVERRIDE;

private:
    QBoxLayout *m_mainBoxLayout;
    QLabel *m_holdSpliter;

    TrayPlugin *m_trayPlugin;

    bool m_expand;

    Dock::Position m_dockPosistion;

    QList<QPointer<FashionTrayWidgetWrapper>> m_holdWrapperList;
};

#endif // FASHIONTRAYHOLDCONTAINER_H
