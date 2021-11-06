#pragma once

#include "BinaryNode.h"

#include <QObject>

class BinaryTreeModel : public QObject
{
    Q_OBJECT
public:
    explicit BinaryTreeModel(QObject *parent = nullptr);

    void setRootItem(BinaryNodePtr rootItem);
    BinaryNodePtr rootItem();
    BinaryNodeConstPtr rootItem() const;

private:
    BinaryNodePtr _rootItem;
};
