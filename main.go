package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", "postgres", "1923247", "testdb")
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return db
}

type Item struct {
	ID      int    `json:"id"`
	Purpose string `json:"purpose"`
	Price   string `json:"price"`
}

func createItem(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodPost {
		w.Header().Add("Allow", "POST")
		http.Error(w, "This method is not allowed!", http.StatusMethodNotAllowed)
		return
	}

	db := setupDB()
	log.Println("Db is working in createItem!")
	defer db.Close()

	item := Item{}
	json.NewDecoder(req.Body).Decode(&item)

	_, err := db.Exec("insert into test_table(purpose, price) values($1, $2)", item.Purpose, item.Price)

	if err != nil {
		log.Fatalln(err.Error())
	}

	var max_id int
	max_id_row := db.QueryRow("select max(id) from test_table")
	check_id_error := max_id_row.Scan(&max_id)

	if check_id_error != nil {
		log.Fatalln(err.Error())
	}

	item.ID = max_id
	json.NewEncoder(w).Encode(&item)
}

func getItem(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodGet {
		w.Header().Add("Allow", "GET")
		http.Error(w, "This method is not allowed!", http.StatusMethodNotAllowed)
		return
	}

	item_id := req.URL.Query().Get("id")

	if item_id == "" {
		http.Error(w, "Wrong id!", 404)
		return
	}
	res_id, err := strconv.Atoi(item_id)

	db := setupDB()
	log.Println("Db is working in getItem by id!")
	defer db.Close()

	var max_id int
	max_id_row := db.QueryRow("select max(id) from test_table")
	check_id_error := max_id_row.Scan(&max_id)

	if err != nil || res_id < 1 || check_id_error != nil || max_id < res_id {
		http.Error(w, "There is no such id!", 404)
		return
	}

	item := Item{}

	row := db.QueryRow("select * from test_table where id = $1", res_id)

	err = row.Scan(&item.ID, &item.Purpose, &item.Price)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "There is no any item with this id!", 404)
		return
	}

	json.NewEncoder(w).Encode(&item)

}

func getAllItems(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodGet {
		w.Header().Add("Allow", "GET")
		http.Error(w, "This method is not allowed!", http.StatusMethodNotAllowed)
		return
	}

	db := setupDB()
	log.Println("Db is working in getAllItems!")
	defer db.Close()

	items := []Item{}

	rows, err := db.Query("select * from test_table order by id")
	if err != nil {
		log.Fatalln(err.Error())
	}

	for rows.Next() {
		var id int
		var purpose string
		var price string

		err = rows.Scan(&id, &purpose, &price)
		if err != nil {
			log.Fatalln(err.Error())
		}

		items = append(items, Item{ID: id, Purpose: purpose, Price: price})
	}

	json.NewEncoder(w).Encode(&items)

}

func updateItem(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodPut {
		w.Header().Add("Allow", "DELETE")
		http.Error(w, "This method is not allowed!", http.StatusMethodNotAllowed)
		return
	}

	db := setupDB()
	log.Println("Db is working in updateItem!")
	defer db.Close()

	item_id := req.URL.Query().Get("id")

	if item_id == "" {
		http.Error(w, "There is no such id", http.StatusNotFound)
		return
	}

	res_id, err := strconv.Atoi(item_id)

	var max_id int
	max_id_row := db.QueryRow("select max(id) from test_table")
	check_id_error := max_id_row.Scan(&max_id)

	if err != nil || res_id < 1 || check_id_error != nil || max_id < res_id {
		http.Error(w, "There is no such id!", 404)
		return
	}

	item := Item{}
	json.NewDecoder(req.Body).Decode(&item)

	_, res_err := db.Exec("update test_table set purpose = $1, price = $2 where id = $3", item.Purpose, item.Price, res_id)
	if res_err != nil {
		log.Fatalln(res_err.Error())
	}

	//fmt.Println(result.RowsAffected())

	row := db.QueryRow("select * from test_table where id = $1", res_id)

	resItem := Item{}
	err = row.Scan(&resItem.ID, &resItem.Purpose, &resItem.Price)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "There is no any item with this id!", 404)
		return
	}

	//resItem.ID = res_id
	json.NewEncoder(w).Encode(&resItem)
}

func deleteItem(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodDelete {
		w.Header().Add("Allow", "DELETE")
		http.Error(w, "This method is not allowed!", http.StatusMethodNotAllowed)
		return
	}

	db := setupDB()
	log.Println("Db is working in deleteItem!")
	defer db.Close()

	item_id := req.URL.Query().Get("id")

	if item_id == "" {
		http.Error(w, "There is no such id", http.StatusNotFound)
		return
	}

	res_id, err := strconv.Atoi(item_id)

	var max_id int
	max_id_row := db.QueryRow("select max(id) from test_table")
	check_id_error := max_id_row.Scan(&max_id)

	if err != nil || res_id < 1 || check_id_error != nil || max_id < res_id {
		http.Error(w, "There is no such id!", 404)
		return
	}

	_, err = db.Exec("delete from test_table where id = $1", res_id)
	if err != nil {
		log.Fatalln(err.Error())
	}

	items := []Item{}

	rows, err := db.Query("select * from test_table order by id")
	if err != nil {
		log.Fatalln(err.Error())
	}

	for rows.Next() {
		var id int
		var purpose string
		var price string

		err = rows.Scan(&id, &purpose, &price)
		if err != nil {
			log.Fatalln(err.Error())
		}

		items = append(items, Item{ID: id, Purpose: purpose, Price: price})
	}

	json.NewEncoder(w).Encode(&items)
}

func checkerFunc(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(w, "Wrong url adress!", 404)
		return
	}

	w.Write([]byte("HomeFunc is working..."))
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", checkerFunc)
	router.HandleFunc("/item/create", createItem)
	router.HandleFunc("/item/get", getItem)
	router.HandleFunc("/item/getAll", getAllItems)
	router.HandleFunc("/item/update", updateItem)
	router.HandleFunc("/item/delete", deleteItem)

	log.Printf("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
