import QtQuick 2.15
import QtQuick.Controls 2.15
import QtQuick.Layouts 1.15

import QtMultimedia 5.15

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

            text: "PIP View here. Feel free to move the video around!"
        }

        Text {
            visible: playerView.videoPlayerFailed

            Layout.fillWidth: true

            text: "Failed playing video. Error: " + playerView.videoPlayerError
            color: "darkred"

            horizontalAlignment: Text.AlignHCenter
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

    PlayerView {
        id: playerView

        z: 100

        x: 100
        y: 100
        width: root.width / 2
        height: width

        MouseArea {
            anchors.fill: parent
            drag.target: parent
            drag.minimumX: 0
            drag.minimumY: 0
            drag.maximumX: root.width - parent.width
            drag.maximumY: root.height - parent.height
        }
    }
}
