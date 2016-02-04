/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef CLEARTRASHDIALOG_H
#define CLEARTRASHDIALOG_H

#include "dbasedialog.h"

class ClearTrashDialog : public DBaseDialog
{
    Q_OBJECT
public:
    explicit ClearTrashDialog(QWidget *parent = 0);
    ~ClearTrashDialog();

signals:

public slots:
    void handleKeyEnter();
};

#endif // CLEARTRASHDIALOG_H
