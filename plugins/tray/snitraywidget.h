/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef SNITRAYWIDGET_H
#define SNITRAYWIDGET_H

#include "constants.h"
#include "abstracttraywidget.h"
#include "util/dockpopupwindow.h"
#include "../../widgets/tipswidget.h"

#include <org_kde_statusnotifieritem.h>
#include <dbusmenu-qt5/dbusmenuimporter.h>

#include <QMenu>
#include <QDBusObjectPath>
#include <QThread>
DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

//using namespace com::deepin::dde;
using namespace org::kde;

enum Command{
    GetHoverTips,
    TraySniAdded
};

class Worker : public QObject
{
    Q_OBJECT

public slots:
    void doWork(Command com,const QString& dbusService,const QString& dbusPath) {
        QProcess p;
        p.start("qdbus", {dbusService});
        if (p.waitForFinished()) {
            switch (com) {
            case GetHoverTips:
            {
                QDBusInterface infc(dbusService, dbusPath);
                QDBusMessage msg = infc.call("Get", "org.kde.StatusNotifierItem", "ToolTip");
                if (msg.type() == QDBusMessage::ReplyMessage) {
                    QDBusArgument arg = msg.arguments().at(0).value<QDBusVariant>().variant().value<QDBusArgument>();
                    DBusToolTip tooltip = qdbus_cast<DBusToolTip>(arg);
                    emit resultReady(tooltip.title);
                    return;
                }
            }
                break;
            case TraySniAdded:
                break;
            default:
                break;
            }
        }
        else {
            qWarning() << "sni dbus service error : " << dbusService;
        }

        emit resultReady("");
    }

signals:
    void resultReady(const QString &);
};

class Controller : public QObject
{
    Q_OBJECT
    QThread workerThread;
public:
    Controller() {
        Worker *worker = new Worker;
        worker->moveToThread(&workerThread);
        connect(&workerThread, &QThread::finished, worker, &QObject::deleteLater);
        connect(this, &Controller::operate, worker, &Worker::doWork);
        connect(worker, &Worker::resultReady, this, &Controller::handleResults);
        workerThread.start();
    }
    ~Controller() {
        workerThread.quit();
        workerThread.wait();
    }
public slots:
    void handleResults(const QString &text)
    {
        if(!text.isEmpty())
            emit resultReady(text);
        this->deleteLater();
    }
signals:
    void operate(Command com,const QString& dbusService,const QString& dbusPath);
    void resultReady(const QString &);
};
class SNITrayWidget : public AbstractTrayWidget
{
    Q_OBJECT

public:
    enum ItemCategory {UnknownCategory = -1, ApplicationStatus, Communications, SystemServices, Hardware};
    enum ItemStatus {Passive, Active, NeedsAttention};
    enum IconType {UnknownIconType = -1, Icon, OverlayIcon, AttentionIcon, AttentionMovieIcon};

public:
    SNITrayWidget(const QString &sniServicePath, QWidget *parent = Q_NULLPTR);
    virtual ~SNITrayWidget();

    QString itemKeyForConfig() override;
    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    void sendClick(uint8_t mouseButton, int x, int y) Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

    bool isValid();
    SNITrayWidget::ItemStatus status();
    SNITrayWidget::ItemCategory category();

    static QString toSNIKey(const QString &sniServicePath);
    static bool isSNIKey(const QString &itemKey);
    static QPair<QString, QString> serviceAndPath(const QString &servicePath);

    void showHoverTips();
    const QPoint topleftPoint() const;
    void showPopupWindow(QWidget * const content, const bool model = false);
    const QPoint popupMarkPoint() const;

    static void setDockPostion(const Dock::Position pos) { DockPosition = pos; }

Q_SIGNALS:
    void statusChanged(SNITrayWidget::ItemStatus status);

private Q_SLOTS:
    void initSNIPropertys();
    void initMenu();
    void refreshIcon();
    void refreshOverlayIcon();
    void refreshAttentionIcon();
    void showContextMenu(int x, int y);
    // SNI property change slot
    void onSNIAttentionIconNameChanged(const QString & value);
    void onSNIAttentionIconPixmapChanged(DBusImageList  value);
    void onSNIAttentionMovieNameChanged(const QString & value);
    void onSNICategoryChanged(const QString & value);
    void onSNIIconNameChanged(const QString & value);
    void onSNIIconPixmapChanged(DBusImageList  value);
    void onSNIIconThemePathChanged(const QString & value);
    void onSNIIdChanged(const QString & value);
    void onSNIMenuChanged(const QDBusObjectPath & value);
    void onSNIOverlayIconNameChanged(const QString & value);
    void onSNIOverlayIconPixmapChanged(DBusImageList  value);
    void onSNIStatusChanged(const QString & value);
    void hidePopup();
    void hideNonModel();
    void popupWindowAccept();
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;

private:
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;
    QPixmap newIconPixmap(IconType iconType);

private:
    StatusNotifierItem *m_sniInter;

    DBusMenuImporter *m_dbusMenuImporter;

    QMenu *m_menu;
    QTimer *m_updateIconTimer;
    QTimer *m_updateOverlayIconTimer;
    QTimer *m_updateAttentionIconTimer;

    QString m_sniServicePath;
    QString m_dbusService;
    QString m_dbusPath;

    QPixmap m_pixmap;
    QPixmap m_overlayPixmap;

    // SNI propertys
    QString m_sniAttentionIconName;
    DBusImageList m_sniAttentionIconPixmap;
    QString m_sniAttentionMovieName;
    QString m_sniCategory;
    QString m_sniIconName;
    DBusImageList m_sniIconPixmap;
    QString m_sniIconThemePath;
    QString m_sniId;
    QDBusObjectPath m_sniMenuPath;
    QString m_sniOverlayIconName;
    DBusImageList m_sniOverlayIconPixmap;
    QString m_sniStatus;
    QTimer *m_popupTipsDelayTimer;
    static Dock::Position DockPosition;
    static QPointer<DockPopupWindow> PopupWindow;
    TipsWidget *m_tipsLabel;
    bool m_popupShown;
};

#endif /* SNIWIDGET_H */
