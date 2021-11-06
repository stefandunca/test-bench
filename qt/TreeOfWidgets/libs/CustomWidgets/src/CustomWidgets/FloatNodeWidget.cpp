#include "FloatNodeWidget.h"

#include <BinaryTreeModel/BinaryNode.h>

#include <QSlider>
#include <QLabel>
#include <QGridLayout>

FloatNodeWidget::FloatNodeWidget(const std::weak_ptr<BinaryNode>& node, QWidget *parent)
    : BaseNodeWidget(node, parent)
    , _valueLabel(new QLabel(this))
    , _valueSlider(new QSlider(Qt::Horizontal, this))
{
    _valueSlider->setMinimum(0);
    _valueSlider->setMaximum(_maxValue * _decimalCount * 10);
    connect(_valueSlider, &QAbstractSlider::valueChanged, this, &onValueChanged);
    auto nodePtr = _node.lock();
    if(nodePtr) {
        connect(nodePtr.get(), &BinaryNode::valueChanged, this, &onDataChanged);
        onDataChanged();
        _valueLabel->setStyleSheet("QLabel { color: " + nodePtr->style().textColor.name() + ";}");
        _valueSlider->setStyleSheet("QSlider::groove::horizontal { color: " + nodePtr->style().primaryColor.name() + "}"
                                    "QSlider::handle::horizontal { color: " + nodePtr->style().secondaryColor.name() + "}");
    }
}

void FloatNodeWidget::addContent(QGridLayout* contentLayout)
{
    contentLayout->addWidget(_valueLabel, 0, 0, 1, 1, Qt::AlignBottom | Qt::AlignHCenter);
    contentLayout->addWidget(_valueSlider, 1, 0, 1, 1, Qt::AlignBottom);
}

void FloatNodeWidget::onValueChanged(int value)
{
    double newVal = static_cast<double>(value)/(_decimalCount * 10);
    _valueLabel->setText(QString::number(newVal, 'f', _decimalCount));
    auto nodePtr = _node.lock();
    if(nodePtr) {
        nodePtr->setValue(newVal);
    }
}

void FloatNodeWidget::onDataChanged()
{
    auto nodePtr = _node.lock();
    if(nodePtr) {
        int newVal = static_cast<int>(nodePtr->floatValue()) * (_decimalCount * 10);
        if(newVal != _valueSlider->value())
            _valueSlider->setValue(newVal);
    }
}
