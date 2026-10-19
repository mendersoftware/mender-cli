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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

const (
	FormatText = "text"
	FormatJSON = "json"
	FormatCSV  = "csv"
	FormatTSV  = "tsv"
)

var Formats = []string{FormatText, FormatJSON, FormatCSV, FormatTSV}

func ValidateFormat(format string) error {
	for _, f := range Formats {
		if format == f {
			return nil
		}
	}
	return fmt.Errorf("invalid output format %q, must be one of: %s",
		format, strings.Join(Formats, ", "))
}

const arraySeparator = ";"

func (d *Device) Lookup(col Column) (Attribute, bool) {
	scopes := Scopes
	if col.Scope != "" {
		scopes = []string{col.Scope}
	}
	for _, scope := range scopes {
		for _, attr := range d.Attributes {
			if attr.Scope == scope && attr.Name == col.Name {
				return attr, true
			}
		}
	}
	return Attribute{}, false
}

func FormatValue(v interface{}) string {
	switch v := v.(type) {
	case nil:
		return ""
	case string:
		return v
	case json.Number:
		return v.String()
	case bool:
		return strconv.FormatBool(v)
	case []interface{}:
		items := make([]string, len(v))
		for i, item := range v {
			items[i] = FormatValue(item)
		}
		return strings.Join(items, arraySeparator)
	default:
		return fmt.Sprint(v)
	}
}

func scopeRank(scope string) int {
	for i, s := range Scopes {
		if s == scope {
			return i
		}
	}
	return len(Scopes)
}

func sortAttributes(attrs []Attribute) {
	sort.SliceStable(attrs, func(i, j int) bool {
		ri, rj := scopeRank(attrs[i].Scope), scopeRank(attrs[j].Scope)
		if ri != rj {
			return ri < rj
		}
		if attrs[i].Scope != attrs[j].Scope {
			return attrs[i].Scope < attrs[j].Scope
		}
		return attrs[i].Name < attrs[j].Name
	})
}

func AllColumns(devices []Device) []Column {
	var attrs []Attribute
	seen := map[Column]bool{}
	for _, dev := range devices {
		for _, attr := range dev.Attributes {
			col := Column{Scope: attr.Scope, Name: attr.Name}
			if !seen[col] {
				seen[col] = true
				attrs = append(attrs, attr)
			}
		}
	}
	sortAttributes(attrs)
	cols := make([]Column, len(attrs))
	for i, attr := range attrs {
		cols[i] = Column{Header: attr.Scope + "/" + attr.Name, Scope: attr.Scope, Name: attr.Name}
	}
	return cols
}

type Writer interface {
	WritePage(devices []Device, raw []json.RawMessage) error
	Close() error
}

func NewWriter(
	w io.Writer,
	format string,
	cols []Column,
	detailLevel int,
	raw bool,
) (Writer, error) {
	if err := ValidateFormat(format); err != nil {
		return nil, err
	}
	switch {
	case raw:
		return &jsonArrayWriter{w: w, element: rawElement}, nil
	case format == FormatText:
		return &textWriter{w: w, cols: cols, detailLevel: detailLevel}, nil
	case format == FormatJSON:
		jw := &jsonArrayWriter{w: w}
		jw.element = func(dev Device, _ json.RawMessage) ([]byte, error) {
			return json.MarshalIndent(flatten(dev, jw.cols), "  ", "  ")
		}
		return &bufferingWriter{cols: cols, Writer: jw, setCols: func(c []Column) {
			jw.cols = c
		}}, nil
	default:
		sep := ','
		if format == FormatTSV {
			sep = '\t'
		}
		cw := &csvWriter{w: csv.NewWriter(w)}
		cw.w.Comma = sep
		return &bufferingWriter{cols: cols, Writer: cw, setCols: func(c []Column) {
			cw.cols = c
		}}, nil
	}
}

func Write(w io.Writer, format string, devices []Device, cols []Column, detailLevel int) error {
	return writeAll(w, format, devices, nil, cols, detailLevel, false)
}

func WriteRaw(w io.Writer, raw []json.RawMessage) error {
	return writeAll(w, FormatJSON, make([]Device, len(raw)), raw, nil, 0, true)
}

func writeAll(
	w io.Writer,
	format string,
	devices []Device,
	raw []json.RawMessage,
	cols []Column,
	detailLevel int,
	rawMode bool,
) error {
	wr, err := NewWriter(w, format, cols, detailLevel, rawMode)
	if err != nil {
		return err
	}
	if err := wr.WritePage(devices, raw); err != nil {
		return err
	}
	return wr.Close()
}

type bufferingWriter struct {
	Writer
	cols     []Column
	setCols  func([]Column)
	started  bool
	buffered []Device
}

func (b *bufferingWriter) WritePage(devices []Device, raw []json.RawMessage) error {
	if len(b.cols) == 0 {
		b.buffered = append(b.buffered, devices...)
		return nil
	}
	if !b.started {
		b.started = true
		b.setCols(b.cols)
	}
	return b.Writer.WritePage(devices, raw)
}

func (b *bufferingWriter) Close() error {
	if len(b.cols) == 0 {
		b.setCols(AllColumns(b.buffered))
		if err := b.Writer.WritePage(b.buffered, nil); err != nil {
			return err
		}
	} else if !b.started {
		b.setCols(b.cols)
	}
	return b.Writer.Close()
}

type textWriter struct {
	w           io.Writer
	cols        []Column
	detailLevel int
}

func (t *textWriter) WritePage(devices []Device, _ []json.RawMessage) error {
	for _, dev := range devices {
		fmt.Fprintf(t.w, "ID: %s\n", dev.ID)
		if t.detailLevel >= 1 && dev.UpdatedTs != "" {
			fmt.Fprintf(t.w, "UpdatedTs: %s\n", dev.UpdatedTs)
		}
		var attrs []Attribute
		var labels []string
		if len(t.cols) > 0 {
			for _, col := range t.cols {
				attr, _ := dev.Lookup(col)
				attrs = append(attrs, attr)
				labels = append(labels, col.Header)
			}
		} else {
			attrs = append(attrs, dev.Attributes...)
			sortAttributes(attrs)
			for _, attr := range attrs {
				labels = append(labels, attr.Scope+"/"+attr.Name)
			}
		}
		for i, attr := range attrs {
			fmt.Fprintf(t.w, "%s: %s\n", labels[i], FormatValue(attr.Value))
			if t.detailLevel >= 1 && attr.Description != "" {
				fmt.Fprintf(t.w, "  Description: %s\n", attr.Description)
			}
			if t.detailLevel >= 1 && attr.Timestamp != "" {
				fmt.Fprintf(t.w, "  Timestamp: %s\n", attr.Timestamp)
			}
		}
		_, err := fmt.Fprintf(t.w,
			"--------------------------------------------------------------------------------\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *textWriter) Close() error {
	return nil
}

func flatten(dev Device, cols []Column) map[string]interface{} {
	obj := map[string]interface{}{"id": dev.ID}
	for _, col := range cols {
		attr, _ := dev.Lookup(col)
		obj[col.Header] = attr.Value
	}
	return obj
}

func rawElement(_ Device, raw json.RawMessage) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "  ", "  "); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type jsonArrayWriter struct {
	w       io.Writer
	cols    []Column
	element func(dev Device, raw json.RawMessage) ([]byte, error)
	count   int
}

func (j *jsonArrayWriter) WritePage(devices []Device, raw []json.RawMessage) error {
	for i, dev := range devices {
		var r json.RawMessage
		if raw != nil {
			r = raw[i]
		}
		data, err := j.element(dev, r)
		if err != nil {
			return err
		}
		prefix := ",\n  "
		if j.count == 0 {
			prefix = "[\n  "
		}
		j.count++
		if _, err := io.WriteString(j.w, prefix); err != nil {
			return err
		}
		if _, err := j.w.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func (j *jsonArrayWriter) Close() error {
	end := "\n]\n"
	if j.count == 0 {
		end = "[]\n"
	}
	_, err := io.WriteString(j.w, end)
	return err
}

type csvWriter struct {
	w             *csv.Writer
	cols          []Column
	headerWritten bool
}

func (c *csvWriter) writeHeader() error {
	if c.headerWritten {
		return nil
	}
	c.headerWritten = true
	record := make([]string, len(c.cols)+1)
	record[0] = "id"
	for i, col := range c.cols {
		record[i+1] = col.Header
	}
	return c.w.Write(record)
}

func (c *csvWriter) WritePage(devices []Device, _ []json.RawMessage) error {
	if err := c.writeHeader(); err != nil {
		return err
	}
	record := make([]string, len(c.cols)+1)
	for _, dev := range devices {
		record[0] = dev.ID
		for i, col := range c.cols {
			attr, _ := dev.Lookup(col)
			record[i+1] = FormatValue(attr.Value)
		}
		if err := c.w.Write(record); err != nil {
			return err
		}
	}
	c.w.Flush()
	return c.w.Error()
}

func (c *csvWriter) Close() error {
	if err := c.writeHeader(); err != nil {
		return err
	}
	c.w.Flush()
	return c.w.Error()
}

func WriteFilterableAttributes(w io.Writer, format string, attrs []FilterAttribute) error {
	switch format {
	case FormatText:
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "SCOPE\tNAME\tCOUNT")
		for _, a := range attrs {
			fmt.Fprintf(tw, "%s\t%s\t%d\n", a.Scope, a.Name, a.Count)
		}
		return tw.Flush()
	case FormatJSON:
		if attrs == nil {
			attrs = []FilterAttribute{}
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(attrs)
	case FormatCSV, FormatTSV:
		cw := csv.NewWriter(w)
		if format == FormatTSV {
			cw.Comma = '\t'
		}
		_ = cw.Write([]string{"scope", "name", "count"})
		for _, a := range attrs {
			_ = cw.Write([]string{a.Scope, a.Name, strconv.Itoa(a.Count)})
		}
		cw.Flush()
		return cw.Error()
	}
	return ValidateFormat(format)
}
