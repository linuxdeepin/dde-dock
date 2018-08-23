#include <QPixmap>
#include <QImageReader>
#include <QApplication>

namespace Utils {
    static QPixmap renderSVG(const QString &path, const QSize &size) {
    QImageReader reader;
    QPixmap pixmap;
    reader.setFileName(path);
    if (reader.canRead()) {
        const qreal ratio = qApp->devicePixelRatio();
        reader.setScaledSize(size * ratio);
        pixmap = QPixmap::fromImage(reader.read());
        pixmap.setDevicePixelRatio(ratio);
    }
    else {
        pixmap.load(path);
    }

    return pixmap;
    }
}