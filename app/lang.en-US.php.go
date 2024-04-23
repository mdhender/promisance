// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/mdhender/promisance/app/cerr"
	"html/template"
	"log"
	"strings"
)

type LanguageManager_t struct {
	DefaultCode string
	DefaultMap  map[string]string
	enUS        map[string]string
}

func NewLanguageManager(code string) (*LanguageManager_t, error) {
	lm := &LanguageManager_t{
		enUS: map[string]string{},
	}
	switch code {
	case "en-US":
		lm.DefaultCode, lm.DefaultMap = code, lm.enUS
		for k, v := range lang_en_US {
			lm.enUS[k] = v
		}
		for k, v := range lang_en_US_admin {
			lm.enUS[k] = v
		}
		for k, v := range lang_en_US_common {
			lm.enUS[k] = v
		}
		for k, v := range lang_en_US_page {
			lm.enUS[k] = v
		}
	default:
		return nil, cerr.ErrUnknownLanguage
	}
	log.Printf("language: code %q: messages %d\n", lm.DefaultCode, len(lm.DefaultMap))
	return lm, nil
}

func (lm *LanguageManager_t) Printf(msg string, args ...any) string {
	xlat, ok := lm.DefaultMap[msg]
	if !ok {
		return fmt.Sprintf("%s %v\n", msg, args)
	}
	if strings.Contains(xlat, "%1$") {
		log.Printf("todo: lang: msg %q: update percent codes\n", msg)
		return fmt.Sprintf("%q %q %v\n", msg, xlat, args)
	}
	return fmt.Sprintf(xlat, args...)
}

func (lm *LanguageManager_t) PrintfHTML(msg string, args ...any) template.HTML {
	return template.HTML(lm.Printf(msg, args...))
}

var (
	lang_en_US = map[string]string{
		// Display name for language (within Preferences)
		`LANG_ID`: `English (United States)`,
		// Language code as reported in <html> tag
		`LANG_CODE`: `en-US`,
		// Text direction as reported in <html> tag
		`LANG_DIR`: `ltr`,
		// Location of additional strings for this language (within lang directory)
		// as well as guide pages (in "guide/" directory within this path)
		`LANG_PATH`: `en-US/`,
		// Various language-dependent string formatting routines
		`FUNC_FORMAT_LIST`:     `en_us_format_list`,
		`FUNC_FORMAT_LABEL`:    `en_us_format_label`,
		`FUNC_FORMAT_NUMBER`:   `en_us_format_number`,
		`FUNC_FORMAT_PRENUM`:   `en_us_format_prenum`,
		`FUNC_FORMAT_MONEY`:    `en_us_format_money`,
		`FUNC_FORMAT_PERCENT`:  `en_us_format_percent`,
		`FUNC_FORMAT_DURATION`: `en_us_format_duration`,
		`FUNC_UNFORMAT_NUMBER`: `en_us_unformat_number`,
		`FUNC_TRUNCATE`:        `en_us_truncate`,
	}
)
