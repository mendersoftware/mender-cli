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
	"fmt"
	"strings"
)

type Column struct {
	Header string
	Scope  string
	Name   string
}

func splitScope(s string) (scope, name string) {
	if i := strings.Index(s, "/"); i > 0 {
		for _, sc := range Scopes {
			if s[:i] == sc {
				return sc, s[i+1:]
			}
		}
	}
	return "", s
}

func ParseColumns(specs []string) ([]Column, error) {
	cols := make([]Column, 0, len(specs))
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		scope, name := splitScope(spec)
		if name == "" {
			return nil, fmt.Errorf("invalid column %q: missing attribute name", spec)
		}
		cols = append(cols, Column{Header: spec, Scope: scope, Name: name})
	}
	return cols, nil
}

func SelectAttributes(cols []Column, sort []SortCriteria) []SelectAttribute {
	if len(cols) == 0 {
		return nil
	}
	var sel []SelectAttribute
	seen := map[SelectAttribute]bool{}
	add := func(a SelectAttribute) {
		if !seen[a] {
			seen[a] = true
			sel = append(sel, a)
		}
	}
	for _, col := range cols {
		if col.Scope != "" {
			add(SelectAttribute{Scope: col.Scope, Attribute: col.Name})
			continue
		}
		for _, scope := range Scopes {
			add(SelectAttribute{Scope: scope, Attribute: col.Name})
		}
	}
	for _, s := range sort {
		add(SelectAttribute{Scope: s.Scope, Attribute: s.Attribute})
	}
	return sel
}

type filterOperator struct {
	token string
	typ   string
}

var filterOperators = []filterOperator{
	{"=nin:", "$nin"},
	{"=in:", "$in"},
	{"=~", "$regex"},
	{"!=", "$ne"},
	{">=", "$gte"},
	{"<=", "$lte"},
	{"=", "$eq"},
	{">", "$gt"},
	{"<", "$lt"},
}

func ParseFilter(expr string) (FilterPredicate, error) {
	if !strings.ContainsAny(expr, "=<>") && strings.HasSuffix(expr, "?") {
		exists := !strings.HasPrefix(expr, "!")
		scope, name := splitScope(strings.TrimSuffix(strings.TrimPrefix(expr, "!"), "?"))
		if name == "" {
			return FilterPredicate{}, fmt.Errorf("invalid filter %q: missing attribute name", expr)
		}
		return FilterPredicate{
			Scope:     scopeOrDefault(scope),
			Attribute: name,
			Type:      "$exists",
			Value:     exists,
		}, nil
	}

	for i := 0; i < len(expr); i++ {
		for _, op := range filterOperators {
			if !strings.HasPrefix(expr[i:], op.token) {
				continue
			}
			scope, name := splitScope(expr[:i])
			if name == "" {
				return FilterPredicate{}, fmt.Errorf(
					"invalid filter %q: missing attribute name", expr)
			}
			pred := FilterPredicate{
				Scope:     scopeOrDefault(scope),
				Attribute: name,
				Type:      op.typ,
			}
			value := expr[i+len(op.token):]
			if op.typ == "$in" || op.typ == "$nin" {
				pred.Value = strings.Split(value, ",")
			} else {
				pred.Value = value
			}
			return pred, nil
		}
	}
	return FilterPredicate{}, fmt.Errorf(
		"invalid filter %q: expected [scope/]name<op>value with op one of "+
			"=, !=, =~, >, >=, <, <=, =in:, =nin: or [!][scope/]name?", expr)
}

func ParseSort(specs []string) ([]SortCriteria, error) {
	sort := make([]SortCriteria, 0, len(specs))
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		attr, order := spec, "asc"
		if i := strings.LastIndex(spec, ":"); i >= 0 {
			switch spec[i+1:] {
			case "asc", "desc":
				attr, order = spec[:i], spec[i+1:]
			}
		}
		scope, name := splitScope(attr)
		if name == "" {
			return nil, fmt.Errorf("invalid sort %q: missing attribute name", spec)
		}
		sort = append(sort, SortCriteria{
			Scope:     scopeOrDefault(scope),
			Attribute: name,
			Order:     order,
		})
	}
	return sort, nil
}

func StableSort(sort []SortCriteria) []SortCriteria {
	for _, key := range []SortCriteria{
		{Scope: ScopeSystem, Attribute: "created_ts", Order: "asc"},
		{Scope: ScopeIdentity, Attribute: "id", Order: "asc"},
	} {
		found := false
		for _, s := range sort {
			if s.Scope == key.Scope && s.Attribute == key.Attribute {
				found = true
				break
			}
		}
		if !found {
			sort = append(sort, key)
		}
	}
	return sort
}

func scopeOrDefault(scope string) string {
	if scope == "" {
		return ScopeInventory
	}
	return scope
}
