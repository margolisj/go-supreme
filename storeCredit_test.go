package main

import (
	"regexp"
	"testing"

	"github.com/levigross/grequests"
	"github.com/stretchr/testify/assert"
)

func TestStoreCreditAPI(t *testing.T) {
	session := grequests.NewSession(nil)
	task := testTask()
	task.Account.Person.Email = "munaweryusuf@mchsupply.com"

	id, err := GetStoreCredits(session, &task)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "18498", id)
}

func TestStoreCreditMatch(t *testing.T) {
	respString := string(`<div id="store_credits" style="display:none">
	<p> You have $60.68 available in store credits.<br /> Do you want to use it?</p>

	<input type="submit" name="commit" value="Use Store Credit" class="button checkout" id="store_credit" store_credit_id="18498" /><input type="submit" name="commit" value="Do Not Use" class="button checkout" id="no_store_credit" />
</div>`)

	re := regexp.MustCompile(`store_credit_id="(?P<id>\d*)"`)
	findResults := re.FindAllStringSubmatch(respString, -1)
	t.Log(findResults[0][1])
}
