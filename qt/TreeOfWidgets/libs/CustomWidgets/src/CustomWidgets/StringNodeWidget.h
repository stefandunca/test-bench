#pragma once

#include "BaseNodeWidget.h"

#include <QObject>

class QLineEdit;

class StringNodeWidget : public BaseNodeWidget
{
    Q_OBJECT
public:
    StringNodeWidget(const std::weak_ptr<BinaryNode>& node,
                     QWidget *parent = nullptr);

    void addContent(QGridLayout* contentLayout) override;

public slots:
    void onValueChanged(const QString& text);
    void onDataChanged();

private:
    QLineEdit* _dataTextEdit;
};
