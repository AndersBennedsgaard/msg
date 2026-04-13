{
  description = "A persistent, filesystem-backed notification manager for Unix-like environments.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {inherit system;};
      version = "0.1.0";
      commit = builtins.substring 0 7 (self.rev  or self.dirtyRev or "unknown");
      # date = toString (self.lastModified or 0);
      date = let
        epoch = self.lastModified or builtins.currentTime;
      in
        builtins.readFile (pkgs.runCommand "date" {} ''
          ${pkgs.coreutils}/bin/date -u -d @${toString epoch} +%Y-%m-%dT%H:%M:%SZ > $out
        '');
    in {
      packages.default = pkgs.buildGoModule {
        pname = "msg";
        version = version;

        src = self;

        vendorHash = "sha256-Om/C7JxRBtvdRDI2NUQocwrff1Qu7Kq0pfhqZgBqMK8=";

        subPackages = ["."];

        ldflags = [
          "-s"
          "-w"
          "-X github.com/AndersBennedsgaard/msg/internal/version.Version=${version}"
          "-X github.com/AndersBennedsgaard/msg/internal/version.Commit=${commit}"
          "-X github.com/AndersBennedsgaard/msg/internal/version.Date=${date}"
        ];
      };

      apps.default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/msg";

        meta = {
          description = "Run the application";
        };
      };

      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
        ];
      };
    });
}
