package esu

import (
	"fmt"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func TestUtils_getConnectionURL(t *testing.T) {
	h := "foobar"
	p := "1234"

	expected := fmt.Sprintf("https://%s:%s", h, p)
	var actual string

	app := cli.NewApp()
	app.HideHelp = true
	setupGlobalFlags(app)

	app.Action = func(ctx *cli.Context) {
		actual = getConnectionURL(ctx).String()
	}

	err := app.Run([]string{"", "--host", h, "--port", p, "--ssl"})
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
