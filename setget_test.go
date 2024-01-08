package copier_test

import (
	"github.com/aacfactory/copier"
	"testing"
	"time"
)

type Getter struct {
	Finished bool      `copy:",GetFinished"`
	Deadline time.Time `copy:",GetDeadline"`
}

func (g Getter) GetFinished() bool {
	return true
}

func (g Getter) GetDeadline() int64 {
	return g.Deadline.UnixMilli()
}

type Setter struct {
	Succeed  bool      `copy:"Finished,SetSucceed"`
	Deadline time.Time `copy:"Deadline,SetDeadline"`
}

func (s *Setter) SetSucceed(b bool) {
	s.Succeed = b
}

func (s *Setter) SetDeadline(n int64) {
	s.Deadline = time.UnixMilli(n)
}

func TestCopy_SetGet(t *testing.T) {
	s := Setter{}
	err := copier.Copy(&s, Getter{
		Finished: false,
		Deadline: time.Now(),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", s)
}
