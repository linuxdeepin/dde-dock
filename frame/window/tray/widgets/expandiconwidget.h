#ifndef EXPANDICONWIDGET_H
#define EXPANDICONWIDGET_H

#include "constants.h"
#include "basetraywidget.h"

class TrayGridView;
class TrayModel;

namespace Dtk { namespace Gui { class DRegionMonitor; } }

class ExpandIconWidget : public BaseTrayWidget
{
    Q_OBJECT

Q_SIGNALS:
    void trayVisbleChanged(bool);

public:
    explicit ExpandIconWidget(QWidget *parent = Q_NULLPTR, Qt::WindowFlags f = Qt::WindowFlags());
    ~ExpandIconWidget() override;
    void setPositonValue(Dock::Position position);

    void sendClick(uint8_t mouseButton, int x, int y) override;
    void setTrayPanelVisible(bool visible);
    QString itemKeyForConfig() override { return "Expand"; }
    void updateIcon() override {}
    QPixmap icon() override;
    TrayGridView *popupTrayView();

private Q_SLOTS:
    void onGlobMousePress(const QPoint &mousePos, const int flag);

protected:
    void paintEvent(QPaintEvent *e) override;
    const QString dropIconFile() const;

    void resetPosition();

private:
    Dtk::Gui::DRegionMonitor *m_regionInter;
    Dock::Position m_position;
    TrayGridView *m_trayView;
};

#endif // EXPANDICONWIDGET_H
