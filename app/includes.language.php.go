// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import "log"

const (
	// php/includes/language.php

	DURATION_SECONDS = 0
	DURATION_MINUTES = 1
	DURATION_HOURS   = 2
	DURATION_DAYS    = 3
)

func (p *PHP) includes_language_php() error {
	if !p.constants.IN_GAME {
		p.die("Access denied")
	}

	p.globals.cur_lang = DEFAULT_LANGUAGE
	p.globals.lang_cache = map[string]any{}

	log.Printf("todo: consider implementing includes/language\n")

	return nil
}

// Selects which language to use when generating the current page
// Returns TRUE if the language is available, FALSE if not
func (p *PHP) setlanguage(name string) {
	panic("not implemented")
}

// Substitutes the supplied string ID for the matching string in the currently selected language
// and substitutes parameters where appropriate
func (p *PHP) lang(id int) {
	panic("not implemented")
}

// Substitutes the supplied string ID for the matching string in the DEFAULT language
// and substitutes parameters where appropriate
func (p *PHP) def_lang(id int) {
	panic("not implemented")
}

// Substitutes the supplied string ID for the matching string in the BASE language
// and substitutes parameters where appropriate
func (p *PHP) base_lang(id int) {
	panic("not implemented")
}

// Locates a file within a language directory
func (p *PHP) langfile(file string) {
	panic("not implemented")
}

// Checks if text is equal to a string defined in any language pack
// Intended for preventing creation of clans named "None"
func (p *PHP) lang_equals_any(compare string, id int) {
	panic("not implemented")
}

// Determines if the specified string is a defined string ID
// Intended for usage in user input validation (to prevent users from entering a string ID as their nickname, clan title, etc.)
func (p *PHP) lang_isset(id int) {
	panic("not implemented")
}

// pluralize a string, with the (comma-formatted) number substituted into the string if requested
// singular/plural forms may be literal strings or language-specific string IDs
func (p *PHP) plural(num int, sing string, plur string, zero string) {
	panic("not implemented")
}

// Substitutes the supplied string ID for the matching string in the currently selected language
// Intended for retrieving language-specific functions
// Resulting function name is cached for faster access until the language is switched
func (p *PHP) lang_func(id int, default_ string) {
	panic("not implemented")
}

func (p *PHP) lang_unavailable() {
	p.warning("Language-specific function not available", 1, "")
}

// Takes a list and separates it with commas and spaces, including "and" before the last entry
// If there are only 2 entries in the list, no commas are used
func (p *PHP) commaList(list ...string) {
	panic("not implemented")
}

// Formats the specified text as a label to be immediately followed by a value
// Label text may be a string ID
func (p *PHP) label(label, value string) {
	panic("not implemented")
}

// Formats the specified text as an ordinary number with thousands separators
func (p *PHP) number(num int) {
	panic("not implemented")
}

// Formats the specified text as a number with a number-sign prefix or suffix
func (p *PHP) prenum(num int) {
	panic("not implemented")
}

// Formats the specified text as currency
func (p *PHP) money(num int) {
	panic("not implemented")
}

// Formats the specified text as a percentage
func (p *PHP) percent(num, decimal int) {
	panic("not implemented")
}

// Formats a number of seconds as "N days, N hours, N minutes, N seconds"
// Precision controls number of decimal places for last token
func (p *PHP) duration(num, precision, min_level, max_level int) {
	panic("not implemented")
}

// Removes all number formatting from the specified text
func (p *PHP) unformat_number(num string) {
	panic("not implemented")
}

func (p *PHP) truncate(str string, num int) {
	panic("not implemented")
}
