#!/usr/bin/env bash

cat << EOB
{"items": [
    {
        "title": "Error",
        "subtitle": "$error_msg",
        "arg": "$error_msg",
        "icon": {
            "path": "/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/AlertStopIcon.icns"
        }
    },
    {
        "title": "Clear credentials to reauthenticate",
        "arg": "clearauth"
    }
]}
EOB
