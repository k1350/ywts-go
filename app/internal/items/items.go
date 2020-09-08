package items

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"app/internal/boards"
	"app/internal/middleware"
	"app/internal/models"
	"app/internal/repository"
)

func GetBoardItems() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		uid := c.Get("uid")

		b := new(models.Board)
		b.Id = c.Param("id")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		authed, err := boards.IsBoardAuthenticated(uid.(string), b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !authed {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		sqlString := `
		SELECT 
		 items.*,
		 users.name as author
		FROM boards 
		INNER JOIN board_members ON boards.id = board_members.board_id 
		LEFT JOIN items on boards.id = items.board_id 
		LEFT JOIN users on users.uid = items.author_uid 
		WHERE board_members.uid = ? and boards.id = ? and items.id is not null 
		ORDER BY items.id DESC`

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

func GetBoardItem() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		uid := c.Get("uid")

		b := new(models.Board)
		b.Id = c.Param("id")
		if err = c.Validate(b); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		i := new(models.Item)
		i.Id = c.Param("iid")
		if err = c.Validate(i); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		authed, err := boards.IsBoardAuthenticated(uid.(string), b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !authed {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		sqlString := `
		SELECT 
		 items.*
		FROM boards 
		INNER JOIN board_members ON boards.id = board_members.board_id 
		LEFT JOIN items on boards.id = items.board_id 
		WHERE board_members.uid = ? and boards.id = ? and items.id = ? 
		limit 1`

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		results, err := repository.Db.Fetch(sqlString, uid, b.Id, i.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, results)
	}
}

func AddBoardItem() echo.HandlerFunc {
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

		i := new(models.Item)
		if err = c.Bind(i); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(i); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		authed, err := boards.IsBoardAuthenticated(uid.(string), b.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !authed {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		created := time.Now().UTC().Format("2006-01-02 15:04:05")

		sqlString := `
		INSERT INTO items 
		(board_id, y, w, t, author_uid, created, updated) 
		VALUES 
		(?, ?, ?, ?, ?, ?, ?) 
		`

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		_, err = repository.Db.Exec(sqlString, b.Id, i.Y, i.W, i.T, uid, created, created)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func UpdateBoardItem() echo.HandlerFunc {
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

		i := new(models.Item)
		if err = c.Bind(i); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		i.Id = c.Param("iid")
		if err = c.Validate(i); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		authed, err := isItemAuthenticated(uid.(string), b.Id, i.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !authed {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		updated := time.Now().UTC().Format("2006-01-02 15:04:05")

		sqlString := `
		UPDATE items 
		SET y = ?, w = ?, t = ?, updated = ? 
		WHERE id = ? `

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		_, err = repository.Db.Exec(sqlString, i.Y, i.W, i.T, updated, i.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func DeleteBoardItem() echo.HandlerFunc {
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

		i := new(models.Item)
		i.Id = c.Param("iid")
		if err = c.Validate(i); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}

		authed, err := isItemAuthenticated(uid.(string), b.Id, i.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !authed {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		sqlString := `
		DELETE FROM items  
		WHERE id = ? `

		err = repository.Db.Connect()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer repository.Db.Close()

		_, err = repository.Db.Exec(sqlString, i.Id)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}

func isItemAuthenticated(uid string, id string, iid string) (bool, error) {
	sqlString := `
		SELECT 
		 count(1)
		FROM items 
		INNER JOIN boards ON items.board_id = boards.id 
		INNER JOIN board_members ON items.board_id = board_members.board_id 
		WHERE items.author_uid = ? and items.board_id = ? and items.id = ? `

	err := repository.Db.Connect()
	if err != nil {
		return false, err
	}
	defer repository.Db.Close()

	count, err := repository.Db.Count(sqlString, uid, id, iid)
	if err != nil {
		return false, err
	}
	if count < 1 {
		return false, nil
	}
	return true, nil
}
