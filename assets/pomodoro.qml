import QtQuick 2.2
import QtMultimedia 5.0

Column {
  Rectangle {
    width: 200
    height: 50
    color: '#cccccc'
    Text {
      anchors {
        fill: parent.fill
        horizontalCenter: parent.horizontalCenter
        verticalCenter: parent.verticalCenter
      }
      text: pomodoro.countdown
      color: 'black'
      font {
        family: 'Arial'
        pointSize: 24
        bold: true
      }
    }
  }
  Image {
    id: pImg

    width: 200
    height: 200
    rotation: 0

    source: 'pomodoro.png'

    state: { return pomodoro.isRinging ? 'ringing' : 'normal' }
    states: [
      State {
        name: 'normal'
        PropertyChanges {
          target: pImg
          rotation: 0
        }
      },
      State {
        name: 'ringing'
      }
    ]

    MouseArea {
      objectName: 'mouseArea'
      anchors.fill: parent
    }

    transitions: [
      Transition { 
        to: 'normal'
        RotationAnimation { target: pImg; property: 'rotation'; duration: 50; to: 0 }
      },
      Transition { 
        to: 'ringing'
        SequentialAnimation { 
          id: ringingAnim
          loops: Animation.Infinite
          RotationAnimation { target: pImg; property: 'rotation'; duration: 50; to: 10 }
          RotationAnimation { target: pImg; property: 'rotation'; duration: 50; to: -10 }
        }
      }
    ]

    SoundEffect {
      objectName: 'tick'
      id: tick
      source: 'tick.wav'
    }

    SoundEffect {
      objectName: 'ring'
      id: ring
      source: 'ring.wav'
    }

  }
}
