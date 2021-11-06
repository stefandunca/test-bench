#include "MainWindow.h"
#include "./ui_MainWindow.h"

#include <BinaryTreeModel/BinaryTreeModel.h>
#include <CustomWidgets/BinaryTreeView.h>

MainWindow::MainWindow(BinaryTreeModel& model, QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
    , _model(model)
{
    ui->setupUi(this);
    ui->binaryTreeView->setModel(_model);
}

MainWindow::~MainWindow()
{
    delete ui;
}
