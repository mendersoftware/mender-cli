// Copyright 2026 Northern.tech AS
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package inventory

import (
	"reflect"
	"testing"
)

func TestParseFilter(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		expr string
		want FilterPredicate
		err  bool
	}{
		"eq default scope": {
			expr: "device_type=precision2",
			want: FilterPredicate{"inventory", "device_type", "$eq", "precision2"},
		},
		"eq with scope": {
			expr: "identity/status=accepted",
			want: FilterPredicate{"identity", "status", "$eq", "accepted"},
		},
		"unknown scope is part of the name": {
			expr: "foo/bar=baz",
			want: FilterPredicate{"inventory", "foo/bar", "$eq", "baz"},
		},
		"name with spaces": {
			expr: "MAC address=00:11",
			want: FilterPredicate{"inventory", "MAC address", "$eq", "00:11"},
		},
		"value with operators": {
			expr: "kernel=Linux version 5.10 >= 4",
			want: FilterPredicate{"inventory", "kernel", "$eq", "Linux version 5.10 >= 4"},
		},
		"empty value": {
			expr: "hostname=",
			want: FilterPredicate{"inventory", "hostname", "$eq", ""},
		},
		"ne": {
			expr: "artifact_name!=release-1",
			want: FilterPredicate{"inventory", "artifact_name", "$ne", "release-1"},
		},
		"regex": {
			expr: "hostname=~pos-[0-9]+",
			want: FilterPredicate{"inventory", "hostname", "$regex", "pos-[0-9]+"},
		},
		"gt": {
			expr: "system/check_in_time>2026-01-01",
			want: FilterPredicate{"system", "check_in_time", "$gt", "2026-01-01"},
		},
		"gte": {
			expr: "system/check_in_time>=2026-01-01",
			want: FilterPredicate{"system", "check_in_time", "$gte", "2026-01-01"},
		},
		"lt": {
			expr: "a<b",
			want: FilterPredicate{"inventory", "a", "$lt", "b"},
		},
		"lte": {
			expr: "a<=b",
			want: FilterPredicate{"inventory", "a", "$lte", "b"},
		},
		"in": {
			expr: "device_type=in:precision2,spectrum",
			want: FilterPredicate{"inventory", "device_type", "$in",
				[]string{"precision2", "spectrum"}},
		},
		"nin": {
			expr: "tags/env=nin:dev",
			want: FilterPredicate{"tags", "env", "$nin", []string{"dev"}},
		},
		"exists": {
			expr: "monitor/alerts?",
			want: FilterPredicate{"monitor", "alerts", "$exists", true},
		},
		"not exists": {
			expr: "!hostname?",
			want: FilterPredicate{"inventory", "hostname", "$exists", false},
		},
		"value ending with question mark": {
			expr: "note=why?",
			want: FilterPredicate{"inventory", "note", "$eq", "why?"},
		},
		"no operator": {
			expr: "hostname",
			err:  true,
		},
		"missing name": {
			expr: "=foo",
			err:  true,
		},
		"missing name after scope": {
			expr: "inventory/=foo",
			err:  true,
		},
		"missing name exists": {
			expr: "!?",
			err:  true,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseFilter(tc.expr)
			if tc.err {
				if err == nil {
					t.Fatalf("expected an error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestParseSort(t *testing.T) {
	t.Parallel()
	got, err := ParseSort([]string{"artifact_name", "system/check_in_time:desc", "tags/x:y:asc"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := []SortCriteria{
		{"inventory", "artifact_name", "asc"},
		{"system", "check_in_time", "desc"},
		{"tags", "x:y", "asc"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
	if _, err := ParseSort([]string{":desc"}); err == nil {
		t.Error("expected an error for a missing attribute name")
	}
}

func TestColumns(t *testing.T) {
	t.Parallel()
	cols, err := ParseColumns([]string{"hostname", "identity/mac", " system/group "})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	wantCols := []Column{
		{"hostname", "", "hostname"},
		{"identity/mac", "identity", "mac"},
		{"system/group", "system", "group"},
	}
	if !reflect.DeepEqual(cols, wantCols) {
		t.Errorf("got %#v, want %#v", cols, wantCols)
	}

	sel := SelectAttributes(
		append(cols, Column{"inventory/hostname", "inventory", "hostname"}),
		[]SortCriteria{{"identity", "mac", "asc"}, {"system", "created_ts", "asc"}})
	wantSel := []SelectAttribute{
		{"inventory", "hostname"},
		{"identity", "hostname"},
		{"system", "hostname"},
		{"tags", "hostname"},
		{"monitor", "hostname"},
		{"identity", "mac"},
		{"system", "group"},
		{"system", "created_ts"},
	}
	if !reflect.DeepEqual(sel, wantSel) {
		t.Errorf("got %#v, want %#v", sel, wantSel)
	}
	if sel := SelectAttributes(nil, []SortCriteria{{"system", "created_ts", "asc"}}); sel != nil {
		t.Errorf("without columns all attributes should be returned, got %#v", sel)
	}

	if _, err := ParseColumns([]string{"tags/"}); err == nil {
		t.Error("expected an error for a missing attribute name")
	}
}

func TestStableSort(t *testing.T) {
	t.Parallel()
	createdTs := SortCriteria{"system", "created_ts", "asc"}
	id := SortCriteria{"identity", "id", "asc"}
	got := StableSort(nil)
	if !reflect.DeepEqual(got, []SortCriteria{createdTs, id}) {
		t.Errorf("got %#v", got)
	}
	sort := []SortCriteria{{"inventory", "artifact_name", "asc"}}
	got = StableSort(sort)
	if !reflect.DeepEqual(got, append(sort, createdTs, id)) {
		t.Errorf("got %#v", got)
	}
	sort = []SortCriteria{{"system", "created_ts", "desc"}, {"inventory", "a", "asc"}}
	got = StableSort(sort)
	if !reflect.DeepEqual(got, append(sort, id)) {
		t.Errorf("an existing created_ts sort should be kept, got %#v", got)
	}
	sort = []SortCriteria{{"identity", "id", "desc"}}
	got = StableSort(sort)
	if !reflect.DeepEqual(got, []SortCriteria{sort[0], createdTs}) {
		t.Errorf("an existing id sort should be kept, got %#v", got)
	}
}
