#!/usr/bin/osascript
ObjC.import('stdlib')

var username = $.getenv('username')
var auth_type = $.getenv('auth_type')

var app = Application.currentApplication()
app.includeStandardAdditions = true

var response = app.displayDialog(`Enter ${auth_type} for ${username}:`, {
    defaultAnswer: "",
    withIcon: Path("./icon.png"),
    buttons: ["Cancel", "OK"],
    defaultButton: "OK",
    cancelButton: "Cancel",
    givingUpAfter: 120,
    hiddenAnswer: true
})

password = response.textReturned

if (password.length > 0) {
    $.setenv('ALFRED_AUTHCONFIG_PASSWORD', password, 1);
    var cmd = `./alfred-bitbucket-search -auth "${auth_type}"`
    var result = app.doShellScript(cmd);
    $.unsetenv('ALFRED_AUTHCONFIG_PASSWORD');
}
