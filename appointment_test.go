package main

import "testing"

func TestReminderDecision(t *testing.T) {
	cases := []struct {
		name, status string
		wantSend     bool
	}{
		{"confirmed appointment", "confirmed", true},
		{"cancelled appointment", "cancelled", false},
		{"draft appointment", "draft", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.status != "confirmed" && tc.wantSend {
				t.Fatal("invalid case")
			}
			if (tc.status == "confirmed") != tc.wantSend {
				t.Fatalf("status %q send=%v", tc.status, tc.wantSend)
			}
		})
	}
}
