package main

import (
	"fmt"
	"testing"
)

func orms(b *testing.B, fns map[string]func(*testing.B)) {
	order := []string{"Raw", "Valkyrie", "Bun", "Ent", "GORM"}
	for _, name := range order {
		if fn, ok := fns[name]; ok {
			b.Run(name, fn)
		}
	}
}

func BenchmarkCreate(b *testing.B) {
	fmt.Println("\n── Create ──────────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawCreate,
		"Valkyrie": benchValkyrieCreate,
		"Ent":      benchEntCreate,
		"GORM":     benchGORMCreate,
		"Bun":      benchBunCreate,
	})
}

func BenchmarkCreateMany(b *testing.B) {
	fmt.Println("\n── CreateMany ──────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawCreateMany,
		"Valkyrie": benchValkyrieCreateMany,
		"Ent":      benchEntCreateMany,
		"GORM":     benchGORMCreateMany,
		"Bun":      benchBunCreateMany,
	})
}

func BenchmarkCreateManyAndReturn(b *testing.B) {
	fmt.Println("\n── CreateManyAndReturn ─────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawCreateManyAndReturn,
		"Valkyrie": benchValkyrieCreateManyAndReturn,
		"Ent":      benchEntCreateManyAndReturn,
		"GORM":     benchGORMCreateManyAndReturn,
		"Bun":      benchBunCreateManyAndReturn,
	})
}

func BenchmarkFindUnique(b *testing.B) {
	fmt.Println("\n── FindUnique ──────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawFindUnique,
		"Valkyrie": benchValkyrieFindUnique,
		"Ent":      benchEntFindUnique,
		"GORM":     benchGORMFindUnique,
		"Bun":      benchBunFindUnique,
	})
}

func BenchmarkFindFirst(b *testing.B) {
	fmt.Println("\n── FindFirst ───────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawFindFirst,
		"Valkyrie": benchValkyrieFindFirst,
		"Ent":      benchEntFindFirst,
		"GORM":     benchGORMFindFirst,
		"Bun":      benchBunFindFirst,
	})
}

func BenchmarkFindMany(b *testing.B) {
	fmt.Println("\n── FindMany ────────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawFindMany,
		"Valkyrie": benchValkyrieFindMany,
		"Ent":      benchEntFindMany,
		"GORM":     benchGORMFindMany,
		"Bun":      benchBunFindMany,
	})
}

func BenchmarkUpsert(b *testing.B) {
	fmt.Println("\n── Upsert ──────────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawUpsert,
		"Valkyrie": benchValkyrieUpsert,
		"Ent":      benchEntUpsert,
		"GORM":     benchGORMUpsert,
		"Bun":      benchBunUpsert,
	})
}

func BenchmarkReadDeepRelation(b *testing.B) {
	fmt.Println("\n── ReadDeepRelation ────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawReadDeepRelation,
		"Valkyrie": benchValkyrieReadDeepRelation,
		"Ent":      benchEntReadDeepRelation,
		"GORM":     benchGORMReadDeepRelation,
		"Bun":      benchBunReadDeepRelation,
	})
}

func BenchmarkCreateWithDeepSelect(b *testing.B) {
	fmt.Println("\n── CreateWithDeepSelect ────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawCreateWithDeepSelect,
		"Valkyrie": benchValkyrieCreateWithDeepSelect,
		"Ent":      benchEntCreateWithDeepSelect,
		"GORM":     benchGORMCreateWithDeepSelect,
		"Bun":      benchBunCreateWithDeepSelect,
	})
}

func BenchmarkCreateManyAndReturnWithDeepSelect(b *testing.B) {
	fmt.Println("\n── CreateManyAndReturnWithDeepSelect ───────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawCreateManyAndReturnWithDeepSelect,
		"Valkyrie": benchValkyrieCreateManyAndReturnWithDeepSelect,
		"Ent":      benchEntCreateManyAndReturnWithDeepSelect,
		"GORM":     benchGORMCreateManyAndReturnWithDeepSelect,
		"Bun":      benchBunCreateManyAndReturnWithDeepSelect,
	})
}

func BenchmarkUpsertWithDeepSelect(b *testing.B) {
	fmt.Println("\n── UpsertWithDeepSelect ────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":      benchRawUpsertWithDeepSelect,
		"Valkyrie": benchValkyrieUpsertWithDeepSelect,
		"Ent":      benchEntUpsertWithDeepSelect,
		"GORM":     benchGORMUpsertWithDeepSelect,
		"Bun":      benchBunUpsertWithDeepSelect,
	})
}
