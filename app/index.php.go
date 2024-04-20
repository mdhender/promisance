// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"log"
	"path/filepath"
)

func (p *PHP) index_php() error {
	// date_default_timezone_set('UTC')
	// ignore_user_abort(1)
	// set_time_limit(0)

	// don't allow CLI invocation

	// prevent recursion
	if p.constants.IN_GAME {
		p.die("Access denied (recursion)")
	}
	p.constants.IN_GAME = true
	p.constants.PROM_BASEDIR = "."

	// Don't allow accessing anything while setup.php is being run
	if p.file_exists(filepath.Join(p.constants.PROM_BASEDIR, "setup.php")) {
		p.die("Access denied (setup running)")
	}
	log.Printf("todo: the setup check is not implemented correctly\n")

	p.require_once(filepath.Join("config.php"))
	p.require_once(filepath.Join("includes/constants.php"))
	p.require_once(filepath.Join("includes/database.php"))
	p.require_once(filepath.Join("includes/language.php"))
	p.require_once(filepath.Join("includes/html.php"))
	p.require_once(filepath.Join("includes/logging.php"))
	p.require_once(filepath.Join("includes/misc.php"))
	p.require_once(filepath.Join("includes/news.php"))
	p.require_once(filepath.Join("includes/permissions.php"))
	p.require_once(filepath.Join("classes/prom_vars.php"))
	p.require_once(filepath.Join("classes/prom_user.php"))
	p.require_once(filepath.Join("classes/prom_empire.php"))
	p.require_once(filepath.Join("classes/prom_clan.php"))
	p.require_once(filepath.Join("classes/prom_session.php"))
	p.require_once(filepath.Join("classes/prom_turns.php"))
	p.require_once(filepath.Join("includes/auth.php"))

	return fmt.Errorf("not implemented")
}
