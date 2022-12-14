## Compile
```bash
go build main.go
```



## Local Usage

Modify config.yml according to your profiles and secrets.

```bash
eval $(./main --config config.yml --profile profile)

```

Or use environment variables to load the config and profile

```bash
export SECRET_PROFILE=default
export CONFIG_FILE=config.yml

eval $(./main)

```

Or just use it with the default values

```bash
eval $(./main)

```

## Usage on Docker

```dockerfile

ENTRYPOINT ["entrypoint.sh"]
CMD ["YOUR_EXECUTABLE","YOUR","PARAMS"]

```