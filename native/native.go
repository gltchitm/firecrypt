package native

import (
	//#cgo release CFLAGS: -DFIRECRYPT_RELEASE
	//#ifdef __APPLE__
	//#cgo darwin CFLAGS: -x objective-c
	//#cgo darwin LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
	//#include "./darwin/webview.m"
	//#else
	//#error "firecrypt only currently supports macOS!"
	//#endif
	"C"
	"runtime"
)

func StartFirecrypt(onMessage func(string, []string) interface{}) {
	runtime.LockOSThread()
	onMessageCallback = onMessage
	C.StartFirecrypt()
}
func RunFirefox(profileName string) {
	C.RunFirefox(C.CString(profileName))
}
