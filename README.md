rsay
====

```
-> AWS SQS -> say command
```

Currently, Mac OSX is only supported.

Usage
-----

```
$ export AWS_ACCESS_KEY_ID=...
$ export AWS_SECRET_ACCESS_KEY=...
$ rsay -q QUEUE_NAME [-v VOICE]
```

From another host or terminal:

```
$ aws sqs send-message --queue-url URL --message-body '{"Message": "Hello"}'
```
