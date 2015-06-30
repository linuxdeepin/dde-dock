/****************************************************************************
** Meta object code from reading C++ file 'docklayout.h'
**
** Created by: The Qt Meta Object Compiler version 67 (Qt 5.4.2)
**
** WARNING! All changes made in this file will be lost!
*****************************************************************************/

#include "../src/Widgets/docklayout.h"
#include <QtCore/qbytearray.h>
#include <QtCore/qmetatype.h>
#if !defined(Q_MOC_OUTPUT_REVISION)
#error "The header file 'docklayout.h' doesn't include <QObject>."
#elif Q_MOC_OUTPUT_REVISION != 67
#error "This file was generated using the moc from 5.4.2. It"
#error "cannot be used with the include files from this version of Qt."
#error "(The moc has changed too much.)"
#endif

QT_BEGIN_MOC_NAMESPACE
struct qt_meta_stringdata_DockLayout_t {
    QByteArrayData data[15];
    char stringdata[154];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_DockLayout_t, stringdata) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_DockLayout_t qt_meta_stringdata_DockLayout = {
    {
QT_MOC_LITERAL(0, 0, 10), // "DockLayout"
QT_MOC_LITERAL(1, 11, 11), // "dragStarted"
QT_MOC_LITERAL(2, 23, 0), // ""
QT_MOC_LITERAL(3, 24, 11), // "itemDropped"
QT_MOC_LITERAL(4, 36, 12), // "slotItemDrag"
QT_MOC_LITERAL(5, 49, 8), // "AppItem*"
QT_MOC_LITERAL(6, 58, 4), // "item"
QT_MOC_LITERAL(7, 63, 15), // "slotItemRelease"
QT_MOC_LITERAL(8, 79, 1), // "x"
QT_MOC_LITERAL(9, 81, 1), // "y"
QT_MOC_LITERAL(10, 83, 15), // "slotItemEntered"
QT_MOC_LITERAL(11, 99, 16), // "QDragEnterEvent*"
QT_MOC_LITERAL(12, 116, 5), // "event"
QT_MOC_LITERAL(13, 122, 14), // "slotItemExited"
QT_MOC_LITERAL(14, 137, 16) // "QDragLeaveEvent*"

    },
    "DockLayout\0dragStarted\0\0itemDropped\0"
    "slotItemDrag\0AppItem*\0item\0slotItemRelease\0"
    "x\0y\0slotItemEntered\0QDragEnterEvent*\0"
    "event\0slotItemExited\0QDragLeaveEvent*"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_DockLayout[] = {

 // content:
       7,       // revision
       0,       // classname
       0,    0, // classinfo
       6,   14, // methods
       0,    0, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       2,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    0,   44,    2, 0x06 /* Public */,
       3,    0,   45,    2, 0x06 /* Public */,

 // slots: name, argc, parameters, tag, flags
       4,    1,   46,    2, 0x08 /* Private */,
       7,    3,   49,    2, 0x08 /* Private */,
      10,    2,   56,    2, 0x08 /* Private */,
      13,    2,   61,    2, 0x08 /* Private */,

 // signals: parameters
    QMetaType::Void,
    QMetaType::Void,

 // slots: parameters
    QMetaType::Void, 0x80000000 | 5,    6,
    QMetaType::Void, QMetaType::Int, QMetaType::Int, 0x80000000 | 5,    8,    9,    6,
    QMetaType::Void, 0x80000000 | 11, 0x80000000 | 5,   12,    6,
    QMetaType::Void, 0x80000000 | 14, 0x80000000 | 5,   12,    6,

       0        // eod
};

void DockLayout::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        DockLayout *_t = static_cast<DockLayout *>(_o);
        switch (_id) {
        case 0: _t->dragStarted(); break;
        case 1: _t->itemDropped(); break;
        case 2: _t->slotItemDrag((*reinterpret_cast< AppItem*(*)>(_a[1]))); break;
        case 3: _t->slotItemRelease((*reinterpret_cast< int(*)>(_a[1])),(*reinterpret_cast< int(*)>(_a[2])),(*reinterpret_cast< AppItem*(*)>(_a[3]))); break;
        case 4: _t->slotItemEntered((*reinterpret_cast< QDragEnterEvent*(*)>(_a[1])),(*reinterpret_cast< AppItem*(*)>(_a[2]))); break;
        case 5: _t->slotItemExited((*reinterpret_cast< QDragLeaveEvent*(*)>(_a[1])),(*reinterpret_cast< AppItem*(*)>(_a[2]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        switch (_id) {
        default: *reinterpret_cast<int*>(_a[0]) = -1; break;
        case 2:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 0:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 3:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 2:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 4:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 1:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 5:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 1:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        void **func = reinterpret_cast<void **>(_a[1]);
        {
            typedef void (DockLayout::*_t)();
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&DockLayout::dragStarted)) {
                *result = 0;
            }
        }
        {
            typedef void (DockLayout::*_t)();
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&DockLayout::itemDropped)) {
                *result = 1;
            }
        }
    }
}

const QMetaObject DockLayout::staticMetaObject = {
    { &QWidget::staticMetaObject, qt_meta_stringdata_DockLayout.data,
      qt_meta_data_DockLayout,  qt_static_metacall, Q_NULLPTR, Q_NULLPTR}
};


const QMetaObject *DockLayout::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *DockLayout::qt_metacast(const char *_clname)
{
    if (!_clname) return Q_NULLPTR;
    if (!strcmp(_clname, qt_meta_stringdata_DockLayout.stringdata))
        return static_cast<void*>(const_cast< DockLayout*>(this));
    return QWidget::qt_metacast(_clname);
}

int DockLayout::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QWidget::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 6)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 6;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 6)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 6;
    }
    return _id;
}

// SIGNAL 0
void DockLayout::dragStarted()
{
    QMetaObject::activate(this, &staticMetaObject, 0, Q_NULLPTR);
}

// SIGNAL 1
void DockLayout::itemDropped()
{
    QMetaObject::activate(this, &staticMetaObject, 1, Q_NULLPTR);
}
QT_END_MOC_NAMESPACE
