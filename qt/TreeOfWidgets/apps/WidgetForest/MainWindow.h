#pragma once

#include <QMainWindow>

QT_BEGIN_NAMESPACE
namespace Ui {
    class MainWindow;
}
QT_END_NAMESPACE

class BinaryTreeModel;

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(BinaryTreeModel& model, QWidget *parent = nullptr);
    ~MainWindow();

private:

    Ui::MainWindow *ui;
    BinaryTreeModel& _model;
};
