#ifndef APPPREVIEWS_H
#define APPPREVIEWS_H

#include <QWidget>
#include <QHBoxLayout>
#include <QLabel>
#include <QDebug>
#include "windowpreview.h"
#include "../dockconstants.h"

class AppPreviews : public QWidget
{
    Q_OBJECT
public:
    explicit AppPreviews(QWidget *parent = 0);

    void addItem(const QString &title,int xid);
    void setTitle(const QString &title);
signals:

public slots:

private:
    QHBoxLayout *m_mainLayout = NULL;
    QList<int> m_xidList;
};

#endif // APPPREVIEWS_H
