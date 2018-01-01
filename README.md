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

```./simplFT --help``` will list all the command line flags with the associated help text.

The **config.json** file contains the following settings:

1. address           - The address on which to serve

2. port              - The port

3. maxDirDepth       - The maximum depth the user can go into directories. A value of 30 means the user can cd into max 30 dirs.

4. absoluteServePath - The path from where to serve the files.

5. pic               - The X and Y max size for the pic command. A value of 0 means original size.

6. upload            - Settings for the upload server.
If one of the settings are changed, the server will reload the configuration.
Except for the absoluteServePath.

## Docker

To build the image run: ```docker build -t simplft .``` and to run the server:

```
docker run -d \
  -it \
  --name devtest \
  --mount type=bind,source="/Users/denis/Downloads",target=/externalVolume \
  -p 8080:8080 -p 8081:8081 \
  simplft
```

* ```-p PORT1:PORT2``` - PORT1 is the host machine's port mapped to the container's PORT2
* ```source="/Users/denis/Downloads"``` - This path should be changed to the path from where you want to serve files.

To stop the server you will first need to identify the running container and the stop it via 
```docker container stop CONTAINER ID```

##### Stopping the server

```
➜  server git:(master) ✗ docker container ls
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                                             NAMES
90b6f00b1331        simplft             "./simplFTP -confi..."   2 minutes ago       Up 2 minutes        0.0.0.0:8081->8081/tcp, 0.0.0.0:32768->8080/tcp   devtest
➜  server git:(master) ✗ docker container stop 90b6f00b1331
90b6f00b1331
```