package models

import (
    "log"
    "reakgo/utility"
    "github.com/jmoiron/sqlx"
)


type FormAddView struct {
    Id int32
    Name string
    Address string
}


type FormAddViewModel struct {
    DB *sqlx.DB
}


func (form FormAddViewModel) Add(name string, address string) {
    utility.Db.MustExec("INSERT INTO simpleForm (name, address) VALUES (?, ?)", name, address)
}

func (form FormAddViewModel) View () ([]FormAddView, error) {
    var resultSet []FormAddView

    rows, err := utility.Db.Query("SELECT * FROM simpleForm")
    if (err != nil){
        log.Println(err)
    } else {
        defer rows.Close()

        for rows.Next(){
            var singleRow FormAddView
            err = rows.Scan(&singleRow.Id, &singleRow.Name, &singleRow.Address)
            if (err != nil){
                log.Println(err)
            } else {
               resultSet = append(resultSet,singleRow) 
            }
        }
    }
    return resultSet, err
}
