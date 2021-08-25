package repo

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ozoncp/ocp-resource-api/internal/models"
	"github.com/rs/zerolog/log"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type repoPostgreSQL struct {
	connectionUrl string
	pool          *pgxpool.Pool
}

func (r *repoPostgreSQL) AddEntity(ctx context.Context, entity models.Resource) (*models.Resource, error) {
	conn, err := r.getConnection()
	if err != nil {
		return nil, err
	}
	sql, err := prepareResourceInsertQuery(entity)
	if err != nil {
		return nil, err
	}
	id := uint64(0)
	err = conn.QueryRow(ctx, sql, entity.UserId, entity.Type, entity.Status).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.DescribeEntity(ctx, id)
}

func (r *repoPostgreSQL) AddEntities(ctx context.Context, addEntityRequests []models.Resource) error {
	if len(addEntityRequests) == 0 {
		return errors.New("addEntityRequests list should not be empty")
	}
	conn, err := r.getConnection()
	if err != nil {
		return err
	}
	batch := pgx.Batch{}
	for i := range addEntityRequests {
		entity := addEntityRequests[i]
		sql, err := prepareResourceInsertQuery(entity)
		if err != nil {
			return err
		}
		batch.Queue(sql, entity.UserId, entity.Type, entity.Status)
	}
	exec, err := conn.SendBatch(ctx, &batch).Exec()
	if err != nil {
		return err
	}
	log.Info().Msgf("Batch completed: %v", exec.RowsAffected())
	return nil
}

func prepareResourceInsertQuery(entity models.Resource) (string, error) {
	sql, _, err := psql.Insert("resource").
		Columns("user_id", "type", "status").
		Values(entity.UserId, entity.Type, entity.Status).
		ToSql()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s RETURNING id", sql), nil
}

func (r *repoPostgreSQL) ListEntities(ctx context.Context, limit uint64, offset uint64) (*[]models.Resource, error) {
	conn, err := r.getConnection()
	if err != nil {
		return nil, err
	}
	queryBuilder := psql.Select("id, user_id, type, status").From("resource").OrderBy("id")
	if offset != 0 {
		queryBuilder = queryBuilder.Offset(offset)
	}
	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit)
	}

	sql, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	result := make([]models.Resource, 0, limit)
	for rows.Next() {
		resource := models.Resource{}
		if err := rows.Scan(&resource.Id, &resource.UserId, &resource.Type, &resource.Status); err != nil {
			continue
		}

		result = append(result, resource)
	}
	return &result, nil
}

func (r *repoPostgreSQL) DescribeEntity(ctx context.Context, entityId uint64) (*models.Resource, error) {
	conn, err := r.getConnection()
	if err != nil {
		return nil, err
	}
	queryBuilder := psql.Select("id, user_id, type, status").From("resource").Where("id = ?")
	sql, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, entityId)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errors.New("resource not found")
	}
	resource := models.Resource{}
	if err := rows.Scan(&resource.Id, &resource.UserId, &resource.Type, &resource.Status); err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *repoPostgreSQL) RemoveEntity(ctx context.Context, entityId uint64) error {
	conn, err := r.getConnection()
	if err != nil {
		return err
	}
	queryBuilder := psql.Delete("resource").Where("id = ?")
	sql, _, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = conn.Query(ctx, sql, entityId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoPostgreSQL) getConnection() (*pgxpool.Pool, error) {
	if r.pool != nil {
		return r.pool, nil
	}

	pool, err := pgxpool.Connect(context.Background(), r.connectionUrl)
	if err != nil {
		return nil, err
	}

	r.pool = pool
	return pool, nil
}

func NewRepoPostgreSQL(connectionUrl string) Repo {
	return &repoPostgreSQL{connectionUrl: connectionUrl}
}
