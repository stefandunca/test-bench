import QtQuick 2.15
import QtQuick.Window 2.15

import custom.Controls 1.0
import "Views" as Views

Window {
    id: root
    width: 720
    height: 1024
    visible: true

    title: qsTr("Test Bench")

    Component {
        id: dynamicViewComponent
        Views.DynamicView {
            onGoBack: stackView.pop();
        }
    }

    StackView {
        id: stackView

        anchors.fill: parent

        initialItem: introView

        Views.PipView {
            id: pipView
            onGoBack: stackView.pop();
            onGoForward: stackView.push(dynamicViewComponent.createObject(stackView))
        }
        Views.MediaView {
            id: mediaView
            onShowPIP: stackView.push(pipView);
            onGoBack: stackView.pop();
            onFullScreenChanged: root.visibility = fullScreen ? Window.FullScreen : Window.Windowed;
        }
        Views.IntroView {
            id: introView
            onShowVideo: stackView.push(mediaView);
        }

        focus: true
        Keys.onReleased: {
            if(event.key === Qt.Key_Back || event.key === Qt.Key_Backspace) {
                event.accepted = true;
                if(stackView.count > 1)
                    stackView.pop();
                else
                    Qt.quit()
            }
        }
    }
}
