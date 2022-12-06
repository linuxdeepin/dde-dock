// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QAccessible>
#include <QAccessibleWidget>
#include <QMap>
#include <QMetaEnum>
#include <QObject>
#include <QWidget>

inline QString getAccessibleName(QWidget *w, QAccessible::Role r, const QString &fallback)
{
#define SEPARATOR "_"
    const QString lowerFallback = fallback.toLower();
    // 避免重复生成
    static QMap<QObject *, QString > objnameMap;
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

class Accessible : public QAccessibleWidget {
public:
    Accessible(QWidget *parent, QAccessible::Role r, const QString &accessibleName)
        : QAccessibleWidget(parent, r)
        , w(parent)
        , accessibleName(accessibleName)
    {}

    // 对于使用dogtail的AT自动化测试工作，实际上只需要使用我们提供的text方法获取控件唯一ID,，然后再通过QAccessibleWidget的rect方法找到其坐标，模拟点击即可
    // rect没必要重新实现，text方法通过getAccessibleName确定返回唯一值
    QString text(QAccessible::Text t) const override {
        switch (t) {
        case QAccessible::Name:
            return getAccessibleName(w, this->role(), accessibleName);
        default:
            return QString();
        }
    }

private:
    QWidget *w;
    QString accessibleName;
};

QAccessibleInterface *accessibleFactory(const QString &classname, QObject *object)
{
    Q_UNUSED(classname);

    static QMap<QString, QAccessible::Role> s_roleMap = {
        {"MainWindow",                      QAccessible::Role::Form}
        , {"MainPanelControl",              QAccessible::Role::Button}
        , {"Dock::TipsWidget",              QAccessible::Role::StaticText}
        , {"DockPopupWindow",               QAccessible::Role::Form}
        , {"LauncherItem",                  QAccessible::Role::Button}
        , {"AppItem",                       QAccessible::Role::Button}
        , {"PreviewContainer",              QAccessible::Role::Button}
        , {"PluginsItem",                   QAccessible::Role::Button}
        , {"TrayPluginItem",                QAccessible::Role::Button}
        , {"PlaceholderItem",               QAccessible::Role::Button}
        , {"AppDragWidget",                 QAccessible::Role::Button}
        , {"AppSnapshot",                   QAccessible::Role::Button}
        , {"FloatingPreview",               QAccessible::Role::Button}
        , {"XEmbedTrayWidget",              QAccessible::Role::Button}
        , {"IndicatorTrayWidget",           QAccessible::Role::Button}
        , {"SNITrayWidget",                 QAccessible::Role::Button}
        , {"AbstractTrayWidget",            QAccessible::Role::Button}
        , {"SystemTrayItem",                QAccessible::Role::Button}
        , {"FashionTrayItem",               QAccessible::Role::Form}
        , {"FashionTrayWidgetWrapper",      QAccessible::Role::Form}
        , {"FashionTrayControlWidget",      QAccessible::Role::Button}
        , {"AttentionContainer",            QAccessible::Role::Form}
        , {"HoldContainer",                 QAccessible::Role::Form}
        , {"NormalContainer",               QAccessible::Role::Form}
        , {"SpliterAnimated",               QAccessible::Role::Form}
        , {"DatetimeWidget",                QAccessible::Role::Form}
        , {"OnboardItem",                   QAccessible::Role::Form}
        , {"TrashWidget",                   QAccessible::Role::Form}
        , {"PopupControlWidget",            QAccessible::Role::Button}
        , {"ShutdownWidget",                QAccessible::Role::Form}
        , {"MultitaskingWidget",            QAccessible::Role::Form}
        , {"ShowDesktopWidget",             QAccessible::Role::Form}
        , {"OverlayWarningWidget",          QAccessible::Role::Form}
        , {"QWidget",                       QAccessible::Role::Form}
        , {"QLabel",                        QAccessible::Role::StaticText}
        , {"Dtk::Widget::DIconButton",      QAccessible::Role::Button}
        , {"Dtk::Widget::DSwitchButton",    QAccessible::Role::Button}
        , {"DesktopWidget",                 QAccessible::Role::Button}
        , {"HorizontalSeperator",           QAccessible::Role::Form}
    };

    static QMap<QString, QString> s_classNameMap = {
        {"Dock::TipsWidget",                "tips"}
        , {"DatetimeWidget",                "plugin-datetime"}
        , {"OnboardItem",                   "plugin-onboard"}
        , {"TrashWidget",                   "plugin-trash"}
        , {"ShutdownWidget",                "plugin-shutdown"}
        , {"MultitaskingWidget",            "plugin-multitasking"}
        , {"ShowDesktopWidget",             "plugin-showdesktop"}
        , {"OverlayWarningWidget",          "plugin-overlaywarningwidget"}
        , {"SoundItem",                     "plugin-sounditem"}
    };

    if (object->isWidgetType())
        return new Accessible(qobject_cast<QWidget *>(object)
                                       , s_roleMap.value(classname, QAccessible::Role::Form)
                                       , s_classNameMap.value(object->metaObject()->className(), object->metaObject()->className()));

    return nullptr;
}
