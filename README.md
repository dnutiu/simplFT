# simplFT
This project was made for the purpose of me to learn and understand Go better and also for the Computer Networking class
that I took in Fall 2017 at UPT.

The scope of this project is to implement a simple server that handles multiple clients and allows the clients to
execute commands on it.

## Configuration

The server can be configured via command line flags with the -ConfigPath option,
specifying a path to the configuration file.
If no configuration file is provided the server will run with the default settings.

Sample Configuration File:
```
{
    "Address": "localhost",
    "Port": "8000",
    "MaxDirDepth": 30,
    "AbsoluteServePath": "./"
}
```

The **config.json** file contains the following settings:

Address           - The address on which to serve
Port              - The port
MaxDirDepth       - The maximum depth the user can go into directories. A value of 30 means the user can cd into max 30 dirs.
AbsoluteServePath - The path from where to serve the files.
