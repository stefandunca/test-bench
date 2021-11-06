#pragma once

#include <QObject>
#include <QVariant>
#include <QColor>

#include <memory>

class BinaryNode;

using BinaryNodePtr = std::shared_ptr<BinaryNode>;
using BinaryNodeConstPtr = std::shared_ptr<const BinaryNode>;

class BinaryNode: public QObject
{
    Q_OBJECT
public:
    explicit BinaryNode(QString name, double data);
    explicit BinaryNode(QString name, QString data);

    BinaryNodeConstPtr left() const;
    BinaryNodePtr left();
    void setLeft(BinaryNodePtr item);
    BinaryNodeConstPtr right() const;
    BinaryNodePtr right();
    void setRight(BinaryNodePtr item);

    const QString& name() const;

    QVariant::Type type() const;

    double floatValue() const;
    void setValue(double newVal);
    QString stringValue() const;
    void setValue(QString newVal);

    int height() const;

    struct Style {
        QColor primaryColor;
        QColor backgroundColor;
        QColor borderColor;
        QColor foregroundColor;
        QColor secondaryColor;
        QColor accentColor;
        QColor titleColor;
        QColor textColor;
        bool fixedSize = true;
    };

    const Style& style() const;
    void setStyle(const Style& style);

signals:
    void valueChanged();

private:
    QString _name;
    QVariant _data;

    std::shared_ptr<Style> _style;
    static const Style _defaultStyle;

    BinaryNodePtr _left;
    BinaryNodePtr _right;
};
