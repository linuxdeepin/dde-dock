#ifndef QUICKSETTINGCONTAINER_H
#define QUICKSETTINGCONTAINER_H

#include "pluginproxyinterface.h"

#include "dtkwidget_global.h"

#include <DListView>

#include <QWidget>

class DockItem;
class QVBoxLayout;
class QuickSettingController;
class MediaWidget;
class VolumeWidget;
class BrightnessWidget;
class QuickSettingItem;
class DockPopupWindow;
class QStackedLayout;
class VolumeDevicesWidget;
class BrightnessMonitorWidget;
class QLabel;
class PluginChildPage;

DWIDGET_USE_NAMESPACE

class QuickSettingContainer : public QWidget
{
    Q_OBJECT

public:
    static DockPopupWindow *popWindow();

protected:
    void mousePressEvent(QMouseEvent *event) override;
    explicit QuickSettingContainer(QWidget *parent = nullptr);
    ~QuickSettingContainer() override;
    void showHomePage();

private Q_SLOTS:
    void onPluginInsert(QuickSettingItem *quickItem);
    void onPluginRemove(QuickSettingItem *quickItem);
    void onItemDetailClick(PluginsItemInterface *pluginInter);
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    // 加载UI
    void initUi();
    // 初始化槽函数
    void initConnection();
    // 调整尺寸
    void resizeView();
    // 调整控件位置
    void resetItemPosition();
    // 初始化控件项目
    void initQuickItem(QuickSettingItem *quickItem);
    // 显示具体的窗体
    void showWidget(QWidget *widget, const QString &title);

private:
    QStackedLayout *m_switchLayout;
    QWidget *m_mainWidget;
    QWidget *m_pluginWidget;
    QVBoxLayout *m_mainlayout;
    QuickSettingController *m_pluginLoader;
    MediaWidget *m_playerWidget;
    VolumeWidget *m_volumnWidget;
    BrightnessWidget *m_brihtnessWidget;

    VolumeDevicesWidget *m_volumeSettingWidget;
    BrightnessMonitorWidget *m_brightSettingWidget;
    PluginChildPage *m_childPage;
};

class CustomMimeData : public QMimeData
{
    Q_OBJECT

public:
    CustomMimeData() : QMimeData(), m_data(nullptr) {}
    ~CustomMimeData() {}
    void setData(void *data) { m_data = data; }
    void *data() { return m_data; }

private:
     void *m_data;
};

class PluginChildPage : public QWidget
{
    Q_OBJECT

Q_SIGNALS:
    void back();
    void closeSelf();

public:
    explicit PluginChildPage(QWidget *parent);
    ~PluginChildPage() override;
    void pushWidget(QWidget *widget);
    void setTitle(const QString &text);
    bool isBack();

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    void initUi();
    void resetHeight();

private:
    QWidget *m_headerWidget;
    QLabel *m_back;
    QLabel *m_title;
    QWidget *m_container;
    QWidget *m_topWidget;
    QVBoxLayout *m_containerLayout;
    bool m_isBack;
};

#endif // PLUGINCONTAINER_H
