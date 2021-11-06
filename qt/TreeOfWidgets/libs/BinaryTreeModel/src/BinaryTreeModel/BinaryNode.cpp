#include "BinaryNode.h"

#include <QDebug>

#include <math.h>

const BinaryNode::Style BinaryNode::_defaultStyle{
    /*primaryColor*/ 0x006374,
    /*backgroundColor*/ 0x363636,
    /*borderColor*/ 0x7b7b7b,
    /*foregroundColor*/ 0x666666,
    /*secondaryColor*/ 0x4491a1,
    /*accentColor*/ 0x265b91,
    /*titleColor*/ 0xb4b4b4,
    /*textColor*/ 0xeeeeee,
};

BinaryNode::BinaryNode(QString name, double data)
    : _name(name)
    , _data(data)
{
}

BinaryNode::BinaryNode(QString name, QString data)
    : _name(name)
    , _data(data)
{
}

BinaryNodeConstPtr BinaryNode::left() const
{
    return _left;
}

BinaryNodePtr
BinaryNode::left()
{
    return _left;
}

void
BinaryNode::setLeft(BinaryNodePtr item)
{
    _left = item;
}

BinaryNodeConstPtr BinaryNode::right() const
{
    return _right;
}

BinaryNodePtr
BinaryNode::right()
{
    return _right;
}

void
BinaryNode::setRight(BinaryNodePtr item)
{
    _right = item;
}

const QString& BinaryNode::name() const
{
    return _name;
}

QVariant::Type BinaryNode::type() const
{
    return _data.type();
}

double BinaryNode::floatValue() const
{
    assert(type() == QVariant::Type::Double);
    return _data.toDouble();
}

void BinaryNode::setValue(double newVal)
{
    if(!qFuzzyCompare(newVal, _data.toDouble())) {
        _data = QVariant::fromValue(newVal);
        emit valueChanged();
    }
}

QString BinaryNode::stringValue() const
{
    assert(type() == QVariant::Type::String);
    return _data.toString();
}

void BinaryNode::setValue(QString newVal)
{
    if(newVal != _data.toString()) {
        _data = QVariant::fromValue(newVal);
        emit valueChanged();
    }
}

int BinaryNode::height() const
{
    return std::max(_left ? _left->height() + 1 : 0, _right ? _right->height() + 1 : 0);
}

const BinaryNode::Style& BinaryNode::style() const
{
    return _style ? *_style : _defaultStyle;
}

void BinaryNode::setStyle(const Style &style)
{
    _style = std::make_shared<Style>(style);
}
