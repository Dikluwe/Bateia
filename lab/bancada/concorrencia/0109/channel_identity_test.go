/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/channel_identity_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"reflect"
	"strings"
	"testing"
)

const identityFixture = `package p;type Gate chan int;type P struct{tokens Gate;other Gate};func(l Gate)release(){<-l};func(p *P)run(){p.tokens<-1;defer p.tokens.release()}`

func TestChannelIdentityContrasts(t *testing.T) {
	for _, c := range []struct{ name, extra, kind string }{
		{"no-event", "", ""},
		{"write", ";func(p *P)configure(){p.tokens=make(Gate,2)}", "FIELD_WRITE"},
		{"alias-write", ";func(p *P)configure(){alias:=p;alias.tokens=make(Gate,2)}", "FIELD_WRITE"},
		{"paren-write", ";func(p *P)configure(){(p.tokens)=make(Gate,2)}", "FIELD_WRITE"},
		{"aggregate", ";func(p *P)configure(){*p=P{}}", "AGGREGATE_OR_BASE_WRITE"},
		{"rebind", ";func(p *P)configure(){p=&P{};_=p}", "AGGREGATE_OR_BASE_WRITE"},
		{"address", ";func(p *P)configure(){_ = &(p.tokens)}", "FIELD_ADDRESS_EXPOSED"},
		{"argument", ";func consume(*P){};func(p *P)configure(){consume(p)}", "RECEIVER_ARGUMENT_EXPOSURE"},
	} {
		t.Run(c.name, func(t *testing.T) {
			code := identityFixture + c.extra
			r, err := Observe(request(code))
			if err != nil {
				t.Fatal(err)
			}
			if len(r.ChannelIdentity) != 1 {
				t.Fatal("audit missing")
			}
			a := r.ChannelIdentity[0]
			if a.BaseLink != "SAME_BASE_OBJECT" || a.State != "INSTANCE_IDENTITY_UNRESOLVED" || a.Authority != "DECLARATION_SCOPED_AUDIT" {
				t.Fatal("identity promoted")
			}
			found := c.kind == "" && len(a.Events) == 0
			for _, e := range a.Events {
				if e.Kind == c.kind {
					found = true
				}
				if e.Location.ContentSHA256 != Digest(code) || e.Location.EndOffset <= e.Location.Offset {
					t.Fatal("event unbound")
				}
			}
			if !found {
				t.Fatalf("event %s missing: %+v", c.kind, a.Events)
			}
		})
	}
}
func TestChannelBasesCapturesAndDistinctField(t *testing.T) {
	code := strings.Replace(identityFixture, "func(p *P)run(){p.tokens<-1;defer p.tokens.release()}", "func(p *P)run(){p.tokens<-1;go func(){defer p.tokens.release()}()}", 1)
	r, err := Observe(request(code))
	if err != nil || len(r.ChannelIdentity) != 1 {
		t.Fatal(err)
	}
	if r.ChannelIdentity[0].BaseLink != "SAME_BASE_OBJECT" || len(r.ChannelIdentity[0].Events) != 1 || r.ChannelIdentity[0].Events[0].Kind != "RECEIVER_CAPTURE" {
		t.Fatal("capture identity lost")
	}
	code = strings.Replace(identityFixture, "func(p *P)run(){p.tokens<-1;defer p.tokens.release()}", "func(p *P)run(){p.tokens<-1};func(q *P)work(){defer q.tokens.release()}", 1)
	r, err = Observe(request(code))
	if err != nil || r.ChannelIdentity[0].BaseLink != "DISTINCT_BASE_OBJECTS" {
		t.Fatal("distinct declarations equated")
	}
	code = strings.Replace(identityFixture, "defer p.tokens.release()", "defer p.other.release()", 1)
	r, err = Observe(request(code))
	if err != nil || len(r.ChannelIdentity) != 0 {
		t.Fatal("different channel fields equated")
	}
	code = identityFixture + ";var bad=missing"
	r, err = Observe(request(code))
	if err != nil || len(r.ChannelIdentity) != 0 {
		t.Fatal("invalid package audited")
	}
	input := request(identityFixture)
	extra := "package p;func(p *P)configure(){p.tokens=make(Gate,2)}"
	input.Files = append(input.Files, SourceFile{Path: "other.go", Content: extra, SHA256: Digest(extra)})
	r, err = Observe(input)
	if err != nil {
		t.Fatal(err)
	}
	input.Files[0], input.Files[1] = input.Files[1], input.Files[0]
	again, err := Observe(input)
	if err != nil || !reflect.DeepEqual(r, again) {
		t.Fatal("audit order changed")
	}
}

func TestUnknownChannelBaseAndAuditLimit(t *testing.T) {
	code := strings.Replace(identityFixture, "func(p *P)run(){p.tokens<-1;defer p.tokens.release()}", "func get()*P{return nil};func run(){get().tokens<-1;defer get().tokens.release()}", 1)
	r, err := Observe(request(code))
	if err != nil || r.ChannelIdentity[0].BaseLink != "UNRESOLVED_BASE_EXPRESSION" {
		t.Fatal("complex base assumed identical")
	}
	code = identityFixture + ";func(p *P)configure(){" + strings.Repeat("p.tokens=nil;", 10001) + "}"
	if _, err = Observe(request(code)); err == nil {
		t.Fatal("audit overflow accepted")
	}
}
