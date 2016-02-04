/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "cleartrashdialog.h"
#include <QVBoxLayout>
#include <QLabel>
#include <QButtonGroup>
#include <QPushButton>

ClearTrashDialog::ClearTrashDialog(QWidget *parent):
    DBaseDialog(parent)
{

    QString icon = ":/images/skin/dialogs/images/user-trash-full.png";
    QString message = tr("Are you sure to empty trash?");
    QString tipMessage = tr("This action cannot be restored");
    QStringList buttons, buttonTexts;
    buttons << "Cancel" << "Empty";
    buttonTexts << tr("Cancel") << tr("Empty");
    initUI(icon, message, tipMessage, buttons, buttons);
    moveCenter();
    getButtonsGroup()->button(1)->setFocus();
    setButtonTexts(buttonTexts);
}

void ClearTrashDialog::handleKeyEnter(){
    handleButtonsClicked(1);
}

ClearTrashDialog::~ClearTrashDialog()
{

}

