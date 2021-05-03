package wsnet

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pion/ice/v2"
	"github.com/pion/webrtc/v3"
)

func TestDialICE(t *testing.T) {
	t.Parallel()

	t.Run("TURN with TLS", func(t *testing.T) {
		t.Parallel()

		addr := createTURNServer(t, ice.SchemeTypeTURNS, "test")
		err := DialICE(webrtc.ICEServer{
			URLs:           []string{fmt.Sprintf("turns:%s", addr)},
			Username:       "example",
			Credential:     "test",
			CredentialType: webrtc.ICECredentialTypePassword,
		}, time.Millisecond)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Protocol mismatch", func(t *testing.T) {
		t.Parallel()

		addr := createTURNServer(t, ice.SchemeTypeTURNS, "test")
		err := DialICE(webrtc.ICEServer{
			URLs:           []string{fmt.Sprintf("turn:%s", addr)},
			Username:       "example",
			Credential:     "test",
			CredentialType: webrtc.ICECredentialTypePassword,
		}, time.Millisecond)
		if !errors.Is(err, ErrMismatchedProtocol) {
			t.Error(err)
		}
	})

	t.Run("Invalid auth", func(t *testing.T) {
		t.Parallel()

		addr := createTURNServer(t, ice.SchemeTypeTURNS, "test")
		err := DialICE(webrtc.ICEServer{
			URLs:           []string{fmt.Sprintf("turns:%s", addr)},
			Username:       "example",
			Credential:     "invalid",
			CredentialType: webrtc.ICECredentialTypePassword,
		}, time.Millisecond)
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Error(err)
		}
	})

	t.Run("Protocol mismatch public", func(t *testing.T) {
		t.Parallel()

		err := DialICE(webrtc.ICEServer{
			URLs: []string{"turn:stun.l.google.com:19302"},
		}, time.Millisecond)
		if !errors.Is(err, ErrMismatchedProtocol) {
			t.Error(err)
		}
	})
}
