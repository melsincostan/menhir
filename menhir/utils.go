package menhir

func sortFunc[T ModuleBase](invert bool) func(a, b T) int {
	return func(a, b T) int {
		if a.Priority() == nil && b.Priority() == nil {
			return 0
		}

		if b.Priority() == nil { // since checked before, a is not nil
			if invert {
				return -1
			} else {
				return 1
			}
		}

		if a.Priority() == nil { // since checked before, b is not nil
			if invert {
				return 1
			} else {
				return -1
			}
		}

		if invert {
			return *b.Priority() - *a.Priority()
		} else {
			return *a.Priority() - *b.Priority()
		}
	}
}
