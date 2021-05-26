import QtQuick 2.15
import QtQuick.Controls 2.15
import QtQuick.Layouts 1.15

import QtMultimedia 5.15

import QtQuick.Dialogs 1.3

Item {
    id: root

    // API
    //
    signal showPIP()

    // Private
    //

    ColumnLayout {
        anchors.fill: parent
        spacing: 10

        Text {
            text: "Media View here"
        }

        Button {
            text: "Move to PIP"
            onClicked: root.showPIP()
        }

        Button {
            text: "Play"
            onClicked: fileDialog.open()
        }

        Rectangle {
            Layout.fillWidth: true
            Layout.fillHeight: true
            Layout.margins: 10

            border.color: "black"

            VideoOutput {
                visible: true
                anchors.fill: parent
                source: mediaPlayer
            }
        }

        FileDialog {
            id: fileDialog
            visible: false
            title: "Please choose a file"
            folder: shortcuts.home
            onAccepted: {
                mediaPlayer.source = fileDialog.fileUrls[0]
                mediaPlayer.play();
            }
        }

        MediaPlayer {
            id: mediaPlayer

            //loops: Audio.Infinite
            //autoPlay: true

            onError: {
                if (MediaPlayer.NoError != error) {
                    console.log("[qmlvideo] VideoItem.onError error " + error + " errorString " + errorString)
                }
            }
        }
    }
}
