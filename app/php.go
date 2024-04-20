// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/syyongx/php2go"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	FALSE = false
	TRUE  = true
)

// PHP implements a wrapper for the manually converted PHP code.
type PHP struct {
	// constant per instance, not per server
	constants struct {
		CUR_TIME     time.Time
		IN_GAME      bool
		IN_SETUP     bool
		PROM_BASEDIR string
		REQUEST_URI  *url.URL
	}

	globals struct {
		banners         []banner_t
		cur_lang        string
		empire_defaults map[string]int
		lang_cache      map[string]any
		notices         string
		required_vars   []string
		sequences       map[string]string
		tables          map[string]string
		timezones       map[int]string
		styles          map[string]css_file_t
		world           *world_t
	}
	required map[string]bool
}
type banner_t struct {
	key   string
	value string
}
type css_file_t struct {
	file string
	name string
}
type world_t struct {
	lotto_current_jackpot   int
	lotto_yesterday_jackpot int
	lotto_last_picked       int
	lotto_last_winner       int
	lotto_jackpot_increase  int
	round_time_begin        time.Time
	round_time_closing      time.Time
	round_time_end          time.Time
	turns_next              time.Time
	turns_next_hourly       time.Time
	turns_next_daily        time.Time
}

func newInstance(w http.ResponseWriter, r *http.Request) (*PHP, error) {
	p := &PHP{
		required: make(map[string]bool),
	}
	if w != nil && r != nil {
		p.constants.REQUEST_URI = r.URL
	}
	return p, nil
}

type array struct {
	data []*arnode
}
type arnode struct {
	key   string
	value any
}

//func (p *PHP) mkarray(kv ...any) *array {
//	a := &array{}
//	if len(kv)%2 != 0 {
//		panic(fmt.Sprintf("php: array: %d args", len(kv)))
//	}
//	for i := 0; i < len(kv); i += 2 {
//		key, ok := kv[i].(string)
//		if !ok {
//			panic(fmt.Sprintf("php: array: key is not a string: %s", kv[i]))
//		}
//		val := kv[i+1]
//		a.data = append(a.data, &arnode{key: key, value: val})
//	}
//	return a
//}

func (p *PHP) die(format string, args ...any) {
	log.Printf(format, args...)
	panic("php: die")
}

// return trues if the file (or directory) exists
func (p *PHP) file_exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func (p *PHP) getcwd() string {
	wd, _ := php2go.Getcwd()
	return wd
}

func (p *PHP) getString(key string) string {
	panic("php: getString: not implemented")
}

func (p *PHP) require_once(path string) {
	// return immediately if we have already loaded the script
	if p.required[path] {
		return
	}
	var err error
	switch path {
	case "config.php":
		err = p.config_php()
	case "classes/prom_clan.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "classes/prom_empire.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "classes/prom_session.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "classes/prom_turns.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "classes/prom_user.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "classes/prom_vars.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "includes/PasswordHash.php":
		err = p.includes_PasswordHash_php()
	case "includes/auth.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "includes/constants.php":
		err = p.includes_constants_php()
	case "includes/database.php":
		err = p.includes_database_php()
	case "includes/html.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "includes/language.php":
		err = p.includes_language_php()
	case "includes/logging.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "includes/misc.php":
		err = p.includes_misc_php()
	case "includes/news.php":
		err = fmt.Errorf("assert(path != %q)", path)
	case "includes/permissions.php":
		err = fmt.Errorf("assert(path != %q)", path)
	default:
		err = fmt.Errorf("assert(path != %q)", path)
	}
	if err != nil {
		log.Printf("php: %s: %v\n", path, err)
		panic(fmt.Sprintf("php: require %q failed", path))
	}
}
