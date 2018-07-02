maybe-slow
====

### Usage

```
$ maybe-slow -h

Usage: maybe-slow [OPTIONS] COMMAND [ARG...]

Run command, maybe it is slow

Options:
  -a	notify and play alert sound
  -h	show help
  -k	kill slow command
  -s int
    	detect slow command, if not finished after N seconds (default 60)
  -v	show version

```

### Example

nop

```
$ maybe-slow -s 10 sleep 5
```

send desktop notification

```
$ maybe-slow -s 3 sleep 5
```

send desktop notification and play alert sound

```
$ maybe-slow -s 3 -a sleep 5
```

kill and send desktop notification

```
$ maybe-slow -s 3 -k sleep 5
```

command can use stdin, and trap any signals

```
$ echo 'sleep 5' | maybe-slow -s 3 bash
```


