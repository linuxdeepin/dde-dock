#include "informationwidget.h"

#include <QVBoxLayout>
#include <QTimer>

InformationWidget::InformationWidget(QWidget *parent)
    : QWidget(parent)

    , m_infoLabel(new QLabel)
{
    m_infoLabel->setStyleSheet("QLabel {"
                               "color: white;"
                               "}");

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addWidget(m_infoLabel);
    centralLayout->setSpacing(0);
    centralLayout->setMargin(0);

    setLayout(centralLayout);

    QTimer::singleShot(1, this, &InformationWidget::refreshInfo);
}

void InformationWidget::refreshInfo()
{
    // TODO: fetch info
    const int remain = 50;
    const int total = 100;

    // update display
    m_infoLabel->setText(QString("Home:\n%1G/%2G").arg(remain).arg(total));
}
