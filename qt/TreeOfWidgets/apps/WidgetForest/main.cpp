#include "MainWindow.h"

#include "BinaryTreeModel/BinaryTreeModel.h"
#include "BinaryTreeModel/BinaryNode.h"

#include <QApplication>

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);

    BinaryTreeModel model;

    //         0
    //      /    \
    //     1      8
    //    /\     / \
    //   2  6   9   14
    //  /\   \   \   \
    // 3 4   7  10    15
    //   \      /
    //   5     11
    //        / \
    //       12 13

    auto rootNode = std::make_shared<BinaryNode>("Root [0]", "Data in multi");
    auto n1 = std::make_shared<BinaryNode>("1", 22.0);
    rootNode->setLeft(n1);
    auto n2 = std::make_shared<BinaryNode>("2 play!", 2.0);
    n1->setLeft(n2);
    n2->setLeft(std::make_shared<BinaryNode>("3", 3.0));
    n2->setRight(std::make_shared<BinaryNode>("4", 4));
    n2->right()->setRight(std::make_shared<BinaryNode>("5", "Change me 5!"));
    n1->setRight(std::make_shared<BinaryNode>("6", 2.0));
    n1->right()->setRight(std::make_shared<BinaryNode>("7", "Change me 7!"));

    auto n8 = std::make_shared<BinaryNode>("8", "December");
    rootNode->setRight(n8);
    n8->setLeft(std::make_shared<BinaryNode>("9", "Text data"));
    n8->left()->setRight(std::make_shared<BinaryNode>("10", "10"));
    auto n11 = std::make_shared<BinaryNode>("11 play!", 11);
    n8->left()->right()->setLeft(n11);
    n11->setLeft(std::make_shared<BinaryNode>("12", "12"));
    n11->setRight(std::make_shared<BinaryNode>("13", 13));
    n8->setRight(std::make_shared<BinaryNode>("14", "2021"));
    n8->right()->setRight(std::make_shared<BinaryNode>("15", "Soon"));
    model.setRootItem(rootNode);

    QObject::connect(n11.get(), &BinaryNode::valueChanged, [n2, n11]() { n2->setValue(n11->floatValue()); });
    QObject::connect(n2.get(), &BinaryNode::valueChanged, [n2, n11]() { n11->setValue(n2->floatValue()); });
    QObject::connect(n2->right()->right().get(), &BinaryNode::valueChanged, [n2, n1]() { n1->right()->right()->setValue(n2->right()->right()->stringValue()); });
    QObject::connect(n1->right()->right().get(), &BinaryNode::valueChanged, [n2, n1]() { n2->right()->right()->setValue(n1->right()->right()->stringValue()); });

    MainWindow w(model);
    w.show();
    return a.exec();
}
