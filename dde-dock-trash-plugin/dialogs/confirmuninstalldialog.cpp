/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "confirmuninstalldialog.h"
#include <QVBoxLayout>
#include <QLabel>
#include <QButtonGroup>
#include <QPushButton>

ConfirmUninstallDialog::ConfirmUninstallDialog(QWidget *parent) : DBaseDialog(parent)
{
    QString icon = ":/images/skin/dialogs/images/user-trash-full.png";
    QString message = "Are you sure to uninstall this application?";
    QString tipMessage = tr("All dependencies will be removed together");
    QStringList buttons, buttonTexts;
    buttons << "Cancel" << "Confirm";
    buttonTexts << tr("Cancel") << tr("Confirm");
    initUI(icon, message, tipMessage, buttons, buttons);
    moveCenter();
    getButtonsGroup()->button(1)->setFocus();
    setButtonTexts(buttonTexts);
}

void ConfirmUninstallDialog::handleKeyEnter(){
    handleButtonsClicked(1);
}

ConfirmUninstallDialog::~ConfirmUninstallDialog()
{

}
