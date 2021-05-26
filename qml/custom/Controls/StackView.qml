import QtQuick 2.15

/// Stack View Control that handles navigating through a list of views in a stack order
//
// Will automatically ingest children views pushing them on the stack in order
//   the first found is visible at the top
//   all other children are hidden and lower in the stack
//
Item {
    id: root

    // API
    //
    function push(view) {
        if(view instanceof Item) {
            _views.push(view);
            root.currentView = _views[_views.length - 1];
        }
        else {
            throw 'Pushed View is not of type Item!';
        }
    }

    function pop() {
        if(root.currentView)
            root.currentView.visible = false

        _views.pop();

        if(_views.length > 1)
            root.currentView = _views[_views.length - 1];
        else
            root.currentView = null
    }

    property Item currentView: null

    // Private
    //

    property var _views: []

    Component.onCompleted: {
        for(var i in root.children) {
            let item = root.children[i]
            if(item instanceof Item) {
                _views.push(item);
                if(root.currentView === null) {
                    root.currentView = item;
                }
                else {
                    item.visible = false;
                }
            }
        }
    }

    // Ensure changes on the currentView are in sync
    Connections {
        target: root
        function onCurrentViewChanged() {
            if(root.currentView !== null) {
                if(_views.indexOf(root.currentView) === -1)
                    root.push(root.currentView);

                if(root.currentView.parent !== root)
                    root.currentView.parent = root;

                root.currentView.anchors.fill = Qt.binding(function() { return root; } )

                if(_views.length > 1)
                    _views[_views.length - 2].visible = false;
                root.currentView.visible = true;
            }
            else {
                root.pop();
            }
        }
    }
}
