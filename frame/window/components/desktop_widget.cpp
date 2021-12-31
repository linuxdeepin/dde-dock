#include "desktop_widget.h"

#include <QDBusInterface>
#include <QDebug>
#include <QPainter>
#include <QProcess>

DesktopWidget::DesktopWidget(QWidget *parent)
    : QWidget (parent)
    , m_isHover(false)
    , m_needRecoveryWin(false)
{
    setMouseTracking(true);
}

void DesktopWidget::paintEvent(QPaintEvent *event)
{
    Q_UNUSED(event);

    QPainter painter(this);
    //描绘桌面区域的颜色
    painter.setOpacity(1);
    QPen pen;
    QColor penColor(0, 0, 0, 25);
    pen.setWidth(2);
    pen.setColor(penColor);
    painter.setPen(pen);
    painter.drawRect(rect());
    if (m_isHover) {
        painter.fillRect(rect(), QColor(255, 255, 255, 51));
    } else {
        painter.fillRect(rect(), QColor(255, 255, 255, 25));
    }
}

void DesktopWidget::enterEvent(QEvent *event)
{
    if (checkNeedShowDesktop()) {
        m_needRecoveryWin = true;
        QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");
    }

    m_isHover = true;
    update();

    return QWidget::enterEvent(event);
}

void DesktopWidget::leaveEvent(QEvent *event)
{
    // 鼠标移入时隐藏了窗口，移出时恢复
    if (m_needRecoveryWin) {
        QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");
    }

    m_isHover = false;
    update();

    return QWidget::leaveEvent(event);
}

void DesktopWidget::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton) {
        if (m_needRecoveryWin) {
            // 手动点击 显示桌面窗口 后，鼠标移出时不再调用显/隐窗口进程，以手动点击设置为准
            m_needRecoveryWin = false;
        } else {
            // 需求调整，鼠标移入，预览桌面时再点击显示桌面保持显示桌面状态，再点击才切换桌面显、隐状态
            QProcess::startDetached("/usr/lib/deepin-daemon/desktop-toggle");
        }
    }

    QWidget::mousePressEvent(event);
}

/**
 * @brief ShowDesktopWidget::checkNeedShowDesktop 根据窗管提供接口（当前是否显示的桌面），提供鼠标
 * 移入 显示桌面窗口 区域时，是否需要显示桌面判断依据
 * @return 窗管返回 当前是桌面 或 窗管接口查询失败 返回false，否则true
 */
bool DesktopWidget::checkNeedShowDesktop()
{
    QDBusInterface wmInter("com.deepin.wm", "/com/deepin/wm", "com.deepin.wm");
    QList<QVariant> argumentList;
    QDBusMessage reply = wmInter.callWithArgumentList(QDBus::Block, QStringLiteral("GetIsShowDesktop"), argumentList);
    if (reply.type() == QDBusMessage::ReplyMessage && reply.arguments().count() == 1) {
        return !reply.arguments().at(0).toBool();
    }

    qDebug() << "wm call GetIsShowDesktop fail, res:" << reply.type();
    return false;
}
