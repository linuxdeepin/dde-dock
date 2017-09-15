#include "pluginwidget.h"

#include <QSvgRenderer>
#include <QPainter>
#include <QMouseEvent>
#include <QApplication>

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
    return QSize(26, 26);
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
            pixmap = loadSvg(":/icons/resources/icons/normal.svg", QSize(16, 16));
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
        const uint percentage = qMin(100.0, qMax(0.0, percentageData.value("Display")));
        const int percent = std::round(percentage / 10.0) * 10;
        const int imageNumber = (percent / 10) & ~0x1;
        const QString image = QString(":/icons/resources/icons/battery_%1%2.svg").arg(imageNumber)
                                                                                 .arg(m_hover || onBattery ? "" : "_plugged");

        pixmap = loadSvg(image, size);
    } while (false);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - pixmap.rect().center() / qApp->devicePixelRatio(), pixmap);
}

void PluginWidget::mousePressEvent(QMouseEvent *e)
{
    if (e->button() != Qt::RightButton)
        return QWidget::mousePressEvent(e);

    const QPoint p(e->pos() - rect().center());
    if (p.manhattanLength() < std::min(width(), height()) * 0.8 * 0.5)
    {
        emit requestContextMenu(SHUTDOWN_KEY);
        return;
    }

    return QWidget::mousePressEvent(e);
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
    const auto ratio = qApp->devicePixelRatio();

    QPixmap pixmap(size * ratio);
    QSvgRenderer renderer(fileName);
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();

    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}
