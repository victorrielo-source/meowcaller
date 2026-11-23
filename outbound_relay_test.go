package meowcaller

import "testing"

// Live case (Miami server calling a phone in Brazil): we were offered mia3c02 (own),
// mia3c03 and gru2c01; the phone probed gru2c01, gig4c01 and mia3c03 only.
func TestPickOutboundRelayUsesARelayTheCalleeProbed(t *testing.T) {
	addr := []relayAddress{{ipv4: "157.240.14.129", port: 3478}}
	rd := &relayData{
		relayTokens: [][]byte{[]byte("t0"), []byte("t1")},
		endpoints: []relayEndpoint{
			{relayName: "gru2c01", tokenID: 1, c2rRTT: 104, addresses: addr},
			{relayName: "mia3c02", tokenID: 0, authTokenID: 1, c2rRTT: 67, addresses: addr},
			{relayName: "mia3c03", tokenID: 0, c2rRTT: 67, addresses: addr},
		},
	}
	probes := []rlProbe{
		{relayName: "gru2c01", latency: 20},
		{relayName: "gig4c01", latency: 1},
		{relayName: "mia3c03", latency: 131},
	}
	if got := pickOutboundRelay(rd, probes); got != "gru2c01" {
		t.Fatalf("picked %q, want gru2c01 (104+20 beats 67+131; gig4c01 is not ours)", got)
	}
	rd.electedRelay = "gru2c01"
	if ep := getMediaRelayEndpoint(rd, false); ep == nil || ep.relayName != "gru2c01" {
		t.Fatalf("outbound endpoint should follow the elected relay, got %+v", ep)
	}
	if got := pickOutboundRelay(rd, []rlProbe{{relayName: "gig4c01", latency: 5}}); got != "" {
		t.Fatalf("no shared relay should return empty, got %q", got)
	}
	if got := pickOutboundRelay(rd, nil); got != "" {
		t.Fatalf("no probes should return empty, got %q", got)
	}
}
