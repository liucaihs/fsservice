package remote

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	memdb  *sql.DB
	wrlock *sync.RWMutex
}

func GetModel() *Model {
	var model Model
	model.wrlock = new(sync.RWMutex)
	model.init()
	return &model
}

func (m *Model) init() {
	var err error
	m.memdb, err = sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		log.Fatalln("open db:", err)
	}
	_, err = m.memdb.Exec(`
                    CREATE TABLE "client_info" (
                    "id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
					"client_id" INTEGER NOT NULL,
                    "phone"  TEXT(16),
                    "imei"  TEXT(128),
					"iccid"  TEXT(128) ,
					"imsi"  TEXT(128),
					"ip"  TEXT(64),
                    "time"  TEXT(64) NOT NULL
                    );`)
	if err != nil {
		log.Fatalln("create table:", err)
	}

}

func (m *Model) sqlErroLog(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (m *Model) rowExists(id uint32) bool {
	m.wrlock.RLock()
	stmt, err := m.memdb.Prepare("SELECT * FROM client_info where client_id = ?")
	if err != nil {
		log.Println(err)
		return false
	}
	rows, errq := stmt.Query(id)
	if errq != nil {
		log.Println(errq)
		return false
	}
	defer func() {
		rows.Close()
		m.wrlock.RUnlock()
	}()
	res := rows.Next()

	return res
}

func (m *Model) DelExists(id uint32) {
	m.wrlock.Lock()
	stmt, err := m.memdb.Prepare("DELETE FROM client_info where client_id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	_, errq := stmt.Exec(id)
	if errq != nil {
		log.Println(errq)
		return
	}
	defer m.wrlock.Unlock()

}

func (m *Model) SaveClientInfo(id uint32, phone interface{}, imei interface{}, iccid interface{}, imsi interface{}, ip interface{}) {
	if id == 0 {
		return
	}

	curtime := time.Now().Format("2006-01-02 15:04:05")
	res := m.rowExists(id)
	m.wrlock.Lock()
	if res {
		stmtu, erru := m.memdb.Prepare("update client_info set phone=? , imei=?, iccid=?, imsi=?, ip=? , time=? where client_id = ?")
		m.sqlErroLog(erru)
		_, erru = stmtu.Exec(phone, imei, iccid, imsi, ip, curtime, id)
		m.sqlErroLog(erru)
	} else {
		stmti, erri := m.memdb.Prepare("INSERT INTO client_info(client_id, phone, imei, iccid, imsi, ip, time) values(?,?,?,?,?,?,?)")
		m.sqlErroLog(erri)
		_, erri = stmti.Exec(id, phone, imei, iccid, imsi, ip, curtime)
		m.sqlErroLog(erri)
	}
	defer m.wrlock.Unlock()
}

func (m *Model) GetClintOne(id uint32) map[string]interface{} {
	sqlstr := "SELECT * FROM client_info where"
	if id == 0 {
		return map[string]interface{}{}
	} else {
		sqlstr = sqlstr + " client_id=" + strconv.Itoa(int(id))
	}

	return m.GetClintRow(sqlstr)
}

func (m *Model) GetClintInfo(phone string, imei string) map[string]interface{} {
	sqlstr := "SELECT * FROM client_info where 1=1"
	if len(phone) == 0 && len(imei) == 0 {
		return map[string]interface{}{}
	} else if len(phone) > 0 {
		sqlstr = sqlstr + " and phone='" + phone + "'"
	} else if len(imei) > 0 {
		sqlstr = sqlstr + " and imei='" + imei + "'"
	}

	return m.GetClintRow(sqlstr)
}

func (m *Model) GetClintRow(sqlstr string) map[string]interface{} {
	result := map[string]interface{}{"id": 0, "client_id": 0, "phone": "", "imei": "", "iccid": "", "imsi": "", "ip": "", "time": ""}

	m.wrlock.RLock()
	log.Println(sqlstr)

	row := m.memdb.QueryRow(sqlstr)

	var id int
	var client_id uint32
	var phones string
	var imeis string
	var time string
	var iccid interface{}
	var imsi interface{}
	var ip string

	row.Scan(&id, &client_id, &phones, &imeis, &iccid, &imsi, &ip, &time)
	log.Println(id, client_id, phones, imeis, iccid, imsi, ip, time)

	result["id"] = id
	result["client_id"] = client_id
	result["phone"] = phones
	result["imei"] = imeis
	if iccid != nil {
		result["iccid"] = string(iccid.([]byte))
	}
	if imsi != nil {
		result["imsi"] = string(imsi.([]byte))
	}

	result["ip"] = ip
	result["time"] = time
	log.Println(result)

	defer func() {
		m.wrlock.RUnlock()
	}()

	return result
}

func (m *Model) GetRowNums(sqlstr string) int32 {
	var nums int32
	m.wrlock.RLock()
	log.Println(sqlstr)
	sql := "SELECT count(*) as ct FROM client_info"
	row := m.memdb.QueryRow(sql)

	row.Scan(&nums)
	log.Println("GetRowNums", nums)

	defer func() {
		m.wrlock.RUnlock()
	}()

	return nums
}

func (m *Model) GetPageRec() (string, int32) {
	var counts int32 = m.GetRowNums("")
	if counts < 1 {
		return "[]", counts
	}
	var reset string = "["
	m.wrlock.RLock()
	rows, err := m.memdb.Query("SELECT * FROM client_info limit 0,10")
	m.sqlErroLog(err)
	result := map[string]interface{}{"id": 0, "client_id": 0, "phone": "", "imei": "", "iccid": "", "imsi": "", "ip": "", "time": ""}
	var id int
	var client_id uint32
	var phones string
	var imeis string
	var time string
	var iccid interface{}
	var imsi interface{}
	var ip string
	for rows.Next() {
		rows.Scan(&id, &client_id, &phones, &imeis, &iccid, &imsi, &ip, &time)
		result["id"] = id
		result["client_id"] = client_id
		result["phone"] = phones
		result["imei"] = imeis
		if iccid != nil {
			result["iccid"] = string(iccid.([]byte))
		}
		if imsi != nil {
			result["imsi"] = string(imsi.([]byte))
		}

		result["ip"] = ip
		result["time"] = time
		log.Println(result)
		datastr, err := json.Marshal(result)
		if err != nil {
			log.Println(id, err)
		}
		reset = reset + string(datastr) + ","

	}
	if len(reset) > 1 {
		reset = reset[0 : len(reset)-1]
	}
	reset = reset + "]"
	log.Println(reset)

	defer func() {
		m.wrlock.RUnlock()
		rows.Close()
	}()
	return reset, counts
}

func (m *Model) CloseConn() {
	if err := m.memdb.Close(); err != nil {
		log.Println("Err from m.memdb.DatabaseClose(): ", err)
	}

}

func (m *Model) ShowAllRec() {
	m.wrlock.RLock()
	rows, err := m.memdb.Query("SELECT * FROM client_info")
	m.sqlErroLog(err)

	for rows.Next() {
		var id int
		var client_id int
		var phone string
		var imei string
		var time string
		var iccid interface{}
		var imsi interface{}
		var ip string
		err = rows.Scan(&id, &client_id, &phone, &imei, &iccid, &imsi, &ip, &time)
		m.sqlErroLog(err)
		log.Println(id, client_id, phone, imei, iccid, imsi, ip, time)

	}
	defer func() {
		rows.Close()
		m.wrlock.RUnlock()
	}()

}
