package main

import (
	"database/sql"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//type album struct {
//	Id     string  `json:"id"`
//	Title  string  `json:"title"`
//	Artist string  `json:"artist"`
//	Price  float64 `json:"price"`
//}
type client struct {
	Id      string  `jsong:"id"`
	Balance float64 `json:"balance"`
}

type usluga struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type operation struct {
	Id_client      string  `json:"Id_client"`
	Id_usluga      string  `json:"usluga"`
	Id_transaction string  `json:"transaction"`
	Price          float64 `json:"price"`
}

var Db *sql.DB
var connStr = "user=postgres password=hbdtkjy2012 dbname=postgres sslmode=disable"
var clients = make([]client, 0)
var reserved = make([]client, 0)
var operations = make([]operation, 0)
var uslugi = make([]usluga, 0)

func parseDB() {
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer Db.Close()
	var newClient client
	var newUsluga usluga
	rows, _ := Db.Query(`SELECT * FROM avito_users ORDER BY id`)
	rows_uslugi, _ := Db.Query(`SELECT * FROM uslugi ORDER BY id`)
	for rows.Next() {
		err := rows.Scan(&newClient.Id, &newClient.Balance)
		clients = append(clients, newClient)
		reserved = append(reserved, client{Id: newClient.Id, Balance: 0})
		Db.Exec(`INSERT INTO reserved_accounts VALUES($1,$2)`, newClient.Id, "0")
		if err != nil {
			log.Fatal(err)
		}
	}
	for rows_uslugi.Next() {
		err := rows_uslugi.Scan(&newUsluga.Id, &newUsluga.Name)
		uslugi = append(uslugi, newUsluga)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func Contains(x []client, y string) bool {
	//if len(x) == 0 {
	//	return false
	//}
	for _, i := range clients {
		if i.Id == y {
			return true
		}
	}
	return false
}
func ContainsOp(x []operation, y string) bool {
	for _, i := range operations {
		if i.Id_transaction == y {
			return true
		}
	}
	return false
}

func getAlbums(c *gin.Context) {
	//var newClient client
	// var newUsluga usluga
	// connStr := "user=postgres password=hbdtkjy2012 dbname=postgres sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	// rows, _ := db.Query(`SELECT * FROM avito_users`)
	// rows_uslugi, _ := db.Query(`SELECT * FROM uslugi`)
	// for rows.Next() {
	// 	err := rows.Scan(&newClient.Id, &newClient.Balance)
	// 	clients = append(clients, newClient)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// for rows_uslugi.Next() {
	// 	err := rows_uslugi.Scan(&newUsluga.Id, &newUsluga.Name)
	// 	uslugi = append(uslugi, newUsluga)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	c.IndentedJSON(http.StatusOK, clients)
}

func addFundsToClient(c *gin.Context) {
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer Db.Close()
	var AddingFunds client
	if err := c.BindJSON(&AddingFunds); err != nil {
		return
	}
	i, _ := strconv.Atoi(AddingFunds.Id)
	if AddingFunds.Balance > 0 {
		if Contains(clients, AddingFunds.Id) {
			clients[i-1].Balance += AddingFunds.Balance
			c.IndentedJSON(http.StatusCreated, clients[i-1])
			Db.Exec(`UPDATE avito_users SET balance=$1 WHERE id=$2 ORDER BY id`, clients[i-1].Balance, AddingFunds.Id)
		} else {
			clients = append(clients, AddingFunds)
			reserved = append(reserved, client{Id: AddingFunds.Id, Balance: 0})
			c.IndentedJSON(http.StatusCreated, AddingFunds)
			Db.Exec(`INSERT INTO avito_users VALUES($1,$2)`, AddingFunds.Id, AddingFunds.Balance)
			Db.Exec(`INSERT INTO reserved_accounts VALUES($1,$2)`, AddingFunds.Id, "0")
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "the amount must be positive"})
	}
}

// func getStats(c *gin.Context) {
// 	id := c.Param("id")
// 	var newTransact operation
// 	for _, a := range clients {
// 		if a.Id == id {
// 			rows, _ := Db.Query(`SELECT * FROM transactions WHERE id_of_client=$1`, id)
// 			for rows.Next() {
// 				err := rows.Scan(&newTransact.Id_transaction, &newTransact.Id_client,
// 					&newTransact.Id_usluga, &newTransact.Price)
// 				operations = append(operations, newTransact)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 			}
// 		}
// 	}
// 	c.IndentedJSON(http.StatusOK, operations)
// }
func reserveOp(c *gin.Context) {
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer Db.Close()
	var newOp operation
	if err = c.BindJSON(&newOp); err != nil {
		log.Fatal(err)
	}
	clientId, _ := strconv.Atoi(newOp.Id_client)
	OpId, _ := strconv.Atoi(newOp.Id_transaction)
	UslId, _ := strconv.Atoi(newOp.Id_usluga)
	if clients[clientId-1].Balance >= newOp.Price {
		clients[clientId-1].Balance -= newOp.Price
		_, err = Db.Exec(`UPDATE avito_users SET balance=$1 WHERE id=$2`, clients[clientId-1].Balance, clientId)
		if err != nil {
			panic(err)
		}
		reserved[clientId-1].Balance += newOp.Price
		_, err = Db.Exec(`UPDATE reserved_accounts SET balance=$1 WHERE id=$2`, reserved[clientId-1].Balance, clientId)
		if err != nil {
			panic(err)
		}

		_, err = Db.Exec(`INSERT INTO transactions VALUES($1,$2,$3,$4)`, OpId, clientId, newOp.Price, UslId)
		operations = append(operations, newOp)
		if err != nil {
			panic(err)
		}
		c.IndentedJSON(http.StatusOK, clients[clientId-1])
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not enough funds"})
	}
}

func acceptOp(c *gin.Context) {
	Db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer Db.Close()
	var newOp operation
	if err := c.BindJSON(&newOp); err != nil {
		panic(err)
	}
	clientId, _ := strconv.Atoi(newOp.Id_client)
	if ContainsOp(operations, newOp.Id_transaction) {
		reserved[clientId-1].Balance -= newOp.Price
		_, err = Db.Exec(`UPDATE reserved_accounts SET balance=$1 WHERE id=$2`, reserved[clientId-1].Balance, clientId)
		c.IndentedJSON(http.StatusOK, reserved[clientId-1])
	}
}

func getClientByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range clients {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "client not found"})
}

func main() {

	parseDB()
	router := gin.Default()

	router.GET("/clients", getAlbums)
	router.POST("/clients/addfunds/", addFundsToClient)
	router.GET("/clients/:id", getClientByID)
	router.POST("/clients/reserve", reserveOp)
	router.POST("/clients/accept", acceptOp)

	router.Run(":8080")
	Db.Close()
}
