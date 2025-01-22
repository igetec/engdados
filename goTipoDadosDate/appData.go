package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
create table if not exists godatas (
	goDatasID bigserial,
	descricao varchar(40),
	dataString date not null,
	dataTime date not null,
	dataStringNull date,
	dataTimeNull date,
	constraint pk_godatas primary key (godatasid));
`

type St_GoDatas struct {
	Godatasid      int64      `db:"godatasid"`
	Descricao      string     `db:"descricao"`
	Datastring     string     `db:"datastring"`
	Datatime       time.Time  `db:"datatime"`
	DatastringNull *string    `db:"datastringnull"`
	DatatimeNull   *time.Time `db:"datatimenull"`
}

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres password=postdba dbname=igetec sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	//-- Salvando valores completos
	data01, _ := time.Parse("2006-01-02", "2025-01-20")
	godatas01 := St_GoDatas{
		Descricao:  "01",
		Datastring: "2025-01-21",
		Datatime:   data01,
	}

	_, _ = db.NamedExec(`INSERT INTO godatas 
        VALUES (default, :descricao, :datastring, :datatime )`, godatas01)

	//-- Recuperando TODOS os dados
	fmt.Println("Recuperando TODOS os dados")
	godatas := []St_GoDatas{}
	db.Select(&godatas, "SELECT * FROM godatas")

	for _, reg := range godatas {
		fmt.Println("Valores:", reg)
	}

	//-- Recuperando os dados por meio de uma data em String
	fmt.Println("Recuperando os dados por meio de uma data em String")
	godatasString := []St_GoDatas{}
	db.Select(&godatasString, `SELECT * 
							 FROM godatas WHERE datastring = $1`, "2025-01-21")

	for _, reg := range godatasString {
		fmt.Println("Valores:", reg)
	}
}
