import QtQuick
import QtQuick.Window
import QtQuick.Layouts

Window {
    width: mainLayout.implicitWidth
    height: mainLayout.implicitHeight

    visible: true
    title: qsTr("Hello World")

    ColumnLayout {
        id: mainLayout

        anchors.fill: parent

        LayoutsExample {
            Layout.margins: 10
        }
    }
}
