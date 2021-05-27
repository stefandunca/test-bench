import QtQml 2.15

pragma Singleton

QtObject {
    property bool isMobilePlatform: ["android", "ios"].indexOf(Qt.platform.os) >= 0
}
