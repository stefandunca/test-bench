import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Item {
    implicitWidth: mainLayout.implicitWidth
    implicitHeight: mainLayout.implicitHeight

    RowLayout {
        id: mainLayout

        anchors.fill: parent

        ColumnLayout {
            GoodControl { text: "control uses hints" }
            GoodControl {
                text: "layout sizing"

                Layout.minimumHeight: 10
                Layout.maximumHeight: 50
            }
        }
        Column {
            GoodControl { text: "positioners" }
        }
    }

    component GoodControl: Item {
        id: ctrlRoot

        width: testLayout.implicitWidth
        height: testLayout.implicitHeight

        property string text: ""

        clip: true

        ColumnLayout {
            id: testLayout

            anchors.fill: parent

            Label { text: `Test ${ctrlRoot.text}` }
            Item { Layout.preferredHeight: 5 }
            Label { text: `1. ...\n2. ...` }
        }
    }
}
