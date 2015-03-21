import QtQuick 2.2
import QtQuick.Controls 1.1
import QtWebKit 3.0

ApplicationWindow {
    id: root
    visible: true
    title: "Browse View"

    width: 650
    height: 500


    ScrollView {
        anchors.fill: parent

        WebView {
            id: browseView
            url: "http://localhost:1055/"
            anchors.fill: parent

            onLoadingChanged: {
                if (loadRequest.status === WebView.LoadFailedStatus) {
                    loadStatusTxt.text = "Load failed.";
                } else {
                    loadStatusTxt.text = "Load worked!";
                }
            }
        }
    }

    Text { 
        id: loadStatusTxt 
        text: "Starting up..."
        anchors.bottom: parent.bottom
    }
}
