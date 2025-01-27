/**
 * This file is generated by dconfig2cpp.
 * Command line arguments: ./dconfig2cpp -p ./dde-dock/toolGenerate/dconfig2cpp ./dde-dock/configs/org.deepin.dde.dock.power.json
 * Generation time: 2025-01-14T10:55:02
 * JSON file version: 1.0
 * 
 * WARNING: DO NOT MODIFY THIS FILE MANUALLY.
 * If you need to change the content, please modify the dconfig2cpp tool.
 */

#ifndef ORG_DEEPIN_DDE_DOCK_POWER_H
#define ORG_DEEPIN_DDE_DOCK_POWER_H

#include <QThread>
#include <QVariant>
#include <QDebug>
#include <QAtomicPointer>
#include <QAtomicInteger>
#include <DConfig>

class org_deepin_dde_dock_power : public QObject {
    Q_OBJECT

    Q_PROPERTY(bool control READ control WRITE setControl NOTIFY controlChanged)
    Q_PROPERTY(bool enable READ enable WRITE setEnable NOTIFY enableChanged)
    Q_PROPERTY(bool menu-enable READ menu-enable WRITE setMenu-enable NOTIFY menu-enableChanged)
    Q_PROPERTY(bool showtimetofull READ showtimetofull WRITE setShowtimetofull NOTIFY showtimetofullChanged)
public:
    explicit org_deepin_dde_dock_power(QThread *thread, const QString &appId, const QString &name, const QString &subpath, QObject *parent = nullptr)
        : QObject(parent) {

        if (!thread->isRunning()) {
            qWarning() << QStringLiteral("Warning: The provided thread is not running.");
        }
        Q_ASSERT(QThread::currentThread() != thread);
        auto worker = new QObject();
        worker->moveToThread(thread);
        QMetaObject::invokeMethod(worker, [=]() {
            auto config = DTK_CORE_NAMESPACE::DConfig::create(appId, name, subpath, nullptr);
            if (!config) {
                qWarning() << QStringLiteral("Failed to create DConfig instance.");
                worker->deleteLater();
                return;
            }
            config->moveToThread(QThread::currentThread());
            initialize(config);
            worker->deleteLater();
        });
    }
    explicit org_deepin_dde_dock_power(QThread *thread, DTK_CORE_NAMESPACE::DConfigBackend *backend, const QString &appId, const QString &name, const QString &subpath, QObject *parent = nullptr)
        : QObject(parent) {

        if (!thread->isRunning()) {
            qWarning() << QStringLiteral("Warning: The provided thread is not running.");
        }
        Q_ASSERT(QThread::currentThread() != thread);
        auto worker = new QObject();
        worker->moveToThread(thread);
        QMetaObject::invokeMethod(worker, [=]() {
            auto config = DTK_CORE_NAMESPACE::DConfig::create(backend, appId, name, subpath, nullptr);
            if (!config) {
                qWarning() << QStringLiteral("Failed to create DConfig instance.");
                worker->deleteLater();
                return;
            }
            config->moveToThread(QThread::currentThread());
            initialize(config);
            worker->deleteLater();
        });
    }
    explicit org_deepin_dde_dock_power(QThread *thread, const QString &name, const QString &subpath, QObject *parent = nullptr)
        : QObject(parent) {

        if (!thread->isRunning()) {
            qWarning() << QStringLiteral("Warning: The provided thread is not running.");
        }
        Q_ASSERT(QThread::currentThread() != thread);
        auto worker = new QObject();
        worker->moveToThread(thread);
        QMetaObject::invokeMethod(worker, [=]() {
            auto config = DTK_CORE_NAMESPACE::DConfig::create(name, subpath, nullptr);
            if (!config) {
                qWarning() << QStringLiteral("Failed to create DConfig instance.");
                worker->deleteLater();
                return;
            }
            config->moveToThread(QThread::currentThread());
            initialize(config);
            worker->deleteLater();
        });
    }
    explicit org_deepin_dde_dock_power(QThread *thread, DTK_CORE_NAMESPACE::DConfigBackend *backend, const QString &name, const QString &subpath, QObject *parent = nullptr)
        : QObject(parent) {

        if (!thread->isRunning()) {
            qWarning() << QStringLiteral("Warning: The provided thread is not running.");
        }
        Q_ASSERT(QThread::currentThread() != thread);
        auto worker = new QObject();
        worker->moveToThread(thread);
        QMetaObject::invokeMethod(worker, [=]() {
            auto config = DTK_CORE_NAMESPACE::DConfig::create(backend, name, subpath, nullptr);
            if (!config) {
                qWarning() << QStringLiteral("Failed to create DConfig instance.");
                worker->deleteLater();
                return;
            }
            config->moveToThread(QThread::currentThread());
            initialize(config);
            worker->deleteLater();
        });
    }
    ~org_deepin_dde_dock_power() {
        if (m_config.loadRelaxed()) {
            m_config.loadRelaxed()->deleteLater();
        }
    }

    bool control() const {
        return p_control;
    }
    void setControl(const bool &value) {
        auto oldValue = p_control;
        p_control = value;
        markPropertySet(0);
        if (auto config = m_config.loadRelaxed()) {
            QMetaObject::invokeMethod(config, [this, value]() {
                m_config.loadRelaxed()->setValue(QStringLiteral("control"), value);
            });
        }
        if (p_control != oldValue) {
            Q_EMIT controlChanged();
        }
    }
    bool enable() const {
        return p_enable;
    }
    void setEnable(const bool &value) {
        auto oldValue = p_enable;
        p_enable = value;
        markPropertySet(1);
        if (auto config = m_config.loadRelaxed()) {
            QMetaObject::invokeMethod(config, [this, value]() {
                m_config.loadRelaxed()->setValue(QStringLiteral("enable"), value);
            });
        }
        if (p_enable != oldValue) {
            Q_EMIT enableChanged();
        }
    }
    bool menu-enable() const {
        return p_menu-enable;
    }
    void setMenu-enable(const bool &value) {
        auto oldValue = p_menu-enable;
        p_menu-enable = value;
        markPropertySet(2);
        if (auto config = m_config.loadRelaxed()) {
            QMetaObject::invokeMethod(config, [this, value]() {
                m_config.loadRelaxed()->setValue(QStringLiteral("menu-enable"), value);
            });
        }
        if (p_menu-enable != oldValue) {
            Q_EMIT menu-enableChanged();
        }
    }
    bool showtimetofull() const {
        return p_showtimetofull;
    }
    void setShowtimetofull(const bool &value) {
        auto oldValue = p_showtimetofull;
        p_showtimetofull = value;
        markPropertySet(3);
        if (auto config = m_config.loadRelaxed()) {
            QMetaObject::invokeMethod(config, [this, value]() {
                m_config.loadRelaxed()->setValue(QStringLiteral("showtimetofull"), value);
            });
        }
        if (p_showtimetofull != oldValue) {
            Q_EMIT showtimetofullChanged();
        }
    }
Q_SIGNALS:
    void controlChanged();
    void enableChanged();
    void menu-enableChanged();
    void showtimetofullChanged();
private:
    void initialize(DTK_CORE_NAMESPACE::DConfig *config) {
        Q_ASSERT(!m_config.loadRelaxed());
        m_config.storeRelaxed(config);
        if (testPropertySet(0)) {
            config->setValue(QStringLiteral("control"), QVariant::fromValue(p_control));
        } else {
            updateValue(QStringLiteral("control"), QVariant::fromValue(p_control));
        }
        if (testPropertySet(1)) {
            config->setValue(QStringLiteral("enable"), QVariant::fromValue(p_enable));
        } else {
            updateValue(QStringLiteral("enable"), QVariant::fromValue(p_enable));
        }
        if (testPropertySet(2)) {
            config->setValue(QStringLiteral("menu-enable"), QVariant::fromValue(p_menu-enable));
        } else {
            updateValue(QStringLiteral("menu-enable"), QVariant::fromValue(p_menu-enable));
        }
        if (testPropertySet(3)) {
            config->setValue(QStringLiteral("showtimetofull"), QVariant::fromValue(p_showtimetofull));
        } else {
            updateValue(QStringLiteral("showtimetofull"), QVariant::fromValue(p_showtimetofull));
        }

        connect(config, &DTK_CORE_NAMESPACE::DConfig::valueChanged, this, [this](const QString &key) {
            updateValue(key);
        }, Qt::DirectConnection);
    }
    void updateValue(const QString &key, const QVariant &fallback = QVariant()) {
        Q_ASSERT(QThread::currentThread() == m_config.loadRelaxed()->thread());
        const QVariant &value = m_config.loadRelaxed()->value(key, fallback);
        if (key == QStringLiteral("control")) {
            auto newValue = qvariant_cast<bool>(value);
            QMetaObject::invokeMethod(this, [this, newValue]() {
                if (p_control != newValue) {
                    p_control = newValue;
                    Q_EMIT controlChanged();
                }
            });
            return;
        }
        if (key == QStringLiteral("enable")) {
            auto newValue = qvariant_cast<bool>(value);
            QMetaObject::invokeMethod(this, [this, newValue]() {
                if (p_enable != newValue) {
                    p_enable = newValue;
                    Q_EMIT enableChanged();
                }
            });
            return;
        }
        if (key == QStringLiteral("menu-enable")) {
            auto newValue = qvariant_cast<bool>(value);
            QMetaObject::invokeMethod(this, [this, newValue]() {
                if (p_menu-enable != newValue) {
                    p_menu-enable = newValue;
                    Q_EMIT menu-enableChanged();
                }
            });
            return;
        }
        if (key == QStringLiteral("showtimetofull")) {
            auto newValue = qvariant_cast<bool>(value);
            QMetaObject::invokeMethod(this, [this, newValue]() {
                if (p_showtimetofull != newValue) {
                    p_showtimetofull = newValue;
                    Q_EMIT showtimetofullChanged();
                }
            });
            return;
        }
    }
    inline void markPropertySet(const int index) {
        if (index < 32) {
            m_propertySetStatus0.fetchAndOrOrdered(1 << (index - 0));
            return;
        }
        Q_UNREACHABLE();
    }
    inline bool testPropertySet(const int index) const {
        if (index < 32) {
            return (m_propertySetStatus0.loadRelaxed() & (1 << (index - 0)));
        }
        Q_UNREACHABLE();
    }
    QAtomicPointer<DTK_CORE_NAMESPACE::DConfig> m_config = nullptr;
    bool p_control { false };
    bool p_enable { true };
    bool p_menu-enable { true };
    bool p_showtimetofull { true };
    QAtomicInteger<quint32> m_propertySetStatus0 = 0;
};

#endif // ORG_DEEPIN_DDE_DOCK_POWER_H
