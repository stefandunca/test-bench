#include "BinaryTreeView.h"

#include "BaseNodeWidget.h"
#include "StringNodeWidget.h"
#include "FloatNodeWidget.h"

#include <BinaryTreeModel/BinaryTreeModel.h>
#include <BinaryTreeModel/BinaryNode.h>

#include <QLabel>
#include <QResizeEvent>
#include <QPainter>

#include <QDebug>

BinaryTreeView::BinaryTreeView(QWidget *parent)
    : QWidget(parent)
{
    setStyleSheet("* { font: 16px; }");
}

void
BinaryTreeView::setModel(BinaryTreeModel &model)
{
    _model = &model;
    _rootWidget = createInPreorder(_model->rootItem(), nullptr);
}

void
BinaryTreeView::resizeEvent(QResizeEvent *event)
{
    _rootWidget->relayout(event->size().width());
}

void
BinaryTreeView::paintEvent(QPaintEvent* /*event*/)
{
    QPainter painter(this);
    painter.save();
    painter.setPen(QPen(QBrush(QColor(0xb96d14)), 5, Qt::SolidLine, Qt::RoundCap));
    painter.setBrush(Qt::SolidPattern);
    std::lock_guard<std::mutex> lock(_mutex);
    for(auto line : _lines)
        painter.drawLine(std::get<0>(line), std::get<1>(line));
    painter.restore();
}

// TODO: make me pretty
BaseNodeWidget*
BinaryTreeView::createInPreorder(const std::shared_ptr<BinaryNode> node, const BaseNodeWidget* parent, const int treeDepth, const bool left)
{
    const int hSpacing = 10;
    const int vSpacing = 25;
    BaseNodeWidget* w = nullptr;
    if(node) {
        switch (node->type()) {
        case QVariant::Type::String:
            w = new StringNodeWidget(node, this);
            break;
        case QVariant::Type::Double:
            w = new FloatNodeWidget(node, this);
            break;
        default:
            qCritical() << "No way!";
            return nullptr;
        }

        w->init();

        auto leftChild = createInPreorder(node->left(), w, treeDepth + 1, true);
        auto rightChild = createInPreorder(node->right(), w, treeDepth + 1, false);

        int offsetFactor = 0;
        if(parent) {
            offsetFactor = (node->left() ? node->left()->height() : 0)
                           + (node->right() ? node->right()->height() : 0);
        }

        _lines.resize(_lines.size() + 1);
        auto& line = _lines.back();

        // Reposition the window based on it's parent
        w->setRelayoutFunction([this, w, parent, left, treeDepth, leftChild, rightChild, offsetFactor, &line](int width) {
            auto parentHCenter = parent ? parent->x() + parent->width()/2 : width/2;
            w->move(parentHCenter + (w->width() + hSpacing) * (left ? -(offsetFactor + 1) : offsetFactor),
                    (treeDepth * (w->height() + vSpacing)));
            if(leftChild)
                leftChild->relayout(width);
            if(rightChild)
                rightChild->relayout(width);
            if(parent) {
                std::lock_guard<std::mutex> lock(_mutex);
                line = std::tuple<QPoint, QPoint>({{parentHCenter, parent->y() + parent->height()},
                                  {w->x() + w->width()/2, w->y()}});
            }
        });
    }
    return w;
}
