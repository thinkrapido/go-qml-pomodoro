package main

import (
  "fmt"
  "os"
  "time"

  "gopkg.in/qml.v1"

  state "github.com/thinkrapido/go-statemachine"
)

const (
  Stop = iota
  Start
  Ring
)

type Pomodoro struct {
  state.Machine

  ringing qml.Object
  ticking qml.Object

  Countdown string

  timer int32
  timerDefault int32

  IsRinging bool
}

func NewPomodoro() *Pomodoro {
  out := &Pomodoro{ timerDefault: 25 * 60 }
  out.init()
  return out
}

func (p *Pomodoro) init() {

  p.Machine.Init()

  p.Learn("waiting", "ticking", "click", p.startClock)
  p.Learn("ticking", "waiting", "click", p.resetClock)

  p.Learn("ticking", "ticking", "tick", p.tick)
  p.Learn("ticking", "ringing", "timeout", p.startRinging)
  p.Learn("ringing", "waiting", "click", p.stopRinging)

  p.SetStartState("waiting")
  p.AddListener(p)

  p.timerDefault = 10
  p.Reset()

}

func (p *Pomodoro) Notify(e *state.Event) {
  switch e.Event {
    case state.StateReachedEvent:
      p.IsRinging = false
      switch p.CurrentState() {
        case "ticking":
          if p.Timer() == 0 {
            p.Trigger("timeout")
          } else {
            time.Sleep(time.Second)
            p.Trigger("tick")
          }
        case "ringing":
          p.IsRinging = true
      }
      qml.Changed(p, &p.IsRinging)
    case state.InconsistencyEvent:
      fmt.Printf("state inconsistent: %s\n", e.Message)
    case state.KillEvent:
      p.Run()
  }
}

func (p *Pomodoro) startClock() {
  p.Reset()
  p.ticking.Set("loops", 1 << 10)
  p.ticking.Call("play")
}

func (p *Pomodoro) resetClock() {
  p.Reset()
  p.ticking.Call("stop")
}

func (p *Pomodoro) tick() {
  p.decTimer()
}

func (p *Pomodoro) startRinging() {
  p.ticking.Call("stop")
  p.ringing.Set("loops", 1 << 10)
  p.ringing.Call("play")
}

func (p *Pomodoro) stopRinging() {
  p.Reset()
  p.ticking.Call("stop")
  p.ringing.Call("stop")
}

func (p *Pomodoro) Reset() *Pomodoro {
  p.setTimer(p.timerDefault)
  return p
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
    pomo.Trigger("click");
  })

  pomo.Run()

  win.Show()
  win.Wait()

  return nil
}