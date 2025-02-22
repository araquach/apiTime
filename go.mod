module github.com/araquach/apiTime

go 1.18

require (
	github.com/araquach/apiAuth v0.0.1
	github.com/araquach/dbService v0.0.1
	github.com/araquach/apiTeam v0.0.1
	github.com/jinzhu/now v1.1.5
	gorm.io/datatypes v1.2.0
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	gorm.io/driver/mysql v1.4.7 // indirect
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gorm.io/driver/postgres v1.5.2 // indirect
	gorm.io/gorm v1.25.1 // indirect
)

replace (
	github.com/araquach/apiAuth => ../apiAuth
	github.com/araquach/dbService => ../dbService
	github.com/araquach/apiTeam => ../apiTeam
)
