#ifndef ABSTRACTCONTAINER_H
#define ABSTRACTCONTAINER_H

#include "constants.h"
#include "trayplugin.h"
#include "../fashiontraywidgetwrapper.h"

#include <QWidget>

class AbstractContainer : public QWidget
{
    Q_OBJECT
public:
    explicit AbstractContainer(TrayPlugin *trayPlugin, QWidget *parent = nullptr);

    virtual bool acceptWrapper(FashionTrayWidgetWrapper *wrapper) = 0;
    virtual void refreshVisible() = 0;

    virtual void addWrapper(FashionTrayWidgetWrapper *wrapper);
    virtual bool removeWrapper(FashionTrayWidgetWrapper *wrapper);
    virtual bool removeWrapperByTrayWidget(AbstractTrayWidget *trayWidget);
    virtual FashionTrayWidgetWrapper *takeWrapper(FashionTrayWidgetWrapper *wrapper);
    virtual void setDockPosition(const Dock::Position pos);
    virtual void setExpand(const bool expand);
    virtual QSize totalSize() const;

    QSize sizeHint() const Q_DECL_OVERRIDE;

    void clearWrapper();
    void saveCurrentOrderToConfig();
    void setWrapperSize(QSize size);
    bool isEmpty();
    bool containsWrapper(FashionTrayWidgetWrapper *wrapper);
    bool containsWrapperByTrayWidget(AbstractTrayWidget *trayWidget);
    FashionTrayWidgetWrapper *wrapperByTrayWidget(AbstractTrayWidget *trayWidget);

Q_SIGNALS:
    void attentionChanged(FashionTrayWidgetWrapper *wrapper, const bool attention);

protected:
    virtual int whereToInsert(FashionTrayWidgetWrapper *wrapper);

    TrayPlugin *trayPlugin() const;
    QList<QPointer<FashionTrayWidgetWrapper>> wrapperList() const;
    QBoxLayout *wrapperLayout() const;
    bool expand() const;
    Dock::Position dockPosition() const;
    QSize wrapperSize() const;

private Q_SLOTS:
    void onWrapperAttentionhChanged(const bool attention);
    void onWrapperDragStart();
    void onWrapperDragStop();
    void onWrapperRequestSwapWithDragging();

private:
    TrayPlugin *m_trayPlugin;
    QBoxLayout *m_wrapperLayout;

    QPointer<FashionTrayWidgetWrapper> m_currentDraggingWrapper;
    QList<QPointer<FashionTrayWidgetWrapper>> m_wrapperList;

    bool m_expand;
    Dock::Position m_dockPosition;

    QSize m_wrapperSize;
};

#endif // ABSTRACTCONTAINER_H
