// Package collections provides small, type-safe, allocation-conscious generic
// collections and slice helpers.
//
// Stack and Queue have useful zero values. Set is backed by map[T]struct{} and
// must be initialized before Add. The collection types are intentionally not
// synchronized; callers that share them across goroutines must provide their
// own synchronization.
package collections
