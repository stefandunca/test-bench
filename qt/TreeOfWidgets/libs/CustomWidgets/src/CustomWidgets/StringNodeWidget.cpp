#include "StringNodeWidget.h"

#include <BinaryTreeModel/BinaryNode.h>

#include <QLineEdit>
#include <QVBoxLayout>

StringNodeWidget::StringNodeWidget(const std::weak_ptr<BinaryNode>& node,
                                   QWidget *parent)
    : BaseNodeWidget(node, parent)
    , _dataTextEdit(new QLineEdit(this))
{
    connect(_dataTextEdit, &QLineEdit::textEdited, this, &onValueChanged);
    auto nodePtr = node.lock();
    if(nodePtr) {
        connect(nodePtr.get(), &BinaryNode::valueChanged, this, &onDataChanged);
        _dataTextEdit->setText(nodePtr->stringValue());
        _dataTextEdit->setStyleSheet("QLineEdit { color: " + nodePtr->style().textColor.name() + ";"
                                               "  background-color: " + nodePtr->style().foregroundColor.name() + ";"
                                               "  border: 0"
                                               "}");
    }
}

void StringNodeWidget::addContent(QGridLayout* contentLayout)
{
    contentLayout->addWidget(_dataTextEdit, 0, 0, 1, 1, Qt::AlignBottom | Qt::AlignHCenter);
}

void StringNodeWidget::onValueChanged(const QString& text)
{
    auto nodePtr = _node.lock();
    if(nodePtr) {
        nodePtr->setValue(text);
    }
}

void StringNodeWidget::onDataChanged()
{
    auto nodePtr = _node.lock();
    if(nodePtr) {
        if(nodePtr->stringValue() != _dataTextEdit->text())
            _dataTextEdit->setText(nodePtr->stringValue());
    }
}
