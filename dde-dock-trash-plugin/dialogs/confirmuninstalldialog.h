/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef CONFIRMUNINSTALLDIALOG_H
#define CONFIRMUNINSTALLDIALOG_H

#include "dbasedialog.h"

class ConfirmUninstallDialog : public DBaseDialog
{
    Q_OBJECT
public:
    explicit ConfirmUninstallDialog(QWidget *parent = 0);
    ~ConfirmUninstallDialog();

signals:

public slots:
    void handleKeyEnter();
};

#endif // CONFIRMUNINSTALLDIALOG_H
