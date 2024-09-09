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

// BenchmarkResult holds the result of a benchmark.
type BenchmarkResult struct {
	Repository string
	Operation  string
	Duration   time.Duration
}

// Create a map to keep track of used emails
var usedEmails = map[string]bool{}

// Create a map to store created author IDs
var createdAuthorIDs = map[uuid.UUID]bool{}

// createRandomAuthor generates a random author for testing.
func createRandomAuthor(rng *rand.Rand) (string, string, string, *time.Time) {
	name := fmt.Sprintf("Author%d", rng.Intn(1000))
	bio := fmt.Sprintf("Bio%d", rng.Intn(1000))

	// Generate a unique email
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

// benchmarkCreate runs the CreateAuthor benchmark.
func benchmarkCreate(repo repositories.AuthorRepository, repoName string, count int, rng *rand.Rand) BenchmarkResult {
	start := time.Now()
	for i := 0; i < count; i++ {
		name, bio, email, dateOfBirth := createRandomAuthor(rng)
		id, err := repo.CreateAuthor(context.Background(), name, bio, email, dateOfBirth)
		if err != nil {
			log.Fatalf("[%s] Failed to create author: %v", repoName, err)
		}
		createdAuthorIDs[id] = true // Store the created ID
	}
	return BenchmarkResult{Repository: repoName, Operation: "CreateAuthor", Duration: time.Since(start)}
}

// benchmarkGet runs the GetAuthor benchmark.
func benchmarkGet(repo repositories.AuthorRepository, repoName string) BenchmarkResult {
	start := time.Now()
	for id := range createdAuthorIDs { // Use IDs that were created
		_, err := repo.GetAuthor(context.Background(), id)
		if err != nil {
			log.Fatalf("[%s] Failed to get author: %v", repoName, err)
		}
	}
	return BenchmarkResult{Repository: repoName, Operation: "GetAuthor", Duration: time.Since(start)}
}

// benchmarkList runs the ListAuthors benchmark.
func benchmarkList(repo repositories.AuthorRepository, repoName string) BenchmarkResult {
	start := time.Now()
	_, err := repo.ListAuthors(context.Background())
	if err != nil {
		log.Fatalf("[%s] Failed to list authors: %v", repoName, err)
	}
	return BenchmarkResult{Repository: repoName, Operation: "ListAuthors", Duration: time.Since(start)}
}

// benchmarkDelete runs the DeleteAuthor benchmark.
func benchmarkDelete(repo repositories.AuthorRepository, repoName string) BenchmarkResult {
	start := time.Now()
	for id := range createdAuthorIDs { // Use IDs that were created
		err := repo.DeleteAuthor(context.Background(), id)
		if err != nil {
			log.Fatalf("[%s] Failed to delete author: %v", repoName, err)
		}
	}
	return BenchmarkResult{Repository: repoName, Operation: "DeleteAuthor", Duration: time.Since(start)}
}

// benchmarkUpdate runs the UpdateAuthor benchmark.
func benchmarkUpdate(repo repositories.AuthorRepository, repoName string, rng *rand.Rand) BenchmarkResult {
	start := time.Now()
	for id := range createdAuthorIDs { // Use IDs that were created
		name, bio, email, dateOfBirth := createRandomAuthor(rng)
		err := repo.UpdateAuthor(context.Background(), id, name, bio, email, dateOfBirth)
		if err != nil {
			log.Fatalf("[%s] Failed to update author: %v", repoName, err)
		}
	}
	return BenchmarkResult{Repository: repoName, Operation: "UpdateAuthor", Duration: time.Since(start)}
}

// benchmarkGetAuthorsByBirthdateRange runs the GetAuthorsByBirthdateRange benchmark.
func benchmarkGetAuthorsByBirthdateRange(repo repositories.AuthorRepository, repoName string, startDate, endDate time.Time) BenchmarkResult {
	start := time.Now()
	_, err := repo.GetAuthorsByBirthdateRange(context.Background(), startDate, endDate)
	if err != nil {
		log.Fatalf("[%s] Failed to get authors by birthdate range: %v", repoName, err)
	}
	return BenchmarkResult{Repository: repoName, Operation: "GetAuthorsByBirthdateRange", Duration: time.Since(start)}
}

func performBenchmarks(mongoRepo, gormRepo repositories.AuthorRepository) {
	// Define the number of test iterations
	testCount := 100 // Adjust the count as needed for your testing

	// Create a new random generator with a seed
	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Create maps to store results for side-by-side comparison
	results := map[string]map[string]BenchmarkResult{
		"MongoDB": {},
		"GORM":    {},
	}

	// Define date range for GetAuthorsByBirthdateRange
	startDate := time.Now().AddDate(-5, 0, 0) // 5 years ago
	endDate := time.Now()

	// Run benchmarks for MongoDB repository
	log.Println("Running benchmarks for MongoDB repository...")
	results["MongoDB"]["CreateAuthor"] = benchmarkCreate(mongoRepo, "MongoDB", testCount, rng)
	results["MongoDB"]["GetAuthor"] = benchmarkGet(mongoRepo, "MongoDB")
	results["MongoDB"]["ListAuthors"] = benchmarkList(mongoRepo, "MongoDB")
	results["MongoDB"]["DeleteAuthor"] = benchmarkDelete(mongoRepo, "MongoDB")
	results["MongoDB"]["UpdateAuthor"] = benchmarkUpdate(mongoRepo, "MongoDB", rng)
	results["MongoDB"]["GetAuthorsByBirthdateRange"] = benchmarkGetAuthorsByBirthdateRange(mongoRepo, "MongoDB", startDate, endDate)

	// Run benchmarks for GORM repository
	log.Println("Running benchmarks for GORM repository...")
	results["GORM"]["CreateAuthor"] = benchmarkCreate(gormRepo, "GORM", testCount, rng)
	results["GORM"]["GetAuthor"] = benchmarkGet(gormRepo, "GORM")
	results["GORM"]["ListAuthors"] = benchmarkList(gormRepo, "GORM")
	results["GORM"]["DeleteAuthor"] = benchmarkDelete(gormRepo, "GORM")
	results["GORM"]["UpdateAuthor"] = benchmarkUpdate(gormRepo, "GORM", rng)
	results["GORM"]["GetAuthorsByBirthdateRange"] = benchmarkGetAuthorsByBirthdateRange(gormRepo, "GORM", startDate, endDate)

	// Log results side by side and determine the winner
	var mongoTotal, gormTotal time.Duration
	for operation := range results["MongoDB"] {
		mongoDuration := results["MongoDB"][operation].Duration
		gormDuration := results["GORM"][operation].Duration
		difference := gormDuration - mongoDuration

		mongoTotal += mongoDuration
		gormTotal += gormDuration

		// Determine winner for each operation
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

	// Summarize overall results
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
	// Set up logging
	logFile, err := pkgs.SetUpLogger("MongoVsGorm.log")
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}
	defer logFile.Close()

	// Set up GORM database connection
	gormDB, err := gorm.Open(postgres.Open("postgresql://postgres:postgresSqlcVsGormSecret@localhost:5434/SqlcVsGorm_GORM"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to GORM DB: %v", err)
	}

	// Auto migrate GORM schema
	err = gormDB.AutoMigrate(&types.Author{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Create the GORM repository
	gormRepo := repositories.NewGORMAuthorRepository(gormDB)

	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	mongoDB := client.Database("MongoDB_Test")
	mongoRepo := repositories.NewMongoAuthorRepository(mongoDB)

	// Perform benchmarks using the repositories
	performBenchmarks(mongoRepo, gormRepo)
}
