package esu

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func TestApp_InitApp(t *testing.T) {
	app := InitApp()
	assert.Equal(t, "esu", app.Name)
	assert.True(t, app.HideHelp)
}

func TestApp_setupGlobalFlags(t *testing.T) {
	app := cli.NewApp()
	app.HideHelp = true
	setupGlobalFlags(app)

	os.Unsetenv("ESU_HOST")
	os.Unsetenv("ESU_PORT")
	os.Unsetenv("ESU_SSL")

	var host, port string
	var ssl bool

	app.Action = func(ctx *cli.Context) {
		host = ctx.GlobalString("host")
		port = ctx.GlobalString("port")
		ssl = ctx.GlobalBool("ssl")
	}

	err := app.Run([]string{""})
	assert.NoError(t, err)
	assert.Equal(t, DefaultHost, host)
	assert.Equal(t, DefaultPort, port)
	assert.False(t, ssl)

	app.Run([]string{"", "-h", "foo", "-p", "bar", "-s"})
	assert.Equal(t, "foo", host)
	assert.Equal(t, "bar", port)
	assert.True(t, ssl)

	os.Setenv("ESU_HOST", "fizz")
	os.Setenv("ESU_PORT", "buzz")
	app.Run([]string{""})
	assert.NoError(t, err)
	assert.Equal(t, "fizz", host)
	assert.Equal(t, "buzz", port)
	assert.False(t, ssl)
}

func TestApp_setupGlobalCommands(t *testing.T) {
	r, w := io.Pipe()

	app := cli.NewApp()
	app.Writer = w
	setupGlobalCommands(app)

	go func() {
		app.Run([]string{"", "help"})
		w.Close()
	}()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	assert.Contains(t, string(b), "help")
}
