package utils

//------------------------------------------------------------------------------
// Here we compare maps
//------------------------------------------------------------------------------

// compareMaps : getting a diff between 2 maps
func compareMaps(currentPath PropPath, map1, map2 map[string]interface{}, options *ComparisonOptions, fromSlice bool) (Comparison, error) {
	// the result from comparing the 2 maps
	thisComparison := map[string]interface{}{}

	// first, let's keep track of the keys we encounter
	checked := map[string]bool{}

	// let's iterate over the first map, and see what we have in the second
	for key1, obj1 := range map1 {
		// this key is being checked
		checked[key1] = true

		// what's in the 2nd map ?
		obj2 := map2[key1]

		// what's the next path ? this map's key from the current path
		nextPath := currentPath.To(key1)

		// if the maps here come from 2 arrays, then the key is not a naturel map key, but an object ID
		if fromSlice {
			if idProp := options.GetIDProp(currentPath); idProp != nil {
				nextPath = currentPath.With(idProp)
			}
		}

		// obj1 and obj2 should be compared
		compObj1Obj2, errComp := compareObjects(nextPath, obj1, obj2, options)
		if errComp != nil {
			return nil, errComp
		}

		// adding only if there's at least 1 difference
		if compObj1Obj2.hasDiffs() {
			thisComparison[key1] = compObj1Obj2
		}
	}

	// now let's iterate over the second map, because there might be stuff not found in the first map
	for key2 := range map2 {
		// we're considering keys that have not been checked yet
		if !checked[key2] {
			// what's the next path ? this map's key from the current path
			nextPath := currentPath.To(key2)

			// if the maps here come from 2 arrays, then the key is not a naturel map key, but an object ID
			if fromSlice {
				if idProp := options.GetIDProp(currentPath); idProp != nil {
					nextPath = currentPath.With(idProp)
				}
			}

			// at this point, obj1 does not exist for this key...
			compObj1Obj2, errComp := compareObjects(nextPath, map1[key2], map2[key2], options)
			if errComp != nil {
				return nil, errComp
			}

			// ... so we're sure to have a difference here
			thisComparison[key2] = compObj1Obj2
		}
	}

	// returning
	return thisComparison, nil
}