#include "dockmodel.h"

DockModel::DockModel(QObject *parent) :
    QAbstractItemModel(parent)
{
}

int DockModel::count()
{
    return this->rowCount(QModelIndex());
}

void DockModel::append(const QJsonObject &dataObj)
{
    insert(count(),dataObj);
}

void DockModel::clear()
{
    this->removeRows(0,count());
}

QJsonObject DockModel::get(int index)
{
    QModelIndex tmpIndex = getIndex(index);
    QJsonObject tmpObj = this->data(tmpIndex,0).toJsonValue().toObject();

    return tmpObj;
}

bool DockModel::insert(int index, const QJsonObject &dataObj)
{
    if (insertRow(index))
    {
        if (setData(getIndex(index),QVariant(dataObj)))
        {
            return true;
        }
        else
            qWarning() << "setData error";
    }
    qWarning() << "insertRow error" ;
    return false;
}

void DockModel::move(int from, int to, int count)
{
    Q_UNUSED(from)
    Q_UNUSED(to)
    Q_UNUSED(count)
}

void DockModel::remove(int index, int count)
{
    this->removeRows(index,count);
}

void DockModel::set(int index, const QJsonObject &dataObj)
{

}

void DockModel::setProperty(int index, const QString &property, const QVariant &value)
{

}

int DockModel::indexOf(const QString &property)
{

}

QModelIndex DockModel::getIndex(int row)
{
    return this->index(row,0,QModelIndex());
}

bool DockModel::setData(const QModelIndex &index, const QVariant &value, int role)
{
    if (index.isValid() && role == Qt::EditRole)
    {
        dataArray.replace(index.row(),QJsonValue(value.toJsonObject()));
        emit dataChanged(index, index);
        return true;
    }
    return false;
}

QModelIndex DockModel::index(int row, int column, const QModelIndex &parent) const
{
    Q_UNUSED(column)
    Q_UNUSED(parent)
    return this->createIndex(row,0);
}

QModelIndex DockModel::parent(const QModelIndex &child) const
{
    return QModelIndex();
}

int DockModel::rowCount(const QModelIndex &parent) const
{
    return dataArray.count();
}

int DockModel::columnCount(const QModelIndex &parent) const
{
    return 1;
}

QVariant DockModel::data(const QModelIndex &index, int role) const
{
    if (!index.isValid())
        return QVariant();
    if (index.row() >= dataArray.count())
        return QVariant();
    if (role == Qt::DisplayRole || role == Qt::EditRole)
    {
        return QVariant(dataArray.at(index.row()));
    }
    else
        return QVariant();
}

bool DockModel::insertRows(int row, int count, const QModelIndex &parent)
{
    beginInsertRows(QModelIndex(), row, row + count-1);
    for (int i = row; i < row + count; i++)
    {
        dataArray.insert(i,QJsonValue());
    }
    endInsertRows();

    return true;
}
