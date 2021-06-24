package models

import (
    "reakgo/utility"
    "github.com/jmoiron/sqlx"
)

type Profile struct {
    Id int32
    Name string
    Location string
    Description string
    FbLink string
    TwitLink string
    GithubLink string
    UnsplshLink string
    DribbleLink string
    InstaLink string
    YtLink string
}

type ProfileModel struct {
	DB *sqlx.DB
}

func (profile ProfileModel) Fetch() (Profile, error) {
    // Locking down to single entry, no multi-user functionality for now
    var p Profile
    rows, err := utility.Db.Queryx("SELECT * FROM profile WHERE id = 1")
    if err != nil{
        return p, err
    }
    for rows.Next() {
        err = rows.StructScan(&p)
    }
    return p, err
}
