package CLF

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type AbsInserter interface {
	Insert(goroutine int)
}

type Inserter struct {
	ParsedLogCh *chan ParsedLog
}

func (pi *Inserter) Insert(goroutine int) {
	log.Printf("Start inserter routine %d\n", goroutine)
	db, err := sql.Open("mysql", MYSQL_CRED)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	logCnt := 0
	tagCnt := 0
	for logItem := range *pi.ParsedLogCh {

		l := &logItem.Log
		tags := &logItem.Tags

		q := "INSERT INTO `log` (_time, log_raw, log_type) VALUES (?, ?, ?);"
		insert, err := db.Prepare(q)
		if err != nil {
			panic(err)
		}

		resp, err := insert.Exec(l.timestampString(), l.LogRaw, l.LogType)
		insert.Close()

		lastInsertLogId, err := resp.LastInsertId()
		if err != nil {
			fmt.Println(err)
			continue
		}

		l.ID = lastInsertLogId
		logCnt++
		for idx, tag := range *tags {
			q := "INSERT INTO `tag` (_key, _value, _type) VALUES (?, ?, ?);"
			insert, _ := db.Prepare(q)
			resp, err = insert.Exec(tag.Key, tag.Value, tag.Type)
			insert.Close()

			var tagId int64 = -1
			if err == nil {
				tagId, _ = resp.LastInsertId()
			} else {
				qu := "SELECT `_id` FROM `tag` WHERE _key = ? AND _value= ?;"
				row := db.QueryRow(qu, tag.Key, tag.Value)
				err = row.Scan(&tagId)
				if err != nil {
					panic(err)
				}
			}
			if tagId == -1 {
				panic("no tag id")
			}

			(*tags)[idx].ID = tagId

			q = "INSERT INTO `r_log_tag` (log_id, tag_id) VALUES (?, ?);"
			insert, _ = db.Prepare(q)
			resp, err = insert.Exec(l.ID, tagId)
			insert.Close()

			if err != nil {

			}
			tagCnt++
		}
	}
	log.Printf("Stop inserter routine %d, insert %d log with %d tag\n", goroutine, logCnt, tagCnt)
}
