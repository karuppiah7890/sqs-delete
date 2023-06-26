# sqs-delete

This is a tool to delete all the messages available in a SQS queue - usually a dead letter queue, given the list of receipt handles to delete the messages

In short, input is - 1 queue and multiple receipt handles of multiple messages, output - all the messages corresponding to the receipt handles getting deleted from the SQS queue

## Running

Just build the tool using `make` and then run it

```bash
make
```

```bash
./sqs-delete
```

`sqs-delete` is just a simple tool and does not run as a service, it just runs once and then exits.

### Setup

#### AWS credentials

Create an IAM user which has access to `sqs:DeleteMessage` action on the required SQS queues
