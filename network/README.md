# Basic Network Config

Note that this basic configuration uses pre-generated certificates and
key material, and also has predefined transactions to initialize a channel named "mychannel".

To regenerate this material, simply run

``` bash
generate.sh
```

To start the network, run

``` bash
start.sh
```

To stop it, run

``` bash
stop.sh
```

To completely remove all incriminating evidence of the network
on your system, run

``` bash
teardown.sh
```