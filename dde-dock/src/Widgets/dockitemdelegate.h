#ifndef DOCKITEMDELEGATE_H
#define DOCKITEMDELEGATE_H

#include <QObject>
#include <QJsonObject>
#include <QItemDelegate>
#include <QStyleOptionViewItem>
#include <QModelIndex>
#include <QAbstractItemModel>
#include <QDebug>
#include "appitem.h"

class DockItemDelegate : public QItemDelegate
{
    Q_OBJECT
public:
    explicit DockItemDelegate(QObject *parent = 0);
    ~DockItemDelegate();

    QWidget *createEditor(QWidget *parent, const QStyleOptionViewItem &option,const QModelIndex &index) const;
    void setEditorData(QWidget * editor, const QModelIndex & index) const;
    void setModelData(QWidget * editor, QAbstractItemModel * model, const QModelIndex & index) const;
    void updateEditorGeometry(QWidget *editor, const QStyleOptionViewItem &option, const QModelIndex &index) const;

signals:

public slots:
};

#endif // DOCKITEMDELEGATE_H
