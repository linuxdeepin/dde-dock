/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     lzm <lizhongming@uniontech.com>
 *
 * Maintainer: lzm <lizhongming@uniontech.com>
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

#ifndef APPEDITDIALOG_H
#define APPEDITDIALOG_H

#include <DLabel>
#include <DDialog>
#include <DLineEdit>

DWIDGET_USE_NAMESPACE

class IconWidget;
class AppEditDialog : public DDialog
{
    Q_OBJECT

    enum ErrorType {
        NoAppNameError = 0,
        NoIconError,

        AppNameError,
        FileSizeError,
        FileTypeError,
        IconSizeError,
    };

public:
    AppEditDialog(const QString& appName, const QString& iconName, QWidget* parent = Q_NULLPTR);
    ~AppEditDialog();

Q_SIGNALS:
    void updateAppInfo(const QString& appName, const QString& iconPath);

private Q_SLOTS:
    void onButtonClicked(int index);
    void onIconClicked();

private:
    void initUi();
    void initConnections();
    void updateErrorPrompt(ErrorType errorType);
    void changeAppInfo();
    ErrorType checkAppIcon(const QString& iconPath);
    ErrorType checkAppName(const QString& appName);

private:
    QString m_appName;
    QString m_iconName;
    QString m_newIconPath;
    DLabel* m_errorLabel = Q_NULLPTR;
    DLineEdit* m_appNameEdit = Q_NULLPTR;
    IconWidget* m_iconEditWidget = Q_NULLPTR;
};

class IconWidget : public QWidget
{
    Q_OBJECT

    enum Status {
        Normal = 0,
        Hover,
        Pressed
    };

public:
    IconWidget(QWidget* parent = Q_NULLPTR);
    ~IconWidget();

Q_SIGNALS:
    void iconClicked();

public Q_SLOTS:
    void updateIcon(const QString& iconName);

protected:
    void mousePressEvent(QMouseEvent *event) Q_DECL_OVERRIDE;
    void enterEvent(QEvent *event) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent *event) Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *event) Q_DECL_OVERRIDE;

private:
    const QPixmap& currentBtnPixmap();

private:
    QPixmap m_appIcon;
    QPixmap m_btnNormal;
    QPixmap m_btnHover;
    QPixmap m_btnPressed;
    Status m_status = Normal;
};

#endif // APPEDITDIALOG_H
