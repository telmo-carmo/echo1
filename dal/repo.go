package dal

import (
	"database/sql"
	"encoding/json"
	"log"

	// db driver
	_ "github.com/mattn/go-sqlite3"
)

type Repo struct {
	db *sql.DB
}

// NewRepo return a new repo
func NewRepo(dbpath string) (*Repo, error) {
	log.Println("SQLITE db:", dbpath)
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	return &Repo{db}, nil

}

// GetVersion return DB current version
func (r *Repo) GetVersion() string {
	row := r.db.QueryRow("select sqlite_version() as v, datetime('now','localtime') as dt")
	var ver string
	var dt string

	row.Scan(&ver, &dt)
	log.Println("SQLITE Version:" + ver + " ; Now(): " + dt)
	return ver
}

// GetBonus returns a list of bonus as a json string
func (r *Repo) GetBonus(id string) string {
	myDb := r.db
	type Bonus struct {
		Ename string `json:"ename"`
		Job   string `json:"job"`
		Sal   int    `json:"sal"`
	}

	log.Printf("[INFO] GetBonus id: %s ", id)

	if id != "" && id[0] != '?' {
		var tempBn Bonus
		err := myDb.QueryRow("SELECT ename, job, sal FROM bonus where ename = ?", id).Scan(&tempBn.Ename, &tempBn.Job, &tempBn.Sal)
		if err != nil {
			if err == sql.ErrNoRows {
				respBy, _ := json.Marshal(map[string]string{"error": "Product not found"})
				return string(respBy)
			}
			log.Printf("[WARN] Query row error: %v", err)
			return "{}"
		} else {
			log.Printf("[INFO] Ename:%s, Job:%s, Sal:%d\n", tempBn.Ename, tempBn.Job, tempBn.Sal)
		}
		ba, err := json.Marshal(tempBn)
		if err != nil {
			log.Printf("[ERROR] json encode error: %v", err)
		}
		return string(ba)
	}

	rows, err := myDb.Query("SELECT ename, job, sal FROM bonus")
	myBonus := []Bonus{}

	for rows.Next() {
		var tempBn Bonus
		rows.Scan(&tempBn.Ename, &tempBn.Job, &tempBn.Sal)
		log.Printf("[INFO] Ename:%s, Job:%s, Sal:%d\n", tempBn.Ename, tempBn.Job, tempBn.Sal)
		myBonus = append(myBonus, tempBn)
	}
	rows.Close()

	if err = rows.Err(); err != nil {
		log.Printf("[WARN] Query rows error: %v", err)
	}

	ba, err := json.Marshal(myBonus)
	if err != nil {
		log.Printf("[ERROR] json encode error: %v", err)
	}
	return string(ba)
}

// Close closes the repo
func (r *Repo) Close() {
	r.db.Close()
}
