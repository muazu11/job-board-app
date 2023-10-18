package db

import (
	"context"
	"fmt"
	"jobboard/backend/services"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string
	Password string
	Port     int
	Host     string

	PageLimit int
}

type DB struct {
	pool      *pgxpool.Pool
	pageLimit int
}

func New(config Config, servicesConfig services.Config) DB {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/postgres",
		config.User, config.Password, config.Host, config.Port,
	)
	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		panic(err)
	}

	db := DB{
		pool:      pool,
		pageLimit: config.PageLimit,
	}
	err = db.migrate(servicesConfig)
	if err != nil {
		panic(err)
	}

	return db
}

func (d DB) Exec(ctx context.Context, query string, args ...any) error {
	_, err := d.pool.Exec(ctx, query, args...)
	return err
}

func (d DB) Query(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, d.pool, dest, query, args...)
}

func (d DB) QueryRow(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Get(ctx, d.pool, dest, query, args...)
}

func (d DB) QueryPage(ctx context.Context, dest Pageable, query, column string, page Page) (Cursors, error) {
	if page.emptyCursor && page.previous {
		return Cursors{Next: page.cursor}, nil
	}

	query = fmt.Sprintf(`
		%s
		WHERE (NOT @previous AND %s > @cursor) OR (@previous AND %s < @cursor)
		ORDER BY %s
		LIMIT @limit`,
		query, column, column, column,
	)
	err := pgxscan.Select(ctx, d.pool, dest, query, page.toArgs(d.pageLimit+1))
	if err != nil {
		return Cursors{}, err
	}

	cursors := d.processPage(dest, page)
	return cursors, nil
}

type ColumnKind int

const (
	TextColumn ColumnKind = iota
	IntColumn
	FloatColumn
)

type Page struct {
	cursor   any
	previous bool

	emptyCursor bool
}

func newPageString(cursor string, previous bool) Page {
	return Page{
		cursor:      cursor,
		previous:    previous,
		emptyCursor: cursor == "",
	}
}

func newPageInt(cursor int, previous bool) Page {
	return Page{
		cursor:      cursor,
		previous:    previous,
		emptyCursor: cursor == 0,
	}
}

func newPageFloat(cursor float64, previous bool) Page {
	return Page{
		cursor:      cursor,
		previous:    previous,
		emptyCursor: cursor == 0.,
	}
}

func NewPage(cursor any, previous bool) Page {
	switch c := cursor.(type) {
	case string:
		return newPageString(c, previous)
	case int:
		return newPageInt(c, previous)
	case float64:
		return newPageFloat(c, previous)
	default:
		panic("unexpected page cursor type")
	}
}

func PageFromContext(c *fiber.Ctx, cursorKind ColumnKind) Page {
	const cursorKey = "page_cursor"
	previous := c.QueryBool("page_previous")
	switch cursorKind {
	case TextColumn:
		return newPageString(c.Query(cursorKey), previous)
	case IntColumn:
		return newPageInt(c.QueryInt(cursorKey), previous)
	case FloatColumn:
		return newPageFloat(c.QueryFloat(cursorKey), previous)
	default:
		panic("unexpected page cursor type")
	}
}

func (p Page) toArgs(limit int) pgx.NamedArgs {
	return pgx.NamedArgs{
		"cursor":   p.cursor,
		"previous": p.previous,
		"limit":    limit,
	}
}

type Cursors struct {
	Previous any
	Next     any
}

type CursorWrap[T any] struct {
	Cursors Cursors
	Data    T
}

func NewCursorWrap[T any](cursors Cursors, data T) CursorWrap[T] {
	return CursorWrap[T]{
		Cursors: cursors,
		Data:    data,
	}
}

type Pageable interface {
	Len() int
	GetCursor(idx int) any
	Slice(start, end int)
}

func (d DB) processPage(dest Pageable, page Page) Cursors {
	destLen := dest.Len()
	canContinue := destLen == d.pageLimit+1

	var cursors Cursors
	if page.previous {
		if destLen != 0 {
			cursors.Next = dest.GetCursor(destLen - 1)
		}
		if canContinue {
			cursors.Previous = dest.GetCursor(1)
			dest.Slice(1, destLen)
		}
	} else {
		if !page.emptyCursor && destLen != 0 {
			cursors.Previous = dest.GetCursor(0)
		}
		if canContinue {
			cursors.Next = dest.GetCursor(destLen - 2)
			dest.Slice(0, destLen-1)
		}
	}

	return cursors
}
