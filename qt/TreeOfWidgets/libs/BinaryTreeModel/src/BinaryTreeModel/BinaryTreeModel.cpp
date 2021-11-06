#include "BinaryTreeModel.h"

#include <queue>

BinaryTreeModel::BinaryTreeModel(QObject *parent)
    : QObject(parent)
{

}

void BinaryTreeModel::setRootItem(BinaryNodePtr rootItem)
{
    _rootItem = rootItem;
}

BinaryNodePtr BinaryTreeModel::rootItem()
{
    return _rootItem;
}

BinaryNodeConstPtr BinaryTreeModel::rootItem() const
{
    return _rootItem;
}
