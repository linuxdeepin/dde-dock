#ifndef _PREVIEWCONTAINER_H
#define _PREVIEWCONTAINER_H

#include <QWidget>
#include <QBoxLayout>

#include "dbus/dbusdockentry.h"
#include "constants.h"

#include <DWindowManagerHelper>

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
    QMap<WId, QWidget *> m_windows;

    DWindowManagerHelper *m_wmHelper;

    QBoxLayout *m_windowListLayout;

};

#endif // _PREVIEWCONTAINER_H
