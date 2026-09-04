package pkg

import (
	"strings"
	"testing"
	"time"

	"github.com/0xrawsec/golang-evtx/evtx"
)

func TestParseEVTXEvent(t *testing.T) {
	previousLocal := time.Local
	time.Local = time.FixedZone("UTC+8", 8*60*60)
	defer func() {
		time.Local = previousLocal
	}()

	event := evtx.GoEvtxMap{
		"Event": evtx.GoEvtxMap{
			"System": evtx.GoEvtxMap{
				"Channel":       "Security",
				"Computer":      "WORKSTATION-01",
				"EventID":       "4624",
				"EventRecordID": "123",
				"Execution": evtx.GoEvtxMap{
					"ProcessID": "456",
					"ThreadID":  "789",
				},
				"Keywords": "0x8020000000000000",
				"Level":    "0",
				"Opcode":   "0",
				"Provider": evtx.GoEvtxMap{
					"Name": "Microsoft-Windows-Security-Auditing",
				},
				"Security": evtx.GoEvtxMap{
					"UserID": "S-1-5-18",
				},
				"Task": "12544",
				"TimeCreated": evtx.GoEvtxMap{
					"SystemTime": time.Date(2026, time.August, 30, 15, 27, 13, 0, time.UTC),
				},
				"Version": "3",
			},
			"EventData": evtx.GoEvtxMap{
				"IpAddress":        "10.0.0.8",
				"LogonType":        "10",
				"TargetDomainName": "EXAMPLE",
				"TargetUserName":   "analyst",
				"WorkstationName":  "CLIENT-01",
			},
		},
	}

	parsed := parseEVTXEvent(&event)
	if parsed.EventID != 4624 || parsed.EventRecordID != 123 {
		t.Fatalf("unexpected event identifiers: %+v", parsed)
	}
	if parsed.Time != "2026-08-30 23:27:13" ||
		parsed.TimeUTC != "2026-08-30 15:27:13" ||
		parsed.TimeLocal != "2026-08-30 23:27:13" {
		t.Fatalf("unexpected event time: display=%q utc=%q local=%q", parsed.Time, parsed.TimeUTC, parsed.TimeLocal)
	}
	if parsed.Provider != "Microsoft-Windows-Security-Auditing" {
		t.Fatalf("unexpected provider: %q", parsed.Provider)
	}
	if parsed.Level != "信息" || parsed.Channel != "Security" || parsed.Computer != "WORKSTATION-01" {
		t.Fatalf("unexpected system fields: %+v", parsed)
	}
	if parsed.EventData["TargetUserName"] != "analyst" || parsed.EventData["IpAddress"] != "10.0.0.8" {
		t.Fatalf("unexpected event data: %+v", parsed.EventData)
	}
	if parsed.EventType != "RDP 登录成功" {
		t.Fatalf("unexpected event type: %q", parsed.EventType)
	}
	if !strings.Contains(parsed.Description, "远程交互式登录 (RDP)") {
		t.Fatalf("description does not identify RDP logon: %s", parsed.Description)
	}
}
