import QtQuick 2.15
import QtQuick.Controls 2.15
import QtQuick.Layouts 1.15

Item {
    id: root

    // API
    //
    signal goBack()

    // Private
    //

    ColumnLayout {
        anchors.fill: parent
        spacing: 10

        // Vertical spacer
        Item {
            Layout.fillHeight: true
        }

        Text {
            Layout.alignment: Qt.AlignHCenter

            text: "The End!"
        }

        // Vertical spacer
        Item {
            Layout.fillHeight: true
        }

        RowLayout {
            Item {
                Layout.fillWidth: true
            }

            Button {
                Layout.margins: 20

                text: "Back"
                onClicked: root.goBack()
            }
        }
    }
}
