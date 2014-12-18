package main

import (
  "fmt"
  "os"
  "time"

  "gopkg.in/qml.v1"
)

type Pomodoro struct {
  IsRinging bool
  ringing qml.Object
  ticking qml.Object
}

func (p *Pomodoro) Run() {

  go func() {
    ticktack := func() {
      p.ticking.Set("loops", 1 << 10)
      p.ticking.Call("play")
    }
    ticktack()
    for {
      time.Sleep(time.Second * 3)
      p.IsRinging = !p.IsRinging
      if (p.IsRinging) {
        p.ticking.Call("stop")
        p.ringing.Set("loops", 1 << 10)
        p.ringing.Call("play")
      } else {
        p.ringing.Call("stop")
        ticktack()
      }
      qml.Changed(p, &p.IsRinging)
    }
  }()

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

  pomo := Pomodoro{}

  context := engine.Context()
  context.SetVar("pomodoro", &pomo)

  win := component.CreateWindow(nil)

  pomo.ticking = win.Root().ObjectByName("tick")
  pomo.ringing = win.Root().ObjectByName("ring")

  pomo.Run()

  win.Show()
  win.Wait()

  return nil
}