#include "CommunitiesPlugin.h"

//#include "ChatModule.h"

void
CommunitiesPlugin::registerTypes(const char *uri)
{
    Q_ASSERT(QString(uri) == "demo");

    //qmlRegisterType<CommunitiesModule>(uri, 1, 0, "ChatModule");
}