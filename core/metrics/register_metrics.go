package metrics

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/televi-go/migrate"
	"github.com/televi-go/televi/core/metrics/grouping"
	"github.com/televi-go/televi/core/metrics/pages"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/util"
	"net/http"
	"strings"
	"time"
)

func SetupDatabaseCommons(db *sql.DB) error {
	return migrate.RunInMemory(db, allMigrations, migrate.MysqlDialect)
}

func SetupRouterCommons(router fiber.Router, db *sql.DB, botInfo string) {
	router.Get("/api/users", func(ctx *fiber.Ctx) error {
		users, err := getAllUsersWithRegisterData(db)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		err = ctx.JSON(users)
		ctx.Set("content-type", "application/json; charset=utf-8")
		return err
	})
	router.Get("/users", func(ctx *fiber.Ctx) error {
		users, err := getAllUsersWithRegisterData(db)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		viewData := pages.JoinedPageViewData{
			Title:  "Users joined",
			Name:   botInfo,
			Groups: makeGroups(users),
		}
		ctx.Set("content-type", "text/html; charset=utf-8")
		return pages.JoinedPageTemplate.Execute(ctx, viewData)
	})
	router.Get("/actions", func(ctx *fiber.Ctx) error {
		action := ctx.Query("action")

		actions, err := getUserActions(db, action)
		uniqueActions, err2 := getUniqueActions(db)
		if err != nil || err2 != nil {
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		ctx.Set("content-type", "text/html; charset=utf-8")
		return pages.ExecuteActionsPage(ctx, actions, action, uniqueActions)
	})
}

type BotNameData struct {
	Name string
}

type UserRowData struct {
	FirstName    string
	LastName     string
	UserName     string
	RegisteredAt string
}

func toJoinedData(in []UserRegisteredAt, timeFormat string) []pages.JoinedAt {
	return util.Map(in, func(elem UserRegisteredAt) pages.JoinedAt {
		return pages.JoinedAt{
			FirstName:   elem.FirstName,
			LastName:    elem.LastName,
			FormattedAt: elem.RegisteredAt.Format(timeFormat),
			UiName:      elem.UiName(),
		}
	})
}

func getUniqueActions(db *sql.DB) ([]pages.ActionInfo, error) {
	rows, err := db.Query("select (CONCAT(domain, '.', action)) from actions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result map[string]int = map[string]int{}

	for rows.Next() {
		var action string
		err = rows.Scan(&action)
		if err != nil {
			return nil, err
		}
		action = strings.TrimPrefix(action, ".")
		result[action] = result[action] + 1
	}

	info := make([]pages.ActionInfo, 0, len(result))
	for a, count := range result {
		info = append(info, pages.ActionInfo{Action: a, Count: count})
	}
	return info, nil
}

func RegisterAction(db *sql.DB, userId int64, domain string, action string) error {
	_, err := db.Exec(
		"insert into actions (user_id, domain, action, committedAt) values (?, ?, ?, ?)",
		userId,
		domain,
		action,
		time.Now(),
	)
	return err
}

func makeGroups(source []UserRegisteredAt) (groups []pages.Group) {
	todayGroupContent, notToday := util.FilterOut(source, func(elem UserRegisteredAt) bool {
		return grouping.IsToday(elem.RegisteredAt)
	})
	if len(todayGroupContent) != 0 {
		groups = append(groups, pages.MakeGroup(
			"Today",
			toJoinedData(todayGroupContent, "15:04"),
		))
	}

	yesterdayContent, notLatterTwoDays := util.FilterOut(notToday, func(elem UserRegisteredAt) bool {
		return grouping.IsYesterday(elem.RegisteredAt)
	})

	if len(yesterdayContent) != 0 {
		groups = append(groups, pages.MakeGroup(
			"Yesterday",
			toJoinedData(yesterdayContent, "15:04"),
		))
	}

	weekContent, earlierContent := util.FilterOut(notLatterTwoDays, func(elem UserRegisteredAt) bool {
		return grouping.IsThisWeek(elem.RegisteredAt)
	})

	if len(weekContent) != 0 {
		groups = append(groups,
			pages.MakeGroup("This week", toJoinedData(weekContent, "Monday 15:04")))
	}

	thisMonthContent, prevMonthsContent := util.FilterOut(earlierContent, func(elem UserRegisteredAt) bool {
		return grouping.IsThisMonth(elem.RegisteredAt)
	})

	if len(thisMonthContent) != 0 {
		groups = append(groups, pages.MakeGroup("This month", toJoinedData(thisMonthContent, "02.01")))
	}

	mi := util.MonthIterFrom(time.Now()).Prev()
	var currMonthContent []UserRegisteredAt
	for len(prevMonthsContent) > 0 {
		currMonthContent, prevMonthsContent = util.FilterOut(prevMonthsContent, func(elem UserRegisteredAt) bool {
			return elem.RegisteredAt.After(mi.Begin()) && elem.RegisteredAt.Before(mi.End())
		})
		if len(currMonthContent) > 0 {
			groups = append(groups, pages.MakeGroup(mi.Begin().Format("January"), toJoinedData(currMonthContent, "02.01")))
		}
	}

	return
}

type UserRegisteredAt struct {
	dto.User
	RegisteredAt time.Time `json:"registered_at"`
}

func (user UserRegisteredAt) UiName() string {
	if user.UserName != "" {
		return user.UserName
	}
	return fmt.Sprintf("@%d", user.ID)
}

func AddRegistered(db *sql.DB, user *dto.User) error {
	row := db.QueryRow("select exists(select id from users where id = ?)", user.ID)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	_, err = tx.Exec(`insert into users(id, is_bot, is_premium, first_name, last_name, username, language_code) VALUES 
    (?, ?, ?, ?, ?, ?, ?)`, user.ID, user.IsBot, user.IsPremium, user.FirstName, user.LastName, user.UserName, user.LanguageCode)
	if err != nil {
		return err
	}

	_, err = tx.Exec("insert into users_joined(user_id, joinedAt) VALUES (?, ?)", user.ID, time.Now())

	return err
}

func getAllUsersWithRegisterData(db *sql.DB) (result []UserRegisteredAt, err error) {
	result = make([]UserRegisteredAt, 0)
	rows, err := db.Query("select id, is_bot, is_premium, first_name, COALESCE(last_name, ''), COALESCE(username, ''), COALESCE(language_code, ''), joinedAt from users left join users_joined on users.id = users_joined.user_id order by users_joined.joinedAt desc ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			user     dto.User
			joinedAt time.Time
		)
		err = rows.Scan(&user.ID, &user.IsBot, &user.IsPremium, &user.FirstName, &user.LastName, &user.UserName, &user.LanguageCode, &joinedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, UserRegisteredAt{RegisteredAt: joinedAt, User: user})
	}

	return
}

func PrintAction(domain, action string) string {
	if domain == "" {
		return action
	}
	return fmt.Sprintf("%s.%s", domain, action)
}

func getUserActions(db *sql.DB, action string) (result []pages.UserWithAction, err error) {
	rows, err := db.Query(`select 
    actions.domain,
    actions.action,
    actions.committedAt,
    users.id,
    users.first_name,
    users.is_premium,
    coalesce(users.last_name, ''),
    coalesce(users.username, '') from actions left join users on actions.user_id = users.id where actions.action = ? or ? = '' order by actions.committedAt desc`, action, action)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			domain      string
			action      string
			committedAt time.Time
			userId      int64
			firstname   string
			isPremium   bool
			lastname    string
			username    string
		)
		err = rows.Scan(&domain, &action, &committedAt, &userId, &firstname, &isPremium, &lastname, &username)
		if err != nil {
			return nil, err
		}
		if username == "" {
			username = fmt.Sprintf("#%d", userId)
		}
		result = append(result, pages.UserWithAction{
			Action:      PrintAction(domain, action),
			CommittedAt: committedAt,
			FirstName:   firstname,
			LastName:    lastname,
			UserName:    username,
			IsPremium:   isPremium,
		})
	}

	return
}

const createUsersTable = `
CREATE TABLE users (
	id BIGINT NOT NULL,
	is_bot BOOLEAN DEFAULT FALSE,
	is_premium BOOLEAN DEFAULT FALSE,
	first_name VARCHAR(64) NOT NULL,
	last_name VARCHAR(64) DEFAULT NULL,
	username VARCHAR(64) DEFAULT NULL,
	language_code VARCHAR(13) DEFAULT NULL,
	PRIMARY KEY (id)
);
`

const createActionsTable = `
CREATE TABLE actions (
    user_id BIGINT NOT NULL,
    domain VARCHAR(64) NOT NULL,
    action VARCHAR(64) NOT NULL,
    committedAt TIMESTAMP NOT NULL,
    id int not null primary key auto_increment,
    index (action),
    index (user_id),
    foreign key fk_user_id (user_id) references users(id) on delete cascade
);
`

const createJoinTable = `
CREATE TABLE users_joined (
    user_id BIGINT NOT NULL,
    joinedAt TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, joinedAt),
    FOREIGN KEY fk_user_id (user_id) references users (id) on delete cascade
);
`

var usersMigration = migrate.InMemory{
	Name:       "users.up.sql",
	Statements: createUsersTable,
}

var joinedMigration = migrate.InMemory{
	Name:       "users-joined.up.sql",
	Statements: createJoinTable,
}

var actionsMigration = migrate.InMemory{
	Name:       "users-with-actions.up.sql",
	Statements: createActionsTable,
}

var allMigrations = []migrate.InMemory{
	usersMigration,
	joinedMigration,
	actionsMigration,
}
