import QtQuick 2.15
import QtQuick.Controls 2.15
import QtQuick.Layouts 1.15

Item {
    id: root

    // API
    //

    signal showVideo()

    // Private
    //

    ColumnLayout {
        id: mainLayout

        spacing: 10
        anchors.fill: parent

        Item {
            Layout.fillHeight: true
        }

        Text {
            Layout.alignment: Qt.AlignHCenter

            id: textLabel
            text: "Intro View here"
        }

        Button {
            Layout.alignment: Qt.AlignHCenter

            text: "See Video"
            onClicked: root.showVideo()
        }

        Item {
            Layout.fillHeight: true
        }
    }
}
