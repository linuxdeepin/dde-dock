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
