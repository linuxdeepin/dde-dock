#ifndef _PREVIEWCONTAINER_H
#define _PREVIEWCONTAINER_H

#include <QWidget>
#include <QBoxLayout>
#include <QTimer>

#include "dbus/dbusdockentry.h"
#include "constants.h"
#include "appsnapshot.h"

#include <DWindowManagerHelper>

#define SNAP_WIDTH       200
#define SNAP_HEIGHT      130

DWIDGET_USE_NAMESPACE

class _PreviewContainer : public QWidget
{
    Q_OBJECT

public:
    explicit _PreviewContainer(QWidget *parent = 0);

signals:
    void requestActivateWindow(const WId wid) const;
    void requestPreviewWindow(const WId wid) const;
    void requestCancelPreview() const;
    void requestHidePreview() const;

public:
    void setWindowInfos(const WindowDict &infos);

public slots:
    void updateLayoutDirection(const Dock::Position dockPos);

private:
    void adjustSize();
    void appendSnapWidget(const WId wid);

private:
    QMap<WId, AppSnapshot *> m_snapshots;

    QBoxLayout *m_windowListLayout;

    DWindowManagerHelper *m_wmHelper;
};

#endif // _PREVIEWCONTAINER_H
