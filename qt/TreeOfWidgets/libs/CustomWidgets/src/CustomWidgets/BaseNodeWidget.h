#pragma once

#include <QWidget>
#include <QColor>

#include <functional>
#include <memory>

class QGridLayout;
class QSpacerItem;
class QLabel;
class BinaryNode;

class BaseNodeWidget : public QWidget
{
    Q_OBJECT
public:
    explicit BaseNodeWidget(const std::weak_ptr<BinaryNode>& node, QWidget *parent = nullptr);

    void init();

    void setRelayoutFunction(std::function<void(int)> fn) { _relayoutFn = fn; }
    void relayout(int width) { if(_relayoutFn) _relayoutFn(width); }

private:
    virtual void addContent(QGridLayout* contentLayout) = 0;

    void paintEvent(QPaintEvent* event) override;

protected:

    QGridLayout* _mainLayout;
    QGridLayout* _contentLayout;
    QSpacerItem* _sideBarSpacer;
    QLabel* _titleLabel;

    std::weak_ptr<BinaryNode> _node;

    std::function<void(int)> _relayoutFn;

    static const int _borderSize;
};
