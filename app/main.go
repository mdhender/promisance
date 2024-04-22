// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements a Promisance server.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/orm"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"log"
)

func main() {
	// default log format to UTC
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	cfg := &config{}

	cobra.CheckErr(Execute(cfg))
}

type config struct{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(cfg *config) error {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&serverArgs.data, "data", "", "path to data files")
	serverCmd.Flags().StringVar(&serverArgs.host, "host", "localhost", "host to bind listener to")
	serverCmd.Flags().StringVar(&serverArgs.port, "port", "8080", "port to bind listener to")
	serverCmd.Flags().StringVar(&serverArgs.templates, "templates", "", "path to template files")
	if err := serverCmd.MarkFlagRequired("data"); err != nil {
		log.Fatalf("setup: markFlagRequired: %v\n", err)
	} else if err = serverCmd.MarkFlagRequired("templates"); err != nil {
		log.Fatalf("setup: markFlagRequired: %v\n", err)
	}

	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().BoolVar(&setupArgs.showTimeZone, "tz", false, "show time zone data")
	setupCmd.Flags().StringVar(&setupArgs.data, "data", ".", "path to data files")
	if err := setupCmd.MarkFlagRequired("data"); err != nil {
		log.Fatalf("setup: markFlagRequired: %v\n", err)
	}

	rootCmd.AddCommand(timeZoneCmd)

	return rootCmd.Execute()
}

var rootArgs struct{}

// rootCmd represents the root command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "app",
	Long:  `app is a command line application for Promisance.`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

var serverArgs struct {
	data      string
	host      string
	port      string
	templates string // path to template files
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the web server",
	Long:  `Start the web server.`,
	Run: func(cmd *cobra.Command, args []string) {
		// verify data path
		if serverArgs.data = strings.TrimSpace(serverArgs.data); serverArgs.data == "" {
			log.Fatal("error: no data path specified\n")
		} else if path, err := filepath.Abs(serverArgs.data); err != nil {
			log.Fatalf("error: data: %v\n", err)
		} else if sb, err := os.Stat(path); err != nil {
			log.Fatalf("error: data: %s: no such directory\n", serverArgs.data)
		} else if !sb.IsDir() {
			log.Fatalf("error: data: %s: not a directory\n", serverArgs.data)
		} else {
			serverArgs.data = path
		}
		log.Printf("server: data %s\n", serverArgs.data)
		// verify templates path
		if serverArgs.templates = strings.TrimSpace(serverArgs.templates); serverArgs.templates == "" {
			log.Fatal("error: no templates path specified\n")
		} else if path, err := filepath.Abs(serverArgs.templates); err != nil {
			log.Fatalf("error: templates: %v\n", err)
		} else if sb, err := os.Stat(path); err != nil {
			log.Fatalf("error: templates: %s: no such directory\n", serverArgs.templates)
		} else if !sb.IsDir() {
			log.Fatalf("error: templates: %s: not a directory\n", serverArgs.templates)
		} else {
			serverArgs.templates = path
		}
		log.Printf("server: templates %s\n", serverArgs.templates)

		s := &server{
			host: "localhost", port: "8080",
			data: serverArgs.data, templates: serverArgs.templates,
			sessions: &SessionManager_t{
				sessions: map[string]*Session_t{},
			},
		}
		s.addr = net.JoinHostPort(s.host, s.port)
		s.tz, _ = time.Now().Zone()

		log.Printf("app: server time zone is %s (logs are UTC)\n", s.tz)

		// temporarily save some routing information. we don't use it, but may.
		// Valid in-game pages - can be specified for 'location' parameter to load corresponding PHP file
		// Values denote any special requirements for loading the page
		valid_locations := map[string]int{
			// 0 - does not require referer or session
			"count":       0, // 0 - does not require referer or session
			"credits":     0, // 0 - does not require referer or session
			"history":     0, // 0 - does not require referer or session
			"login":       0, // 0 - does not require referer or session
			"pguide":      0, // 0 - does not require referer or session
			"playerstats": 0, // 0 - does not require referer or session
			"signup":      0, // 0 - does not require referer or session
			"topclans":    0, // 0 - does not require referer or session
			"topempires":  0, // 0 - does not require referer or session
			"topplayers":  0, // 0 - does not require referer or session
			"relogin":     0, // redirect from login page load; redirects don't set referer, and this could be a bookmark

			// 1 - requires referer from any site
			"game": 1, // redirect from login page submission; redirects don't set referer

			// 2 - requires referer from in-game, also requires active session
			"banner":     2, // 2 - requires referer from in-game, also requires active session
			"guide":      2, // 2 - requires referer from in-game, also requires active session
			"messages":   2, // 2 - requires referer from in-game, also requires active session
			"revalidate": 2, // 2 - requires referer from in-game, also requires active session
			"validate":   2, // 2 - requires referer from in-game, also requires active session
			"main":       2, // both "relogin" and "game" redirect to here

			// Information
			"clanstats": 2, // 2 - requires referer from in-game, also requires active session
			"contacts":  2, // 2 - requires referer from in-game, also requires active session
			"graveyard": 2, // 2 - requires referer from in-game, also requires active session
			"news":      2, // 2 - requires referer from in-game, also requires active session
			"scores":    2, // 2 - requires referer from in-game, also requires active session
			"search":    2, // 2 - requires referer from in-game, also requires active session
			"status":    2, // 2 - requires referer from in-game, also requires active session

			// Use Turns
			"build":    2, // 2 - requires referer from in-game, also requires active session
			"cash":     2, // 2 - requires referer from in-game, also requires active session
			"demolish": 2, // 2 - requires referer from in-game, also requires active session
			"farm":     2, // 2 - requires referer from in-game, also requires active session
			"land":     2, // 2 - requires referer from in-game, also requires active session

			// Finances
			"bank":          2, // 2 - requires referer from in-game, also requires active session
			"lottery":       2, // 2 - requires referer from in-game, also requires active session
			"pubmarketbuy":  2, // 2 - requires referer from in-game, also requires active session
			"pubmarketsell": 2, // 2 - requires referer from in-game, also requires active session
			"pvtmarketbuy":  2, // 2 - requires referer from in-game, also requires active session
			"pvtmarketsell": 2, // 2 - requires referer from in-game, also requires active session

			// Foreign Affairs
			"aid":       2, // 2 - requires referer from in-game, also requires active session
			"clan":      2, // 2 - requires referer from in-game, also requires active session
			"clanforum": 2, // 2 - requires referer from in-game, also requires active session
			"magic":     2, // 2 - requires referer from in-game, also requires active session
			"military":  2, // 2 - requires referer from in-game, also requires active session

			// Management
			"delete":        2, // 2 - requires referer from in-game, also requires active session
			"manage/clan":   2, // 2 - requires referer from in-game, also requires active session
			"manage/empire": 2, // 2 - requires referer from in-game, also requires active session
			"manage/user":   2, // 2 - requires referer from in-game, also requires active session

			// Administration
			"admin/clans":       2, // 2 - requires referer from in-game, also requires active session
			"admin/empedit":     2, // 2 - requires referer from in-game, also requires active session
			"admin/empires":     2, // 2 - requires referer from in-game, also requires active session
			"admin/history":     2, // 2 - requires referer from in-game, also requires active session
			"admin/log":         2, // 2 - requires referer from in-game, also requires active session
			"admin/market":      2, // 2 - requires referer from in-game, also requires active session
			"admin/messages":    2, // 2 - requires referer from in-game, also requires active session
			"admin/permissions": 2, // 2 - requires referer from in-game, also requires active session
			"admin/round":       2, // 2 - requires referer from in-game, also requires active session
			"admin/users":       2, // 2 - requires referer from in-game, also requires active session

			// Logout
			"logout": 2, // 2 - requires referer from in-game, also requires active session
		}
		log.Printf("server: todo: implement valid_locations referer logic (%d pages)\n", len(valid_locations))

		handler := s.routes(valid_locations)

		dbFile := filepath.Join(serverArgs.data, "promisance.sqlite")
		log.Printf("server: connecting to database: %s\n", dbFile)
		var err error
		s.db, err = orm.OpenSqliteDatabase(dbFile)
		if err != nil {
			log.Fatalf("server: database: %v\n", err)
		}
		defer func() {
			if err := s.db.Close(); err != nil {
				log.Printf("server: db close failed: %v\n", err)
			}
			log.Printf("server: db closed\n")
		}()

		// world data is a one time load that is shared with all the handlers
		s.world, err = s.db.WorldVarsFetch()
		if err != nil {
			log.Fatalf("server: failed to fetch vars: %v\n", err)
		} else if s.world == nil {
			log.Fatalf("server: database: world_vars undefined\n")
		} else if s.world.Id != 1 {
			log.Fatalf("server: database: world_vars corrupted\n")
		}
		log.Printf("server: fetched world variables\n")

		// If we're configured for cronless turn updates, check them now
		//turns, err := s.db.P prom_turns();
		//if !TURNS_CRONTAB {
		//	turns.doUpdate()
		//}

		log.Printf("server: serving on http://%s\n", s.addr)
		log.Fatalln(http.ListenAndServe(s.addr, handler))
	},
}

var setupArgs struct {
	data         string
	showTimeZone bool
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup a new server",
	Long:  `Setup a new server database and files.`,
	Run: func(cmd *cobra.Command, args []string) {
		startedAt := time.Now()
		log.Printf("setup: started\n")

		serverTimeZone, serverTimeOffset := time.Now().Zone()
		if setupArgs.showTimeZone {
			log.Printf("setup: server time zone %q\n", serverTimeZone)
			if serverTimeZone != "UTC" {
				localNow := time.Now()
				log.Printf("setup: server_time %v\n", localNow)
				log.Printf("setup: utc_time    %v\n", localNow.UTC())
				// calculate the offset in hours and minutes
				offsetInMinutes := serverTimeOffset / 60
				hours, minutes := offsetInMinutes/60, offsetInMinutes%60
				// if offset is negative, then hours and minutes both should be negative.
				if offsetInMinutes < 0 {
					minutes = -minutes
				}
				if hours != 0 && minutes != 0 {
					log.Printf("setup: offset between %s and UTC is currently %4d hours and %4d minutes\n", serverTimeZone, hours, minutes)
				} else if hours != 0 {
					log.Printf("setup: offset between %s and UTC is currently %4d hours\n", serverTimeZone, hours)
				} else if minutes != 0 {
					log.Printf("setup: offset between %s and UTC is currently %4d minutes\n", serverTimeZone, minutes)
				}
			}
		}

		// verify data path
		if setupArgs.data = strings.TrimSpace(setupArgs.data); setupArgs.data == "" {
			log.Fatal("error: no data path specified\n")
		} else if path, err := filepath.Abs(setupArgs.data); err != nil {
			log.Fatalf("error: data: %v\n", err)
		} else if sb, err := os.Stat(path); err != nil {
			log.Fatalf("error: data: %s: no such directory\n", setupArgs.data)
		} else if !sb.IsDir() {
			log.Fatalf("error: data: %s: not a directory\n", setupArgs.data)
		} else {
			setupArgs.data = path
		}
		log.Printf("setup: data %s\n", setupArgs.data)

		// make sure we're not running setup twice into the same path
		dbFile := filepath.Join(setupArgs.data, "promisance.sqlite")
		if _, err := os.Stat(dbFile); err == nil {
			if err := os.Remove(dbFile); err != nil {
				log.Fatalf("error: %s: %v\n", dbFile, err)
			}
			log.Printf("setup: todo: remove force delete of existing database\n")
			// log.Fatalf("setup: this server has already been setup\n")
		}

		// load the configuration file
		var cfg struct {
			Administrator struct {
				UserName   string `json:"username"`
				Password   string `json:"password"`
				Nickname   string `json:"nickname"`
				Email      string `json:"email"`
				EmpireName string `json:"empire_name"`
			} `json:"administrator"`
			Site struct {
				Language          string `json:"language,omitempty"`
				DefaultTimezone   string `json:"default_timezone,omitempty"`
				DefaultDateFormat string `json:"default_date_format,omitempty"`
				RoundBegin        string `json:"round_begin,omitempty"`
				RoundClosing      string `json:"round_closing,omitempty"`
				RoundEnd          string `json:"round_end,omitempty"`
			} `json:"site,omitempty"`
		}
		cfg.Site.DefaultTimezone = serverTimeZone
		defaultRoundBegin := time.Now().Add(2 * 24 * time.Hour) // default to 48 hours
		cfg.Site.RoundBegin = defaultRoundBegin.Format(time.RFC3339)
		cfg.Site.RoundClosing = defaultRoundBegin.Add(21 * 24 * time.Hour).Format(time.RFC3339)
		cfg.Site.RoundEnd = defaultRoundBegin.Add(28*24*time.Hour - time.Second).Format(time.RFC3339)

		cfgFile := filepath.Join(setupArgs.data, "config.json")
		log.Printf("setup: config %s\n", cfgFile)
		if buf, err := os.ReadFile(cfgFile); err != nil {
			log.Fatalf("setup: config: %v\n", err)
		} else if err = json.Unmarshal(buf, &cfg); err != nil {
			log.Fatalf("setup: config: %v\n", err)
		} else {
			log.Printf("setup: config: loaded successfully\n")
		}

		if cfg.Site.DefaultTimezone == "" {
			cfg.Site.DefaultTimezone, _ = time.Now().Zone()
		}

		var roundBegin, roundClosing, roundEnd time.Time
		if cfg.Site.RoundBegin == "" {
			// by default, start in 48 hours
			roundBegin = defaultRoundBegin.Add(2 * 24 * time.Hour)
		} else {
			// todo: parse the round values from the configuration
			roundBegin = defaultRoundBegin.Add(2 * 24 * time.Hour)
		}
		roundClosing = roundBegin.Add(21 * 24 * time.Hour)
		roundEnd = roundBegin.Add(28*24*time.Hour - time.Second)
		if roundBegin.Before(time.Now().Add(5 * time.Minute)) {
			// the round shouldn't start too soon after setting up the server
			log.Printf("setup: def_time    %v\n", defaultRoundBegin)
			log.Printf("setup: round_begin %v\n", roundBegin)
			log.Fatalf("setup: bad round_begin\n")
		}
		log.Printf("setup: round_begin   %v\n", roundBegin)
		log.Printf("setup: round_closing %v\n", roundClosing)
		log.Printf("setup: round_end     %v\n", roundEnd)

		// validate inputs
		var inputErrors []error
		if cfg.Administrator.UserName == "" {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.username: missing"))
		} else if len(cfg.Administrator.UserName) < 6 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.username: must be at least 6 characters"))
		} else if len(cfg.Administrator.UserName) >= 255 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.username: must be less than 255 characters"))
		}
		if cfg.Administrator.Password == "" {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.password: missing"))
		} else if len(cfg.Administrator.Password) < 6 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.password: must be at least 6 characters"))
		} else if len(cfg.Administrator.Password) >= 255 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.password: must be less than 255 characters"))
		}
		if cfg.Administrator.Nickname == "" {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.nickname: missing"))
		} else if len(cfg.Administrator.Nickname) < 6 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.nickname: must be at least 6 characters"))
		} else if len(cfg.Administrator.Nickname) >= 255 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.nickname: must be less than 255 characters"))
		}
		if cfg.Administrator.Email == "" {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.email: missing"))
		} else if len(cfg.Administrator.Email) < 6 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.email: must be at least 6 characters"))
		} else if len(cfg.Administrator.Email) >= 255 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.email: must be less than 255 characters"))
		} else if _, err := mail.ParseAddress(cfg.Administrator.Email); err != nil {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.email: %v", err))
		}
		if cfg.Administrator.EmpireName == "" {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.empire_name: missing"))
		} else if len(cfg.Administrator.EmpireName) < 6 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.empire_name: must be at least 6 characters"))
		} else if len(cfg.Administrator.EmpireName) >= 255 {
			inputErrors = append(inputErrors, fmt.Errorf("administrator.empire_name: must be less than 255 characters"))
		}
		if roundBegin.Before(time.Now().Add(5 * time.Minute)) {
			inputErrors = append(inputErrors, fmt.Errorf("site.round_begin: round must begin in future\n"))
		} else if !roundBegin.Before(roundClosing) {
			inputErrors = append(inputErrors, fmt.Errorf("site.round_closing: round must close after it begins\n"))
		} else if roundClosing.Add(VACATION_START + VACATION_LIMIT).After(roundEnd) {
			inputErrors = append(inputErrors, fmt.Errorf("site.round_closing: cooldown breaks round_end?\n"))
		}
		if len(inputErrors) != 0 {
			for _, err := range inputErrors {
				log.Printf("setup: input error: %v\n", err)
			}
			log.Fatalf("setup: please correct errors and restart\n")
		}
		log.Printf("setup: administrator.username    %q\n", cfg.Administrator.UserName)
		log.Printf("setup: administrator.password    %q\n", cfg.Administrator.Password)
		log.Printf("setup: administrator.nickname    %q\n", cfg.Administrator.Nickname)
		log.Printf("setup: administrator.email       %q\n", cfg.Administrator.Email)
		log.Printf("setup: administrator.empire_name %q\n", cfg.Administrator.EmpireName)
		log.Printf("setup: site.round_begin          %v\n", roundBegin)
		log.Printf("setup: site.round_closing        %v\n", roundClosing)
		log.Printf("setup: site.round_end            %v\n", roundEnd)

		// load the database initialization script and run it
		log.Printf("setup: running database initialization\n")
		db, err := orm.CreateSqliteDatabase(dbFile)
		if err != nil {
			log.Fatalf("setup: failed to open database: %v\n", err)
		}
		defer func() {
			_ = db.Close()
			log.Printf("setup: database closed %q\n", dbFile)
		}()

		world := &model.World_t{
			LottoCurrentJackpot:   LOTTERY_JACKPOT,
			LottoYesterdayJackpot: LOTTERY_JACKPOT,
			LottoLastPicked:       0,
			LottoLastWinner:       0,
			LottoJackpotIncrease:  0,
			RoundTimeBegin:        roundBegin,
			RoundTimeClosing:      roundClosing,
			RoundTimeEnd:          roundEnd,
			TurnsNext:             roundBegin.Add(TURNS_OFFSET * 60 * time.Second),
			TurnsNextHourly:       roundBegin.Add(TURNS_OFFSET_HOURLY * 60 * time.Hour),
			TurnsNextDaily:        roundBegin.Add(TURNS_OFFSET_DAILY * 24 * 60 * time.Hour),
		}
		if err := db.WorldVarsInitialize(world); err != nil {
			log.Fatalf("setup: failed to initialize world variables: %v\n", err)
		}
		log.Printf("setup: created world variables\n")

		// persist these to the database
		user, err := db.UserCreate(cfg.Administrator.UserName, cfg.Administrator.Email)
		if err != nil {
			log.Fatalf("setup: failed to create administrator: %v\n", err)
		}
		user.Password = cfg.Administrator.Password
		if err := db.UserPasswordUpdate(user); err != nil {
			log.Fatalf("setup: failed to update administrator password: %v\n", err)
		}
		user.Nickname = cfg.Administrator.Nickname
		user.Flags = model.UserFlag_t{Admin: true, Mod: true}
		user.LastIP = "localhost"
		if err := db.UserAttributesUpdate(user); err != nil {
			log.Fatalf("setup: failed to update administrator: %v\n", err)
		}
		log.Printf("setup: created administrator %q\n", user.UserName)

		empire, err := db.EmpireCreate(user, cfg.Administrator.EmpireName, "HUMAN")
		if err != nil {
			log.Fatalf("setup: failed to create empire: %v\n", err)
		}
		empire.Flags = model.EmpireFlag_t{Admin: true, Mod: true}
		if err := db.EmpireAttributesUpdate(empire); err != nil {
			log.Fatalf("setup: failed to update empire: %v\n", err)
		}
		log.Printf("setup: created empire %q\n", empire.Name)

		log.Printf("setup: completed in %v\n", time.Now().Sub(startedAt))
	},
}

// timeZoneCmd implements a command to show the current time zone data
var timeZoneCmd = &cobra.Command{
	Use:   "tz",
	Short: "tz data for the server",
	Long:  `tz shows time zone data for the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverTimeZone, serverTimeOffset := time.Now().Zone()
		log.Printf("tz: server time zone %q\n", serverTimeZone)
		if serverTimeZone != "UTC" {
			localNow := time.Now()
			log.Printf("tz: server_time %v\n", localNow)
			log.Printf("tz: utc_time    %v\n", localNow.UTC())
			// calculate the offset in hours and minutes
			offsetInMinutes := serverTimeOffset / 60
			hours, minutes := offsetInMinutes/60, offsetInMinutes%60
			// if offset is negative, then hours and minutes both should be negative.
			if offsetInMinutes < 0 {
				minutes = -minutes
			}
			if hours != 0 && minutes != 0 {
				log.Printf("tz: offset between %s and UTC is currently %4d hours and %4d minutes\n", serverTimeZone, hours, minutes)
			} else if hours != 0 {
				log.Printf("tz: offset between %s and UTC is currently %4d hours\n", serverTimeZone, hours)
			} else if minutes != 0 {
				log.Printf("tz: offset between %s and UTC is currently %4d minutes\n", serverTimeZone, minutes)
			}
		}
	},
}
