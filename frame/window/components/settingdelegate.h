#ifndef SETTINGDELEGATE_H
#define SETTINGDELEGATE_H

#include <DStyledItemDelegate>

DWIDGET_USE_NAMESPACE

static const int itemCheckRole = Dtk::UserRole + 1;
static const int itemDataRole = Dtk::UserRole + 2;

class SettingDelegate : public DStyledItemDelegate
{
    Q_OBJECT

Q_SIGNALS:
    void selectIndexChanged(const QModelIndex &);

public:
    explicit SettingDelegate(QAbstractItemView *parent = nullptr);
    ~SettingDelegate() override;

protected:
    void paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const override;
    bool editorEvent(QEvent *event, QAbstractItemModel *model, const QStyleOptionViewItem &option, const QModelIndex &index) override;
};

#endif // SETTINGDELEGATE_H
