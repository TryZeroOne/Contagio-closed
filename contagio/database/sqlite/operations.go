package sqlite

import (
	"contagio/contagio/cnc/utils"
	"contagio/contagio/config/logging"
)

func RemoveUser(login string) {
	statement, err := Db.Prepare("DELETE from users where login=?")
	if err != nil {
		logging.PrintWarning("Can't remove user: " + err.Error())
		return
	}
	statement.Exec(utils.Sha3(login))
}

func RemoveIp(ip string) {

	statement, err := Db.Prepare("DELETE from allowed where ip=?")
	if err != nil {
		return
	}
	statement.Exec(utils.Sha3(ip))
}

func AddIp(ip string) {

	statement, err := Db.Prepare("INSERT INTO allowed(ip) VALUES (?)")

	if err != nil {
		logging.PrintWarning("Can't add ip: " + err.Error())
		return
	}

	statement.Exec(utils.Sha3(ip))

}

func GetCount() int {
	query := "SELECT COUNT(*) FROM users"

	var count int
	err := Db.QueryRow(query).Scan(&count)
	if err != nil {
		logging.PrintError("Can't get count: " + err.Error())
		return 0
	}

	return count
}

func CheckUser(login, password string) bool {
	rows, err := Db.Query("SELECT id, login, password FROM users WHERE login=? AND password=?", login, password)

	if err != nil {
		logging.PrintWarning("Can't check user: " + err.Error())
		return false
	}

	defer rows.Close()
	var id int
	var _login string
	var _password string

	for rows.Next() {
		rows.Scan(&id, &_login, &_password)
		return true
	}
	return false
}

func CheckIp(ip string) bool {
	rows, err := Db.Query("SELECT id, ip FROM allowed WHERE ip=?", utils.Sha3(ip))

	if err != nil {
		logging.PrintWarning("Can't check ip: " + err.Error())
		return false
	}

	defer rows.Close()
	var id int
	var _ip string

	for rows.Next() {
		rows.Scan(&id, &_ip)
		return true
	}
	return false
}

func AddUser(login, password string) {
	statement, err := Db.Prepare("INSERT INTO users(login, password) VALUES (?, ?)")

	if err != nil {
		logging.PrintWarning("Can't add user: " + err.Error())
		return
	}

	statement.Exec(utils.Sha3(login), utils.Sha3(password))
}

func SetPid(pid int) {
	var old string
	err := Db.QueryRow("SELECT pid FROM pids").Scan(&old)
	if err == nil {
		statement, err := Db.Prepare("DELETE from pids where pid=? ")
		if err != nil {
			logging.PrintWarning("Can't remove old pid: " + err.Error())
			return
		}
		statement.Exec(old)

	}

	statement, err := Db.Prepare("INSERT INTO pids(pid) VALUES(?) ")
	if err != nil {
		logging.PrintWarning("Can't set pid: " + err.Error())
		return
	}
	statement.Exec(pid)
}

func GetPid() string {
	var result string
	err := Db.QueryRow("SELECT pid FROM pids").Scan(&result)
	if err != nil {
		logging.PrintWarning("Can't get pid: " + err.Error())
		return ""
	}

	return result
}

func AddSession(sess int) {
	var c string
	err := Db.QueryRow("SELECT count FROM sessions").Scan(&c)
	if err != nil {
		statement, err := Db.Prepare("INSERT INTO sessions(count) VALUES(?) ")
		if err != nil {
			logging.PrintWarning("Can't set session: " + err.Error())
			return
		}
		statement.Exec(sess)
		return
	}

	statement, err := Db.Prepare("UPDATE sessions SET count=?")
	if err != nil {
		logging.PrintWarning("Can't update session: " + err.Error())
		return
	}
	statement.Exec(sess)
}

func GetSessions() string {
	var result string
	err := Db.QueryRow("SELECT count FROM sessions").Scan(&result)
	if err != nil {
		logging.PrintWarning("Can't get count: " + err.Error())
		return ""
	}

	return result

}

func AddBot(bot int) {
	var c string
	err := Db.QueryRow("SELECT count FROM bots").Scan(&c)
	if err != nil {
		statement, err := Db.Prepare("INSERT INTO bots(count) VALUES(?) ")
		if err != nil {
			logging.PrintWarning("Can't set bot: " + err.Error())
			return
		}
		statement.Exec(bot)
		return
	}

	statement, err := Db.Prepare("UPDATE bots SET count=?")
	if err != nil {
		logging.PrintWarning("Can't update bot: " + err.Error())
		return
	}
	statement.Exec(bot)
}

func GetBots() string {
	var result string
	err := Db.QueryRow("SELECT count FROM bots").Scan(&result)
	if err != nil {
		logging.PrintWarning("Can't get count: " + err.Error())
		return ""
	}

	return result

}

func SetStats(inc, out string) {
	var c string
	err := Db.QueryRow("SELECT inc FROM stats").Scan(&c)
	if err != nil {
		statement, err := Db.Prepare("INSERT INTO stats(inc, out) VALUES(?, ?) ")
		if err != nil {
			logging.PrintWarning("Can't set stats: " + err.Error())
			return
		}
		statement.Exec(inc, out)
		return
	}

	statement, err := Db.Prepare("UPDATE stats SET inc=?, out=?")
	if err != nil {
		logging.PrintWarning("Can't update stats: " + err.Error())
		return
	}
	statement.Exec(inc, out)
}

func GetStats() (string, string) {

	rows, err := Db.Query("SELECT inc, out FROM stats")
	if err != nil {
		logging.PrintWarning("Can't get stats: " + err.Error())
		return "", ""
	}

	defer rows.Close()
	var inc string
	var out string

	for rows.Next() {
		rows.Scan(&inc, &out)
	}

	return inc, out
}
