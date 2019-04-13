//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"time"
)

const (
	maxSeconds = 10
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

type UserStrategy interface {
	Process(process func()) bool
}

type PremiumUserFlow struct {
	User *User
}

func (pu PremiumUserFlow) Process (process func()) bool {
	process()
	return true
}

type FreeUserFlow struct {
	User *User
}


func (fu FreeUserFlow) Process (process func()) bool {
	timeLeft := fu.getTimeLeft()
	limiter := time.Tick(time.Duration(timeLeft) * time.Second)
	start := time.Now()
	defer fu.updateUsageTime(start)
	done := make(chan bool)
	go func(done chan bool) {
		// by default time.sleep can not be interrupted
		// so we have to put it into the separate goroutine so
		// ticker can work in current routine
		process()
		done <- true
	}(done)
	select {
	case <- limiter :

		return false
	case <- done:
		return true
	}
}


func (fu FreeUserFlow) updateUsageTime(start time.Time)  {
		t := time.Now()
		secondsConsumed := int64(t.Sub(start).Seconds())
		fu.User.TimeUsed += secondsConsumed
}


func (fu FreeUserFlow) getTimeLeft() int64 {
	timeLeft := maxSeconds - fu.User.TimeUsed
	if timeLeft < 0 {
		timeLeft = 0
	}
	return timeLeft
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	strategy := GetUserStrategy(u)
	return strategy.Process(process)
}

func GetUserStrategy(u *User) UserStrategy  {
	if u.IsPremium {
		return PremiumUserFlow{User: u}
	}
	return FreeUserFlow{User: u}
}

func main() {
	RunMockServer()
}
