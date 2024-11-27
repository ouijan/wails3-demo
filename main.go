package main

import (
	"embed"
	_ "embed" // Required for the 'embed' package to work
	"log"
	"runtime"
	"time"

	"github.com/ouijan/wails3-demo/src/icons"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
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

	systemTray := app.NewSystemTray()

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	} else {
		systemTray.SetDarkModeIcon(icons.SystrayDark)
		systemTray.SetIcon(icons.SystrayLight)
	}

	// Support for menu
	myMenu := app.NewMenu()
	myMenu.Add("Show / Hide").OnClick(func(_ *application.Context) {
		if mainWindow.IsVisible() {
			mainWindow.Hide()
		} else {
			mainWindow.Show()
			mainWindow.Focus()
		}
	})

	myMenu.AddSeparator()
	myMenu.Add("Quit").OnClick(func(_ *application.Context) {
		app.Quit()
	})
	systemTray.SetMenu(myMenu)
	// systemTray.SetLabel("Wails3 Demo")
	// systemTray.SetIcon(appIcon) // Doesn't work

	// This will center the window to the systray icon with a 5px offset
	// It will automatically be shown when the systray icon is clicked
	// and hidden when the window loses focus
	// systemTray.AttachWindow(mainWindow).WindowOffset(5)

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
