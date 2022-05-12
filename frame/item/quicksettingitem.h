#ifndef QUICKSETTINGITEM_H
#define QUICKSETTINGITEM_H

#include "dockitem.h"

class PluginsItemInterface;

class QuickSettingItem : public DockItem
{
    Q_OBJECT

    friend class QuickSettingController;

Q_SIGNALS:
    void detailClicked(PluginsItemInterface *);

public:
    PluginsItemInterface *pluginItem() const;
    ItemType itemType() const override;
    const QPixmap dragPixmap();
    const QString itemKey() const;

protected:
    QuickSettingItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~QuickSettingItem() override;

    void paintEvent(QPaintEvent *e) override;
    QRect iconRect();
    QColor foregroundColor() const;
    QColor backgroundColor() const;
    QColor shadowColor() const;

    void mouseReleaseEvent(QMouseEvent *event) override;

private:
    int xMarginSpace();
    int yMarginSpace();

private:
    PluginsItemInterface *m_pluginInter;
    QString m_itemKey;
};

#endif // QUICKSETTINGITEM_H
