package appstate

import (
    "ghershon/internal/storage"
)

type AppServices struct {
    DatabaseSrv *sql_l.DatabaseService
    KeySecret   []byte
}
