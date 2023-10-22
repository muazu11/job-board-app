package db

import (
	"context"
	"fmt"
	"jobboard/backend/services"
	jsonutil "jobboard/backend/util/json"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sanggonlee/gosq"
)

type CursorError error

var (
	ErrInvalidCursor CursorError = fmt.Errorf("invalid cursor type")
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

func (d DB) QueryPage(
	ctx context.Context, dest Pageable, query, column string, page Page, args ...any,
) (Cursors, error) {

	if page.emptyCursor && page.previous {
		return Cursors{Next: page.cursor}, nil
	}

	query, err := gosq.Compile(`
		{{ .Query }}
		WHERE {{ .Column }} {{ [if] .Previous [then] < [else] > }} $1
		ORDER BY {{ .Column }}
		LIMIT $2`,
		map[string]any{
			"Query":    query,
			"Column":   column,
			"Previous": page.previous,
		},
	)
	args = append([]any{page.cursor, d.pageLimit + 1}, args...)
	err = pgxscan.Select(ctx, d.pool, dest, query, args...)
	if err != nil {
		return Cursors{}, err
	}

	cursors := d.processPage(dest, page)
	return cursors, nil
}

type Page struct {
	cursor   any
	previous bool

	emptyCursor bool
}

func NewPage(cursor any, previous bool) (Page, error) {
	page := Page{
		cursor:   cursor,
		previous: previous,
	}
	switch c := cursor.(type) {
	case float64:
		page.emptyCursor = c == 0.
	case string:
		page.emptyCursor = c == ""
	default:
		return page, ErrInvalidCursor
	}
	return page, nil
}

func DecodePage(data jsonutil.Value) (page Page, err error) {
	page.previous, err = data.Get("pagePrevious").Bool()
	if err != nil {
		return
	}
	page.cursor, err = data.Get("pageCursor").Float()
	if err == nil {
		page.emptyCursor = page.cursor == 0.
		return
	}
	page.cursor, err = data.Get("pageCursor").String()
	if err != nil {
		return
	}
	page.emptyCursor = page.cursor == ""
	return
}

func (p Page) toArgs(limit int) pgx.NamedArgs {
	return pgx.NamedArgs{
		"cursor":   p.cursor,
		"previous": p.previous,
		"limit":    limit,
	}
}

type Cursors struct {
	Previous any `json:"previous"`
	Next     any `json:"next"`
}

type CursorWrap[T any] struct {
	Cursors Cursors `json:"cursors"`
	Data    T       `json:"data"`
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
