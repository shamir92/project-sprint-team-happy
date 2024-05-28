package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity/pb"
	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/handler"
	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/repository"
	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/service"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting server...")
	fmt.Println("Generating NIPs...")
	itNIPs := generateNIPs("615")
	nurseNIPs := generateNIPs("303")
	fmt.Println("Generated NIPs!")

	if _, err := os.Stat("usedAccount.db"); os.IsNotExist(err) {
		file, err := os.Create("usedAccount.db")
		if err != nil {
			handleErr(err)
		}
		defer file.Close()
		fmt.Println("Created usedAccount.db file")
	}

	db, err := sql.Open("sqlite3", "usedAccount.db")
	handleErr(err)
	migrateDb(db)
	resetCount(db)
	runDbJob(db)

	var nipMutex sync.Mutex
	var itMutex sync.Mutex

	handler := handler.NewRequestHandler(service.NewNipService(
		itNIPs,
		nurseNIPs,
		&itMutex,
		&nipMutex,
		repository.NewRepository(db),
	))

	server := grpc.NewServer()
	pb.RegisterNIPServiceServer(server, handler)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func runDbJob(db *sql.DB) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			updateCount(db)
			fmt.Println("Updated count")
		}
	}()
}

func migrateDb(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS used_it_account`)
	handleErr(err)
	_, err = db.Exec(`DROP TABLE IF EXISTS used_nurse_account`)
	handleErr(err)
	_, err = db.Exec(`DROP TABLE IF EXISTS meta_data`)
	handleErr(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS used_it_account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nip TEXT NOT NULL,
		password TEXT NOT NULL
	)`)
	handleErr(err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS used_nurse_account (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nip TEXT NOT NULL,
		password TEXT NOT NULL
	)`)
	handleErr(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS meta_data (
		key TEXT PRIMARY KEY,
		value INT
	);`)
	handleErr(err)

	_, err = db.Exec(`INSERT OR IGNORE INTO meta_data (key, value) VALUES ('itIndex', 0)`)
	handleErr(err)
	_, err = db.Exec(`INSERT OR IGNORE INTO meta_data (key, value) VALUES ('nurseIndex', 0)`)
	handleErr(err)
}

func resetCount(db *sql.DB) {
	_, err := db.Exec(`UPDATE meta_data SET value = 0 WHERE key = 'itIndex'`)
	handleErr(err)
	_, err = db.Exec(`UPDATE meta_data SET value = 0 WHERE key = 'nurseIndex'`)
	handleErr(err)
}

func updateCount(db *sql.DB) {
	var itCount, nurseCount int
	err := db.QueryRow("SELECT COUNT(1) FROM used_it_account").Scan(&itCount)
	handleErr(err)
	err = db.QueryRow("SELECT COUNT(1) FROM used_nurse_account").Scan(&nurseCount)
	handleErr(err)

	_, err = db.Exec(`UPDATE meta_data SET value = ? WHERE key = 'itIndex'`, itCount)
	handleErr(err)
	_, err = db.Exec(`UPDATE meta_data SET value = ? WHERE key = 'nurseIndex'`, nurseCount)
	handleErr(err)
}

func generateNIPs(prefix string) []uint64 {
	currentYear := time.Now().Year()
	res := []uint64{}

	for year := 2000; year <= currentYear; year++ {
		for month := 1; month <= 12; month++ {
			yearStr := fmt.Sprintf("%d", year)
			monthStr := fmt.Sprintf("%02d", month)
			for gender := 1; gender <= 2; gender++ {
				genderStr := strconv.Itoa(gender)
				for randomDigits := 0; randomDigits <= 9999; randomDigits++ {
					randomPart := strconv.Itoa(randomDigits)
					if len(randomPart) < 3 {
						randomPart = fmt.Sprintf("%03s", randomPart)
					}
					nipStr := prefix + genderStr + yearStr + monthStr + randomPart
					nip, err := strconv.ParseUint(nipStr, 10, 64)
					handleErr(err)
					res = append(res, nip)
				}
			}
		}
	}
	return res
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
