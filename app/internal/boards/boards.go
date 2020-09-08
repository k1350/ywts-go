package boards

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"crypto/md5"
	"encoding/hex"

	"github.com/labstack/echo/v4"

	"app/internal/middleware"
	"app/internal/models"
	"app/internal/repository"
)

const constAdmin = 1
const constUser = 2

func GetBoards() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		uid := c.Get("uid")

		sqlString := `
		SELECT 
		 boards.id, 
		 boards.title, 
		 boards.created, 
		 boards.updated, 
		 CASE roles.description WHEN "admin" THEN 1 ELSE 0 END AS is_admin, 
		 boards.code 
		FROM boards 
		INNER JOIN board_members ON boards.id = board_members.board_id 
		LEFT JOIN roles ON board_members.role = roles.role 
		WHERE board_members.uid = ? 
		ORDER BY boards.id DESC 
		`

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		results, err := repository.Db.Fetch(sqlString, uid)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, results)
	}
}

func GetBoardIdByCode() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		b := new(models.BoardForJoin)
		b.Code = c.Param("code")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		cSqlString := `
		SELECT 
		 count(1)
		FROM boards 
		WHERE boards.code = ? `

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		count, err := repository.Db.Count(cSqlString, b.Code)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if count < 1 {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		sqlString := `
		SELECT 
		 id 
		FROM boards 
		WHERE boards.code = ? 
		ORDER BY id 
		LIMIT 1 
		`

		results, err := repository.Db.Fetch(sqlString, b.Code)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, results)
	}
}

func GetBoard() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		uid := c.Get("uid")

		b := new(models.Board)
		b.Id = c.Param("id")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		authed, err := IsBoardAuthenticated(uid.(string), b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !authed {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		sqlString := `
		SELECT 
		 boards.*
		FROM boards 
		INNER JOIN board_members ON boards.id = board_members.board_id 
		WHERE board_members.uid = ? and boards.id = ? 
		limit 1`

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		results, err := repository.Db.Fetch(sqlString, uid, b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, results)
	}
}

func AddBoard() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		isLegitimate := middleware.CheckSession(c)
		if !isLegitimate {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		uid := c.Get("uid")

		b := new(models.BoardForAdd)
		if err = c.Bind(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		err = repository.Db.Transaction(func(tx *sql.Tx) error {
			created := time.Now().UTC().Format("2006-01-02 15:04:05")

			sqlString := `
			INSERT INTO boards
			(title, code, created, updated)
			VALUES
			(?, "-1", ?, ?)
			`

			result, err := tx.Exec(sqlString, b.Title, created, created)
			if err != nil {
				return err
			}
			bid, err := result.LastInsertId()
			if err != nil {
				return err
			}

			b := []byte(strconv.FormatInt(bid, 10) + created)
			md5 := md5.Sum(b)
			code := hex.EncodeToString(md5[:])

			bSqlString := `
			UPDATE boards 
			SET code = ? 
			WHERE id = ? 
			`

			_, err = tx.Exec(bSqlString, code, bid)
			if err != nil {
				return err
			}

			mSqlString := `
			INSERT INTO board_members
			(board_id, uid, role)
			VALUES
			(?, ?, ?)
			`

			_, err = tx.Exec(mSqlString, bid, uid, constAdmin)

			return err
		})
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func DeleteBoard() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		isLegitimate := middleware.CheckSession(c)
		if !isLegitimate {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		uid := c.Get("uid")

		b := new(models.Board)
		b.Id = c.Param("id")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		authed, err := IsBoardAdmin(uid.(string), b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !authed {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		err = repository.Db.Transaction(func(tx *sql.Tx) error {
			mSqlString := `
			DELETE FROM board_members 
			WHERE board_id = ? 
			`

			_, err = tx.Exec(mSqlString, b.Id)
			if err != nil {
				return err
			}

			iSqlString := `
			DELETE FROM items  
			WHERE board_id = ? 
			`

			_, err = tx.Exec(iSqlString, b.Id)
			if err != nil {
				return err
			}

			sqlString := `
			DELETE FROM boards 
			WHERE id = ? 
			`

			_, err := tx.Exec(sqlString, b.Id)

			return err
		})
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func UpdateBoard() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		isLegitimate := middleware.CheckSession(c)
		if !isLegitimate {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		uid := c.Get("uid")

		b := new(models.Board)
		b.Id = c.Param("id")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		ba := new(models.BoardForAdd)
		if err = c.Bind(ba); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(ba); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		authed, err := IsBoardAdmin(uid.(string), b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !authed {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		updated := time.Now().UTC().Format("2006-01-02 15:04:05")

		sqlString := `
		UPDATE boards  
		SET title = ?, updated = ? 
		WHERE id = ? `

		_, err = repository.Db.Exec(sqlString, ba.Title, updated, b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func JoinBoard() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		isLegitimate := middleware.CheckSession(c)
		if !isLegitimate {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		uid := c.Get("uid")

		b := new(models.Board)
		b.Id = c.Param("id")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		bj := new(models.BoardForJoin)
		if err = c.Bind(bj); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// URLパラメータのコードとPOST内容のidが正しいことを確認する
		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		cSqlString := `
		SELECT count(1) 
		FROM boards 
		WHERE id = ? and code = ? 
		`
		count, err := repository.Db.Count(cSqlString, b.Id, bj.Code)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if count != 1 {
			return c.NoContent(http.StatusBadRequest)
		}
		repository.Db.Close()

		authed, err := IsBoardAuthenticated(uid.(string), b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		// もう既にメンバーだった場合はエラーにはしない
		if authed {
			return c.NoContent(http.StatusOK)
		}

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		sqlString := `
		INSERT INTO board_members 
		(board_id, uid, role)
		VALUES 
		(?, ?, ?) 
		`

		_, err = repository.Db.Exec(sqlString, b.Id, uid, constUser)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func GetBoardMembers() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		b := new(models.Board)
		b.Id = c.Param("id")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		sqlString := `
		SELECT 
		 boards.title, 
		 board_members.id, 
		 board_members.uid, 
		 users.name,  
		 board_members.role,
		 roles.description as role_description
		FROM boards  
		LEFT JOIN board_members ON boards.id = board_members.board_id 
		LEFT JOIN roles ON board_members.role = roles.role 
		LEFT JOIN users ON board_members.uid = users.uid 
		WHERE boards.id = ? 
		ORDER BY board_members.id 
		`

		results, err := repository.Db.Fetch(sqlString, b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, results)
	}
}

func DeleteBoardMember() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		isLegitimate := middleware.CheckSession(c)
		if !isLegitimate {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		uid := c.Get("uid")

		b := new(models.Board)
		b.Id = c.Param("id")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		u := new(models.User)
		u.Uid = c.Param("uid")
		if err = c.Validate(u); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		// 自分自身を除外する場合以外は管理者のみ実行可能
		if uid.(string) != u.Uid {
			authed, err := IsBoardAdmin(uid.(string), b.Id)
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}
			if !authed {
				return echo.NewHTTPError(http.StatusForbidden)
			}
		}

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		sqlString := `
		DELETE FROM board_members 
		WHERE board_id = ? AND uid = ?`

		_, err = repository.Db.Exec(sqlString, b.Id, u.Uid)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

// ボードのメンバーであるかどうか
func IsBoardAuthenticated(uid string, id string) (bool, error) {
	sqlString := `
		SELECT 
		 count(1)
		FROM boards 
		INNER JOIN board_members ON boards.id = board_members.board_id 
		WHERE board_members.uid = ? and boards.id = ? `

	err := repository.Db.Connect()
	if err != nil {
		return false, err
	}
	defer repository.Db.Close()

	count, err := repository.Db.Count(sqlString, uid, id)
	if err != nil {
		return false, err
	}

	if count < 1 {
		return false, nil
	}
	return true, nil
}

// ボードの管理者であるかどうか
func IsBoardAdmin(uid string, id string) (bool, error) {
	sqlString := `
		SELECT 
		 count(1)
		FROM boards 
		INNER JOIN board_members ON boards.id = board_members.board_id 
		LEFT JOIN roles ON board_members.role = roles.role 
		WHERE board_members.uid = ? and boards.id = ? and roles.role = ? `

	err := repository.Db.Connect()
	if err != nil {
		return false, err
	}
	defer repository.Db.Close()

	count, err := repository.Db.Count(sqlString, uid, id, constAdmin)
	if err != nil {
		return false, err
	}

	if count < 1 {
		return false, nil
	}
	return true, nil
}
