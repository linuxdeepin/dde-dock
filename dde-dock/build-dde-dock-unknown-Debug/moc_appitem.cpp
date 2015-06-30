/****************************************************************************
** Meta object code from reading C++ file 'appitem.h'
**
** Created by: The Qt Meta Object Compiler version 67 (Qt 5.4.2)
**
** WARNING! All changes made in this file will be lost!
*****************************************************************************/

#include "../src/Widgets/appitem.h"
#include <QtCore/qbytearray.h>
#include <QtCore/qmetatype.h>
#if !defined(Q_MOC_OUTPUT_REVISION)
#error "The header file 'appitem.h' doesn't include <QObject>."
#elif Q_MOC_OUTPUT_REVISION != 67
#error "This file was generated using the moc from 5.4.2. It"
#error "cannot be used with the include files from this version of Qt."
#error "(The moc has changed too much.)"
#endif

QT_BEGIN_MOC_NAMESPACE
struct qt_meta_stringdata_AppItem_t {
    QByteArrayData data[20];
    char stringdata[187];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_AppItem_t, stringdata) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_AppItem_t qt_meta_stringdata_AppItem = {
    {
QT_MOC_LITERAL(0, 0, 7), // "AppItem"
QT_MOC_LITERAL(1, 8, 9), // "dragStart"
QT_MOC_LITERAL(2, 18, 0), // ""
QT_MOC_LITERAL(3, 19, 8), // "AppItem*"
QT_MOC_LITERAL(4, 28, 4), // "item"
QT_MOC_LITERAL(5, 33, 11), // "dragEntered"
QT_MOC_LITERAL(6, 45, 16), // "QDragEnterEvent*"
QT_MOC_LITERAL(7, 62, 5), // "event"
QT_MOC_LITERAL(8, 68, 10), // "dragExited"
QT_MOC_LITERAL(9, 79, 16), // "QDragLeaveEvent*"
QT_MOC_LITERAL(10, 96, 4), // "drop"
QT_MOC_LITERAL(11, 101, 11), // "QDropEvent*"
QT_MOC_LITERAL(12, 113, 12), // "mouseEntered"
QT_MOC_LITERAL(13, 126, 11), // "mouseExited"
QT_MOC_LITERAL(14, 138, 10), // "mousePress"
QT_MOC_LITERAL(15, 149, 1), // "x"
QT_MOC_LITERAL(16, 151, 1), // "y"
QT_MOC_LITERAL(17, 153, 12), // "mouseRelease"
QT_MOC_LITERAL(18, 166, 16), // "mouseDoubleClick"
QT_MOC_LITERAL(19, 183, 3) // "pos"

    },
    "AppItem\0dragStart\0\0AppItem*\0item\0"
    "dragEntered\0QDragEnterEvent*\0event\0"
    "dragExited\0QDragLeaveEvent*\0drop\0"
    "QDropEvent*\0mouseEntered\0mouseExited\0"
    "mousePress\0x\0y\0mouseRelease\0"
    "mouseDoubleClick\0pos"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_AppItem[] = {

 // content:
       7,       // revision
       0,       // classname
       0,    0, // classinfo
       9,   14, // methods
       1,  100, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       9,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    1,   59,    2, 0x06 /* Public */,
       5,    2,   62,    2, 0x06 /* Public */,
       8,    2,   67,    2, 0x06 /* Public */,
      10,    2,   72,    2, 0x06 /* Public */,
      12,    1,   77,    2, 0x06 /* Public */,
      13,    1,   80,    2, 0x06 /* Public */,
      14,    3,   83,    2, 0x06 /* Public */,
      17,    3,   90,    2, 0x06 /* Public */,
      18,    1,   97,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void, 0x80000000 | 3,    4,
    QMetaType::Void, 0x80000000 | 6, 0x80000000 | 3,    7,    4,
    QMetaType::Void, 0x80000000 | 9, 0x80000000 | 3,    7,    4,
    QMetaType::Void, 0x80000000 | 11, 0x80000000 | 3,    7,    4,
    QMetaType::Void, 0x80000000 | 3,    4,
    QMetaType::Void, 0x80000000 | 3,    4,
    QMetaType::Void, QMetaType::Int, QMetaType::Int, 0x80000000 | 3,   15,   16,    4,
    QMetaType::Void, QMetaType::Int, QMetaType::Int, 0x80000000 | 3,   15,   16,    4,
    QMetaType::Void, 0x80000000 | 3,    4,

 // properties: name, type, flags
      19, QMetaType::QPoint, 0x00095003,

       0        // eod
};

void AppItem::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        AppItem *_t = static_cast<AppItem *>(_o);
        switch (_id) {
        case 0: _t->dragStart((*reinterpret_cast< AppItem*(*)>(_a[1]))); break;
        case 1: _t->dragEntered((*reinterpret_cast< QDragEnterEvent*(*)>(_a[1])),(*reinterpret_cast< AppItem*(*)>(_a[2]))); break;
        case 2: _t->dragExited((*reinterpret_cast< QDragLeaveEvent*(*)>(_a[1])),(*reinterpret_cast< AppItem*(*)>(_a[2]))); break;
        case 3: _t->drop((*reinterpret_cast< QDropEvent*(*)>(_a[1])),(*reinterpret_cast< AppItem*(*)>(_a[2]))); break;
        case 4: _t->mouseEntered((*reinterpret_cast< AppItem*(*)>(_a[1]))); break;
        case 5: _t->mouseExited((*reinterpret_cast< AppItem*(*)>(_a[1]))); break;
        case 6: _t->mousePress((*reinterpret_cast< int(*)>(_a[1])),(*reinterpret_cast< int(*)>(_a[2])),(*reinterpret_cast< AppItem*(*)>(_a[3]))); break;
        case 7: _t->mouseRelease((*reinterpret_cast< int(*)>(_a[1])),(*reinterpret_cast< int(*)>(_a[2])),(*reinterpret_cast< AppItem*(*)>(_a[3]))); break;
        case 8: _t->mouseDoubleClick((*reinterpret_cast< AppItem*(*)>(_a[1]))); break;
        default: ;
        }
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        switch (_id) {
        default: *reinterpret_cast<int*>(_a[0]) = -1; break;
        case 0:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 0:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 1:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 1:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 2:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 1:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 3:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 1:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 4:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 0:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 5:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 0:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 6:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 2:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 7:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 2:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        case 8:
            switch (*reinterpret_cast<int*>(_a[1])) {
            default: *reinterpret_cast<int*>(_a[0]) = -1; break;
            case 0:
                *reinterpret_cast<int*>(_a[0]) = qRegisterMetaType< AppItem* >(); break;
            }
            break;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        void **func = reinterpret_cast<void **>(_a[1]);
        {
            typedef void (AppItem::*_t)(AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::dragStart)) {
                *result = 0;
            }
        }
        {
            typedef void (AppItem::*_t)(QDragEnterEvent * , AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::dragEntered)) {
                *result = 1;
            }
        }
        {
            typedef void (AppItem::*_t)(QDragLeaveEvent * , AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::dragExited)) {
                *result = 2;
            }
        }
        {
            typedef void (AppItem::*_t)(QDropEvent * , AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::drop)) {
                *result = 3;
            }
        }
        {
            typedef void (AppItem::*_t)(AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::mouseEntered)) {
                *result = 4;
            }
        }
        {
            typedef void (AppItem::*_t)(AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::mouseExited)) {
                *result = 5;
            }
        }
        {
            typedef void (AppItem::*_t)(int , int , AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::mousePress)) {
                *result = 6;
            }
        }
        {
            typedef void (AppItem::*_t)(int , int , AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::mouseRelease)) {
                *result = 7;
            }
        }
        {
            typedef void (AppItem::*_t)(AppItem * );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&AppItem::mouseDoubleClick)) {
                *result = 8;
            }
        }
    }
}

const QMetaObject AppItem::staticMetaObject = {
    { &DockItem::staticMetaObject, qt_meta_stringdata_AppItem.data,
      qt_meta_data_AppItem,  qt_static_metacall, Q_NULLPTR, Q_NULLPTR}
};


const QMetaObject *AppItem::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *AppItem::qt_metacast(const char *_clname)
{
    if (!_clname) return Q_NULLPTR;
    if (!strcmp(_clname, qt_meta_stringdata_AppItem.stringdata))
        return static_cast<void*>(const_cast< AppItem*>(this));
    return DockItem::qt_metacast(_clname);
}

int AppItem::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = DockItem::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 9)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 9;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 9)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 9;
    }
#ifndef QT_NO_PROPERTIES
      else if (_c == QMetaObject::ReadProperty) {
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< QPoint*>(_v) = pos(); break;
        default: break;
        }
        _id -= 1;
    } else if (_c == QMetaObject::WriteProperty) {
        void *_v = _a[0];
        switch (_id) {
        case 0: move(*reinterpret_cast< QPoint*>(_v)); break;
        default: break;
        }
        _id -= 1;
    } else if (_c == QMetaObject::ResetProperty) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 1;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 1;
    } else if (_c == QMetaObject::RegisterPropertyMetaType) {
        if (_id < 1)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 1;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void AppItem::dragStart(AppItem * _t1)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)) };
    QMetaObject::activate(this, &staticMetaObject, 0, _a);
}

// SIGNAL 1
void AppItem::dragEntered(QDragEnterEvent * _t1, AppItem * _t2)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)), const_cast<void*>(reinterpret_cast<const void*>(&_t2)) };
    QMetaObject::activate(this, &staticMetaObject, 1, _a);
}

// SIGNAL 2
void AppItem::dragExited(QDragLeaveEvent * _t1, AppItem * _t2)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)), const_cast<void*>(reinterpret_cast<const void*>(&_t2)) };
    QMetaObject::activate(this, &staticMetaObject, 2, _a);
}

// SIGNAL 3
void AppItem::drop(QDropEvent * _t1, AppItem * _t2)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)), const_cast<void*>(reinterpret_cast<const void*>(&_t2)) };
    QMetaObject::activate(this, &staticMetaObject, 3, _a);
}

// SIGNAL 4
void AppItem::mouseEntered(AppItem * _t1)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)) };
    QMetaObject::activate(this, &staticMetaObject, 4, _a);
}

// SIGNAL 5
void AppItem::mouseExited(AppItem * _t1)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)) };
    QMetaObject::activate(this, &staticMetaObject, 5, _a);
}

// SIGNAL 6
void AppItem::mousePress(int _t1, int _t2, AppItem * _t3)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)), const_cast<void*>(reinterpret_cast<const void*>(&_t2)), const_cast<void*>(reinterpret_cast<const void*>(&_t3)) };
    QMetaObject::activate(this, &staticMetaObject, 6, _a);
}

// SIGNAL 7
void AppItem::mouseRelease(int _t1, int _t2, AppItem * _t3)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)), const_cast<void*>(reinterpret_cast<const void*>(&_t2)), const_cast<void*>(reinterpret_cast<const void*>(&_t3)) };
    QMetaObject::activate(this, &staticMetaObject, 7, _a);
}

// SIGNAL 8
void AppItem::mouseDoubleClick(AppItem * _t1)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)) };
    QMetaObject::activate(this, &staticMetaObject, 8, _a);
}
QT_END_MOC_NAMESPACE
