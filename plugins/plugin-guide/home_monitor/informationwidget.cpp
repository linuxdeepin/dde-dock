#include "informationwidget.h"

#include <QVBoxLayout>
#include <QTimer>
#include <QDebug>

InformationWidget::InformationWidget(QWidget *parent)
    : QWidget(parent)
    , m_infoLabel(new QLabel)
    , m_refreshTimer(new QTimer(this))
    // 使用 "/home" 初始化 QStorageInfo
    // 如果 "/home" 没有挂载到一个单独的分区上，QStorageInfo 收集的数据将会是根分区的
    , m_storageInfo(new QStorageInfo("/home"))
{
    m_infoLabel->setStyleSheet("QLabel {"
                               "color: white;"
                               "}");
    m_infoLabel->setAlignment(Qt::AlignCenter);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addWidget(m_infoLabel);
    centralLayout->setSpacing(0);
    centralLayout->setMargin(0);

    setLayout(centralLayout);

    // 连接 Timer 超时的信号到更新数据的槽上
    connect(m_refreshTimer, &QTimer::timeout, this, &InformationWidget::refreshInfo);

    // 设置 Timer 超时为 10s，即每 10s 更新一次控件上的数据，并启动这个定时器
    m_refreshTimer->start(10000);

    refreshInfo();
}

void InformationWidget::refreshInfo()
{
    // 获取分区总容量
    const double total = m_storageInfo->bytesTotal();
    // 获取可用总容量
    const double available = m_storageInfo->bytesAvailable();
    // 得到可用百分比
    const int percent = qRound(available / total * 100);

    // 更新内容
    m_infoLabel->setText(QString("Home:\n%1\%").arg(percent));
}
