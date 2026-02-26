# digSignalGen

## Util for generating mock digital signal

### Build

In root directory of project run command

```bash
    make build
```

### Usage

#### Example

```bash
    ./dsg gen
```

#### Docs

| command | parametr               | default     | reqired | desc                                                                  |
| ------- | ---------------------- | ----------- | ------- | --------------------------------------------------------------------- |
| `gen`   |                        |             |         | generate digital wave with outputing into graph                       |
|         | `--alpha=<string>`     | `0.2\|0.7`  | `N`     | configuring alpha's for calc's in format `alpha1\|alpha2`             |
|         |                        |             |         | example : `--alpha=0.1\|0.5`                                          |
|         | `--length=<uint64>`    | `1000`      | `N`     | configuring length of generating digital wave 1 value = 1 x cordiante |
|         | `--init-cond=<string>` | `0.6\|-0.1` | `N`     | configuring start X's for calc's in format `startX[n-1]\|startX[n-2]` |
|         |                        |             |         | example : `--init-cond=-0.1\|0.6`                                     |
|         | `--disable-output`     | `false`     | `N`     | disable gui graph output                                              |
|         | `--save-file=<string>` | `""`        | `N`     | path for saving generated graph if empty dsg will not save nothing    |
