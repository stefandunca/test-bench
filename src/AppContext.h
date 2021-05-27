#ifndef APPCONTEXT_H
#define APPCONTEXT_H

#include <QObject>

class AppContext : public QObject
{
    Q_OBJECT

    Q_PROPERTY(QString pwd READ pwd CONSTANT)
public:
    explicit AppContext(QObject *parent = nullptr);

    QString pwd() const;
};

#endif // APPCONTEXT_H
