# How to use
At the server side, you need to import package to run initializer
```
import _ "github.com/un000/grpc-snappy"
```

At the client use DialOptions:
```
    var dialOpts = []grpc.DialOption{
        grpc.WithInsecure(),
    }

    if *UseCompressor {
        // here it is
        dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor(snappy.Name)))
    }

    conn, err := grpc.Dial(
        "you-grpc-address:port",
        dialOpts...,
    )
```

For more information visit the official docs:
    https://github.com/grpc/grpc-go/blob/master/Documentation/compression.md