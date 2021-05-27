#include "AppContext.h"

#include <QCoreApplication>

AppContext::AppContext(QObject *parent) : QObject(parent)
{

}

QString AppContext::pwd() const
{
    return QCoreApplication::applicationDirPath();
}
