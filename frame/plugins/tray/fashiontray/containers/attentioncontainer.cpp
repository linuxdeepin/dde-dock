#include "attentioncontainer.h"

AttentionContainer::AttentionContainer(TrayPlugin *trayPlugin, QWidget *parent)
    : AbstractContainer(trayPlugin, parent)
{

}

FashionTrayWidgetWrapper *AttentionContainer::takeAttentionWrapper()
{
    if (isEmpty()) {
        return nullptr;
    }

    return takeWrapper(wrapperList().first());
}

bool AttentionContainer::acceptWrapper(FashionTrayWidgetWrapper *wrapper)
{
    return true;
}

void AttentionContainer::refreshVisible()
{
    setVisible(!isEmpty());
}

void AttentionContainer::addWrapper(FashionTrayWidgetWrapper *wrapper)
{
    if (!isEmpty()) {
        qDebug() << "Reject! Already contains a attention wrapper!";
        return;
    }

    AbstractContainer::addWrapper(wrapper);
}
