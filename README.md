# gator

boot.dev project written in Go.

## Usage

```bash
gator reset                     # Deletes everything from db
gator register <user_name>      # Registers user in db
gator login <user_name>         # Logs into registered user
gator users                     # Prints list of registered users
gator agg <time_between_reps>   # Aggregates posts, looping after time_between_reps (e.g. '10s', '1m')
```

## Development Setup

### Nix

1. Install `direnv` with `nix-direnv`:

```nix
  programs.direnv = {
    enable = true;
    nix-direnv.enable = true;
  };
```

2. Enter the project directory:

```bash
  cd my-project
  direnv allow
```

A `nix develop` environment will automatically be created from the `flake.nix` and `flake.lock` files.

Manual option: run `nix develop` each time you navigate to the project directory in a shell or editor.

### Other

1. Install `go` and `postgres`
2. Run `go install github.com/corygyarmathy/gator`
3. Refer to Usage section for using.
