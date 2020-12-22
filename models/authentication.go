package models

import (
    "database/sql"
    "reakgo/utility"
)

type Authentication struct {
    Id int32
    Email string
    Password string
}

type AuthenticationModel struct {
	DB *sql.DB
}

func (auth AuthenticationModel) All() ([]Authentication, error) {
    rows, err := utility.Db.Query("SELECT * FROM authentication")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var bks []Authentication

    for rows.Next() {
        var bk Authentication

        err := rows.Scan(&bk.Id, &bk.Email, &bk.Password)
        if err != nil {
            return nil, err
        }

        bks = append(bks, bk)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return bks, nil
}
