# Test-bed Large File Service

The Test-bed Large File service allows uploading large files. An uploaded file may be specified to be public or private. Public files are directory listed and retrievable by their original name. Private files are not directory listed, and can only be retrieved by using their obfuscated name.

# Installation and running

* Install Go: https://golang.org/doc/install
* Install/Update Large File Service: `go get github.com/driver-eu/large-file-service` This installs or updates the service to your `%GOPATH%\bin` directory.
* Navigate to your `%GOPATH/src/github.com/driver-eu/large-file-service` directory.
* Run `go build .`
* Execute `large-file-service` or `large-file-service.exe` for windows binary from your command line

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

# Example Usage

To upload a file publically via cURL:

`curl -X POST "http://localhost:9090/upload" -H "accept: application/json" -H "Content-Type: multipart/form-data" -F "uploadFile=@publicfile.pdf;type=application/pdf" -F "private=false"`

To upload a file privately via cURL:

`curl -X POST "http://localhost:9090/upload" -H "accept: application/json" -H "Content-Type: multipart/form-data" -F "uploadFile=@privatefile.pdf;type=application/pdf" -F "private=true"`

Alternatively, you may use functionality provided by the programming language of your choice to POST multipart/form-data.

After uploading your file you will receive a JSON response, with either an obfuscated or public URL:

```json
{"FileURL":"http://localhost:9090/private/B4FE4514-02C1-3FCE-4677-3373C7D5E914.pdf"}
```
This FileURL can then be filled in the [large_data_update](https://github.com/DRIVER-EU/avro-schemas/blob/master/core/large-data/system_large_data_update-value.avsc) message that can be sent to the DRIVER+ Test-bed to notify other Applications that are conntected to the Test-bed that the file can be downloaded.

