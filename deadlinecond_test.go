package deadlinecond

import (
  "testing"
  "time"
  "sync"
  "math"
)

func TestCond(t *testing.T){
  cond := NewCond(nil)
  go func(){
    time.Sleep(10 * time.Millisecond)
    cond.Signal()
  }()
  n := time.Now()
  cond.L.Lock()
  cond.SetDeadline(n.Add(20 * time.Millisecond))
  cond.Wait()
  e := time.Now().Sub(n)
  if e > 11 * time.Millisecond {
    t.Fatal("Waited too long")
  }
}

func TestCondWait2(t *testing.T){
  cond := NewCond(nil)
  go func(){
    time.Sleep(10 * time.Millisecond)
    cond.Signal()
  }()
  n := time.Now()
  cond.L.Lock()
  cond.SetDeadline(n.Add(20 * time.Millisecond))
  w := cond.Wait()
  if w {
    t.Fatal("Timedout althought it shouldn't")
  }
}

func TestCondWait2Timeout(t *testing.T){
  cond := NewCond(nil)
  go func(){
    time.Sleep(20 * time.Millisecond)
    cond.Signal()
  }()
  n := time.Now()
  cond.L.Lock()
  cond.SetDeadline(n.Add(10 * time.Millisecond))
  w := cond.Wait()
  if !w {
    t.Fatal("Not timedout althought it should")
  }
}

func TestCondEdgeCase(t *testing.T){
  cond := &Cond{Cond: *sync.NewCond(&sync.Mutex{}), waitTag: math.MaxInt32 }
  go func(){
    time.Sleep(20 * time.Millisecond)
    cond.Signal()
  }()
  n := time.Now()
  cond.L.Lock()
  cond.SetDeadline(n.Add(10 * time.Millisecond))
  w := cond.Wait()
  if !w {
    t.Fatal("Not timedout althought it should")
  }
}


