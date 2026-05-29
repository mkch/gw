package gw

import "github.com/mkch/gw/app"

// Run initializes an [app.App], call initUI and runs the message loop by calling [app.App.Run].
// The initUI function is called after an [app.App] is created and before the message loop starts,
// so it can be used to create windows and do other UI operations before the message loop starts.
// The cleanup function is called after the message loop ends and before the app is destroyed,
// so it can be used to do cleanup operations.
// If cleanup is nil, it will be ignored.
// Run returns the return value of [app.App.Run].
// After Run returns the app passed to initUI is automatically destroyed, so no gw operation should
// be performed after that.
func Run(initUI func(app *app.App), cleanup func(app *app.App)) int {
	app := app.New()
	defer app.Destroy()
	initUI(app)
	ret := app.Run()
	if cleanup != nil {
		cleanup(app)
	}
	return ret
}
