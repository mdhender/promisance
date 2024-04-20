// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import "fmt"

func (p *PHP) install_setup_php() error {
	// don't allow CLI invocation
	if p.constants.REQUEST_URI == nil {
		p.die("Access denied (cli forbidden)")
	}

	// prevent recursion
	if !p.constants.IN_GAME {
		p.die("Access denied (recursion)")
	}

	p.constants.IN_GAME = true
	p.constants.IN_SETUP = true

	return fmt.Errorf("not implemented")
}
