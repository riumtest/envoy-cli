// Package planner translates a differ.Result into an ordered execution Plan.
//
// Each Change in a diff is mapped to a typed Action:
//   - Added keys become ActionSet
//   - Removed keys become ActionDelete
//   - Changed keys become ActionUpdate
//   - Unchanged keys become ActionNoop
//
// Example:
//
//	result := differ.Compare(base, target)
//	plan := planner.Build(result)
//	for _, action := range plan.Actions {
//		fmt.Println(action)
//	}
package planner
