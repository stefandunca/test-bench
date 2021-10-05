import QtQuick 2.0

import QtMultimedia 5.15

Item {
    id: root

    // API
    //
    property bool videoPlayerFailed: mediaPlayer.error !== MediaPlayer.NoError
    property string videoPlayerError: root.videoPlayerFailed ? ("MediaPlayer error code" + mediaPlayer.error + "; Error message " + mediaPlayer.errorString) : ""

    Rectangle {
        anchors.fill: parent

        border.color: "black"

        VideoOutput {
            id: videoOutput

            visible: !root.videoPlayerFailed

            anchors.fill: parent
            anchors.margins: parent.border.width

            source: mediaPlayer
        }

        Image {
            visible: !videoOutput.visible

            anchors.fill: parent
            anchors.margins: parent.border.width

            source: "qrc:/walking.gif"
        }
    }

    MediaPlayer {
        id: mediaPlayer

        loops: MediaPlayer.Infinite
        autoPlay: true

        source: "qrc:/walking.gif"

        onError: {
            if (MediaPlayer.NoError !== error) {
                console.log("[qmlvideo] " + root.videoPlayerError)
                source = backupFsPath
            }
        }

        readonly property string backupFsPath: "file:/" + AppContext.pwd + "/walking.gif"
    }
}
