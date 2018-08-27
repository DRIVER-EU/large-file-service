# Test-bed Large File Service

The Test-bed Large File service allows uploading large files. An uploaded file may be specified to be public or private. Public files are directory listed and retrievable by their original name. Private files are not directory listed, and can only be retrieved by using their obfuscated name.

# Installation and running

* Install Go: https://golang.org/doc/install
* Install Large File Service: `go get github.com/driver-eu/large-file-service`.
* Run: `large-file-service` from your command line.

Alternatively download a pre-built binary and run that:

* Windows: https://github.com/DRIVER-EU/large-file-service/releases/download/0.0.1/large-file-service.exe
* Linux: Not available yet

# Usage (API Spec)

## GET /upload
A very simple HTML form can be found at `http://<host>/upload` where a file can be selected and the upload can be indicated as private or public.

## POST /upload
Allows uploading a large file as either private or public.

Input Parameters (type multipart/form-data):
| Parameter Name     | Type             | Value               |
|--------------------|------------------|---------------------|
| uploadFile         | file             | file to be uploaded |
| private            | checkbox         | private             |

Return Value (type JSON):
```json
{ "FileURL": <url where file can be downloaded> }
```

## GET /private/{fileName}

Return Value (type application/octet-stream binary):
Requested file if it exists, 404 otherwise.

## GET /public/
Standard file-server directory listing all public files

## GET /public/{fileName}

Return Value (type application/octet-stream binary):
Requested file if it exists, 404 otherwise.

# Configuration

The following Environment variables can be set for configuring the service:

| Variable           | Description                                                         | Default value  |
|--------------------|---------------------------------------------------------------------|----------------|
| HOST               | hostname:port that the service will listen on                       | localhost:9090 |
| PUBLIC_UPLOAD_DIR  | relative                                                            | public         |
| PRIVATE_UPLOAD_DIR | relative location of the directory where private files are uploaded | private        |
| WRITE_TIMEOUT_SECS | timeout limit in seconds for a file upload                          | 120            |
| READ_TIMEOUT_SECS  | timeout limit in seconds for a file download                        | 120            |


