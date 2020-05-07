#ifndef TIPSWIDGET_H
#define TIPSWIDGET_H

#include <QFrame>

class TipsWidget : public QFrame
{
    Q_OBJECT
    enum ShowType
    {
        SingleLine,
        MultiLine
    };
public:
    explicit TipsWidget(QWidget *parent = nullptr);

    const QString& text(){return m_text;}
    const QStringList &textList() { return  m_textList; }
    void setText(const QString &text);
    void setTextList(const QStringList &textList);
    
protected:
    void paintEvent(QPaintEvent *event) override;

private:
    QString m_text;
    QStringList m_textList;
    int m_width;
    ShowType m_type;
};

#endif // TIPSWIDGET_H
