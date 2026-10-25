package repository

import "testing"

func TestRequestListOrderBy(t *testing.T) {
	got := requestListOrderBy(AIRequestListFilter{SortBy: "input", SortDesc: true}, true)
	if got != "r.input_tokens DESC, r.id DESC" {
		t.Fatalf("got %q", got)
	}
	got = requestListOrderBy(AIRequestListFilter{SortBy: "tag", SortDesc: false}, true)
	if got != "r.tag ASC, r.id ASC" {
		t.Fatalf("tag: got %q", got)
	}
	got = requestListOrderBy(AIRequestListFilter{SortBy: "tag", SortDesc: false}, false)
	if got != "r.started_at ASC, r.id ASC" {
		t.Fatalf("tag without column: got %q", got)
	}
	got = requestListOrderBy(AIRequestListFilter{SortBy: "nope", SortDesc: false}, true)
	if got != "r.started_at ASC, r.id ASC" {
		t.Fatalf("default: got %q", got)
	}
}
