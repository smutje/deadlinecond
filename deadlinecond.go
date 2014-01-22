package deadlinecond

import (
  "sync"
  "time"
)

type Cond struct {
  sync.Cond

  waitTag   int32
  timer *time.Timer
}

func NewCond(l sync.Locker) *Cond {
  if l == nil {
    l = &sync.Mutex{}
  }
  return &Cond{ Cond: *sync.NewCond(l) }
}

func (d *Cond) timeout(){
  d.Cond.L.Lock()
  defer d.Cond.L.Unlock()
  d.waitTag++
  d.Cond.Broadcast()
}

func (d *Cond) SetDeadline(t time.Time) error {
  if t == (time.Time{}) {
    if d.timer == nil {
      return nil
    }
    d.timer.Stop()
  }else{
    if d.timer == nil {
      d.timer = time.AfterFunc(t.Sub(time.Now()),d.timeout)
    }else{
      d.timer.Reset(t.Sub(time.Now()))
    }
  }
  return nil
}

func (d *Cond) Wait() bool {
  waitTag := d.waitTag
  d.Cond.Wait()
  return waitTag != d.waitTag
}

