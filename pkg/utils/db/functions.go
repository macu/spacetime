package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"spacetime/pkg/utils/types"
)

// Returns a placeholder representing the given arg in args,
// adding the arg to args if not already present.
func Arg(args *[]interface{}, arg interface{}) string {
	for i := 0; i < len(*args); i++ {
		if (*args)[i] == arg {
			return "$" + strconv.Itoa(i+1)
		}
	}
	*args = append(*args, arg)
	return "$" + strconv.Itoa(len(*args))
}

// Returns the "= $N" or "IS NULL" part of an equality condition
// where the operand may be null.
func Eq(col string, args *[]interface{}, arg interface{}) string {
	if types.IsNil(arg) {
		return col + " IS NULL"
	}
	return col + " = " + Arg(args, arg)
}

func In(col string, args *[]interface{}, values interface{}) string {
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice {
		panic("values must be a slice")
	}

	var out strings.Builder
	out.WriteString(col + " IN (")

	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			out.WriteString(",")
		}
		out.WriteString(Arg(args, v.Index(i).Interface()))
	}

	out.WriteString(")")
	return out.String()
}

// create a VALUES (), (), ... postgres string using argument placeholders
func ArgValuesMap(args *[]interface{}, values [][]interface{}) string {
	var out = `VALUES `

	for i := 0; i < len(values); i++ {
		if i > 0 {
			out += `,`
		}

		out += `(`

		for j := 0; j < len(values[i]); j++ {
			if j > 0 {
				out += `,`
			}

			var v = values[i][j]
			out += Arg(args, v)

			// include type casts with placeholders
			switch v.(type) {
			case int, uint, int64:
				out += `::int`
			case string:
				out += `::text`
			}
		}

		out += `)`
	}

	return out
}

func InTransaction(db *sql.DB, f func(*sql.Tx) error) error {

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	err = f(tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("rollback: %v; on run function: %w", rbErr, err)
		}
		return fmt.Errorf("run function in transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("rollback: %v; on commit: %w", rbErr, err)
		}
		return fmt.Errorf("commit: %w", err)
	}

	return nil

}
