# ok cli

aws cli abstraction

```shell
go build -o ok main.go && mv ok $HOME/go/bin/ok
mkdir -p ~/.ok
cp .ok* ~/.ok
```

## ok cli

```shell
ok help
```

## ok prep

```shell
ok prep help

ok prep docker create --private -f .ok.prep.prototype
ok prep docker destroy --private -f .ok.prep.prototype
ok prep docker create --public -f .ok.prep.prototype
ok prep docker destroy --public -f .ok.prep.prototype

ok prep helm create --private -f .ok.prep.prototype
ok prep helm destroy --private -f .ok.prep.prototype
ok prep helm create --public -f .ok.prep.prototype
ok prep helm destroy --public -f .ok.prep.prototype
```

## ok tidy

```shell
ok tidy help

ok tidy
ok tidy -f .ok.tidy
```

## ok whoami

```shell
ok whoami help

ok whoami --account 000000000000 --region us-west-2 --environment prototype --version v1 --organization stxkxs --name team001 --alias development
ok whoami --account 000000000000 --region us-west-2 --environment production --version v1 --organization stxkxs --name team001 --alias live
```
