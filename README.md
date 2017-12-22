# simplFT

[![Go Report Card](https://goreportcard.com/badge/github.com/Metonimie/simplFT)](https://goreportcard.com/report/github.com/Metonimie/simplFT)
[![Build Status](https://travis-ci.org/metonimie/simplFT.svg?branch=master)](https://travis-ci.org/metonimie/simplFT)

This project was made for the purpose of me to learn and understand Go better and also for the Computer Networking class
that I took in Fall 2017 at UPT.

The scope of this project is to implement a simple server that handles multiple clients and allows the clients to
execute commands on it.

## Commands

The server accepts the following commands:

```
get <filename> - Download the requested filename.
ls             - List the files in the current directory.
cd             - Changes the directory.
clear          - Clear the screen.
exit           - Close the connection with the server.c
pic <filename> - Returns the ascii art of an image. :-)
```

#### Sending commands via netcat

To grab a file the following command can be send:

```echo "get file.txt" | nc ip port > lfile.txt```

If someone wishes to send multiple commands, the following syntax
can be used:

```(echo "get file1.txt" & echo "get file2.txt") | nc ip port > concatenated.txt```

#### The upload server

If the upload server is running, the user will be able to put files
on the **absoluteServePath**. After the file is uploaded successfully,
if the timeout is not reached, the user will get back the filename.

To send data to the upload server, the following command can be used:

```nc ip port < gopher.png```

## Configuration

The server can be configured via command line flags with the -config option,
specifying a path to the configuration file.
If no configuration file is provided the server will run with the default settings.

Sample Configuration File:
```
{
    "address": "localhost",
    "port": 8080,
    "maxDirDepth": 30,
    "absoluteServePath": "/Users/denis/Dropbox/Pictures/CuteAvatars",
    "pic": {
        "color": true,
        "x": 197,
        "y": 50
    },
    "upload": {
        "enabled": false,
        "directory": "upload",
        "timeout": 5,
        "address": "localhost",
        "port": 8081
    }
}
```

The **config.json** file contains the following settings:

1. address           - The address on which to serve

2. port              - The port

3. maxDirDepth       - The maximum depth the user can go into directories. A value of 30 means the user can cd into max 30 dirs.

4. absoluteServePath - The path from where to serve the files.

5. pic               - The X and Y max size for the pic command. A value of 0 means original size.

6. upload            - Settings for the upload server.
If one of the settings are changed, the server will reload the configuration.
Except for the absoluteServePath.
