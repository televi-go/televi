package metrics

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/televi-go/migrate"
	"github.com/televi-go/televi/core/metrics/pages"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/util"
	"net/http"
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

func isSameDay(val time.Time, point time.Time) bool {
	yearP, monthP, dayP := point.Date()
	yearV, monthV, dayV := val.Date()
	return yearP == yearV && monthP == monthV && dayP == dayV
}

func isToday(val time.Time) bool {
	return isSameDay(val, time.Now())
}

func isYesterday(val time.Time) bool {
	return isSameDay(val, time.Now().Add(-time.Hour*24)) && !isToday(val)
}

func isThisWeek(val time.Time) bool {
	currYear, currWeek := time.Now().ISOWeek()
	valYear, valWeek := val.ISOWeek()
	return currYear == valYear && currWeek == valWeek
}

func isThisMonth(val time.Time) bool {
	yearP, monthP, _ := time.Now().Date()
	yearV, monthV, _ := val.Date()
	return yearP == yearV && monthP == monthV
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

func makeGroups(source []UserRegisteredAt) (groups []pages.Group) {
	todayGroupContent, notToday := util.FilterOut(source, func(elem UserRegisteredAt) bool {
		return isToday(elem.RegisteredAt)
	})
	if len(todayGroupContent) != 0 {
		groups = append(groups, pages.Group{
			Title: "Today",
			Users: toJoinedData(todayGroupContent, "15:04"),
		})
	}

	yesterdayContent, notLatterTwoDays := util.FilterOut(notToday, func(elem UserRegisteredAt) bool {
		return isYesterday(elem.RegisteredAt)
	})

	if len(yesterdayContent) != 0 {
		groups = append(groups, pages.Group{
			Title: "Yesterday",
			Users: toJoinedData(yesterdayContent, "15:04"),
		})
	}

	weekContent, earlierContent := util.FilterOut(notLatterTwoDays, func(elem UserRegisteredAt) bool {
		return isThisWeek(elem.RegisteredAt)
	})

	if len(weekContent) != 0 {
		groups = append(groups, pages.Group{Title: "This week", Users: toJoinedData(weekContent, "Monday 15:04")})
	}

	if len(earlierContent) != 0 {
		groups = append(groups, pages.Group{Title: "Earlier", Users: toJoinedData(earlierContent, "02.01 15:04")})
	}

	return
}

const userRowStyle = `
<style>
	.user-row-wrap {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding-block: 8px;
	}

	.joined_at {
		width: 7rem;
	}

	.username {
		width: 7rem;
	}
	.first_name {
		width: 7rem;
	}

</style>
`

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
    primary key (user_id, domain, action),
    index (user_id),
    foreign key fk_user_id (user_id) references users(id)
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
