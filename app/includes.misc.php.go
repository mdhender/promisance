// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"log"
)

func (p *PHP) includes_misc_php() error {
	if !p.constants.IN_GAME {
		p.die("Access denied")
	}

	p.require_once("includes/PasswordHash.php")

	// Special variables used by functionss in this file

	// Notices accumulated using notice() and displayed using notices()
	p.globals.notices = ""

	log.Printf("todo: includes/misc does not implement code, just stubs functions\n")

	return nil
}

// Sends mail - returns 0 for success, >0 for error
func (p *PHP) prom_mail(to, subj, msg string) {
	panic("not implemented")
	//return !mail($to, $subj, $msg, 'From: '. lang('EMAIL_FROM') .' <'. MAIL_VALIDATE .">\nX-Mailer: PHP/". phpversion());
}

// checks if a form is being submitted via GET or POST
func (p *PHP) isFormPost() {
	panic("not implemented")
	//return ($_SERVER['REQUEST_METHOD'] == 'POST');
}

// Converts non-breaking spaces into normal ones
func (p *PHP) remove_nbsp(str string) {
	panic("not implemented")
	//return str_replace(html_entity_decode('&nbsp;', ENT_COMPAT, 'UTF-8'), ' ', $str);
}

// fetches a field from the posted form, trims whitespace
func (p *PHP) getFormVar(var_, default_ string) {
	panic("not implemented")
	//if (isset($_REQUEST[$var]))
	//	return trim(remove_nbsp($_REQUEST[$var]));
	//else	return $default;
}

// fetches an array from the posted form
func (p *PHP) getFormArr(var_ string, default_ any) {
	panic("not implemented")
	//if (isset($_REQUEST[$var]))
	//	return $_REQUEST[$var];
	//else	return $default;
}

// remove any special punctuation (thousands separators), allow positive integers only
func (p *PHP) fixInputNum(num string) {
	panic("not implemented")
	//$result = floor(unformat_number($num));
	//if ($result > 0)
	//	return $result;
	//else	return 0;
}

// remove any special punctuation (thousands separators), allow positive or negative integers
func (p *PHP) fixInputNumSigned(num string) {
	panic("not implemented")
	//$result = floor(unformat_number($num));
	//return $result;
}

// restrict input to a boolean value
func (p *PHP) fixInputBool(bool_ string) {
	panic("not implemented")
	//if ($bool)
	//	return TRUE;
	//else	return FALSE;
}

// Builds a string "(?,?,?,...,?)" for a given array of parameters
func (p *PHP) sqlArgList(args string, fill ...string) {
	panic("not implemented")
	//return '('. implode(',', array_fill(0, count($args), $fill)) .')';
}

// Returns one of several frequently used lookup tables
func (p *PHP) lookup(id int) {
	panic("not implemented")
	//static $lookup = array(
	//	// unit lists of various types
	//	'list_mil'	=> array('trparm', 'trplnd', 'trpfly', 'trpsea'),
	//	'list_mkt'	=> array('trparm', 'trplnd', 'trpfly', 'trpsea', 'food'),
	//	'list_aid'	=> array('trparm', 'trplnd', 'trpfly', 'trpsea', 'cash', 'runes', 'food'),
	//
	//	// private market property lookups
	//	'pvtmkt_name_id'	=> array('trparm' => 'e_mktarm', 'trplnd' => 'e_mktlnd', 'trpfly' => 'e_mktfly', 'trpsea' => 'e_mktsea', 'food' => 'e_mktfood'),
	//	'pvtmkt_name_limit'	=> array('trparm' => 'e_mktperarm', 'trplnd' => 'e_mktperlnd', 'trpfly' => 'e_mktperfly', 'trpsea' => 'e_mktpersea'),
	//	'pvtmkt_name_cost'	=> array('trparm' => PVTM_TRPARM, 'trplnd' => PVTM_TRPLND, 'trpfly' => PVTM_TRPFLY, 'trpsea' => PVTM_TRPSEA, 'food' => PVTM_FOOD),
	//
	//	// public market property lookups
	//	'pubmkt_name_id'	=> array('trparm' => MARKET_TRPARM, 'trplnd' => MARKET_TRPLND, 'trpfly' => MARKET_TRPFLY, 'trpsea' => MARKET_TRPSEA, 'food' => MARKET_FOOD),
	//	'pubmkt_id_name'	=> array(MARKET_TRPARM => 'trparm', MARKET_TRPLND => 'trplnd', MARKET_TRPFLY => 'trpfly', MARKET_TRPSEA => 'trpsea', MARKET_FOOD => 'food'),
	//	'pubmkt_id_cost'	=> array(MARKET_TRPARM => PVTM_TRPARM, MARKET_TRPLND => PVTM_TRPLND, MARKET_TRPFLY => PVTM_TRPFLY, MARKET_TRPSEA => PVTM_TRPSEA, MARKET_FOOD => PVTM_FOOD),
	//);
	//
	//if (!isset($lookup[$id]))
	//{
	//	warning('Attempted to fetch undefined lookup table '. $id, 1);
	//	return NULL;
	//}
	//
	//$data = $lookup[$id];
	//return $data;
}

// Generates a gaussian random number within a particular range
// Mean defaults to the center of the range
// Standard deviation defaults to 1/6th of the range (such that only 0.3% of values will get clipped)
func (p *PHP) gauss_rand(min, max, dev, mean float64) {
	panic("not implemented")
	//if ($mean == 0)
	//	$mean = ($max + $min) / 2;
	//if ($dev == 0)
	//	$dev = ($max - $min) / 6;
	//
	//$randmax = mt_getrandmax();
	//
	//while (1)
	//{
	//	do
	//	{
	//		$x1 = mt_rand();
	//		$x2 = mt_rand();
	//	} while ($x1 == 0);
	//
	//	$x1 /= $randmax;
	//	$x2 /= $randmax;
	//
	//	// can also change cos() into sin(), but that's not necessary
	//	$y1 = sqrt(-2 * log($x1)) * cos(2 * M_PI * $x2);
	//
	//	$val = $y1 * $dev + $mean;
	//	// if the value is out of range, reroll
	//	if (($val >= $min) && ($val <= $max))
	//		break;
	//}
	//return $val;
}

// Shortcut for always getting an integer value
func (p *PHP) gauss_intrand(min, max, dev, mean float64) {
	panic("not implemented")
	//return intval(gauss_rand($min, $max, $dev, $mean));
}

func (p *PHP) notice(format string, args ...any) {
	panic("not implemented")
	//global $notices;
	//if (strlen($notices) > 0)
	//	$notices .= "<br />\n";
	//$notices .= $msg;
	if len(p.globals.notices) > 0 {
		p.globals.notices += "<br />"
	}
	p.globals.notices += fmt.Sprintf(format, args...)
}

func (p *PHP) notices(style int) {
	panic("not implemented")
	//global $notices;
	//if (!empty($notices))
	//{
	//	switch ($style)
	//	{
	//	case 2:
	//		echo '<h4 class="cwarn">'. $notices .'</h4>';
	//		break;
	//	case 1:
	//		echo '<h4>'. $notices .'</h4>';
	//		break;
	//	case 0:
	//	default:
	//		echo $notices .'<hr />';
	//		break;
	//	}
	//	$notices = '';
	//}
}

func (p *PHP) redirect(newurl string) {
	panic("not implemented")
	//header('Location: '. $newurl);
	//exit;
}

func (p *PHP) validate_email(email string) {
	panic("not implemented")
	//return preg_match('/[a-z0-9!#$%&\'*+\/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&\'*+\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?/i', $email);
}

func (p *PHP) validate_url(url string) {
	panic("not implemented")
	//return preg_match('/https?:\/\/(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?(\/[^\s]*)?/i', $url);
}

// require_once(PROM_BASEDIR .'includes/PasswordHash.php');

func (p *PHP) enc_password(pass string) {
	panic("not implemented")
	//$pwh = new PasswordHash(10, FALSE);
	//return $pwh->HashPassword($pass);
}

func (p *PHP) chk_password(pass, hash string) {
	panic("not implemented")
	//$pwh = new PasswordHash(10, FALSE);
	//return $pwh->CheckPassword($pass, $hash);
}
