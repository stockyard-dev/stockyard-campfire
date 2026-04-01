package server

import "github.com/stockyard-dev/stockyard-campfire/internal/license"

type Limits struct{ MaxCategories int; MaxThreads int }

var freeLimits = Limits{MaxCategories: 3, MaxThreads: 50}
var proLimits = Limits{MaxCategories: 0, MaxThreads: 0}

func LimitsFor(info *license.Info) Limits {
	if info != nil && info.IsPro() { return proLimits }
	return freeLimits
}
