package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/lordofthemind/mongoVsGorm_GO/internals/repositories"
	"github.com/lordofthemind/mongoVsGorm_GO/internals/types"
	"github.com/lordofthemind/mongoVsGorm_GO/pkgs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/rand"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BenchmarkResult struct {
	Repository string
	Operation  string
	Duration   time.Duration
}

var usedEmails = map[string]bool{}
var createdAuthorIDs = map[uuid.UUID]bool{}

func createRandomAuthor(rng *rand.Rand) (string, string, string, *time.Time) {
	name := fmt.Sprintf("Author%d", rng.Intn(1000))
	bio := fmt.Sprintf("Bio%d", rng.Intn(1000))

	var email string
	for {
		email = fmt.Sprintf("author%d@example.com", rng.Intn(1000))
		if !usedEmails[email] {
			usedEmails[email] = true
			break
		}
	}

	dateOfBirth := time.Now().AddDate(-rng.Intn(60), 0, 0)
	return name, bio, email, &dateOfBirth
}

func benchmarkCreate(repo repositories.AuthorRepository, repoName string, count int, rng *rand.Rand) BenchmarkResult {
	start := time.Now()
	for i := 0; i < count; i++ {
		name, bio, email, dateOfBirth := createRandomAuthor(rng)
		id, err := repo.CreateAuthor(context.Background(), name, bio, email, dateOfBirth)
		if err != nil {
			log.Fatalf("[%s] Failed to create author: %v", repoName, err)
		}
		log.Printf("[%s] Created author with ID: %v", repoName, id) // Log the created UUID
		createdAuthorIDs[id] = true
	}
	return BenchmarkResult{Repository: repoName, Operation: "CreateAuthor", Duration: time.Since(start)}
}

func benchmarkGet(repo repositories.AuthorRepository, repoName string) BenchmarkResult {
	start := time.Now()
	for id := range createdAuthorIDs {
		_, err := repo.GetAuthor(context.Background(), id)
		if err != nil {
			log.Printf("[%s] Failed to get author (ID: %v): %v", repoName, id, err) // Handle missing records gracefully
			continue
		}
		log.Printf("[%s] Retrieved author with ID: %v", repoName, id) // Log retrieved UUID
	}
	return BenchmarkResult{Repository: repoName, Operation: "GetAuthor", Duration: time.Since(start)}
}

func benchmarkList(repo repositories.AuthorRepository, repoName string) BenchmarkResult {
	start := time.Now()
	_, err := repo.ListAuthors(context.Background())
	if err != nil {
		log.Fatalf("[%s] Failed to list authors: %v", repoName, err)
	}
	return BenchmarkResult{Repository: repoName, Operation: "ListAuthors", Duration: time.Since(start)}
}

func benchmarkDelete(repo repositories.AuthorRepository, repoName string) BenchmarkResult {
	start := time.Now()
	for id := range createdAuthorIDs {
		err := repo.DeleteAuthor(context.Background(), id)
		if err != nil {
			log.Printf("[%s] Failed to delete author (ID: %v): %v", repoName, id, err) // Handle missing records gracefully
			continue
		}
		log.Printf("[%s] Deleted author with ID: %v", repoName, id) // Log deleted UUID
	}
	return BenchmarkResult{Repository: repoName, Operation: "DeleteAuthor", Duration: time.Since(start)}
}

func benchmarkUpdate(repo repositories.AuthorRepository, repoName string, rng *rand.Rand) BenchmarkResult {
	start := time.Now()
	for id := range createdAuthorIDs {
		name, bio, email, dateOfBirth := createRandomAuthor(rng)
		err := repo.UpdateAuthor(context.Background(), id, name, bio, email, dateOfBirth)
		if err != nil {
			log.Printf("[%s] Failed to update author (ID: %v): %v", repoName, id, err) // Handle missing records gracefully
			continue
		}
		log.Printf("[%s] Updated author with ID: %v", repoName, id) // Log updated UUID
	}
	return BenchmarkResult{Repository: repoName, Operation: "UpdateAuthor", Duration: time.Since(start)}
}

func benchmarkGetAuthorsByBirthdateRange(repo repositories.AuthorRepository, repoName string, startDate, endDate time.Time) BenchmarkResult {
	start := time.Now()
	_, err := repo.GetAuthorsByBirthdateRange(context.Background(), startDate, endDate)
	if err != nil {
		log.Fatalf("[%s] Failed to get authors by birthdate range: %v", repoName, err)
	}
	return BenchmarkResult{Repository: repoName, Operation: "GetAuthorsByBirthdateRange", Duration: time.Since(start)}
}

func performBenchmarks(mongoRepo, gormRepo repositories.AuthorRepository) {
	testCount := 100

	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	results := map[string]map[string]BenchmarkResult{
		"MongoDB": {},
		"GORM":    {},
	}

	startDate := time.Now().AddDate(-5, 0, 0)
	endDate := time.Now()

	log.Println("Running benchmarks for MongoDB repository...")
	results["MongoDB"]["CreateAuthor"] = benchmarkCreate(mongoRepo, "MongoDB", testCount, rng)
	results["MongoDB"]["GetAuthor"] = benchmarkGet(mongoRepo, "MongoDB")
	results["MongoDB"]["ListAuthors"] = benchmarkList(mongoRepo, "MongoDB")
	results["MongoDB"]["DeleteAuthor"] = benchmarkDelete(mongoRepo, "MongoDB")
	results["MongoDB"]["UpdateAuthor"] = benchmarkUpdate(mongoRepo, "MongoDB", rng)
	results["MongoDB"]["GetAuthorsByBirthdateRange"] = benchmarkGetAuthorsByBirthdateRange(mongoRepo, "MongoDB", startDate, endDate)

	log.Println("Running benchmarks for GORM repository...")
	results["GORM"]["CreateAuthor"] = benchmarkCreate(gormRepo, "GORM", testCount, rng)
	results["GORM"]["GetAuthor"] = benchmarkGet(gormRepo, "GORM")
	results["GORM"]["ListAuthors"] = benchmarkList(gormRepo, "GORM")
	results["GORM"]["DeleteAuthor"] = benchmarkDelete(gormRepo, "GORM")
	results["GORM"]["UpdateAuthor"] = benchmarkUpdate(gormRepo, "GORM", rng)
	results["GORM"]["GetAuthorsByBirthdateRange"] = benchmarkGetAuthorsByBirthdateRange(gormRepo, "GORM", startDate, endDate)

	var mongoTotal, gormTotal time.Duration
	for operation := range results["MongoDB"] {
		mongoDuration := results["MongoDB"][operation].Duration
		gormDuration := results["GORM"][operation].Duration
		difference := gormDuration - mongoDuration

		mongoTotal += mongoDuration
		gormTotal += gormDuration

		winner := "MongoDB"
		if gormDuration < mongoDuration {
			winner = "GORM"
		}

		log.Printf("Operation: %s\n", operation)
		log.Printf("  MongoDB Duration: %v\n", mongoDuration)
		log.Printf("  GORM Duration    : %v\n", gormDuration)
		log.Printf("  Difference       : %v\n", difference)
		log.Printf("  Winner           : %s\n", winner)
		log.Println()
	}

	log.Println("Summary:")
	log.Printf("  Total MongoDB Time: %v\n", mongoTotal)
	log.Printf("  Total GORM Time   : %v\n", gormTotal)
	if mongoTotal < gormTotal {
		log.Println("Overall Winner: MongoDB")
	} else {
		log.Println("Overall Winner: GORM")
	}
}

func main() {
	logFile, err := pkgs.SetUpLogger("MongoVsGorm.log")
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}
	defer logFile.Close()

	gormDB, err := gorm.Open(postgres.Open("postgresql://postgres:MongoVsGormSecret@localhost:5432/MongoVsGorm_PGDB"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to GORM DB: %v", err)
	}

	err = pkgs.CheckAndEnableUUIDExtension(gormDB)
	if err != nil {
		log.Fatalf("Failed to check or enable UUID extension: %v", err)
	}

	err = gormDB.AutoMigrate(&types.Author{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	gormRepo := repositories.NewGORMAuthorRepository(gormDB)

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	mongoDB := client.Database("MongoVsGorm_MGDB")
	mongoRepo := repositories.NewMongoAuthorRepository(mongoDB)

	performBenchmarks(mongoRepo, gormRepo)
}
