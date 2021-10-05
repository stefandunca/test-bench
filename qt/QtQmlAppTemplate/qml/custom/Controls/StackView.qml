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
        if(!pushAnimation.running && !popAnimation.running) {
            if(view instanceof Item) {
                _nextView = view;
                pushAnimation.start()
            }
            else
                throw 'Pushed View is not of type Item!';
        }
    }

    function pop() {
        if(!pushAnimation.running && !popAnimation.running) {
            if(_currentView)
                popAnimation.restart()
        }
    }

    property Item initialItem: null
    property int animationDurationMs: 200
    readonly property alias currentView: root._currentView
    readonly property alias leftView: root._leftView
    readonly property alias count: root._count

    // Private
    //
    property var _views: []
    property Item _currentView: null
    property Item _leftView: null
    property Item _nextView: null
    property int _count: 0

    function viewsChanged() {
        _currentView = _views.length > 0 ? _views[_views.length - 1] : null
        if(_currentView) {
            _currentView.visible = true;
            _currentView.parent = root
            _currentView.anchors.fill = root
        }

        _leftView = _views.length > 1 ? _views[_views.length - 2] : null
        if(_leftView) {
            _leftView.visible = false;
            _leftView.anchors.fill = undefined
        }
        _count = _views.length;
    }

    Component.onCompleted: {
        let foundInitialItem = false;
        for(let i = root.children.length - 1; i >= 0; i--) {
            let item = root.children[i]
            if(item instanceof Item) {
                item.visible = false
                if(!foundInitialItem) {
                    _views.push(item);
                    foundInitialItem = (initialItem === item)
                }
            }
        }
        viewsChanged();
    }

    NumberAnimation {
        id: pushAnimation
        target: _currentView
        property: "x"
        from: 0
        to: -root.width
        duration: animationDurationMs
        easing.type: Easing.InOutQuad
        onStarted: {
            _currentView.anchors.fill = undefined
            _nextView.x = _currentView.x + _currentView.width
            _nextView.y = _currentView.y
            _nextView.width = root.width
            _nextView.height = root.height
            _nextView.visible = true;
        }
        onFinished: {
            _views.push(_nextView);
            viewsChanged();
        }
    }

    Connections {
        target: _currentView
        function onXChanged() { _nextView.x = _currentView.x + _currentView.width }
        enabled: _nextView !== null && pushAnimation.running
    }

    NumberAnimation {
        id: popAnimation
        target: _currentView
        property: "x"
        from: 0
        to: root.width
        duration: animationDurationMs
        easing.type: Easing.InOutQuad
        onStarted: {
            _currentView.anchors.fill = undefined
            if(_leftView) {
                _leftView.x = _currentView.x - _currentView.width
                _leftView.y = _currentView.y
                _leftView.width = root.width
                _leftView.height = root.height
                _leftView.visible = true
            }
        }
        onFinished: {
            _currentView.visible = false
            _views.pop();
            viewsChanged();
        }
    }

    Connections {
        target: _currentView
        function onXChanged() { _leftView.x = _currentView.x - _currentView.width }
        enabled: _leftView !== null && popAnimation.running
    }
}
