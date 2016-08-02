#include "pluginwidget.h"

#include <QSvgRenderer>
#include <QPainter>

PluginWidget::PluginWidget(QWidget *parent)
    : QWidget(parent),
      m_hover(false),
      m_powerInter(new DBusPower(this))
{
    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, static_cast<void (PluginWidget::*)()>(&PluginWidget::update));
    connect(m_powerInter, &DBusPower::BatteryStateChanged, this, static_cast<void (PluginWidget::*)()>(&PluginWidget::update));
    connect(m_powerInter, &DBusPower::OnBatteryChanged, this, static_cast<void (PluginWidget::*)()>(&PluginWidget::update));
}

QSize PluginWidget::sizeHint() const
{
    return QSize(24, 24);
}

void PluginWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPixmap pixmap;
    do
    {
        const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

        if (displayMode == Dock::Efficient)
        {
            pixmap = loadSvg(":/icons/resources/icons/normal.svg", QSize(18, 18));
            break;
        }

        const int iconSize = std::min(width(), height()) * 0.8;
        const QSize size = QSize(iconSize, iconSize);
        const BatteryPercentageMap percentageData = m_powerInter->batteryPercentage();
        if (percentageData.isEmpty())
        {
            pixmap = loadSvg(":/icons/resources/icons/fashion.svg", size);
            break;
        }

        const BatteryStateMap stateData = m_powerInter->batteryState();
        if (stateData.isEmpty())
        {
            pixmap = loadSvg(":/icons/resources/icons/battery_unknow.svg", size);
            break;
        }

        // battery full, charged
        if (stateData.value("Display") == 4)
        {
            if (!m_hover)
                pixmap = loadSvg(":/icons/resources/icons/battery_plugged.svg", size);
            else
                pixmap = loadSvg(":/icons/resources/icons/battery_10.svg", size);
            break;
        }

        const bool onBattery = m_powerInter->onBattery();
        const int percent = percentageData.value("Display");
        const int imageNumber = (percent / 10) & ~0x1;
        const QString image = QString(":/icons/resources/icons/battery_%1%2.svg").arg(imageNumber)
                                                                                 .arg(m_hover || onBattery ? "" : "_plugged");

        pixmap = loadSvg(image, size);
    } while (false);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - pixmap.rect().center(), pixmap);
}

void PluginWidget::enterEvent(QEvent *)
{
    m_hover = true;
}

void PluginWidget::leaveEvent(QEvent *)
{
    m_hover = false;
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
