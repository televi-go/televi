package metrics

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/televi-go/migrate"
	"github.com/televi-go/televi/gopage"
	"github.com/televi-go/televi/telegram/dto"
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
		writer := gopage.NewHtmlWriter(ctx)
		return writer.WritePage(usersPage(users, botInfo))
	})
}

func headMixin(context gopage.Context, title string) {
	context.OpenTag("head")
	context.OpenSelfClosing("meta")
	context.Attributes(gopage.Attr{Key: "charset", Value: "UTF-8"})
	context.CloseTag()
	context.OpenTag("title")
	context.Content(title)
	context.CloseTag()
	context.OpenTag("link")
	gopage.WriteAttribute(context, "rel", "stylesheet")
	gopage.WriteAttribute(context, "href", "https://www.unpkg.com/televi_assets_x@latest/css/main.css")
	context.CloseTag()
	context.CloseTag()
}

type HeaderData struct {
	Title string
}

var headerComponent = gopage.MakeComponent[HeaderData](`
	<div class="heading" >
        <div class="content-wrap" style="position:relative; height: 80px">
            <img src="https://www.unpkg.com/televi_assets_x@latest/images/logo.png"
                 style="display: block; position: absolute; height:100%; top:0; left:-90px"
                 alt="">
            <h1 style="margin:auto 0">
                {Title}
            </h1>
        </div>
    </div>
	`)

func writeHeader(ctx gopage.Context, title string) {
	headerComponent(HeaderData{Title: title}, ctx)
}

func bodyWrap(content gopage.RenderAction) gopage.RenderAction {
	return func(ctx gopage.Context) {
		ctx.OpenTag("div")
		gopage.WriteAttribute(ctx, "class", "content-wrap")
		gopage.WriteAttribute(ctx, "style", "display:block")
		content(ctx)
		ctx.CloseTag()
	}
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

var usersRowComponent = gopage.MakeComponent[UserRowData](`
<div class="user-row-wrap">
	<div class="joined_at">{RegisteredAt}</div>
	<div class="username">{UserName}</div>
	<div class="first_name">{FirstName}</div>
	<div class="last_name">{LastName}</div>
</div>
`)

var descriptionComponent = gopage.MakeComponent[BotNameData](`
<h2>
Bot <a href="https://t.me/{Name}">@{Name}</a>
</h2>
`)

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

func usersPage(users []UserRegisteredAt, botInfo string) gopage.Page {
	return func(context gopage.Context) {
		headMixin(context, "users")
		context.OpenTag("body")

		writeHeader(context, "Clients")

		bodyWrap(func(ctx gopage.Context) {
			descriptionComponent.Mount(ctx, BotNameData{Name: botInfo})
			context.OpenTag("div")
			gopage.WriteAttribute(ctx, "class", "users-table")
			context.Content(userRowStyle)
			for _, user := range users {
				usersRowComponent.Mount(ctx, UserRowData{
					FirstName:    user.FirstName,
					LastName:     user.LastName,
					UserName:     user.UiName(),
					RegisteredAt: user.RegisteredAt.Format("02.01.06 15:04"),
				})
			}
			context.CloseTag()
		}).Mount(context)

		context.CloseTag()
	}
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

const createJoinTable = `
CREATE TABLE users_joined (
    user_id BIGINT NOT NULL,
    joinedAt TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, joinedAt),
    FOREIGN KEY fk_user_id (user_id) references users (id) on delete cascade
);
`

var usersMigration = migrate.InMemory{
	Name:       "users-joined.up.sql",
	Statements: createUsersTable,
}

var joinedMigration = migrate.InMemory{
	Name:       "users.up.sql",
	Statements: createJoinTable,
}

var allMigrations = []migrate.InMemory{
	usersMigration,
	joinedMigration,
}
