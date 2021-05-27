import QtQuick 2.15
import QtQuick.Controls 2.15
import QtQuick.Layouts 1.15

import custom.Controls 1.0

Item {
    id: root

    // API
    //
    signal showPIP()
    signal goBack()
    property bool fullScreen: false

    // Private
    //

    property real swipeRatio: 0

    onSwipeRatioChanged: {
        if(root.fullScreen) {
            if((swipeView.currentIndex === 0 && root.swipeRatio <= 0))
                root.fullScreen = false
        }
        else {
            if(swipeView.currentIndex === 1 && root.swipeRatio >= 1)
                root.fullScreen = true
        }
    }

    ColumnLayout {
        anchors.fill: parent
        spacing: 10

        Text {
            Layout.alignment: Qt.AlignHCenter

            visible: !root.fullScreen

            text: qsTr("Video image")
        }

        // Vertical spacer
        Item {
            Layout.fillHeight: true
        }

        PlayerView {
            Layout.preferredWidth: ((root.width < root.height ? root.width : root.height) * (0.5 + (root.swipeRatio * 0.5))).toFixed(0)
            Layout.preferredHeight: Layout.preferredWidth
            Layout.alignment: Qt.AlignHCenter
        }

        // Vertical spacer
        Item {
            Layout.fillHeight: true
        }

        RowLayout {
            visible: !root.fullScreen

            Button {
                Layout.margins: 20

                text: "Move to PIP"

                onClicked: root.showPIP()
            }

            Item {
                Layout.fillWidth: true
            }

            Button {
                Layout.margins: 20

                text: "Back"

                //visible: !Context.isMobilePlatform

                onClicked: root.goBack()
            }
        }
    }

    // Quick workaround for swipe gesture detection
    SwipeView {
        id: swipeView

        anchors.fill: parent
        z: -1

        Timer {
            interval: 1000/60; running: true; repeat: true
            onTriggered: {
                root.swipeRatio = (upSwipable.mapFromItem(root, 0, 0).y)/(upSwipable.height != 0 ? upSwipable.height : 1)
            }
        }

        currentIndex: 0
        orientation: Qt.Vertical

        // Non-fullscreen item
        Rectangle {
            id: upSwipable
        }
        // Fullscreen item
        Item {
        }
    }
}
