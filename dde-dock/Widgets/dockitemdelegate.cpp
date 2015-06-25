#include "dockitemdelegate.h"

DockItemDelegate::DockItemDelegate(QObject *parent) : QItemDelegate(parent)
{

}

QWidget * DockItemDelegate::createEditor(QWidget *parent, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    AppItem * editor = new AppItem(parent);
    editor->resize(50,50);

    return editor;
}

void DockItemDelegate::setEditorData(QWidget *editor, const QModelIndex &index) const
{
    QJsonObject dataObj = index.model()->data(index, Qt::DisplayRole).toJsonValue().toObject();

    if (dataObj.isEmpty())
    {
        return;
    }

    AppItem *appItem = static_cast<AppItem*>(editor);
    if (dataObj.contains("itemIconPath"))
        appItem->setIcon(dataObj.value("itemIconPath").toString());
    if (dataObj.contains("itemTitle"))
        appItem->setTitle(dataObj.value("itemTitle").toString());
}

void DockItemDelegate::setModelData(QWidget *editor, QAbstractItemModel *model, const QModelIndex &index) const
{

    AppItem *appItem = static_cast<AppItem*>(editor);
//    appItem->interpretText();
//    int value = appItem->value();

//    model->setData(index, value);
}

void DockItemDelegate::updateEditorGeometry(QWidget *editor, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    editor->setGeometry(option.rect);
}

DockItemDelegate::~DockItemDelegate()
{

}

