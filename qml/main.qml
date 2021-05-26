import QtQuick 2.15
import QtQuick.Window 2.15

import custom.Controls 1.0
import "Views" as Views

Window {
    width: 720
    height: 1024
    visible: true

    title: qsTr("Test Bench")

    Component {
        id: mediaView
        Views.MediaView {
            onShowPIP: stackView.push(pipView.createObject(stackView));
        }
    }
    Component {
        id: pipView
        Views.PipView {
            onGoBack: stackView.pop();
        }
    }

    Loader {
        id: viewLoader
    }

    // TODO: if desktop add a drawer and a navigation bar
    // TODO: if mobile use the system's navigation bar
    StackView {
        id: stackView

        anchors.fill: parent

        Views.IntroView {
            onShowVideo: stackView.push(mediaView.createObject(stackView))
        }
    }
}
