#include "quicksettingcontainer.h"
#include "quicksettingcontroller.h"
#include "pluginsiteminterface.h"
#include "quicksettingitem.h"
#include "mediawidget.h"
#include "dockpopupwindow.h"
#include "brightnesswidget.h"
#include "volumewidget.h"
#include "volumedeviceswidget.h"
#include "brightnessmonitorwidget.h"

#include <DListView>
#include <DStyle>
#include <QDrag>

#include <QVBoxLayout>
#include <QMetaObject>
#include <QStackedLayout>

DWIDGET_USE_NAMESPACE

static const int QuickItemRole = Dtk::UserRole + 10;

#define ITEMWIDTH 70
#define ITEMHEIGHT 60
#define CTRLHEIGHT 56
#define ITEMSPACE 10
#define COLUMNCOUNT 4

QuickSettingContainer::QuickSettingContainer(QWidget *parent)
    : QWidget(parent)
    , m_switchLayout(new QStackedLayout(this))
    , m_mainWidget(new QWidget(this))
    , m_pluginWidget(new QWidget(m_mainWidget))
    , m_mainlayout(new QVBoxLayout(m_mainWidget))
    , m_pluginLoader(QuickSettingController::instance())
    , m_playerWidget(new MediaWidget(m_mainWidget))
    , m_volumnWidget(new VolumeWidget(m_mainWidget))
    , m_brihtnessWidget(new BrightnessWidget(m_mainWidget))
    , m_volumeSettingWidget(new VolumeDevicesWidget(m_volumnWidget->model(), this))
    , m_brightSettingWidget(new BrightnessMonitorWidget(m_brihtnessWidget->model(), this))
    , m_childPage(new PluginChildPage(this))
{
    initUi();
    initConnection();
}

QuickSettingContainer::~QuickSettingContainer()
{
}

void QuickSettingContainer::showHomePage()
{
    m_switchLayout->setCurrentIndex(0);
}

DockPopupWindow *QuickSettingContainer::popWindow()
{
    static DockPopupWindow *popView = nullptr;
    if (!popView) {
        popView = new DockPopupWindow;
        popView->setWindowFlag(Qt::Popup);
        popView->setShadowBlurRadius(20);
        popView->setRadius(18);
        popView->setShadowYOffset(2);
        popView->setShadowXOffset(0);
        popView->setArrowWidth(18);
        popView->setArrowHeight(10);

        QuickSettingContainer *container = new QuickSettingContainer(popView);
        popView->setContent(container);
    } else {
        QuickSettingContainer *container = static_cast<QuickSettingContainer *>(popView->getContent());
        container->showHomePage();
    }

    return popView;
}

void QuickSettingContainer::initQuickItem(QuickSettingItem *quickItem)
{
    int pluginWidth = ITEMWIDTH;
    if (quickItem->pluginItem()->isPrimary())
        pluginWidth = ITEMWIDTH * 2 + ITEMSPACE;

    quickItem->setFixedSize(pluginWidth, ITEMHEIGHT);
    quickItem->setParent(m_pluginWidget);
    quickItem->setMouseTracking(true);
    quickItem->installEventFilter(this);
    connect(quickItem, &QuickSettingItem::detailClicked, this, &QuickSettingContainer::onItemDetailClick);
}

void QuickSettingContainer::onItemDetailClick(PluginsItemInterface *pluginInter)
{
    QuickSettingItem *quickItemWidget = static_cast<QuickSettingItem *>(sender());
    if (!quickItemWidget)
        return;

    QWidget *widget = pluginInter->itemWidget(quickItemWidget->itemKey());
    if (!widget)
        return;

    showWidget(widget, pluginInter->pluginDisplayName());
}

bool QuickSettingContainer::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_childPage && event->type() == QEvent::Resize)
        resizeView();

    return QWidget::eventFilter(watched, event);
}

void QuickSettingContainer::showWidget(QWidget *widget, const QString &title)
{
    m_childPage->setTitle(title);
    m_childPage->pushWidget(widget);
    m_switchLayout->setCurrentWidget(m_childPage);
}

void QuickSettingContainer::onPluginInsert(QuickSettingItem *quickItem)
{
    initQuickItem(quickItem);
    resetItemPosition();
    resizeView();
}

void QuickSettingContainer::onPluginRemove(QuickSettingItem *quickItem)
{
    QObjectList childrens = m_pluginWidget->children();
    for (QObject *child : childrens) {
        if (child != quickItem)
            continue;

        disconnect(quickItem, &QuickSettingItem::detailClicked, this, &QuickSettingContainer::onItemDetailClick);
        quickItem->setParent(nullptr);
        quickItem->removeEventFilter(this);
        quickItem->setMouseTracking(false);
        quickItem->hide();
        break;
    }
    //调整子控件的位置
    resetItemPosition();
    resizeView();
}

void QuickSettingContainer::mousePressEvent(QMouseEvent *event)
{
    if (event->button() != Qt::LeftButton)
        return;

    QuickSettingItem *moveItem = qobject_cast<QuickSettingItem *>(childAt(event->pos()));
    if (!moveItem || moveItem->pluginItem()->isPrimary())
        return QWidget::mousePressEvent(event);

    QDrag *drag = new QDrag(this);
    CustomMimeData *mimedata = new CustomMimeData;
    mimedata->setData(moveItem);
    drag->setMimeData(mimedata);
    QPixmap dragPixmap = moveItem->dragPixmap();
    drag->setPixmap(dragPixmap);
    drag->setHotSpot(QPoint(dragPixmap.width() / 2, dragPixmap.height() / 2));

    drag->exec(Qt::MoveAction | Qt::CopyAction);
}

void QuickSettingContainer::resetItemPosition()
{
    QObjectList childrens = m_pluginWidget->children();
    QList<QuickSettingItem *> primaryQuickItems;
    QList<QuickSettingItem *> quickItems;
    for (QObject *child : childrens) {
        QuickSettingItem *quickItem = qobject_cast<QuickSettingItem *>(child);
        if (!quickItem || quickItem->isHidden())
            continue;

        if (quickItem->pluginItem()->isPrimary())
            primaryQuickItems << quickItem;
        else
            quickItems << quickItem;
    }
    static QStringList existKeys = {"network-item-key", "sound-item-key", "VPN", "PROJECTSCREEN"};
    qSort(primaryQuickItems.begin(), primaryQuickItems.end(), [ = ](QuickSettingItem *item1, QuickSettingItem *item2) {
        int index1 = existKeys.indexOf(item1->itemKey());
        int index2 = existKeys.indexOf(item2->itemKey());
        if (index1 >= 0 || index2 >= 0) {
            return index1 < index2;
        }
        return true;
    });
    int primaryColumnCount = COLUMNCOUNT / 2;
    for (int i = 0; i < primaryQuickItems.size(); i++) {
        QuickSettingItem *quickitem = primaryQuickItems[i];
        QPoint ptItem(ITEMSPACE + (ITEMWIDTH + ITEMSPACE) * 2 * (i % primaryColumnCount),
                      (ITEMHEIGHT + ITEMSPACE) * (static_cast<int>(i / primaryColumnCount)));
        quickitem->move(ptItem);
    }
    int startCount = primaryQuickItems.size() * 2;
    for (int i = 0; i < quickItems.size(); i++) {
        QuickSettingItem *qsi = quickItems[i];
        int columnIndex = (startCount + i) % COLUMNCOUNT;
        int rowIndex = (startCount + i) / COLUMNCOUNT;
        int x = (ITEMWIDTH + ITEMSPACE) * columnIndex + ITEMSPACE;
        int y = (ITEMHEIGHT + ITEMSPACE) * rowIndex;
        qsi->move(x, y);
    }
}

void QuickSettingContainer::initUi()
{
    auto setWidgetStyle = [](DBlurEffectWidget *widget) {
        widget->setMaskColor(QColor(239, 240, 245));
        widget->setBlurRectXRadius(8);
        widget->setBlurRectYRadius(8);
    };

    // 添加音乐播放插件
    m_playerWidget->setFixedHeight(CTRLHEIGHT);
    m_volumnWidget->setFixedHeight(CTRLHEIGHT);
    m_brihtnessWidget->setFixedHeight(CTRLHEIGHT);

    setWidgetStyle(m_playerWidget);
    setWidgetStyle(m_volumnWidget);
    setWidgetStyle(m_brihtnessWidget);

    m_mainlayout->setSpacing(0);
    m_mainlayout->setContentsMargins(0, ITEMSPACE, 0, ITEMSPACE);

    m_mainlayout->addWidget(m_pluginWidget);

    QWidget *ctrlWidget = new QWidget(m_mainWidget);
    QVBoxLayout *ctrlLayout = new QVBoxLayout(ctrlWidget);
    ctrlLayout->setContentsMargins(ITEMSPACE, ITEMSPACE, ITEMSPACE, ITEMSPACE);

    ctrlLayout->addSpacing(ITEMSPACE);
    ctrlLayout->addWidget(m_playerWidget);
    ctrlLayout->setSpacing(ITEMSPACE);
    ctrlLayout->addWidget(m_volumnWidget);
    ctrlLayout->setSpacing(ITEMSPACE);
    ctrlLayout->addWidget(m_brihtnessWidget);

    m_mainlayout->addWidget(ctrlWidget);
    // 加载所有的插件
    QList<QuickSettingItem *> pluginItems = m_pluginLoader->settingItems();
    for (QuickSettingItem *quickItem: pluginItems)
        initQuickItem(quickItem);

    m_switchLayout->addWidget(m_mainWidget);
    m_switchLayout->addWidget(m_childPage);

    m_volumeSettingWidget->hide();
    m_brightSettingWidget->hide();

    setMouseTracking(true);
    setAcceptDrops(true);

    QMetaObject::invokeMethod(this, [ = ] {
        if (pluginItems.size() > 0)
            resetItemPosition();
        // 设置当前窗口的大小
        resizeView();
        setFixedWidth(ITEMWIDTH * 4 + (ITEMSPACE * 5));
    }, Qt::QueuedConnection);
}

void QuickSettingContainer::initConnection()
{
    connect(m_pluginLoader, &QuickSettingController::pluginInsert, this, &QuickSettingContainer::onPluginInsert);
    connect(m_pluginLoader, &QuickSettingController::pluginRemove, this, &QuickSettingContainer::onPluginRemove);
    connect(m_playerWidget, &MediaWidget::visibleChanged, this, [ this ] { resizeView(); });
    connect(m_volumnWidget, &VolumeWidget::visibleChanged, this, [ this ] { resizeView(); });
    connect(m_volumnWidget, &VolumeWidget::rightIconClick, this, [ this ] {
        showWidget(m_volumeSettingWidget, tr("voice"));
        resizeView();
    });
    connect(m_brihtnessWidget, &BrightnessWidget::visibleChanged, this, [ this ] { resizeView(); });
    connect(m_brihtnessWidget, &BrightnessWidget::rightIconClicked, this, [ this ] {
        showWidget(m_brightSettingWidget, tr("brightness"));
        resizeView();
    });
    connect(m_childPage, &PluginChildPage::back, this, [ this ] {
        m_switchLayout->setCurrentWidget(m_mainWidget);
    });
    connect(m_childPage, &PluginChildPage::closeSelf, this, [ this ] {
        if (!m_childPage->isBack())
            topLevelWidget()->hide();
    });
    m_childPage->installEventFilter(this);
}

void QuickSettingContainer::resizeView()
{
    if (m_switchLayout->currentWidget() == m_mainWidget) {
        QList<QuickSettingItem *> pluginItems = m_pluginLoader->settingItems();
        int selfPluginCount = 0;
        for (QuickSettingItem *item : pluginItems) {
            // 如果是置顶的插件，则认为它占用两个普通插件的位置
            int increCount = (item->pluginItem()->isPrimary() ? 2 : 1);
            selfPluginCount += increCount;
        }
        int rowCount = selfPluginCount / COLUMNCOUNT;
        if (selfPluginCount % COLUMNCOUNT > 0)
            rowCount++;

        m_pluginWidget->setFixedHeight(ITEMHEIGHT * rowCount + ITEMSPACE * (rowCount - 1));

        int panelCount = 0;
        if (m_playerWidget->isVisible())
            panelCount++;
        if (m_volumnWidget->isVisible())
            panelCount++;
        if (m_brihtnessWidget->isVisible())
            panelCount++;

        int topHeight = ((ITEMHEIGHT + ITEMSPACE) * rowCount) + ITEMSPACE;
        setFixedHeight(topHeight + (CTRLHEIGHT + ITEMSPACE) * panelCount + ITEMSPACE);
    } else if (m_switchLayout->currentWidget() == m_childPage) {
        setFixedHeight(m_childPage->height());
    }
}

/**
 * @brief PluginChildPage::PluginChildPage
 * @param parent
 */
PluginChildPage::PluginChildPage(QWidget *parent)
    : QWidget(parent)
    , m_headerWidget(new QWidget(this))
    , m_back(new QLabel(m_headerWidget))
    , m_title(new QLabel(m_headerWidget))
    , m_container(new QWidget(this))
    , m_topWidget(nullptr)
    , m_containerLayout(new QVBoxLayout(m_container))
    , m_isBack(false)
{
    initUi();
    m_back->installEventFilter(this);
}

PluginChildPage::~PluginChildPage()
{
}

void PluginChildPage::pushWidget(QWidget *widget)
{
    // 首先将界面其他的窗体移除
    for (int i = m_containerLayout->count() - 1; i >= 0; i--) {
        QLayoutItem *item = m_containerLayout->itemAt(i);
        item->widget()->removeEventFilter(this);
        item->widget()->hide();
        m_containerLayout->removeItem(item);
    }
    m_topWidget = widget;
    widget->installEventFilter(this);
    m_containerLayout->addWidget(widget);
    widget->show();
    m_isBack = false;
    QMetaObject::invokeMethod(this, &PluginChildPage::resetHeight, Qt::QueuedConnection);
}

void PluginChildPage::setTitle(const QString &text)
{
    m_title->setText(text);
}

bool PluginChildPage::isBack()
{
    return m_isBack;
}

void PluginChildPage::initUi()
{
    m_back->setFixedWidth(24);
    m_title->setAlignment(Qt::AlignCenter);
    QHBoxLayout *headerLayout = new QHBoxLayout(m_headerWidget);
    headerLayout->setContentsMargins(11, 12, 24 + 11, 12);
    headerLayout->setSpacing(0);
    headerLayout->addWidget(m_back);
    headerLayout->addWidget(m_title);
    m_headerWidget->setFixedHeight(48);

    QVBoxLayout *mainLayout = new QVBoxLayout(this);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->setSpacing(0);

    mainLayout->addWidget(m_headerWidget);
    mainLayout->addWidget(m_container);
    m_containerLayout->setContentsMargins(11, 0, 11, 0);
    m_containerLayout->setSpacing(0);
}

bool PluginChildPage::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_back && event->type() == QEvent::MouseButtonRelease) {
        m_isBack = true;
        Q_EMIT back();
        return true;
    }
    if (watched == m_topWidget) {
        if (event->type() == QEvent::Hide) {
            Q_EMIT closeSelf();
            return true;
        }
        if (event->type() == QEvent::Resize)
            resetHeight();
    }
    return QWidget::eventFilter(watched, event);
}

void PluginChildPage::resetHeight()
{
    QMargins m = m_containerLayout->contentsMargins();
    m_container->setFixedHeight(m.top() + m.bottom() + m_topWidget->height());
    setFixedHeight(m_headerWidget->height() + m_container->height());
}
