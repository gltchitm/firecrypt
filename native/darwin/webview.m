#pragma once

#import <Foundation/Foundation.h>
#import <Webkit/Webkit.h>
#import <Cocoa/Cocoa.h>

#define CRASH(message) \
    fprintf(stderr, "CRASH: %s\n", message); \
    fflush(stderr); \
    abort();

extern char* onMessage(char*);

@interface AppDelegate : NSObject<NSApplicationDelegate, WKScriptMessageHandler>
    @property (retain) NSWindow* window;
    @property (retain) WKWebView* webView;
@end

bool firefoxRunning = false;
NSImage* icon;

@implementation AppDelegate
    -(instancetype)init {
        [super init];

        NSUInteger windowStyleMask = NSWindowStyleMaskTitled |
            NSWindowStyleMaskClosable |
            NSWindowStyleMaskMiniaturizable;

        self.window = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, 330, 300)
                                                  styleMask:windowStyleMask
                                                    backing:NSBackingStoreBuffered
                                                      defer:NO];
        [self.window setLevel:NSFloatingWindowLevel];
        [self.window setCollectionBehavior:NSWindowCollectionBehaviorManaged];
        [self.window setBackgroundColor:[NSColor colorWithDeviceRed:0.106
                                                              green:0.149
                                                               blue:0.173
                                                              alpha:1]];
        [self.window setTitle:@"Firecrypt"];
        [self.window center];

        self.webView = [[WKWebView alloc] init];

#ifdef FIRECRYPT_RELEASE
        NSString* scriptSource = [NSString stringWithFormat:@"window.__FIRECRYPT_RELEASE = true"];
        [self.webView evaluateJavaScript:scriptSource
                       completionHandler:^(NSString* result, NSError* error) {
                           if (error != nil){
                               CRASH("failed to set release flag!");
                           }
                       }];
#endif

#ifdef FIRECRYPT_RELEASE
        NSBundle* bundle = [NSBundle mainBundle];
        [self.webView loadFileURL:[bundle URLForResource:@"firecrypt"
                                           withExtension:@"html"]
          allowingReadAccessToURL:[bundle resourceURL]];
#else
        NSURL* pwd = [NSURL fileURLWithPath:[[NSFileManager defaultManager] currentDirectoryPath]];

        [self.webView loadFileURL:[NSURL fileURLWithPath:@"ui/firecrypt.html"
                                           relativeToURL:pwd]
          allowingReadAccessToURL:[NSURL fileURLWithPath:@"ui/"
                                           relativeToURL:pwd]];
        [[[self.webView configuration] preferences] setValue:@YES
                                                      forKey:@"developerExtrasEnabled"];
#endif

        [self.webView setValue:@NO forKey:@"drawsBackground"];
        [[[self.webView configuration] userContentController] addScriptMessageHandler:self
                                                                                 name:@"firecrypt"];

        [self.window setContentView:self.webView];


        return self;
    }

    -(void)userContentController:(WKUserContentController*)userContentController
         didReceiveScriptMessage:(WKScriptMessage*)message {
        id payload = [message body];
        if (![payload isKindOfClass:[NSString class]]) {
            CRASH("bad payload!");
        }

        char* response = onMessage((char*) [payload UTF8String]);

        NSString* responseString = [NSString stringWithUTF8String:response];
        if ([responseString rangeOfString:@"[^a-zA-Z0-9=+\\/]"
                                  options:NSRegularExpressionSearch].location != NSNotFound) {
            CRASH("response malformed");
        }

        NSString* scriptSource = [NSString stringWithFormat:@"__resolveMessage(`%s`)", response];
        [self.webView evaluateJavaScript:scriptSource
                       completionHandler:^(NSString* result, NSError* error) {
                           if (error != nil) {
                               CRASH("failed to resolve message!");
                           }
                       }];
    }

    -(BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication*)sender {
        return YES;
    }
    -(void)applicationDidFinishLaunching:(NSNotification*)notification {
        [self.window orderFront:nil];
        [[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
    }
    -(void)observeValueForKeyPath:(NSString*)keyPath
                         ofObject:(id)object
                           change:(NSDictionary<NSKeyValueChangeKey, id>*)change
                          context:(void*)firefox {
        if (![change valueForKey:@"terminated"]) {
            [object removeObserver:self forKeyPath:keyPath];
            [(NSRunningApplication*) firefox release];

            NSApplication* application = [NSApplication sharedApplication];
            [application setActivationPolicy:NSApplicationActivationPolicyRegular];
            [application activateIgnoringOtherApps:YES];

            #ifndef FIRECRYPT_RELEASE
            [application setApplicationIconImage:icon];
            #endif

            firefoxRunning = false;
        }
    }
@end

bool started = false;

void StartFirecrypt() {
    if (started) {
        CRASH("already started!");
    }

    #ifndef FIRECRYPT_RELEASE
    NSLog(@"this is a debug build. issues not present in release builds may occur.\n");
    #endif

    started = true;

    NSApplication* application = [NSApplication sharedApplication];

    AppDelegate* applicationDelegate = [[AppDelegate alloc] init];
    [application setActivationPolicy:NSApplicationActivationPolicyRegular];
    [application setDelegate:applicationDelegate];

#ifndef FIRECRYPT_RELEASE
    icon = [[NSImage alloc] initWithContentsOfFile:@"icon/darwin/icon.png"];
    [application setApplicationIconImage:icon];
#endif

    [application run];
}

void RunFirefox(char* profileName) {
    if (!started) {
        CRASH("not started!");
    } else if (firefoxRunning) {
        CRASH("firefox is already running!");
    }

    firefoxRunning = true;

    NSApplication* application = [NSApplication sharedApplication];
    [application setActivationPolicy:NSApplicationActivationPolicyProhibited];

    NSWorkspaceOpenConfiguration* configuration = [NSWorkspaceOpenConfiguration configuration];
    [configuration setActivates:YES];
    [configuration setCreatesNewApplicationInstance:YES];
    [configuration setArguments:@[@"-p", [NSString stringWithUTF8String:profileName]]];

    [[NSWorkspace sharedWorkspace] openApplicationAtURL:[NSURL fileURLWithPath:@"/Applications/Firefox.app"]
                                          configuration:configuration
                                      completionHandler:^(NSRunningApplication* firefox, NSError* error) {
                                          if (error) {
                                              CRASH("launch firefox failed!");
                                          }

                                          [firefox retain];

                                          [firefox addObserver:[application delegate]
                                                    forKeyPath:@"terminated"
                                                       options:0
                                                       context:firefox];
                                      }];
}
