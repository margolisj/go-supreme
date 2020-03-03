package main

func testAccount() Account {
	p := Person{
		"Jax",
		"Blax",
		"none@none.com",
		"215-834-1857",
	}

	a := Address{
		"102 Broad Street",
		"",
		"12345",
		"Philadeliphia",
		"PA",
		"USA",
	}

	c := Card{
		"1285 4827 5948 2017",
		"02",
		"2019",
		"847",
		"",
	}

	return Account{p, a, c}
}

// testTask is a tester task. It is missing API and refresh rate settings.
func testTask() Task {
	item := taskItem{
		[]string{"shaolin"},
		"shirts",
		"",
		"orange",
	}

	return Task{
		TaskName: "Task1",
		Item:     item,
		Account:  testAccount(),
		API:      "mobile",
	}
}
