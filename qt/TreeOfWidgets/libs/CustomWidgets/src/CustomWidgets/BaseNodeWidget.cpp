#include "BaseNodeWidget.h"

#include<BinaryTreeModel/BinaryNode.h>

#include <QPainter>
#include <QGridLayout>
#include <QLabel>
#include <QSpacerItem>

const int BaseNodeWidget::_borderSize{1};

BaseNodeWidget::BaseNodeWidget(const std::weak_ptr<BinaryNode>& node, QWidget *parent)
    : QWidget(parent)
    , _mainLayout(new QGridLayout(this))
    , _contentLayout(new QGridLayout())
    , _titleLabel(new QLabel(this))
    , _sideBarSpacer(new QSpacerItem(10, 100, QSizePolicy::Fixed, QSizePolicy::MinimumExpanding))
    , _node(node)
{
    // Setup layout as 2x2 with left cells ocupied by side bar
    // right-top cell by the title
    // right-bottom cell by the content
    _mainLayout->setSpacing(0);
    _mainLayout->setContentsMargins(10, 10, 10, 10);

    _mainLayout->addItem(_sideBarSpacer, 0, 0, 2, 1);
    _mainLayout->addWidget(_titleLabel, 0, 1, Qt::AlignLeft);

    _mainLayout->addItem(_contentLayout, 1, 1, Qt::AlignCenter);

    // Populate data
    auto nodePtr = _node.lock();
    if(nodePtr) {
        _titleLabel->setText(nodePtr->name());

        // Style object
        auto& style = nodePtr->style();
        if(style.fixedSize)
            setFixedSize(150, 150);
        setAttribute(Qt::WA_StyledBackground, true);
        setStyleSheet("BaseNodeWidget { border: " + QString::number(style.borderColor.isValid() ? _borderSize : 0) + " solid " +
                                                (style.borderColor.isValid() ? style.borderColor.name() : QString()) + ";"
                                     "  color: " + style.textColor.name() + ";"
                                     "  background-color: " + style.backgroundColor.name() + ";"
                                     "}"
        );
        _titleLabel->setStyleSheet("QLabel { color: " + style.titleColor.name() + ";"
                                            "font: 26px;"
                                          "}");
    }
}

void BaseNodeWidget::init()
{
    addContent(_contentLayout);
}

void BaseNodeWidget::paintEvent(QPaintEvent* /*event*/)
{
    auto nodePtr = _node.lock();
    if(nodePtr) {
        auto& style = nodePtr->style();
        QPainter painter(this);
        painter.save();
        painter.setPen(Qt::NoPen);
        painter.setBrush(QBrush(style.accentColor, Qt::SolidPattern));
        painter.drawRect(QRect(_borderSize, _borderSize, _sideBarSpacer->geometry().width(), height() - _borderSize - 1));
        painter.restore();
    }
}
