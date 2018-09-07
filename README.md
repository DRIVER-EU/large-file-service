# Test-bed Large File Service

The Test-bed Large File service allows uploading large files. An uploaded file may be specified to be public or private. Public files are directory listed and retrievable by their original name. Private files are not directory listed, and can only be retrieved by using their obfuscated name.

# Installation and running

* Install Go: https://golang.org/doc/install
* Add your GOPATH bin folder (default: `%USERPROFILE%\go\bin`) to your PATH.
* Install Large File Service: `go get github.com/driver-eu/large-file-service` This installs the service to your `%GOPATH%\bin` directory.
* Run: `large-file-service` from your command line

Alternatively download a pre-built binary and run:

* Windows: https://github.com/DRIVER-EU/large-file-service/releases/download/0.0.1/large-file-service.exe
* Linux: Not available yet

# Usage (API Spec)

Once the Service is Running you can find a Swagger UI containing API definitions and allowing you to try out the API at http://<hostname>/api. This is http://localhost:9090/api by default.

# Configuration

The DRIVER+ Large File Service may be configured via either Environment variables or by providing command line parameters. If no configuration is provided, the defaults will be applied.

Order of application of configuration is: command line parameters > environment variables > defaults.

## Environment Variable Configuration

| Variable           | Description                                                         | Default value  |
|--------------------|---------------------------------------------------------------------|----------------|
| HOST               | hostname:port that the service will listen on                       | localhost:9090 |
| PUBLIC_UPLOAD_DIR  | path of the directory where public files are uploaded               | ./public       |
| PRIVATE_UPLOAD_DIR | path of the directory where private files are uploaded              | ./private      |
| WRITE_TIMEOUT_SECS | timeout limit in seconds for a GET request / download               | 120            |
| READ_TIMEOUT_SECS  | timeout limit in seconds for a POST request / upload                | 120            |

## Command Line Parameter Configuration

| Command Line Parameter             | Description                                                         | Default value  |
|------------------------------------|---------------------------------------------------------------------|----------------|
| `-hostname=<hostname>`             | hostname:port that the service will listen on                       | localhost:9090 |
| `-publicDir=<path>`                | relative                                                            | ./public       |
| `-privateDir=<path>`               | relative location of the directory where private files are uploaded | ./private      |
| `-writeTimout=<secs>`              | timeout limit in seconds for a GET request / download               | 120            |
| `-readTimeout=<secs>`              | timeout limit in seconds for a POST request / upload                | 120            |
| `-h or -help`                      | lists available command line parameters and their default values    | 120            |
