// Package core: object façade (specs_v2_3.md §5). core.Object is the
// immutable generic carrier; its authoritative fields are unexported kernel
// state, the zero value is invalid, and no mutator exists. The public core
// package deliberately does NOT forward the mint entry point — minting stays
// with fixed domain constructors, ops results, session actions, and the
// hypothesis package, all calling internal/kernel.MintObject (§5.4).
package core

import "github.com/PithomLabs/phys/internal/kernel"

// Aliases over the kernel object representation (specs_v2_3.md §5).
type (
	// Object is the one immutable generic carrier accepted by every
	// generic symbolic operation.
	Object = kernel.Object
	// Kind is the closed physical-kind enum (§6, exactly 18 values).
	Kind = kernel.Kind
	// CorpusStatus is the human-curated corpus-status axis (§13.4); the
	// library never computes or revises it.
	CorpusStatus = kernel.CorpusStatus
)

// The exact 18 MVP kinds in listed order with stable ordinals (§6).
const (
	KindMass          = kernel.KindMass
	KindTime          = kernel.KindTime
	KindPosition      = kernel.KindPosition
	KindVelocity      = kernel.KindVelocity
	KindAcceleration  = kernel.KindAcceleration
	KindForce         = kernel.KindForce
	KindMomentum      = kernel.KindMomentum
	KindEnergy        = kernel.KindEnergy
	KindKineticEnergy = kernel.KindKineticEnergy
	KindRestMass      = kernel.KindRestMass
	KindThreeMomentum = kernel.KindThreeMomentum
	KindFourMomentum  = kernel.KindFourMomentum
	KindSpeedOfLight  = kernel.KindSpeedOfLight
	KindSpacetime     = kernel.KindSpacetime
	KindMinkowski     = kernel.KindMinkowski
	KindExpression    = kernel.KindExpression
	KindRelation      = kernel.KindRelation
	KindBranchSet     = kernel.KindBranchSet
)

// The exact five corpus statuses (§13.4).
const (
	CorpusNone        = kernel.CorpusNone
	CorpusEstablished = kernel.CorpusEstablished
	CorpusContested   = kernel.CorpusContested
	CorpusSuperseded  = kernel.CorpusSuperseded
	CorpusFalsified   = kernel.CorpusFalsified
)
