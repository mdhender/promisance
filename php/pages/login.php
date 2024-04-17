<?php
/* QM Promisance - Turn-based strategy game
 * Copyright (C) QMT Productions
 *
 * $Id: login.php 1983 2014-10-01 15:18:43Z quietust $
 */

if (!defined('IN_GAME'))
	die('Access denied');

if ($action == 'login') do
{
	if (!isFormPost())
		break;
	$username = getFormVar('login_username');
	$password = getFormVar('login_password');

	if (empty($username))
	{
		notice(lang('INPUT_NEED_USERNAME'));
		break;
	}

	if (empty($password))
	{
		notice(lang('INPUT_NEED_PASSWORD'));
		break;
	}

	$user1 = new prom_user();
	if (!$user1->findName($username))
	{
		notice(lang('LOGIN_USER_NOT_FOUND'));
		logmsg(E_USER_NOTICE, 'failed (username) - '. $username);
		break;
	}
	if (!$user1->load())
	{
		// this should be impossible if findName() succeeded
		notice(lang('LOGIN_USER_NOT_FOUND'));
		logmsg(E_USER_NOTICE, 'failed (load) - '. $username);
		break;
	}
	if (!$user1->checkPassword($password))
	{
		notice(lang('INPUT_INCORRECT_PASSWORD'));
		logmsg(E_USER_NOTICE, 'failed (password) - '. $username);
		break;
	}
	if ($user1->u_flags & UFLAG_CLOSED)
	{
		notice(lang('LOGIN_USER_CLOSED'));
		logmsg(E_USER_NOTICE, 'failed (closed) - '. $username);
		break;
	}

	$q = $db->prepare('SELECT e_id FROM '. EMPIRE_TABLE .' WHERE u_id = ? AND e_flags & ? = 0 ORDER BY e_id ASC');
	$q->bindIntValue(1, $user1->u_id);
	$q->bindIntValue(2, EFLAG_DELETE);
	$q->execute() or error_500('ERROR_TITLE', 'Failed to check for registered empires');
	$data = $q->fetchAll();
	$emplist = array();
	foreach ($data as $row)
		$emplist[] = $row['e_id'];

	// if they've signed up before but don't have an empire, bounce them over to the signup page
	if (count($emplist) == 0)
	{
		if (!ROUND_SIGNUP)
		{
			notice(lang('LOGIN_NO_EMPIRE'));
			break;
		}
		if (SIGNUP_CLOSED_EMPIRE)
		{
			notice(lang('LOGIN_NO_EMPIRE_CLOSED'));
			break;
		}
		redirect(URL_BASE .'?location=signup&amp;registered='. urlencode($user1->u_username));
	}
	// load the first empire owned by the user
	$empire = $emplist[0];

	$emp1 = new prom_empire($empire);
	$emp1->load();

	// set lock for logging purposes
	$lock['emp1'] = $emp1->e_id;
	$lock['user1'] = $user1->u_id;
	db_lockentities(array($user1, $emp1), $user1->u_id);
	logevent(varlist(array('username', 'emplist'), get_defined_vars()));

	// convert account password if necessary
	$user1->checkPassword($password, TRUE);

	prom_session::start();
	$_SESSION['user'] = $user1->u_id;
	$_SESSION['empire'] = $emp1->e_id;

	$user1->u_lastip = $_SERVER['REMOTE_ADDR'];
	$user1->u_lastdate = CUR_TIME;

	// only set them online if the round has actually started
	if (ROUND_STARTED)
		$emp1->setFlag(EFLAG_ONLINE);

	$user1->save();
	$emp1->save();
	$db->releaseLocks();

	redirect(URL_BASE .'?location=game');
} while (0);
else
{
	// If the session cookie is even set, then jump to the main page
	// If it turns out to be invalid, it'll be unset and bounced back to here
	if (prom_session::check())
		redirect(URL_BASE .'?location=relogin');
}
// destroy any active session
if (prom_session::check())
{
	prom_session::start();
	prom_session::kill();
}

$html = new prom_html_compact();
$html->begin('LOGIN_TITLE');

$num = $db->queryCell('SELECT COUNT(*) FROM '. EMPIRE_TABLE .' WHERE u_id != 0');
$num = str_pad($num, COUNTER_DIGITS, '0', STR_PAD_LEFT);
if (strlen(COUNTER_TEMPLATE) > 0)
{
	$counter = getimagesize(PROM_BASEDIR .'images/'. COUNTER_TEMPLATE);
	$count_data = '<img src="?location=count" alt="'. $num .'" style="width:'. ($counter[0] / 10 * strlen($num)) .'px;height:'. $counter[1] .'px" />';
}
else	$count_data = '<b>'. $num .'</b>';
?>
<h2><?php echo GAME_TITLE; ?></h2>
<?php echo lang('LOGIN_VERSION', GAME_VERSION); ?><br />
<?php echo lang('LOGIN_DATE_RANGE', gmdate('F j', $world->round_time_begin), gmdate('F j', $world->round_time_end)); ?><br />
<?php echo lang('LOGIN_COUNTER', $count_data); ?><br />
<?php
notices(1);
?>
<form method="post" action="?location=login">
<div>
<?php echo lang('LABEL_USERNAME'); ?> <input type="text" name="login_username" size="8" /><br />
<?php echo lang('LABEL_PASSWORD'); ?> <input type="password" name="login_password" size="8" /><br />
<input type="hidden" name="action" value="login" /><input type="submit" value="<?php echo lang('LOGIN_SUBMIT'); ?>" />
</div>
</form>
<?php
if (ROUND_SIGNUP && !(SIGNUP_CLOSED_USER && SIGNUP_CLOSED_EMPIRE))
	echo '<a href="?location=signup"><b>'. lang('LOGIN_SIGNUP') .'</b></a><br />';
else	echo '<b>'. lang('LOGIN_SIGNUP_CLOSED') .'</b><br />';
echo '<a href="?location=topempires"><b>'. lang('LOGIN_TOPEMPIRES') .'</b></a><br />';
if (CLAN_ENABLE)
	echo '<a href="?location=topclans"><b>'. lang('LOGIN_TOPCLANS') .'</b></a><br />';
echo '<br />';
echo '<a href="?location=topplayers"><b>'. lang('LOGIN_TOPPLAYERS') .'</b></a><br />';
echo '<a href="?location=history"><b>'. lang('LOGIN_HISTORY') .'</b></a><br />';
echo '<a href="?location=pguide"><b>'. lang('LOGIN_GUIDE') .'</b></a><br />';

$html->end();
?>
