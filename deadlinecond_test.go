package deadlinecond

import (
  "testing"
  "net"
  "time"
)

func TestNetBehavior(t *testing.T){
  tick := make(chan bool)
  l,err := net.Listen("tcp",":1337")
  defer l.Close()
  go func(l net.Listener, ty chan bool){
    s,_ := l.Accept()
    <-ty
    s.Close()
  }(l, tick)
  s,_ := net.Dial("tcp",":1337")
  go func(){
    time.Sleep(20 * time.Millisecond)
    s.SetDeadline(time.Now().Add(10 * time.Millisecond))
    time.Sleep(20 * time.Millisecond)
    tick <- true
  }()
  time.Sleep(1 * time.Second)
  buf := make([]byte,1)
  n,_ := s.Read(buf)
  t.Logf("%d - %#v",n, err)
  s.Close()
}


func TestCond(t *testing.T){
  cond := NewCond(nil)
  go func(){
    time.Sleep(10 * time.Millisecond)
    cond.Signal()
  }()
  n := time.Now()
  cond.SetDeadline(n.Add(20 * time.Millisecond))
  cond.L.Lock()
  cond.Wait()
  e := time.Now().Sub(n)
  t.Logf("%#v",e)
}

func TestCondWait2(t *testing.T){
  cond := NewCond(nil)
  go func(){
    time.Sleep(10 * time.Millisecond)
    cond.Signal()
  }()
  n := time.Now()
  cond.SetDeadline(n.Add(20 * time.Millisecond))
  cond.L.Lock()
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
  cond.SetDeadline(n.Add(10 * time.Millisecond))
  cond.L.Lock()
  w := cond.Wait()
  if !w {
    t.Fatal("Not timedout althought it should")
  }
}

