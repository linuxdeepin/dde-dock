#ifndef DOCKMODEL_H
#define DOCKMODEL_H

#include <QObject>
#include <QVariant>
#include <QJsonObject>
#include <QJsonArray>
#include <QAbstractItemModel>
#include <QDebug>

class DockModel : public QAbstractItemModel
{
    Q_OBJECT
public:
    explicit DockModel(QObject *parent = 0);

    int count();
    void append(const QJsonObject &dataObj);
    void clear();
    QJsonObject get(int index);
    bool insert(int index, const QJsonObject &dataObj);
    void move(int from, int to, int count);
    void remove(int index, int count = 1);
    void set(int index, const QJsonObject &dataObj);
    void setProperty(int index, const QString &property, const QVariant &value);
    int indexOf(const QString &property);
    QModelIndex getIndex(int row);

    bool setData(const QModelIndex &index, const QVariant &value, int role = Qt::EditRole);
    QModelIndex index(int row, int column, const QModelIndex &parent = QModelIndex()) const;
    QModelIndex parent(const QModelIndex &child) const;
    int rowCount(const QModelIndex &parent) const;
    int columnCount(const QModelIndex &parent) const;
    QVariant data(const QModelIndex &index, int role = Qt::EditRole) const;
    bool insertRows(int row, int count, const QModelIndex &parent);
signals:

public slots:

private:
    QJsonArray dataArray;
};

#endif // DOCKMODEL_H
