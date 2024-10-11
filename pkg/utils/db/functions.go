package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/logging"
	"treetime/pkg/utils/types"
)

// Returns a placeholder representing the given arg in args,
// adding the arg to args if not already present.
func ArgPlaceholder(arg interface{}, args *[]interface{}) string {
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
func EqCond(col string, arg interface{}, args *[]interface{}) string {
	if types.IsNil(arg) {
		return col + " IS NULL"
	}
	return col + " = " + ArgPlaceholder(arg, args)
}

// Returns the "= $N" or "IS NULL" part of an equality condition
// where the operand may be null, and the placeholder index is given.
func EqCondIndexed(col string, arg interface{}, index int) string {
	if types.IsNil(arg) {
		return col + " IS NULL"
	}
	return col + " = $" + strconv.Itoa(index)
}

func CreateArgsList(args *[]interface{}, values ...interface{}) string {
	out := ``
	for i := 0; i < len(values); i++ {
		if i > 0 {
			out += `,`
		}
		out += ArgPlaceholder(values[i], args)
	}
	return out
}

func CreateArgsListInt64s(args *[]interface{}, values ...int64) string {
	out := ``
	for i := 0; i < len(values); i++ {
		if i > 0 {
			out += `,`
		}
		out += ArgPlaceholder(values[i], args)
	}
	return out
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
			out += ArgPlaceholder(v, args)

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

func InTransaction(r *http.Request, db *sql.DB, f func(*sql.Tx) error) error {

	c := r.Context()
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("rollback: %v; on begin transaction: %w", rbErr, err)
		}
		return fmt.Errorf("begin transaction: %w", err)
	}

	err = f(tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("rollback: %v; on run function: %w", rbErr, err)
		}
		return fmt.Errorf("run function: %w", err)
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

func HandleInTransaction(r *http.Request, db *sql.DB, auth *ajax.Auth,
	f func(*sql.Tx) (interface{}, int, error)) (interface{}, int) {

	c := r.Context()
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			logging.LogError(r, auth, fmt.Errorf("rollback: %v; on begin transaction: %w", rbErr, err))
			return nil, http.StatusInternalServerError
		}
		logging.LogError(r, auth, fmt.Errorf("begin transaction: %w", err))
		return nil, http.StatusInternalServerError
	}

	response, statusCode, err := f(tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			logging.LogError(r, auth, fmt.Errorf("rollback: %v; on run function: %w", rbErr, err))
			return nil, http.StatusInternalServerError
		}
		logging.LogError(r, auth, fmt.Errorf("run function: %w", err))
		return nil, statusCode
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			logging.LogError(r, auth, fmt.Errorf("rollback: %v; on commit: %w", rbErr, err))
			return nil, http.StatusInternalServerError
		}
		logging.LogError(r, auth, fmt.Errorf("commit: %w", err))
		return nil, http.StatusInternalServerError
	}

	return response, statusCode
}
