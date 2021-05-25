/*
 * Copyright (C) 2018 ~ 2028 Deepin Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng_cm@deepin.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng_cm@deepin.com>
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

// 为了方便使用,把相关定义独立出来,如有需要,直接包含这个头文件,然后使用SET_*的宏去设置,USE_*宏开启即可
// 注意：对项目中出现的所有的QWidget的派生类都要再启用一次accessiblity，包括qt的原生控件[qt未限制其标记名称为空的情况]
// 注意：使用USE_ACCESSIBLE_BY_OBJECTNAME开启accessiblity的时候，一定要再对这个类用一下USE_ACCESSIBLE，否则标记可能会遗漏

#ifndef ACCESSIBLEINTERFACE_H
#define ACCESSIBLEINTERFACE_H

#include <QAccessible>
#include <QAccessibleWidget>
#include <QEvent>
#include <QMap>
#include <QString>
#include <QWidget>
#include <QObject>
#include <QMetaEnum>
#include <QMouseEvent>
#include <QApplication>

#define SEPARATOR "_"

inline QString getAccessibleName(QWidget *w, QAccessible::Role r, const QString &fallback)
{
    const QString lowerFallback = fallback.toLower();
    // 避免重复生成
    static QMap< QObject *, QString > objnameMap;
    if (!objnameMap[w].isEmpty())
        return objnameMap[w];

    static QMap< QAccessible::Role, QList< QString > > accessibleMap;
    QString oldAccessName = w->accessibleName().toLower();
    oldAccessName.replace(SEPARATOR, "");

    // 按照类型添加固定前缀
    QMetaEnum metaEnum = QMetaEnum::fromType<QAccessible::Role>();
    QByteArray prefix = metaEnum.valueToKeys(r);
    switch (r) {
    case QAccessible::Button:       prefix = "Btn";         break;
    case QAccessible::StaticText:   prefix = "Label";       break;
    default:                        break;
    }

    // 再加上标识
    QString accessibleName = QString::fromLatin1(prefix) + SEPARATOR;
    QString objectName = w->objectName().toLower();
    accessibleName += oldAccessName.isEmpty() ? (objectName.isEmpty() ?lowerFallback : objectName) : oldAccessName;
    // 检查名称是否唯一
    if (accessibleMap[r].contains(accessibleName)) {
        if (!objnameMap.key(accessibleName)) {
            objnameMap.remove(objnameMap.key(accessibleName));
            objnameMap.insert(w, accessibleName);
            return accessibleName;
        }
        // 获取编号，然后+1
        int pos = accessibleName.indexOf(SEPARATOR);
        int id = accessibleName.mid(pos + 1).toInt();

        QString newAccessibleName;
        do {
            // 一直找到一个不重复的名字
            newAccessibleName = accessibleName + SEPARATOR + QString::number(++id);
        } while (accessibleMap[r].contains(newAccessibleName));

        accessibleMap[r].append(newAccessibleName);
        objnameMap.insert(w, newAccessibleName);

        // 对象销毁后移除占用名称
        QObject::connect(w, &QWidget::destroyed, [ = ] (QObject *obj) {
            objnameMap.remove(obj);
            accessibleMap[r].removeOne(newAccessibleName);
         });
        return newAccessibleName;
    } else {
        accessibleMap[r].append(accessibleName);
        objnameMap.insert(w, accessibleName);

        // 对象销毁后移除占用名称
        QObject::connect(w, &QWidget::destroyed, [ = ] (QObject *obj) {
            objnameMap.remove(obj);
            accessibleMap[r].removeOne(accessibleName);
        });
        return accessibleName;
    }
}

// 公共的功能
#define FUNC_CREATE(classname,accessibletype,accessdescription) explicit Accessible##classname(classname *w) \
    : QAccessibleWidget(w,accessibletype,#classname)\
    , m_w(w)\
    , m_description(accessdescription)\
{}\

#define FUNC_TEXT(classname,accessiblename) QString Accessible##classname::text(QAccessible::Text t) const{\
    switch (t) {\
    case QAccessible::Name:\
    return getAccessibleName(m_w, this->role(), accessiblename);\
    case QAccessible::Description:\
    return m_description;\
    default:\
    return QString();\
    }\
    }\

// button控件特有功能
#define FUNC_ACTIONNAMES(classname) QStringList Accessible##classname::actionNames() const{\
    if(!m_w->isEnabled())\
    return QStringList();\
    return QStringList() << pressAction()<< showMenuAction();\
    }\

#define FUNC_DOACTION(classname) void Accessible##classname::doAction(const QString &actionName){\
    if(actionName == pressAction())\
{\
    QPointF localPos = m_w->geometry().center();\
    QMouseEvent event(QEvent::MouseButtonPress,localPos,Qt::LeftButton,Qt::LeftButton,Qt::NoModifier);\
    qApp->sendEvent(m_w,&event);\
    }\
    else if(actionName == showMenuAction())\
{\
    QPointF localPos = m_w->geometry().center();\
    QMouseEvent event(QEvent::MouseButtonPress,localPos,Qt::RightButton,Qt::RightButton,Qt::NoModifier);\
    qApp->sendEvent(m_w,&event);\
    }\
    }\

// Label控件特有功能
#define FUNC_TEXT_(classname) QString Accessible##classname::text(int startOffset, int endOffset) const{\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    return m_w->text();\
    }\

// Slider控件特有功能
#define FUNC_CURRENTVALUE(classname) QVariant Accessible##classname::currentValue() const{\
    return m_w->value();\
    }\

#define FUNC_SETCURRENTVALUE(classname) void Accessible##classname::setCurrentValue(const QVariant &value){\
    return m_w->setValue(value.toInt());\
    }\

#define FUNC_MAXMUMVALUE(classname) QVariant Accessible##classname::maximumValue() const{\
    return QVariant(m_w->maximum());\
    }\

#define FUNC_FUNC_MINIMUMVALUE(classname) QVariant Accessible##classname::minimumValue() const{\
    return QVariant(m_w->minimum());\
    }\

//  DSlider控件特有功能函数
#define FUNC_FUNC_MINIMUMSTEPSIZE(classname) QVariant Accessible##classname::minimumStepSize() const{\
    return QVariant(m_w->pageStep());\
    }\

#define SET_FORM_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,accessdescription)  class Accessible##classname : public QAccessibleWidget\
{\
    public:\
    FUNC_CREATE(classname,QAccessible::Form,accessdescription)\
    QString text(QAccessible::Text t) const override;\
    void *interface_cast(QAccessible::InterfaceType t) override{\
    switch (t) {\
    case QAccessible::ActionInterface:\
    return static_cast<QAccessibleActionInterface*>(this);\
    default:\
    return nullptr;\
    }\
    }\
    private:\
    classname *m_w;\
    QString m_description;\
    };\
    FUNC_TEXT(classname,accessiblename)\

#define SET_BUTTON_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,accessdescription)  class Accessible##classname : public QAccessibleWidget\
{\
    public:\
    FUNC_CREATE(classname,QAccessible::Button,accessdescription)\
    QString text(QAccessible::Text t) const override;\
    void *interface_cast(QAccessible::InterfaceType t) override{\
    switch (t) {\
    case QAccessible::ActionInterface:\
    return static_cast<QAccessibleActionInterface*>(this);\
    default:\
    return nullptr;\
    }\
    }\
    QStringList actionNames() const override;\
    void doAction(const QString &actionName) override;\
    private:\
    classname *m_w;\
    QString m_description;\
    };\
    FUNC_TEXT(classname,accessiblename)\
    FUNC_ACTIONNAMES(classname)\
    FUNC_DOACTION(classname)\

#define SET_LABEL_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,accessdescription)  class Accessible##classname : public QAccessibleWidget, public QAccessibleTextInterface\
{\
    public:\
    FUNC_CREATE(classname,QAccessible::StaticText,accessdescription)\
    QString text(QAccessible::Text t) const override;\
    void *interface_cast(QAccessible::InterfaceType t) override{\
    switch (t) {\
    case QAccessible::ActionInterface:\
    return static_cast<QAccessibleActionInterface*>(this);\
    case QAccessible::TextInterface:\
    return static_cast<QAccessibleTextInterface*>(this);\
    default:\
    return nullptr;\
    }\
    }\
    QString text(int startOffset, int endOffset) const override;\
    void selection(int selectionIndex, int *startOffset, int *endOffset) const override {\
    Q_UNUSED(selectionIndex)\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    }\
    int selectionCount() const override { return 0; }\
    void addSelection(int startOffset, int endOffset) override {\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    }\
    void removeSelection(int selectionIndex) override {\
    Q_UNUSED(selectionIndex)\
    }\
    void setSelection(int selectionIndex, int startOffset, int endOffset) override {\
    Q_UNUSED(selectionIndex)\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    }\
    int cursorPosition() const override { return 0; }\
    void setCursorPosition(int position) override {\
    Q_UNUSED(position)\
    }\
    int characterCount() const override { return 0; }\
    QRect characterRect(int offset) const override {\
    Q_UNUSED(offset)\
    return QRect();\
    }\
    int offsetAtPoint(const QPoint &point) const override {\
    Q_UNUSED(point)\
    return 0;\
    }\
    void scrollToSubstring(int startIndex, int endIndex) override {\
    Q_UNUSED(startIndex)\
    Q_UNUSED(endIndex)\
    }\
    QString attributes(int offset, int *startOffset, int *endOffset) const override {\
    Q_UNUSED(offset)\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    return QString();\
    }\
    private:\
    classname *m_w;\
    QString m_description;\
    };\
    FUNC_TEXT(classname,accessiblename)\
    FUNC_TEXT_(classname)\

#define SET_SLIDER_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,accessdescription)  class Accessible##classname : public QAccessibleWidget, public QAccessibleValueInterface\
{\
    public:\
    FUNC_CREATE(classname,QAccessible::Slider,accessdescription)\
    QString text(QAccessible::Text t) const override;\
    void *interface_cast(QAccessible::InterfaceType t) override{\
    switch (t) {\
    case QAccessible::ActionInterface:\
    return static_cast<QAccessibleActionInterface*>(this);\
    case QAccessible::ValueInterface:\
    return static_cast<QAccessibleValueInterface*>(this);\
    default:\
    return nullptr;\
    }\
    }\
    QVariant currentValue() const override;\
    void setCurrentValue(const QVariant &value) override;\
    QVariant maximumValue() const override;\
    QVariant minimumValue() const override;\
    QVariant minimumStepSize() const override;\
    private:\
    classname *m_w;\
    QString m_description;\
    };\
    FUNC_TEXT(classname,accessiblename)\
    FUNC_CURRENTVALUE(classname)\
    FUNC_SETCURRENTVALUE(classname)\
    FUNC_MAXMUMVALUE(classname)\
    FUNC_FUNC_MINIMUMVALUE(classname)\
    FUNC_FUNC_MINIMUMSTEPSIZE(classname)\

#define SET_EDITABLE_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,accessdescription)  class Accessible##classname : public QAccessibleWidget, public QAccessibleEditableTextInterface, public QAccessibleTextInterface\
{\
    public:\
    FUNC_CREATE(classname,QAccessible::EditableText,accessdescription)\
    QString text(QAccessible::Text t) const override;\
    QAccessibleInterface *child(int index) const override { Q_UNUSED(index); return nullptr; }\
    void *interface_cast(QAccessible::InterfaceType t) override{\
    switch (t) {\
    case QAccessible::ActionInterface:\
    return static_cast<QAccessibleActionInterface*>(this);\
    case QAccessible::TextInterface:\
    return static_cast<QAccessibleTextInterface*>(this);\
    case QAccessible::EditableTextInterface:\
    return static_cast<QAccessibleEditableTextInterface*>(this);\
    default:\
    return nullptr;\
    }\
    }\
    QString text(int startOffset, int endOffset) const override;\
    void selection(int selectionIndex, int *startOffset, int *endOffset) const override {\
    Q_UNUSED(selectionIndex)\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    }\
    int selectionCount() const override { return 0; }\
    void addSelection(int startOffset, int endOffset) override {\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    }\
    void removeSelection(int selectionIndex) override { Q_UNUSED(selectionIndex);}\
    void setSelection(int selectionIndex, int startOffset, int endOffset) override {\
    Q_UNUSED(selectionIndex)\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    }\
    int cursorPosition() const override { return 0; }\
    void setCursorPosition(int position) override {\
    Q_UNUSED(position)\
    }\
    int characterCount() const override { return 0; }\
    QRect characterRect(int offset) const override { \
    Q_UNUSED(offset)\
    return QRect(); }\
    int offsetAtPoint(const QPoint &point) const override {\
    Q_UNUSED(point)\
    return 0; }\
    void scrollToSubstring(int startIndex, int endIndex) override {\
    Q_UNUSED(startIndex)\
    Q_UNUSED(endIndex)\
    }\
    QString attributes(int offset, int *startOffset, int *endOffset) const override {\
    Q_UNUSED(offset)\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    return QString(); }\
    void insertText(int offset, const QString &text) override {\
    Q_UNUSED(offset)\
    Q_UNUSED(text)\
    }\
    void deleteText(int startOffset, int endOffset) override {\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    }\
    void replaceText(int startOffset, int endOffset, const QString &text) override {\
    Q_UNUSED(startOffset)\
    Q_UNUSED(endOffset)\
    Q_UNUSED(text)\
    }\
    private:\
    classname *m_w;\
    QString m_description;\
    };\
    FUNC_TEXT(classname,accessiblename)\
    FUNC_TEXT_(classname)\

#define USE_ACCESSIBLE(classnamestring,classname)    if (classnamestring == QLatin1String(#classname) && object && object->isWidgetType())\
{\
    interface = new Accessible##classname(static_cast<classname *>(object));\
    }\

#define ELSE_USE_ACCESSIBLE(classnamestring,classname)    else if (classnamestring == QLatin1String(#classname) && object && object->isWidgetType())\
{\
    interface = new Accessible##classname(static_cast<classname *>(object));\
    }\


// [指定objectname]---适用同一个类，但objectname不同的情况
#define USE_ACCESSIBLE_BY_OBJECTNAME(classnamestring,classname,objectname)    if (classnamestring == QLatin1String(#classname) && object && (object->objectName() == objectname) && object->isWidgetType())\
{\
    interface = new Accessible##classname(static_cast<classname *>(object));\
    }\

#define ELSE_USE_ACCESSIBLE_BY_OBJECTNAME(classnamestring,classname,objectname)    else if (classnamestring == QLatin1String(#classname) && object && (object->objectName() == objectname) && object->isWidgetType())\
{\
    interface = new Accessible##classname(static_cast<classname *>(object));\
    }\

/*******************************************简化使用*******************************************/
#define SET_FORM_ACCESSIBLE(classname,accessiblename)                          SET_FORM_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,"")

#define SET_BUTTON_ACCESSIBLE(classname,accessiblename)                        SET_BUTTON_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,"")

#define SET_LABEL_ACCESSIBLE(classname,accessiblename)                         SET_LABEL_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,"")

#define SET_SLIDER_ACCESSIBLE(classname,accessiblename)                        SET_SLIDER_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,"")

#define SET_EDITABLE_ACCESSIBLE(classname,accessiblename)                      SET_EDITABLE_ACCESSIBLE_WITH_DESCRIPTION(classname,accessiblename,"")
/************************************************************************************************/

#endif // ACCESSIBLEINTERFACE_H
