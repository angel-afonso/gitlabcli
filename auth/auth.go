package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"

	"gitlab.com/angel-afonso/gitlabcli/utils"
	"go.etcd.io/bbolt"

	color "gopkg.in/gookit/color.v1"
)

const (
	applicationid = "fa19a133b14bbcaf20a6e5bf6a4e5666cbdf19d0e8ad4f106ba3dea235a1e16b"
	callback      = "http://localhost:7890"
)

// Session has a bbolt db with the authentication data
type Session struct {
	Token string
	Type  string
}

// OpenSession opens session database and returns session struct
func OpenSession() *Session {
	executableDir, _ := os.Executable()
	db, err := bbolt.Open(path.Join(path.Dir(executableDir), "session"), 0600, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	session, err := lookUpSession(db)

	if err != nil {
		session = storeToken(db, login())
	}

	return session
}

// LookUpSession search token in the database
// and return session struct
func lookUpSession(db *bbolt.DB) (session *Session, err error) {
	db.Update(func(tx *bbolt.Tx) error {
		bucket := new(bbolt.Bucket)
		bucket, err = tx.CreateBucketIfNotExists([]byte("session"))

		if err != nil {
			return err
		}

		token := bucket.Get([]byte("access_token"))
		tokenType := bucket.Get([]byte("token_type"))

		if len(token) == 0 || len(tokenType) == 0 {
			err = errors.New("no tokens")
			return err
		}

		session = &Session{
			Token: string(token),
			Type:  string(tokenType),
		}

		return err
	})

	return
}

func storeToken(db *bbolt.DB, data map[string]string) *Session {
	err := db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("session"))
		if err != nil {
			log.Fatal(err.Error())
		}

		err = bucket.Put([]byte("access_token"), []byte(data["access_token"]))
		err = bucket.Put([]byte("token_type"), []byte(data["token_type"]))

		return err
	})

	if err != nil {
		log.Fatal(err.Error())
	}
	color.Green.Light().Println("Login successful!\n")

	return &Session{
		Token: data["access_token"],
		Type:  data["token_type"],
	}
}

func login() map[string]string {
	color.Cyan.Println("Logging with gitlab")
	srv := &http.Server{Addr: "0.0.0.0:7890"}

	var data map[string]string
	spinner := utils.ShowSpinner()

	openBrowser()

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		spinner.Stop()
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		srv.Shutdown(context.TODO())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
		<script>
			function getHashParams() {
				var hashParams = {};
				var e,
					a = /\+/g,  // Regex for replacing addition symbol with a space
					r = /([^&;=]+)=?([^&;]*)/g,
					d = function (s) { return decodeURIComponent(s.replace(a, " ")); },
					q = window.location.hash.substring(1);

				while (e = r.exec(q))
				hashParams[d(e[1])] = d(e[2]);

				return hashParams;
			}

			var xhr = new XMLHttpRequest();
			xhr.open('POST', 'http://localhost:7890/token', true);
			xhr.setRequestHeader('Content-type', 'application/json');
			xhr.send(JSON.stringify(getHashParams()));
			
			window.close();
		</script>
	`)
	})

	srv.ListenAndServe()
	return data
}

func openBrowser() error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args,
		fmt.Sprintf("https://gitlab.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=token&scope=api",
			applicationid, callback),
	)

	return exec.Command(cmd, args...).Start()
}

// fmt.Sprintf("https://gitlab.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=token&state=%s&scope=api",
// 	applicationid, callback, "asd"),
