#ifndef NORMALCONTAINER_H
#define NORMALCONTAINER_H

#include "abstractcontainer.h"

class NormalContainer : public AbstractContainer
{
    Q_OBJECT
public:
    explicit NormalContainer(TrayPlugin *trayPlugin, QWidget *parent = nullptr);

    // AbstractContainer interface
public:
    bool acceptWrapper(FashionTrayWidgetWrapper *wrapper) override;
    void addWrapper(FashionTrayWidgetWrapper *wrapper) override;
    void refreshVisible() override;
    void setExpand(const bool expand) override;
    int itemCount() override;
    QSize sizeHint() const override;
    void updateSize();

protected:
    int whereToInsert(FashionTrayWidgetWrapper *wrapper) override;
    void resizeEvent(QResizeEvent *event) override;

private:
    int whereToInsertByDefault(FashionTrayWidgetWrapper *wrapper) const;
    int whereToInsertAppTrayByDefault(FashionTrayWidgetWrapper *wrapper) const;
    int whereToInsertSystemTrayByDefault(FashionTrayWidgetWrapper *wrapper) const;
    void compositeChanged();
    void adjustMaxSize(const QSize size);

private:
    mutable QVariantAnimation *m_sizeAnimation;
};

#endif // NORMALCONTAINER_H
