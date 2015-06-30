/****************************************************************************
** Meta object code from reading C++ file 'screenmask.h'
**
** Created by: The Qt Meta Object Compiler version 67 (Qt 5.4.2)
**
** WARNING! All changes made in this file will be lost!
*****************************************************************************/

#include "../src/Widgets/screenmask.h"
#include <QtCore/qbytearray.h>
#include <QtCore/qmetatype.h>
#if !defined(Q_MOC_OUTPUT_REVISION)
#error "The header file 'screenmask.h' doesn't include <QObject>."
#elif Q_MOC_OUTPUT_REVISION != 67
#error "This file was generated using the moc from 5.4.2. It"
#error "cannot be used with the include files from this version of Qt."
#error "(The moc has changed too much.)"
#endif

QT_BEGIN_MOC_NAMESPACE
struct qt_meta_stringdata_TransformLabel_t {
    QByteArrayData data[2];
    char stringdata[25];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_TransformLabel_t, stringdata) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_TransformLabel_t qt_meta_stringdata_TransformLabel = {
    {
QT_MOC_LITERAL(0, 0, 14), // "TransformLabel"
QT_MOC_LITERAL(1, 15, 9) // "transform"

    },
    "TransformLabel\0transform"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_TransformLabel[] = {

 // content:
       7,       // revision
       0,       // classname
       0,    0, // classinfo
       0,    0, // methods
       1,   14, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       0,       // signalCount

 // properties: name, type, flags
       1, QMetaType::QTransform, 0x00095103,

       0        // eod
};

void TransformLabel::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    Q_UNUSED(_o);
    Q_UNUSED(_id);
    Q_UNUSED(_c);
    Q_UNUSED(_a);
}

const QMetaObject TransformLabel::staticMetaObject = {
    { &QLabel::staticMetaObject, qt_meta_stringdata_TransformLabel.data,
      qt_meta_data_TransformLabel,  qt_static_metacall, Q_NULLPTR, Q_NULLPTR}
};


const QMetaObject *TransformLabel::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *TransformLabel::qt_metacast(const char *_clname)
{
    if (!_clname) return Q_NULLPTR;
    if (!strcmp(_clname, qt_meta_stringdata_TransformLabel.stringdata))
        return static_cast<void*>(const_cast< TransformLabel*>(this));
    return QLabel::qt_metacast(_clname);
}

int TransformLabel::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QLabel::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    
#ifndef QT_NO_PROPERTIES
     if (_c == QMetaObject::ReadProperty) {
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< QTransform*>(_v) = getTransform(); break;
        default: break;
        }
        _id -= 1;
    } else if (_c == QMetaObject::WriteProperty) {
        void *_v = _a[0];
        switch (_id) {
        case 0: setTransform(*reinterpret_cast< QTransform*>(_v)); break;
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
struct qt_meta_stringdata_ScreenMask_t {
    QByteArrayData data[6];
    char stringdata[51];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_ScreenMask_t, stringdata) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_ScreenMask_t qt_meta_stringdata_ScreenMask = {
    {
QT_MOC_LITERAL(0, 0, 10), // "ScreenMask"
QT_MOC_LITERAL(1, 11, 11), // "itemDropped"
QT_MOC_LITERAL(2, 23, 0), // ""
QT_MOC_LITERAL(3, 24, 3), // "pos"
QT_MOC_LITERAL(4, 28, 11), // "itemEntered"
QT_MOC_LITERAL(5, 40, 10) // "itemExited"

    },
    "ScreenMask\0itemDropped\0\0pos\0itemEntered\0"
    "itemExited"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_ScreenMask[] = {

 // content:
       7,       // revision
       0,       // classname
       0,    0, // classinfo
       3,   14, // methods
       0,    0, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       3,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    1,   29,    2, 0x06 /* Public */,
       4,    0,   32,    2, 0x06 /* Public */,
       5,    0,   33,    2, 0x06 /* Public */,

 // signals: parameters
    QMetaType::Void, QMetaType::QPoint,    3,
    QMetaType::Void,
    QMetaType::Void,

       0        // eod
};

void ScreenMask::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        ScreenMask *_t = static_cast<ScreenMask *>(_o);
        switch (_id) {
        case 0: _t->itemDropped((*reinterpret_cast< QPoint(*)>(_a[1]))); break;
        case 1: _t->itemEntered(); break;
        case 2: _t->itemExited(); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        void **func = reinterpret_cast<void **>(_a[1]);
        {
            typedef void (ScreenMask::*_t)(QPoint );
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&ScreenMask::itemDropped)) {
                *result = 0;
            }
        }
        {
            typedef void (ScreenMask::*_t)();
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&ScreenMask::itemEntered)) {
                *result = 1;
            }
        }
        {
            typedef void (ScreenMask::*_t)();
            if (*reinterpret_cast<_t *>(func) == static_cast<_t>(&ScreenMask::itemExited)) {
                *result = 2;
            }
        }
    }
}

const QMetaObject ScreenMask::staticMetaObject = {
    { &QWidget::staticMetaObject, qt_meta_stringdata_ScreenMask.data,
      qt_meta_data_ScreenMask,  qt_static_metacall, Q_NULLPTR, Q_NULLPTR}
};


const QMetaObject *ScreenMask::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *ScreenMask::qt_metacast(const char *_clname)
{
    if (!_clname) return Q_NULLPTR;
    if (!strcmp(_clname, qt_meta_stringdata_ScreenMask.stringdata))
        return static_cast<void*>(const_cast< ScreenMask*>(this));
    return QWidget::qt_metacast(_clname);
}

int ScreenMask::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QWidget::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 3)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 3;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 3)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 3;
    }
    return _id;
}

// SIGNAL 0
void ScreenMask::itemDropped(QPoint _t1)
{
    void *_a[] = { Q_NULLPTR, const_cast<void*>(reinterpret_cast<const void*>(&_t1)) };
    QMetaObject::activate(this, &staticMetaObject, 0, _a);
}

// SIGNAL 1
void ScreenMask::itemEntered()
{
    QMetaObject::activate(this, &staticMetaObject, 1, Q_NULLPTR);
}

// SIGNAL 2
void ScreenMask::itemExited()
{
    QMetaObject::activate(this, &staticMetaObject, 2, Q_NULLPTR);
}
QT_END_MOC_NAMESPACE
