#include "pluginwidget.h"

#include <QSvgRenderer>
#include <QPainter>

PluginWidget::PluginWidget(QWidget *parent)
    : QWidget(parent),
      m_powerInter(new DBusPower(this))
{
    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, &PluginWidget::refershIconPixmap);
    connect(m_powerInter, &DBusPower::OnBatteryChanged, this, &PluginWidget::refershIconPixmap);
}

void PluginWidget::displayModeChanged()
{
    refershIconPixmap();
}

QSize PluginWidget::sizeHint() const
{
    return QSize(24, 24);
}

void PluginWidget::resizeEvent(QResizeEvent *e)
{
    QWidget::resizeEvent(e);

    refershIconPixmap();
}

void PluginWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_iconPixmap.rect().center(), m_iconPixmap);
}

const QPixmap PluginWidget::loadSvg(const QString &fileName, const QSize &size) const
{
    QPixmap pixmap(size);
    QSvgRenderer renderer(fileName);
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();

    return pixmap;
}

void PluginWidget::refershIconPixmap()
{
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    if (displayMode == Dock::Efficient)
    {
        m_iconPixmap = loadSvg(":/icons/resources/icons/normal.svg", QSize(24 * 0.8, 24 * 0.8));
        return;
    }

    const int iconSize = std::min(width(), height()) * 0.8;
    const QSize size = QSize(iconSize, iconSize);
    const BatteryPercentageMap percentageData = m_powerInter->batteryPercentage();
    if (percentageData.isEmpty())
    {
        m_iconPixmap = loadSvg(":/icons/resources/icons/fashion.svg", size);
        return;
    }

    const BatteryStateMap stateData = m_powerInter->batteryState();
    if (stateData.isEmpty())
    {
        m_iconPixmap = loadSvg(":/icons/resources/icons/battery_unknow.svg", size);
        return;
    }

    // battery full, charged
    if (stateData.value("Display") == 4)
    {
        m_iconPixmap = loadSvg(":/icons/resources/icons/battery_plugged.svg", size);
        return;
    }

    const bool onBattery = m_powerInter->onBattery();
    const int percent = percentageData.value("Display");
    const int imageNumber = (percent / 10) & ~0x1;
    const QString image = QString(":/icons/resources/icons/battery_%1%2.svg").arg(imageNumber)
                                                                             .arg(onBattery ? "_plugged" : "");

    m_iconPixmap = loadSvg(image, size);

}
