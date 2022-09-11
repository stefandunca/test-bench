#include <QGuiApplication>
#include <QQmlApplicationEngine>

#include <QtDebug>

int main(int argc, char *argv[])
{
    QGuiApplication app(argc, argv);

    QQmlApplicationEngine engine;

    auto pluginDir = QCoreApplication::applicationDirPath() + "/../../libs/Chat";
    qDebug() << pluginDir;

    engine.addPluginPath(pluginDir);

    const QUrl url("qrc:/main.qml");
    QObject::connect(&engine, &QQmlApplicationEngine::objectCreated,
                     &app, [url](QObject *obj, const QUrl &objUrl) {
        if (!obj && url == objUrl)
            QCoreApplication::exit(-1);
    }, Qt::QueuedConnection);
    engine.load(url);

    return app.exec();
}
