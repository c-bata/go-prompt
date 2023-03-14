//go:build !windows
// +build !windows

package prompt

import (
	"reflect"
	"testing"
)

func emptyCompleter(in Document) []Suggest {
	return []Suggest{}
}

func TestFormatCompletion(t *testing.T) {
	scenarioTable := []struct {
		scenario      string
		completions   []Suggest
		prefix        string
		suffix        string
		expected      []Suggest
		maxWidth      int
		expectedWidth int
	}{
		{
			scenario: "",
			completions: []Suggest{
				{Text: "select"},
				{Text: "from"},
				{Text: "insert"},
				{Text: "where"},
			},
			prefix: " ",
			suffix: " ",
			expected: []Suggest{
				{Text: " select "},
				{Text: " from   "},
				{Text: " insert "},
				{Text: " where  "},
			},
			maxWidth:      20,
			expectedWidth: 8,
		},
		{
			scenario: "",
			completions: []Suggest{
				{Text: "select", Description: "select description"},
				{Text: "from", Description: "from description"},
				{Text: "insert", Description: "insert description"},
				{Text: "where", Description: "where description"},
			},
			prefix: " ",
			suffix: " ",
			expected: []Suggest{
				{Text: " select ", Description: " select description "},
				{Text: " from   ", Description: " from description   "},
				{Text: " insert ", Description: " insert description "},
				{Text: " where  ", Description: " where description  "},
			},
			maxWidth:      40,
			expectedWidth: 28,
		},
	}

	for _, s := range scenarioTable {
		ac, width := formatSuggestions(s.completions, s.maxWidth)
		if !reflect.DeepEqual(ac, s.expected) {
			t.Errorf("Should be %#v, but got %#v", s.expected, ac)
		}
		if width != s.expectedWidth {
			t.Errorf("Should be %#v, but got %#v", s.expectedWidth, width)
		}
	}
}

func TestBreakLineCallback(t *testing.T) {
	var i int
	r := &Render{
		prefix:                       "> ",
		out:                          &PosixWriter{},
		livePrefixCallback:           func() (string, bool) { return "", false },
		prefixTextColor:              Blue,
		prefixBGColor:                DefaultColor,
		inputTextColor:               DefaultColor,
		inputBGColor:                 DefaultColor,
		previewSuggestionTextColor:   Green,
		previewSuggestionBGColor:     DefaultColor,
		suggestionTextColor:          White,
		suggestionBGColor:            Cyan,
		selectedSuggestionTextColor:  Black,
		selectedSuggestionBGColor:    Turquoise,
		descriptionTextColor:         Black,
		descriptionBGColor:           Turquoise,
		selectedDescriptionTextColor: White,
		selectedDescriptionBGColor:   Cyan,
		scrollbarThumbColor:          DarkGray,
		scrollbarBGColor:             Cyan,
		col:                          1,
	}
	b := NewBuffer()
	l := NewLexer()
	r.BreakLine(b, l)

	if i != 0 {
		t.Errorf("i should initially be 0, before applying a break line callback")
	}

	r.breakLineCallback = func(doc *Document) {
		i++
	}
	r.BreakLine(b, l)
	r.BreakLine(b, l)
	r.BreakLine(b, l)

	if i != 3 {
		t.Errorf("BreakLine callback not called, i should be 3")
	}
}

func TestLinesToTracebackRender(t *testing.T) {
	scenarios := []struct {
		previousText     string
		nextText         string
		linesToTraceBack int
		lastKey          Key
	}{
		{previousText: "select..", nextText: "", linesToTraceBack: 0, lastKey: Enter},
		{previousText: "select.. \n from.. \n where..", nextText: "", linesToTraceBack: 0, lastKey: Enter},
		{previousText: "select.. \n from.. \n where..", nextText: "select..", linesToTraceBack: 2, lastKey: Tab},
		{previousText: "select.. \n from.. \n where..", nextText: "select.. \n from.. \n where field = 2", linesToTraceBack: 2, lastKey: Tab},
		{previousText: "select.. \n from.. \n where..", nextText: "select.. \n from.. \n where field = 2", linesToTraceBack: 2, lastKey: Right},
		{previousText: "select.. \n from.. ", nextText: "previous statement", linesToTraceBack: 1, lastKey: Up},
		{previousText: "select.. \n from.. ", nextText: "next statement", linesToTraceBack: 1, lastKey: Down},
		{previousText: "select.. \n from.. ", nextText: "next statement", linesToTraceBack: 1, lastKey: ControlDown},
		{previousText: "select.. \n from.. ", nextText: "", linesToTraceBack: 1, lastKey: Down},
		{previousText: "select.. \n from.. ", nextText: "", linesToTraceBack: 1, lastKey: ControlDown},
	}

	r := &Render{
		prefix:                       "> ",
		out:                          &PosixWriter{},
		livePrefixCallback:           func() (string, bool) { return "", false },
		prefixTextColor:              Blue,
		prefixBGColor:                DefaultColor,
		inputTextColor:               DefaultColor,
		inputBGColor:                 DefaultColor,
		previewSuggestionTextColor:   Green,
		previewSuggestionBGColor:     DefaultColor,
		suggestionTextColor:          White,
		suggestionBGColor:            Cyan,
		selectedSuggestionTextColor:  Black,
		selectedSuggestionBGColor:    Turquoise,
		descriptionTextColor:         Black,
		descriptionBGColor:           Turquoise,
		selectedDescriptionTextColor: White,
		selectedDescriptionBGColor:   Cyan,
		scrollbarThumbColor:          DarkGray,
		scrollbarBGColor:             Cyan,
		col:                          1,
	}

	for _, s := range scenarios {
		b := NewBuffer()
		b.InsertText(s.nextText, false, true)
		l := NewLexer()

		tracedBackLines := r.Render(b, s.previousText, s.lastKey, NewCompletionManager(emptyCompleter, 0), l)

		if tracedBackLines != s.linesToTraceBack {
			t.Errorf("Should've traced back %d lines before rendering, but got %d", s.linesToTraceBack, tracedBackLines)
		}
	}
}
