#pragma once

#include "BaseNodeWidget.h"

#include <QObject>

class QSlider;
class QLabel;

class FloatNodeWidget : public BaseNodeWidget
{
    Q_OBJECT
public:
    FloatNodeWidget(const std::weak_ptr<BinaryNode>& node, QWidget *parent = nullptr);

    void addContent(QGridLayout* contentLayout) override;

private slots:
    void onValueChanged(int value);
    void onDataChanged();

private:
    QLabel* _valueLabel;
    QSlider* _valueSlider;

    const int _decimalCount = 2;
    const int _maxValue = 20;
};
