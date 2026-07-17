package main

import (
	"fmt"
	"testing"
)

func BenchmarkCreate(b *testing.B) {
	fmt.Println("\n── Create ──────────────────────────────────────")
	b.Run("Ent", benchEntCreate)
	b.Run("GORM", benchGORMCreate)
	b.Run("Valkyrie", benchValkyrieCreate)
	b.Run("Raw", benchRawCreate)
}

func BenchmarkCreateMany(b *testing.B) {
	fmt.Println("\n── CreateMany ──────────────────────────────────")
	b.Run("Ent", benchEntCreateMany)
	b.Run("GORM", benchGORMCreateMany)
	b.Run("Valkyrie", benchValkyrieCreateMany)
	b.Run("Raw", benchRawCreateMany)
}

func BenchmarkCreateManyAndReturn(b *testing.B) {
	fmt.Println("\n── CreateManyAndReturn ─────────────────────────")
	b.Run("Ent", benchEntCreateManyAndReturn)
	b.Run("GORM", benchGORMCreateManyAndReturn)
	b.Run("Valkyrie", benchValkyrieCreateManyAndReturn)
	b.Run("Raw", benchRawCreateManyAndReturn)
}

func BenchmarkFindUnique(b *testing.B) {
	fmt.Println("\n── FindUnique ──────────────────────────────────")
	b.Run("Ent", benchEntFindUnique)
	b.Run("GORM", benchGORMFindUnique)
	b.Run("Valkyrie", benchValkyrieFindUnique)
	b.Run("Raw", benchRawFindUnique)
}

func BenchmarkFindFirst(b *testing.B) {
	fmt.Println("\n── FindFirst ───────────────────────────────────")
	b.Run("Ent", benchEntFindFirst)
	b.Run("GORM", benchGORMFindFirst)
	b.Run("Valkyrie", benchValkyrieFindFirst)
	b.Run("Raw", benchRawFindFirst)
}

func BenchmarkFindMany(b *testing.B) {
	fmt.Println("\n── FindMany ────────────────────────────────────")
	b.Run("Ent", benchEntFindMany)
	b.Run("GORM", benchGORMFindMany)
	b.Run("Valkyrie", benchValkyrieFindMany)
	b.Run("Raw", benchRawFindMany)
}

func BenchmarkUpsert(b *testing.B) {
	fmt.Println("\n── Upsert ──────────────────────────────────────")
	b.Run("Ent", benchEntUpsert)
	b.Run("GORM", benchGORMUpsert)
	b.Run("Valkyrie", benchValkyrieUpsert)
	b.Run("Raw", benchRawUpsert)
}
