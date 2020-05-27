#include "appeditdialog.h"
#include "../util/themeappicon.h"

#include <QDebug>
#include <QHBoxLayout>
#include <DLabel>
#include <DButtonBox>
#include <DWidget>
#include <QPaintEvent>
#include <QStandardPaths>
#include <QFileDialog>
#include <DFileDialog>

#define UnsupportedFileFormat tr("Unsupported file format")
#define UnsupportedDimensions tr("Unsupported dimensions")
#define TheFileIsTooLarge tr("The file is too large")

constexpr int MaxAppIconByteSize = 1024 * 1024 * 1;
static const int DialogWidth = 380;
static const int DIalogHeight = 390;

AppEditDialog::AppEditDialog(const QString& appName, const QString& iconName,  QWidget *parent)
    : DDialog(parent)
    , m_appName(appName)
    , m_iconName(iconName)
{
    initUi();
    initConnections();
}

AppEditDialog::~AppEditDialog()
{

}

void AppEditDialog::initUi()
{
    setFixedSize(DialogWidth, DIalogHeight);
    setWindowFlags(windowFlags() | Qt::WindowStaysOnTopHint);

    addContent(new DLabel(tr("Edit"), this), Qt::AlignHCenter);

    addSpacing(15);

    m_iconEditWidget = new IconWidget(this);
    m_iconEditWidget->setFixedSize(108, 108);
    m_iconEditWidget->updateIcon(m_iconName);
    addContent(m_iconEditWidget, Qt::AlignHCenter);

    addSpacing(20);

    addContent(new DLabel(tr("Change icon"), this), Qt::AlignHCenter);

    addSpacing(6);

    addContent(new DLabel(tr("SVG only; dimensions: 96*96; size: ≤1 MB"), this), Qt::AlignHCenter);

    addSpacing(10);

    DWidget* edit = new DWidget(this);
    QVBoxLayout* editLayout = new QVBoxLayout(edit);

    QHBoxLayout* lineEditLayout = new QHBoxLayout(edit);
    lineEditLayout->addWidget(new DLabel(tr("App name")));
    m_appNameEdit = new DLineEdit(edit);
    m_appNameEdit->setText(m_appName);
    m_appNameEdit->setFixedHeight(36)   ;
    lineEditLayout->addWidget(m_appNameEdit);
    editLayout->addLayout(lineEditLayout);

    lineEditLayout->addSpacing(10);

    m_errorLabel = new DLabel(this);
    m_errorLabel->setFixedHeight(18);
    m_errorLabel->setForegroundRole(DPalette::TextWarning);
    m_errorLabel->setAlignment(Qt::AlignHCenter);
    editLayout->addWidget(m_errorLabel, Qt::AlignHCenter);

    addContent(edit);

    addButton(tr("Cancel"));
    addButton(tr("Confirm"), true, ButtonType::ButtonRecommend);
    setOnButtonClickedClose(false);
}

void AppEditDialog::initConnections()
{
    connect(this, &AppEditDialog::buttonClicked, this, &AppEditDialog::onButtonClicked);
    connect(m_iconEditWidget, &IconWidget::iconClicked, this, &AppEditDialog::onIconClicked);
    connect(m_appNameEdit, &DLineEdit::textChanged, this, [this](const QString &text){
        m_appName = text;
        updateErrorPrompt(checkAppName(text));
    });
}

void AppEditDialog::updateErrorPrompt(ErrorType errorType)
{
    switch (errorType) {
        case NoAppNameError: {
            m_appNameEdit->setAlert(false);
        }
        break;
        case NoIconError: {
            m_errorLabel->clear();
        }
        break;
        case FileTypeError: {
            m_errorLabel->setText(UnsupportedFileFormat);
        }
        break;
        case FileSizeError: {
            m_errorLabel->setText(TheFileIsTooLarge);
        }
        break;
        case IconSizeError: {
            m_errorLabel->setText(tr("SVG only; dimensions: 96*96; size: ≤1 MB"));
        }
        break;
        case AppNameError: {
            m_appNameEdit->setAlert(true);
            m_appNameEdit->showAlertMessage(tr("\\/:*?\"<>| are not allowed"));
        }
        break;
        default: return;
    }
}

void AppEditDialog::changeAppInfo()
{
    ErrorType err = checkAppName(m_appName);
    if (err != NoAppNameError) {
        updateErrorPrompt(err);
        return;
    }

    //if ()

    emit updateAppInfo(m_appName, m_newIconPath);

    close();
}

AppEditDialog::ErrorType AppEditDialog::checkAppIcon(const QString &iconPath)
{
    qDebug() << __FUNCTION__ << "iconPath: " << iconPath;

    QFileInfo iconFile(iconPath);
    if (!iconFile.exists() || iconFile.size() > MaxAppIconByteSize) {
        return FileSizeError;
    }

    QPixmap temp(iconPath);
    if (temp.isNull()) {
        return FileTypeError;
    }

    if (temp.size() != QSize(96, 96)) {
        return IconSizeError;
    }

    return NoIconError;
}

AppEditDialog::ErrorType AppEditDialog::checkAppName(const QString &appName)
{
    if (appName.isEmpty() || appName.contains(QRegularExpression(".*[\\/:*?\"<>|].*"))) {
        return AppNameError;
    }

    return NoAppNameError;
}

void AppEditDialog::onButtonClicked(int index)
{
    switch (index) {
        case 0: close();         break;
        case 1: changeAppInfo(); break;
        default:                 break;
    }
}

void AppEditDialog::onIconClicked()
{
    QString openDir;
    QStringList directorys = QStandardPaths::standardLocations(QStandardPaths::PicturesLocation);
    if (!directorys.isEmpty()) {
        openDir = directorys.first();
    }

    m_newIconPath = QFileDialog::getOpenFileName(this, "", openDir, tr("Images") + "(*.svg)");
    if (m_newIconPath.isEmpty()) {
        return;
    }

    ErrorType err = checkAppIcon(m_newIconPath);
    if (NoIconError != err) {
        m_newIconPath.clear();
        updateErrorPrompt(err);
        return;
    }

    m_iconEditWidget->updateIcon(m_newIconPath);
}

IconWidget::IconWidget(QWidget *parent)
    : QWidget(parent)
    , m_btnNormal(":/window/resources/edit_btn_normal.svg")
    , m_btnHover(":/window/resources/edit_btn_hover.svg")
    , m_btnPressed(":/window/resources/edit_btn_press.svg")
{
    setWindowFlag(Qt::FramelessWindowHint);
}

IconWidget::~IconWidget()
{

}

void IconWidget::updateIcon(const QString &iconName)
{
    qDebug() << __FUNCTION__ << "iconName: " << iconName;
    m_appIcon = ThemeAppIcon::getIcon(iconName, width(), 1.0);
    update();
}

void IconWidget::mousePressEvent(QMouseEvent *event)
{
    Q_UNUSED(event);

    m_status = Pressed;
    update();

    emit iconClicked();
}

void IconWidget::enterEvent(QEvent *event)
{
    Q_UNUSED(event);

    m_status = Hover;
    update();
}

void IconWidget::leaveEvent(QEvent *event)
{
    Q_UNUSED(event);

    m_status = Normal;
    update();
}

void IconWidget::paintEvent(QPaintEvent *event)
{
    QPainter painter(this);
    painter.drawPixmap(event->rect(), m_appIcon);
    painter.drawPixmap(event->rect(), currentBtnPixmap());
}

const QPixmap &IconWidget::currentBtnPixmap()
{
    switch (m_status) {
        case Normal:  return m_btnNormal;
        case Hover:   return m_btnHover;
        case Pressed: return m_btnPressed;
        default:      return m_btnNormal;
    }
}
