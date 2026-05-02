package session

import "testing"

func TestUserSession_Continue(t *testing.T) {
	userSession := NewUserSession(1)

	userSession.Continue()
	userSession.Continue()
	userSession.Continue()

	if len(userSession.continueChan) != 1 {
		t.Errorf("expected 1 signal in channel, got %d", len(userSession.continueChan))
	}
}

func TestUserSession_WaitContinue(t *testing.T) {
	userSession := NewUserSession(1)
	userSession.Continue()
	userSession.Continue()
	userSession.Continue()

	select {
	case <-userSession.WaitContinue():
	default:
		t.Errorf("expected a continue signal")
	}
}
