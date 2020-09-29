package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"

	"app/internal/middleware"
	"app/internal/models"
	"app/internal/repository"
	"app/internal/sessionstore"
)

type User struct {
	Uid string `json:"uid" validate:"gte=1,lte=128"`
}

func SetupFirebase() (*auth.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// Firebase SDK のセットアップ
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_JSON_PATH"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app.Auth(context.Background())
}

func GetCsrfToken() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// 既にログイン済みセッションがあるかどうか確認し、あったらログイン済みユーザー情報を返す
		aSession, err := sessionstore.GetInstance().Store.Get(c.Request(), "auth-session")
		if err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, "session error.")
		}
		if aSession.Values["csrfToken"] != nil {
			client, err := SetupFirebase()
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			// ユーザー情報の取得
			uid := aSession.Values["uid"].(string)
			user, err := client.GetUser(context.Background(), uid)
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			// ユーザー名を最新に更新
			err = updateName(uid, user.UserInfo.DisplayName)
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			// セッション内で使うCSRFトークンを作る
			rand.Seed(time.Now().UnixNano())
			byteToken, err := getBinaryBySHA256WithKey(strconv.Itoa(rand.Intn(2147483647)), uid)
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}
			gToken := hex.EncodeToString(byteToken)

			// 既存のセッションを破棄して張りなおす
			aSession.Options.MaxAge = -1
			err = aSession.Save(c.Request(), c.Response())
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, "failed to delete session.")
			}

			// セッションを張る
			aSession2, err := sessionstore.GetInstance().Store.Get(c.Request(), "auth-session")
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, "session error.")
			}
			// 有効期間を1週間にする
			aSession2.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 7,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			aSession2.Values["csrfToken"] = gToken
			aSession2.Values["uid"] = uid
			err = aSession2.Save(c.Request(), c.Response())
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			// 新しいCSRFトークンとユーザー情報を返す
			m := map[string]interface{}{
				"csrfToken": gToken,
				"uid":       uid,
				"name":      user.UserInfo.DisplayName,
				"email":     user.UserInfo.Email,
			}
			jresults, err := json.Marshal(m)
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, string(jresults))
		}

		rand.Seed(time.Now().UnixNano())
		byteToken, err := getBinaryBySHA256WithKey(strconv.Itoa(rand.Intn(2147483647)), strconv.Itoa(rand.Intn(2147483647)))
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		token := hex.EncodeToString(byteToken)

		session, err := sessionstore.GetInstance().Store.Get(c.Request(), "login-session")
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "session error.")
		}
		// ログイン用の一時的なセッションなので1時間で切れるようにする
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		session.Values["csrfToken"] = token

		err = session.Save(c.Request(), c.Response())
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		m := map[string]interface{}{
			"csrfToken": token,
		}
		jresults, err := json.Marshal(m)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, string(jresults))

	}
}

func Signin() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(u); err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// セッションを確認してセッションを破棄する
		session, err := sessionstore.GetInstance().Store.Get(c.Request(), "login-session")
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "session error.")
		}
		sToken := session.Values["csrfToken"]
		session.Options.MaxAge = -1
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "failed to delete session.")
		}

		// ログイン用の一時CSRFトークンを検証する
		tokenHeader := c.Request().Header.Get("X-CSRF-Token")
		if tokenHeader == "" {
			log.Print("token not found.")
			return c.String(http.StatusBadRequest, "token not found.")
		}
		if tokenHeader != sToken {
			log.Print("session error.")
			return c.String(http.StatusBadRequest, "session error.")
		}

		// idTokenを検証する
		// クライアントから送られてきた JWT 取得
		authHeader := c.Request().Header.Get("Authorization")
		idToken := strings.Replace(authHeader, "Bearer ", "", 1)

		err = godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		// Firebase SDK のセットアップ
		client, err := SetupFirebase()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// JWT の検証
		_, err = client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// ユーザー情報の取得
		user, err := client.GetUser(context.Background(), u.Uid)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// ユーザー名を最新に更新
		err = updateName(u.Uid, user.UserInfo.DisplayName)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// セッション内で使うCSRFトークンを作る
		rand.Seed(time.Now().UnixNano())
		byteToken, err := getBinaryBySHA256WithKey(strconv.Itoa(rand.Intn(2147483647)), u.Uid)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		gToken := hex.EncodeToString(byteToken)

		// セッションを張る
		aSession, err := sessionstore.GetInstance().Store.Get(c.Request(), "auth-session")
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "session error.")
		}
		// 有効期間を1週間にする
		aSession.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		aSession.Values["csrfToken"] = gToken
		aSession.Values["uid"] = u.Uid

		err = aSession.Save(c.Request(), c.Response())
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		m := map[string]interface{}{
			"csrfToken": gToken,
		}
		jresults, err := json.Marshal(m)
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, string(jresults))
	}
}

func UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		isLegitimate := middleware.CheckSession(c)
		if !isLegitimate {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		uid := c.Get("uid")

		u := new(models.UserForUpdate)
		if err = c.Bind(u); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(u); err != nil {
			log.Print(err)
			return c.String(http.StatusBadRequest, err.Error())
		}
		// 本人以外は更新できない
		if uid != u.Uid {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		// Firebase SDK のセットアップ
		client, err := SetupFirebase()
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		// ユーザー情報の取得
		user, err := client.GetUser(context.Background(), u.Uid)

		// Firebase上の名前と送信されてきた名前が違う場合は更新する
		if user.UserInfo.DisplayName != u.Name {
			// Firebase上のユーザー名を更新
			params := (&auth.UserToUpdate{}).
				DisplayName(u.Name)
			_, err := client.UpdateUser(context.Background(), u.Uid, params)
			if err != nil {
				log.Print("error updating user name: %v\n", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}
			// ローカルDB上のユーザー名を更新
			err = updateName(u.Uid, u.Name)
			if err != nil {
				log.Print(err)
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
		// Firebase上のメールアドレスと送信されてきたメールアドレスが違う場合は更新する
		if user.UserInfo.Email != u.Email {
			// Firebase上のメールアドレスを更新
			params := (&auth.UserToUpdate{}).
				Email(u.Email)
			_, err := client.UpdateUser(context.Background(), u.Uid, params)
			if err != nil {
				log.Print("error updating user email: %v\n", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
		return c.NoContent(http.StatusOK)
	}
}

func Signout() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// セッションを確認してセッションを破棄する
		session, err := sessionstore.GetInstance().Store.Get(c.Request(), "auth-session")
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "session error.")
		}
		session.Options.MaxAge = -1
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "failed to delete session.")
		}
		return c.NoContent(http.StatusOK)
	}
}

func updateName(uid string, name string) error {
	err := repository.Db.Connect()
	if err != nil {
		return err
	}
	defer repository.Db.Close()

	count, err := repository.Db.Count("SELECT count(1) FROM users where uid = ?", uid)
	if err != nil {
		return err
	}

	if count > 0 {
		_, err := repository.Db.Exec("UPDATE users SET name = ? WHERE uid = ?", name, uid)
		if err != nil {
			return err
		}
	} else {
		_, err := repository.Db.Exec("INSERT INTO users (uid, name) VALUES (?, ?)", uid, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func getBinaryBySHA256WithKey(msg, key string) ([]byte, error) {
	r := sha256.Sum256([]byte(key))
	mac := hmac.New(sha256.New, r[:])
	_, err := mac.Write([]byte(msg))
	return mac.Sum(nil), err
}
