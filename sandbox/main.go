package main

import (
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("./.env")
	// fileBytes, _ := os.ReadFile("./schema.prisma")

	// rawString := string(fileBytes)

	// tokens := schema.ExtractTokens(rawString)
	// for i, token := range tokens {
	// 	fmt.Printf("%d: %+v\n", i, token)
	// }

	// schema, errs := schema.ParseSchema(rawString)
	// if len(errs) > 0 {
	// 	for _, err := range errs {
	// 		fmt.Println(err)
	// 	}
	// 	os.Exit(1)

	// }
	// mig, err := migration.GenerateMigration(schema)
	// if err != nil {
	// 	panic(err)
	// }
	// b, _ := json.MarshalIndent(schema, "", "  ")

	// os.WriteFile("result.json", b, 0644)
	// if strings.Contains(rawString, `provider = "sqlite"`) {
	// 	os.WriteFile("migrate_Sqlite.sql", []byte(mig), 0644)
	// } else {
	// 	os.WriteFile("migrate_Postgres.sql", []byte(mig), 0644)
	// }
	// fmt.Println(string(b))
	// fmt.Println(os.Getenv("DATABASE_DIRECT_URL"))

}
