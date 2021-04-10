# go-mefs-api


> A go interface to mefs's HTTP API

## install

```sh
go get -u github.com/memoio/mefs-go-http-client
```

## prepare

To interact with the API, you need to have a local daemon running. It needs to be open on the right port. `5001` is the default, and is used in the examples below, but it can be set to whatever you need.

```shell
# Show the mefs config API port
> mefs config Addresses.API
/ip4/127.0.0.1/tcp/5001
# set api port and binding to all ip 
> mefs config Addresses.API /ip4/0.0.0.0/tcp/5001
# Restart the daemon after changing the config
> mefs shutdown
# Run the daemon
> mefs daemon
```

### CORS

In a web browser mefs HTTP client (either browserified or CDN-based) might encounter an error saying that the origin is not allowed. This would be a CORS ("Cross Origin Resource Sharing") failure: mefs servers are designed to reject requests from unknown domains by default. You can whitelist the domain that you are calling from by changing your mefs config like this:

```shell
> mefs config --json API.HTTPHeaders.Access-Control-Allow-Origin  '["http://example.com"]'
> mefs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "POST", "GET"]'
# Restart the daemon after changing the config
> mefs shutdown
# Run the daemon
> mefs daemon
```

## example

see example directory

```go
package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/memoio/mefs-go-http-client"
)

func main() {
	
	// your local address
	p := path.Join(os.Getenv("HOME"), "test")
	file, err := os.Open(p)
	ob, err := sh.PutObject(file, path.Base(file.Name()), "bucket01")
	fmt.Println(ob, err)
	bks, err := sh.ListBuckets()
	fmt.Println(bks, err)
	obs, err := sh.ListObjects(bks.Buckets[0].BucketName)
	fmt.Println(obs, err)
}
```

## Usage

See [mefs docs](https://github.com/memoio/docs)

### LFS

The API enables users to use the LFS abstraction of MEFS. 

#### StartUser

> start user's lfs service 

##### `CreateBucket(addr,ops...)`


`addr` AddressID. Initialize user service with the given address. Type is `string`.

`options` is an optional object argument that might include the following keys:

- `pwd` PassWord. Password of the actual user that you want to execute. Type is `string`.
- `sk` SecreteKey. Private key of the actual user that you want to execute. Type is `string`. If `sk` is not `nil`, mefs will store `sk` in the keystore with `pwd`; otherwise, mefs tries to load the private key from the keystore use the `addr` and `pwd`.  

```go
	sh = shell.NewShell("localhost:5001")
	
	// set multiple parameters
	op1 := shell.SetOp(option1, opsValue1)
	op2 := shell.SetOp(option2, opsValue2)
	...
	// start user
	sh.StartUser(addr, op1, op2,...)
```

#### CreateBucket

> create a bucket in lfs.

##### `CreateBucket(bucketName,ops...)`

`bucketName` is a string of the bucket name we want.

`options` is an optional object argument that might include the following keys:

- `addr` AddressID. The actual user's addressid that you want to execute. Type is `string`.
- `pl` Policy. The Storage policy you want to use. Type is `bool`. `true`for erasure coding and `false` for copyset. `default` is `true`.
- `dc` DataCount. `default` is 3.
- `pc` ParityCount. `default` is 2.


#### `DeleteBucket`

> delete a bucket in lfs.

##### `DeleteBucket(bucketName,ops...)`

`bucketname` is a string of the bucket name .

`options` is an optional object argument that might include the following keys:

- `addr` AddressID. The actual user's addressid that you want to execute. Type is `string`.


#### `PutObject`

> put an object to a bucket

##### `PutObject(data, objectName, bucketName, ops...)`


`data` is the data we want to store.
`bucketName` is a string of the bucket name.
`objectName` is a string of the object name.

`options` is an optional object argument that might include the following keys:

- `addr` AddressID. The actual user's addressid that you want to execute. Type is `string`.


#### `GetObject`

> get an object in a bucket.

##### `GetObject(objectName, bucketName, ops...)`

`bucketName` is a string of the bucket name.
`objectName` is a string of the object name.

`options` is an optional object argument that might include the following keys:

- `addr` AddressID. The actual user's addressid that you want to execute. Type is `string`.


#### `DeleteObject`

> delete an object in a bucket.

##### `DeleteObject(bucketName,objectName)`

`bucketName` is a string of the bucket name.
`objectName` is a string of the object name.



