import QtQuick 2.15
import QtQuick.Controls 2.15

Item {
    id: root

    // API
    //

    signal showVideo()

    // Private
    //

    Text {
        id: textLabel
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: parent.top
        text: "Intro View here"
    }

    Button {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: textLabel.bottom
        anchors.margins: 10

        text: "See Video"
        onClicked: root.showVideo()
    }
}
