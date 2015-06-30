/****************************************************************************
** Meta object code from reading C++ file 'appicon.h'
**
** Created by: The Qt Meta Object Compiler version 67 (Qt 5.4.2)
**
** WARNING! All changes made in this file will be lost!
*****************************************************************************/

#include "../src/Widgets/appicon.h"
#include <QtCore/qbytearray.h>
#include <QtCore/qmetatype.h>
#if !defined(Q_MOC_OUTPUT_REVISION)
#error "The header file 'appicon.h' doesn't include <QObject>."
#elif Q_MOC_OUTPUT_REVISION != 67
#error "This file was generated using the moc from 5.4.2. It"
#error "cannot be used with the include files from this version of Qt."
#error "(The moc has changed too much.)"
#endif

QT_BEGIN_MOC_NAMESPACE
struct qt_meta_stringdata_AppIcon_t {
    QByteArrayData data[5];
    char stringdata[34];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_AppIcon_t, stringdata) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_AppIcon_t qt_meta_stringdata_AppIcon = {
    {
QT_MOC_LITERAL(0, 0, 7), // "AppIcon"
QT_MOC_LITERAL(1, 8, 10), // "getSysIcon"
QT_MOC_LITERAL(2, 19, 0), // ""
QT_MOC_LITERAL(3, 20, 8), // "iconName"
QT_MOC_LITERAL(4, 29, 4) // "size"

    },
    "AppIcon\0getSysIcon\0\0iconName\0size"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_AppIcon[] = {

 // content:
       7,       // revision
       0,       // classname
       0,    0, // classinfo
       2,   14, // methods
       0,    0, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       0,       // signalCount

 // slots: name, argc, parameters, tag, flags
       1,    2,   24,    2, 0x0a /* Public */,
       1,    1,   29,    2, 0x2a /* Public | MethodCloned */,

 // slots: parameters
    QMetaType::QString, QMetaType::QString, QMetaType::Int,    3,    4,
    QMetaType::QString, QMetaType::QString,    3,

       0        // eod
};

void AppIcon::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        AppIcon *_t = static_cast<AppIcon *>(_o);
        switch (_id) {
        case 0: { QString _r = _t->getSysIcon((*reinterpret_cast< const QString(*)>(_a[1])),(*reinterpret_cast< int(*)>(_a[2])));
            if (_a[0]) *reinterpret_cast< QString*>(_a[0]) = _r; }  break;
        case 1: { QString _r = _t->getSysIcon((*reinterpret_cast< const QString(*)>(_a[1])));
            if (_a[0]) *reinterpret_cast< QString*>(_a[0]) = _r; }  break;
        default: ;
        }
    }
}

const QMetaObject AppIcon::staticMetaObject = {
    { &QLabel::staticMetaObject, qt_meta_stringdata_AppIcon.data,
      qt_meta_data_AppIcon,  qt_static_metacall, Q_NULLPTR, Q_NULLPTR}
};


const QMetaObject *AppIcon::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *AppIcon::qt_metacast(const char *_clname)
{
    if (!_clname) return Q_NULLPTR;
    if (!strcmp(_clname, qt_meta_stringdata_AppIcon.stringdata))
        return static_cast<void*>(const_cast< AppIcon*>(this));
    return QLabel::qt_metacast(_clname);
}

int AppIcon::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QLabel::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 2)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 2;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 2)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 2;
    }
    return _id;
}
QT_END_MOC_NAMESPACE
