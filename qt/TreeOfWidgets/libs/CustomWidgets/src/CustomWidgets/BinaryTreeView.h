#pragma once

#include <QWidget>

#include <tuple>
#include <mutex>
#include <list>

class BinaryTreeModel;
class BinaryNode;
class BaseNodeWidget;

class BinaryTreeView : public QWidget
{
    Q_OBJECT
public:
    explicit BinaryTreeView(QWidget *parent = nullptr);

    void setModel(BinaryTreeModel& model);

private:
    void resizeEvent(QResizeEvent *event) override;
    void paintEvent(QPaintEvent* /*event*/);

    // Ugly factory function that also does layouting, drawing and what not
    BaseNodeWidget* createInPreorder(const std::shared_ptr<BinaryNode> node, const BaseNodeWidget* parent = nullptr,
                                     const int treeDepth = 0, const bool left = true);

    BaseNodeWidget* _rootWidget = nullptr;

    BinaryTreeModel* _model;

    std::mutex _mutex;
    std::list<std::tuple<QPoint, QPoint>> _lines;
};
