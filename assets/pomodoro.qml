import QtQuick 2.2
import QtMultimedia 5.0

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