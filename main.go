package main

import (
  "fmt"
  "os"
  "time"

  "gopkg.in/qml.v1"
)

const (
  Stop = iota
  Start
  Ring
)

type Pomodoro struct {
  IsRinging bool
  isTicking bool
  ringing qml.Object
  ticking qml.Object
  Countdown string
  timer int32
  timerDefault int32
  event chan int8
}

func NewPomodoro() *Pomodoro {
  out := &Pomodoro{ timerDefault: 25 * 60 }
  out.Reset()
  out.event = make(chan int8)
  return out
}

func (p *Pomodoro) Run() {

  p.timerDefault = 10
  go func() {
    p.event <- Stop
  }()

  go func() {
    for event := range p.event {
      switch event {
        case Start:
          p.Reset()
          p.ticking.Set("loops", 1 << 10)
          p.ticking.Call("play")
          p.IsRinging = false
          p.isTicking = true
          qml.Changed(p, &p.IsRinging)
          go func() {
            for p.isTicking {
              time.Sleep(time.Second)
              if p.isTicking && p.decTimer().Timer() == 0 {
                p.event <- Ring
              }
            }
          }()
        case Stop:
          p.Reset()
          p.ticking.Call("stop")
          p.ringing.Call("stop")
          p.IsRinging = false
          p.isTicking = false
          qml.Changed(p, &p.IsRinging)
        case Ring:
          p.ringing.Set("loops", 1 << 10)
          p.ringing.Call("play")
          p.IsRinging = true
          p.isTicking = false
          qml.Changed(p, &p.IsRinging)
      }
    }
  }()

}
func (p *Pomodoro) ToggleRun() {
  if p.IsRinging || p.isTicking {
    p.event <- Stop
  } else if !p.isTicking {
    p.event <- Start
  }
}

func (p *Pomodoro) SetTimerDefault(timer int32) *Pomodoro {
  p.timerDefault = timer
  return p
}
func (p *Pomodoro) setTimer(timer int32) *Pomodoro {
  p.timer = timer
  p.Countdown = fmt.Sprintf("%02d:%02d", p.timer / 60, p.timer % 60)
  qml.Changed(p, &p.Countdown)
  return p
}
func (p *Pomodoro) Timer() int32 {
  return p.timer
}
func (p *Pomodoro) decTimer() *Pomodoro {
  p.setTimer(p.Timer() - 1)
  return p
}
func (p *Pomodoro) Reset() *Pomodoro {
  p.setTimer(p.timerDefault)
  return p
}

func main() {
  if err := qml.Run(run); err != nil {
    fmt.Fprintf(os.Stderr, "error: %v\n", err)
    os.Exit(1)
  }
}

func run() error {

  engine := qml.NewEngine()
  component, err := engine.LoadFile("assets/pomodoro.qml")
  if err != nil {
    return err
  }

  pomo := NewPomodoro()

  context := engine.Context()
  context.SetVar("pomodoro", pomo)

  win := component.CreateWindow(nil)

  pomo.ticking = win.Root().ObjectByName("tick")
  pomo.ringing = win.Root().ObjectByName("ring")
  mouse := win.Root().ObjectByName("mouseArea")
  mouse.On("clicked", func(event qml.Object) {
    pomo.ToggleRun();
  })

  pomo.Run()

  win.Show()
  win.Wait()

  return nil
}