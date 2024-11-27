package app

import (
	"embed"
	"runtime"
	"time"

	"github.com/ouijan/wails3-demo/backend/icons"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func Bootstrap(assets embed.FS) *application.App {
	app := newApp(assets)
	mainWindow := newMainWindow(app)

	systemTrayMenu := newSystemTrayMenu(app, mainWindow)
	systemTray := newSystemTray(app)
	systemTray.SetMenu(systemTrayMenu)

	newTimeEvent(app)

	return app
}

func newApp(assets embed.FS) *application.App {
	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	return application.New(application.Options{
		Name:        "wails3-demo",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
			// ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
}

func newMainWindow(app *application.App) *application.WebviewWindow {
	// Create a new mainWindow with the necessary options.
	// 'Title' is the title of the mainWindow.
	// 'Mac' options tailor the mainWindow when running on macOS.
	// 'BackgroundColour' is the background colour of the mainWindow.
	// 'URL' is the URL that will be loaded into the webview.
	mainWindow := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:         400,
		Height:        600,
		Frameless:     false,
		AlwaysOnTop:   true,
		DisableResize: false,
		Hidden:        true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		MinimiseButtonState: application.ButtonHidden,
		MaximiseButtonState: application.ButtonHidden,
		CloseButtonState:    application.ButtonEnabled,

		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		// BackgroundColour: application.NewRGB(27, 38, 54),
		// URL:              "/",
	})

	// Disable window closing by canceling the event
	mainWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		e.Cancel()
		mainWindow.Hide()
	})

	return mainWindow
}

func newSystemTray(app *application.App) *application.SystemTray {
	systemTray := app.NewSystemTray()

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	} else {
		systemTray.SetDarkModeIcon(icons.SystrayDark)
		systemTray.SetIcon(icons.SystrayLight)
	}

	// This will center the window to the systray icon with a 5px offset
	// It will automatically be shown when the systray icon is clicked
	// and hidden when the window loses focus
	// systemTray.AttachWindow(mainWindow).WindowOffset(5)

	return systemTray
}

func newSystemTrayMenu(app *application.App, mainWindow *application.WebviewWindow) *application.Menu {
	// Support for menu
	myMenu := app.NewMenu()
	myMenu.Add("Show / Hide").OnClick(func(_ *application.Context) {
		if mainWindow.IsVisible() {
			mainWindow.Hide()
		} else {
			mainWindow.Show().Focus()
		}
	})

	myMenu.AddSeparator()
	myMenu.Add("Quit").OnClick(func(_ *application.Context) {
		app.Quit()
	})

	return myMenu
}

func newTimeEvent(app *application.App) {
	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()
}
