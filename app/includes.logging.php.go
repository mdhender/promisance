// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

// Generates a warning message
// If in the setup script, display as formatted XHTML
// If in turns script, display as plain text both to STDOUT (to be logged to disk) and to STDERR (to be emailed to the admin)
// If in a utility script, display as plain text to STDOUT only
// If in-game, store in database (along with other relevant information)
// Specify level = 0 to indicate the file/line where warning() itself was called,
// level = 1 for the caller, level = 2 for the caller's caller, etc.
// If a description is specified, a message will be delivered to the in-game moderator mailbox
func (p *PHP) warning(msg string, level int, desc string) {
	panic("not implemented")
}
