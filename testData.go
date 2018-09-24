package main

func testData() account {
	p := person{
		"Jax",
		"Blax",
		"none@none.com",
		"215-834-1857",
	}

	a := address{
		"102 Broad Street",
		"",
		"12345",
		"Philadeliphia",
		"PA",
		"USA",
	}

	c := card{
		"visa",
		"5482 4827 5948 1254",
		"02",
		"2022",
		"847",
	}

	return account{p, a, c}
}

func testTask() Task {
	i := taskItem{
		[]string{"shaolin"},
		"hats",
		"",
		"orange",
	}

	return Task{
		"Task1",
		i,
	}
}

func testFullTask() FullTask {
	return FullTask{testTask(), testData()}
}

func realData() account {
	p := person{
		"",
		"",
		"",
		"",
	}

	a := address{
		"",
		"",
		"",
		"",
		"",
		"",
	}

	c := card{
		"",
		"",
		"",
		"",
		"",
	}

	return account{p, a, c}
}

func realTask() Task {
	i := taskItem{
		[]string{"hanes", "boxer"},
		"accessories",
		"Medium",
		"white",
	}

	return Task{
		"Checkout Task",
		i,
	}
}
