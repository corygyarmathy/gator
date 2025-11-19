{
  description = "Gator - RSS aggregator";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };

        # PostgreSQL helper scripts
        pgScripts = pkgs.writeShellScriptBin "pg-helpers" ''
          export PGDATA="$PWD/.postgres"
          export PGHOST="$PWD/.postgres"
        '';

        pgstart = pkgs.writeShellScriptBin "pgstart" ''
          source ${pgScripts}/bin/pg-helpers
          if pg_ctl status >/dev/null 2>&1; then
            echo "PostgreSQL is already running"
          else
            pg_ctl -o "-k $PGHOST" -l "$PGDATA/logfile" start
            echo "PostgreSQL started"
          fi
        '';

        pgstop = pkgs.writeShellScriptBin "pgstop" ''
          source ${pgScripts}/bin/pg-helpers
          if pg_ctl status >/dev/null 2>&1; then
            pg_ctl stop -m fast
            echo "PostgreSQL stopped"
          else
            echo "PostgreSQL is not running"
          fi
        '';

        pgstatus = pkgs.writeShellScriptBin "pgstatus" ''
          source ${pgScripts}/bin/pg-helpers
          pg_ctl status
        '';

        pglogs = pkgs.writeShellScriptBin "pglogs" ''
          source ${pgScripts}/bin/pg-helpers
          tail -f "$PGDATA/logfile"
        '';

      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Go toolchain
            go
            golangci-lint
            gotools
            delve

            # Database
            postgresql
            goose

            # PostgreSQL helpers
            pgstart
            pgstop
            pgstatus
            pglogs
          ];

          shellHook = ''
            export PGDATA="$PWD/.postgres"
            export PGHOST="$PWD/.postgres"
            export PGDATABASE="gator"
            export DATABASE_URL="postgres:///$PGDATABASE?host=$PGHOST"

            # Initialise PostgreSQL if needed
            if [ ! -d "$PGDATA" ]; then
              echo "Initialising PostgreSQL..."
              initdb --auth=trust --no-locale --encoding=UTF8 >/dev/null
            fi

            # Check if PostgreSQL is running
            if pg_ctl status >/dev/null 2>&1; then
              echo "✓ gator dev environment ready (PostgreSQL running)"
            else
              echo "✓ gator dev environment ready"
              echo "  Run 'pgstart' to start PostgreSQL"
            fi

            echo "  Commands: pgstart, pgstop, pgstatus, pglogs"
          '';
        };
      }
    );
}
