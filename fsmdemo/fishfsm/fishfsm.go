package fishfsm

func Match(data string) int {
	var idx = ^uint(0)

l0:
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] != 'f' {
		goto l0
	}

l1: // e.g. "f"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'f' {
		goto l1
	}

	if data[idx] != 'i' {
		goto l0
	}

l2: // e.g. "fi"
	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'f' {
		goto l1
	}

	if data[idx] == 'i' {
		goto l2
	}

	if data[idx] != 's' {
		goto l0
	}

	if idx++; idx >= uint(len(data)) {
		return -1
	}

	if data[idx] == 'f' {
		goto l1
	}

	if data[idx] != 'h' {
		goto l0
	}

	{
		return 4
	}

}
