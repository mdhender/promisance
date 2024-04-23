// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import "html/template"

func (s *server) notices(style int, notices []string) template.HTML {
	var result string
	switch style {
	case 2:
		result = `<h4 class="cwarn">`
	case 1:
		result = `<h4>`
	}
	for n, notice := range notices {
		if n > 0 {
			result += `<br />`
		}
		result += notice
	}
	switch style {
	case 1, 2:
		result += `</h4>`
	default:
		result += `<hr />`
	}
	return template.HTML(result)
}
