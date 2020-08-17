# gotmpl

A lightweight tool for rendering Go templates from the command line.

## Command Line Usage

Pass the template file(s) and they will be rendered to standard out.

```console
$ gotmpl /app/etc/wp-config.php.tmpl
```

Templates may use [sprig](https://masterminds.github.io/sprig/) functions.

```go-template
define('DB_HOST', {{ squote ( env "WORDPRESS_DB_HOST" ) }});
```

Inline templates may be used with the `--template` option.

```console
$ gotmpl \
  --template '{{ now | date "2006-01-02" }}'
```

External data from files or URLs may also be used (JSON, YAML, or raw).

```console
$ gotmpl \
  --json data.json \
  --template '{{ .location }}'
```

Multiple data sources may be used by prefixing the source with a key.

```console
$ gotmpl \
  --json jdat=data.json \
  --yaml ydat=data.yaml \
  --template '{{ .jdat.location }} vs {{ .ydat.location }}'
```

Output may be redirected elsewhere (or to multiple copies).

```console
$ gotmpl template.tmpl \
  --output /mnt/mirror-1/result.txt \
  --output /mnt/mirror-2/result.txt
```

Outputs may specify a [`block` name](https://golang.org/pkg/text/template/) to render different things to files.

```console
$ gotmpl templates/*.tmpl \
  --output summary.txt=summary \
  --output results.csv=csv
```

## Alternatives

 * [`consul-template`](https://github.com/hashicorp/consul-template) – similar for Consul and Vault integration

## License

[MIT License](LICENSE)
