# 12fact

Simple app demonstrating the 12Factor app principles (<https://12factor.net/>)

This service allows you to manipulate images and apply the following
transformations:

-   Rotate
-   Blur
-   Crop

Supported formats:

    * `.jpg`, `.jpeg`
    * `.png`
    * `.tif`, `.tiff`
    * `.bmp`
    * `.gif`

## Getting started

### Using go

Simply `go get` the repository:

    go get -u github.com/ganitzsh/12fact

### Building

You can build the project, but you need to have `protoc` installed

## General error codes

The API returns the following error codes:

| Code      | Description                           |
| --------- | ------------------------------------- |
| some_code | What it means and what could cause it |

## Server mode

By default the program will start in server mode and run a gRPC server on port
`8080` by default.

You can also run an HTTP server in parallel by setting
`http_on: true` in the config file or by using the `--http`, it will run on port
`8081` by default

By default the server will operate in dev mode. For now the dv mode is simply
displaying more log and disabling the authentication.

Server mode gRPC only:

    ./12fact
    WARN[0000] This instance is running in dev mode
    INFO[0000] RPC server started on port 8080

Server mode with both gRPC and HTTP&#x3A;

    ./12fact --http
    WARN[0000] This instance is running in dev mode
    INFO[0000] Starting HTTP Server on port 8081
    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in
    production.
    - using env:	export GIN_MODE=release
    - using code:	gin.SetMode(gin.ReleaseMode)
    [GIN-debug] POST   /api/v1/images/rotate
    [GIN-debug] POST   /api/v1/images/blur
    [GIN-debug] POST   /api/v1/images/crop
    INFO[0000] RPC server started on port 8080

#### Flags

A set of flags is available to you in server mode:

    -a, --addr string   The server's address
    -c, --config string Config file to use (default is ./config.yml)
    --http              In server mode, starts the HTTP server
                        (default port 8081)
    --http-port int32   (default 8081)
    -p, --port int32    port on which the server will listen

### HTTP Server

The HTTP server has routes to manipulate images and transit over http. All of
them are prefixed by `/api/v1`.

##### HTTP-Specific error codes

| Code      | Description                           |
| --------- | ------------------------------------- |
| some_code | What it means and what could cause it |

#### Rotate

**Method**: `POST`
**URI**: `/api/v1/images/rotate`

**Request**:

    Content-Type: multipart/form-data

    file: <some multipart file>
    angle: <float or int>
    clockWise: true | false

**Response**:

    Content-Type: <identical_to_file>
    Status: 200

    <bytes containing the transformed image>

#### Blur

**Method**: `POST`
**URI**: `/api/v1/images/blur`

**Request**:

    Content-Type: multipart/form-data

    file: <some multipart file>
    sigma: <the factor of blur>

**Response**:

    Content-Type: <identical_to_file>
    Status: 200

    <bytes containing the transformed image>

#### Crop

**Method**: `POST`
**URI**: `/api/v1/images/blur`

**Request**:

    Content-Type: multipart/form-data

    file: <some multipart file>
    topLeftX: <The position of the topLeft point on the X axis>
    topLeftY: <The position of the topLeft point on the Y axis>
    width: <The desired width of the final image>
    height: <The desired height of the final image>

_NOTE: Going out of bounds will return an error with a 500 Status Code_

**Response**:

    Content-Type: <identical_to_file>
    Status: 200

    <bytes containing the transformed image>

## Client mode

Here is an overview of the different command available to you in client mode:

    blur    Blurs the given image with a factor of [sigma]
    crop    Crops the given [file] starting at topLeft of size [width]x[height]
    rotate  Rotates the given [file] with the given [angle]. Clockwise by
            default

All these commands are taking a file as a parameter and creating a new one with
the modifications. By default the new file will be saved at the root of the
executable with the following pattern: `out.<extension>`. This can be overridden
by specifying the `-o` or `--out` flag when running the command

### Blur

Blurs the given image with a factor of [sigma]

    12fact blur [file] [sigma]

### Rotate

Rotates the given [file] with the given [angle]. Counter-clockwise by default.

    12fact rotate [file] [angle]

To rotate clockwise, you can add the `--cw` to the command:

    12fact rotate --cw image.jpg 90

### Crop

Crops the given [file] starting at topLeft with size [width]x[height]

    12fact crop [file] [topLeftX] [topLeftY] [width] [height]

## Library

The service can also be used as a library, a simple program tanking an image and
applying transformations would look like this:

package main

    ```go
    import (
    	"io"
    	"os"
    	"path/filepath"

    	trans "github.com/ganitzsh/12fact/service"
    )

    func main() {
    	image := "./image.png"
    	ext := filepath.Ext(image)
    	f, err := os.Open(image)
    	if err != nil {
    		panic(err)
    	}

    	// You can now start applying transformations
    	rotation := trans.NewRotate().SetAngle(90).SetClockWise(true)
    	ret, err := trans.SingleTransformImage(f, ext, rotation)
    	if err != nil {
    		panic(err)
    	}

    	blur := trans.NewBlur().SetSigma(4.5)
    	ret, err = trans.SingleTransformImage(ret, ext, blur)
    	if err != nil {
    		panic(err)
    	}

    	// You can also apply multiple transforamtions at once
    	transformations := []trans.Transformation{
    		trans.NewCrop().
    			SetTopLeftX(40).SetTopLeftY(150).SetWidth(400).SetHeight(400),
    		trans.NewRotate().SetAngle(45).SetClockWise(true),
    		trans.NewBlur().SetSigma(3.7),
    	}

    	ret, err := trans.TransformImage(f, ext, transformations)
    	if err != nil {
    		panic(err)
    	}

    	out, err := os.Create("./out.png")
    	if err != nil {
    		panic(err)
    	}
    	defer out.Close()

    	if _, err := io.Copy(out, ret); err != nil {
    		panic(err)
    	}
    }
    ```
