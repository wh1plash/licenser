package store

import (
	"context"
	"database/sql"
	"licenser/server/types"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (p PostgresStore) GetAppList(ctx context.Context) ([]*types.App, error) {
	time.Sleep(time.Second * 5)
	rows, err := p.db.QueryContext(ctx, "select * from apps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []*types.App{}
	for rows.Next() {
		app := new(types.App)
		err := rows.Scan(
			&app.ID,
			&app.Name,
			&app.CreatedAt,
			&app.Until,
		)
		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}

	if len(apps) == 0 {
		return nil, sql.ErrNoRows
	}

	return apps, nil

}

func (p PostgresStore) GetApp(ctx context.Context, name string) (*types.App, error) {
	rows, err := p.db.QueryContext(ctx, "select * from apps where name=$1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	apps := &types.App{}
	if err := rows.Scan(
		&apps.ID,
		&apps.Name,
		&apps.CreatedAt,
		&apps.Until,
	); err != nil {
		return nil, err
	}
	return apps, nil
}

func (p PostgresStore) Init() error {
	return p.createAppTable()
}

func (p PostgresStore) CreateApp() error {
	rows, _ := p.db.Query("select * from apps where name=$1", "App")
	if rows.Next() {
		return nil
	}

	app := &types.App{
		Name:      "App",
		CreatedAt: time.Now(),
		Until:     time.Now().AddDate(0, 1, 0),
	}

	_, err := p.InsertApp(context.Background(), app)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresStore) InsertApp(ctx context.Context, app *types.App) (*types.App, error) {
	query := `insert into apps
		(name, created_at, until)
		values($1, $2, $3)
		returning id, name, created_at, until
	`
	insApp := &types.App{}
	err := p.db.QueryRowContext(
		ctx,
		query,
		app.Name,
		app.CreatedAt,
		app.Until,
	).Scan(
		&insApp.ID,
		&insApp.Name,
		&insApp.CreatedAt,
		&insApp.Until,
	)
	if err != nil {
		return nil, err
	}

	return insApp, nil
}

func (p PostgresStore) createAppTable() error {
	query := `create table if not exists apps (
		id serial primary key,
		name varchar(50) Not Null,
		created_at timestamp,
		until timestamp,
		CONSTRAINT unique_app_name UNIQUE (name)
	);`

	_, err := p.db.Exec(query)
	if err != nil {
		return err
	}

	return err
}
