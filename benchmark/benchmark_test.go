package main

import (
	"fmt"
	"math/rand"
	"testing"
)

// Mirrors the shuffle technique used by github.com/efectn/go-orm-benchmarks.
func orms(b *testing.B, fns map[string]func(*testing.B)) {
	order := []string{"Raw", "Phi", "Bun", "Ent", "GORM"}
	rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
	for _, name := range order {
		if fn, ok := fns[name]; ok {
			b.Run(name, fn)
		}
	}
}

func BenchmarkCreate(b *testing.B) {
	fmt.Println("\n── Create ──────────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawCreate,
		"Phi":  benchPhiCreate,
		"Ent":  benchEntCreate,
		"GORM": benchGORMCreate,
		"Bun":  benchBunCreate,
	})
}

func BenchmarkCreateMany(b *testing.B) {
	fmt.Println("\n── CreateMany ──────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawCreateMany,
		"Phi":  benchPhiCreateMany,
		"Ent":  benchEntCreateMany,
		"GORM": benchGORMCreateMany,
		"Bun":  benchBunCreateMany,
	})
}

func BenchmarkCreateManyAndReturn(b *testing.B) {
	fmt.Println("\n── CreateManyAndReturn ─────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawCreateManyAndReturn,
		"Phi":  benchPhiCreateManyAndReturn,
		"Ent":  benchEntCreateManyAndReturn,
		"GORM": benchGORMCreateManyAndReturn,
		"Bun":  benchBunCreateManyAndReturn,
	})
}

func BenchmarkFindUnique(b *testing.B) {
	fmt.Println("\n── FindUnique ──────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawFindUnique,
		"Phi":  benchPhiFindUnique,
		"Ent":  benchEntFindUnique,
		"GORM": benchGORMFindUnique,
		"Bun":  benchBunFindUnique,
	})
}

func BenchmarkFindFirst(b *testing.B) {
	fmt.Println("\n── FindFirst ───────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawFindFirst,
		"Phi":  benchPhiFindFirst,
		"Ent":  benchEntFindFirst,
		"GORM": benchGORMFindFirst,
		"Bun":  benchBunFindFirst,
	})
}

func BenchmarkFindMany(b *testing.B) {
	fmt.Println("\n── FindMany ────────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawFindMany,
		"Phi":  benchPhiFindMany,
		"Ent":  benchEntFindMany,
		"GORM": benchGORMFindMany,
		"Bun":  benchBunFindMany,
	})
}

func BenchmarkUpsert(b *testing.B) {
	fmt.Println("\n── Upsert ──────────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawUpsert,
		"Phi":  benchPhiUpsert,
		"Ent":  benchEntUpsert,
		"GORM": benchGORMUpsert,
		"Bun":  benchBunUpsert,
	})
}

func BenchmarkFindUniqueSelfRef2L(b *testing.B) {
	fmt.Println("\n── FindUniqueSelfRef2L ────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawFindUniqueSelfRef2L,
		"Phi":  benchPhiFindUniqueSelfRef2L,
		"Ent":  benchEntFindUniqueSelfRef2L,
		"GORM": benchGORMFindUniqueSelfRef2L,
		"Bun":  benchBunFindUniqueSelfRef2L,
	})
}

func BenchmarkCreateWithDeepSelect(b *testing.B) {
	fmt.Println("\n── CreateWithDeepSelect ────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawCreateWithDeepSelect,
		"Phi":  benchPhiCreateWithDeepSelect,
		"Ent":  benchEntCreateWithDeepSelect,
		"GORM": benchGORMCreateWithDeepSelect,
		"Bun":  benchBunCreateWithDeepSelect,
	})
}

func BenchmarkCreateManyAndReturnWithDeepSelect(b *testing.B) {
	fmt.Println("\n── CreateManyAndReturnWithDeepSelect ───────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawCreateManyAndReturnWithDeepSelect,
		"Phi":  benchPhiCreateManyAndReturnWithDeepSelect,
		"Ent":  benchEntCreateManyAndReturnWithDeepSelect,
		"GORM": benchGORMCreateManyAndReturnWithDeepSelect,
		"Bun":  benchBunCreateManyAndReturnWithDeepSelect,
	})
}

func BenchmarkUpsertWithDeepSelect(b *testing.B) {
	fmt.Println("\n── UpsertWithDeepSelect ────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawUpsertWithDeepSelect,
		"Phi":  benchPhiUpsertWithDeepSelect,
		"Ent":  benchEntUpsertWithDeepSelect,
		"GORM": benchGORMUpsertWithDeepSelect,
		"Bun":  benchBunUpsertWithDeepSelect,
	})
}

func BenchmarkDelete(b *testing.B) {
	fmt.Println("\n── Delete ────────────────────────────────────")
	orms(b, map[string]func(*testing.B){
		"Raw":  benchRawDelete,
		"Phi":  benchPhiDelete,
		"Ent":  benchEntDelete,
		"GORM": benchGORMDelete,
		"Bun":  benchBunDelete,
	})
}
