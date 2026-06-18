//file to setup all necessary environment init: db, observability, web launcher

import (
	"context"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)
