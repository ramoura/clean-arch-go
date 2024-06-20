package test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ramoura/clean-arch-go/pkg/infra/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

type TestDatabase struct {
	DbAddress  string
	Connection database.Connection
	container  testcontainers.Container
}

var db *sql.DB

//func TestMain(m *testing.M) {
//	var err error
//	db, err = setupDB()
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//
//	// Criação das tabelas e dados de teste
//	setupTestData(db)
//
//	// Run tests
//	code := m.Run()
//
//	// Cleanup
//	teardownTestData(db)
//
//	os.Exit(code)
//}

func SetupTestDatabase() *TestDatabase {

	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, connection, dbAddr, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	setupTestData(connection)
	cancel()

	return &TestDatabase{
		container:  container,
		Connection: connection,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, database.Connection, string, error) {

	dbName := "users"
	dbUser := "user"
	dbPass := "password"

	var env = map[string]string{
		"POSTGRES_PASSWORD": dbPass,
		"POSTGRES_USER":     dbUser,
		"POSTGRES_DB":       dbName,
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())

	connection, err := database.NewPgAdapter(dbUser, dbPass, dbAddr, dbName)
	if err != nil {
		return container, connection, dbAddr, fmt.Errorf("failed to establish database Connection: %v", err)
	}

	return container, connection, dbAddr, nil
}

func setupTestData(db database.Connection) {
	// Limpa a tabela de pessoas e insere dados de teste
	db.Exec("DROP TABLE IF EXISTS pessoas")
	db.Exec(`CREATE TABLE IF NOT EXISTS PESSOAS (
					ID VARCHAR(36)  CONSTRAINT ID_PK PRIMARY KEY,
					APELIDO VARCHAR(32),
					NOME VARCHAR(100),
					NASCIMENTO CHAR(10),
					STACK VARCHAR(1024),
					BUSCA_TRGM TEXT GENERATED ALWAYS AS (NOME || APELIDO || STACK) STORED
				)`)
}
