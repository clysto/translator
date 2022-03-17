package main

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Foundation -framework Cocoa

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

void init(void* window) {
	[(id)window setLevel:NSFloatingWindowLevel];
	[(id)window setHasShadow:YES];
}
*/
import "C"
import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/url"
	"os"

	"github.com/clysto/translator/clipboard"
	"github.com/clysto/translator/webview"
)

var view webview.WebView

//go:embed index.html
var html string

var appid = os.Getenv("appid")
var appkey = os.Getenv("appkey")

func main() {
	view = webview.New(true)
	defer view.Destroy()
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	changed := clipboard.Watch(context.Background(), clipboard.FmtText)
	view.SetTitle("Translator")
	view.SetSize(350, 500, webview.HintMin)
	view.SetSize(400, 800, webview.HintNone)
	view.Navigate("data:text/html," + url.PathEscape(html))
	C.init(view.Window())
	go func() {
		for {
			content := string(<-changed)
			r, err := translate(content, appid, appkey, "zh")
			if err == nil {
				view.Eval(fmt.Sprintf(`document.querySelector("#src").innerText="%s";`,
					template.JSEscapeString(content)))
				view.Eval(fmt.Sprintf(`document.querySelector("#dst").innerText="%s";`,
					template.JSEscapeString(r.Results[0].Dst)))
			} else {
				view.Eval(fmt.Sprintf(`alert("%s")`, template.JSEscapeString(err.Error())))
			}
		}
	}()
	view.Run()
}
