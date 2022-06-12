#define _GNU_SOURCE

#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
#include <JavaScriptCore/JavaScript.h>
#include <sys/wait.h>
#include <regex.h>

#define CRASH(message) \
    fprintf(stderr, "CRASH: %s\n", message); \
    fflush(stderr); \
    abort();

extern char* onMessage(char*);

bool started = false;

GtkWidget* window;
GtkWidget* webview;

void message(WebKitUserContentManager* mananger, WebKitJavascriptResult* value, gpointer ctx) {
    char* payload = jsc_value_to_string(webkit_javascript_result_get_js_value(value));

    regex_t regex;
    if (regcomp(&regex, "[^a-zA-Z0-9=+\\/]", 0)) {
        CRASH("compile regex failed!");
    }

    char* response = onMessage(payload);
    if (regexec(&regex, response, 0, NULL, 0) != REG_NOMATCH) {
        CRASH("malformed response!");
    }

    regfree(&regex);

    char* script;
    if (asprintf(&script, "__resolveMessage(`%s`)", response) == -1) {
        CRASH("resolve message failed!");
    }

    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(webview), script, NULL, NULL, NULL);
    free(script);
}
void load_changed(WebKitWebView* webview, WebKitLoadEvent load_event, gpointer load_data) {
#ifdef FIRECRYPT_RELEASE
    if (load_event == WEBKIT_LOAD_COMMITTED) {
        webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(webview), "window.__FIRECRYPT_RELEASE = true", NULL, NULL, NULL);
    }
#endif
}
void StartFirecrypt() {
    if (started) {
        CRASH("already started!");
    }

    started = true;

#ifndef FIRECRYPT_RELEASE
    printf("this is a debug build. issues not present in release builds may occur.\n");
#endif

    gtk_init(0, NULL);

    window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_title(GTK_WINDOW(window), "Firecrypt");
    gtk_window_set_default_size(GTK_WINDOW(window), 300, 300);
    gtk_window_set_keep_above(GTK_WINDOW(window), TRUE);
    gtk_window_set_resizable(GTK_WINDOW(window), FALSE);

    g_signal_connect(window, "destroy", G_CALLBACK(gtk_main_quit), NULL);

    webview = webkit_web_view_new_with_context(webkit_web_context_new_ephemeral());

#ifndef FIRECRYPT_RELEASE
    WebKitSettings* settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webview));
    webkit_settings_set_enable_developer_extras(settings, TRUE);
#endif

    WebKitUserContentManager* manager = webkit_web_view_get_user_content_manager(WEBKIT_WEB_VIEW(webview));

    g_signal_connect(manager, "script-message-received", G_CALLBACK(message), NULL);
    g_signal_connect(webview, "load-changed", G_CALLBACK(load_changed), NULL);

    webkit_user_content_manager_register_script_message_handler(manager, "firecrypt");

    char cwd[PATH_MAX];
    if (getcwd(cwd, sizeof(cwd)) == NULL) {
        CRASH("can't get cwd!");
    }

    char* path;
    if (asprintf(&path, "file://%s/ui/firecrypt.html", cwd) == -1) {
        CRASH("format path failed!");
    }

    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webview), path);

    gtk_container_add(GTK_CONTAINER(window), GTK_WIDGET(webview));

    gtk_widget_show_all(window);

    gtk_main();
}

void sigchld_handler() {
    gtk_widget_show(window);
    signal(SIGCHLD, SIG_DFL);
}
void RunFirefox(char* profileName) {
    if (!started) {
        CRASH("not started!");
    }

    gtk_widget_hide(window);

    pid_t pid = fork();
    if (pid < 0) {
        CRASH("fork failed!");
    } else if (pid == 0) {
        char* firefoxPath = "/usr/bin/firefox";
        char* args[] = { firefoxPath, "-p", profileName, NULL };
        execv(firefoxPath, args);
        return;
    }

    signal(SIGCHLD, sigchld_handler);
}
