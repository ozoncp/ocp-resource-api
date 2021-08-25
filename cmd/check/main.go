package main

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func main() {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	//sql, i, err := psql.Select("id, user_id, type, status").From("resource").Where("id == ?").ToSql()
	//if err != nil {
	//	return
	//}
	//fmt.Printf("%v %v", sql, i)

	entityIds := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 102}
	query := psql.Select("id, user_id, type, status").From("resource").Where("id IN (?)", entityIds)
	sql, _, _ := query.ToSql()
	fmt.Printf("%v", sql)
}
