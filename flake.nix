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
    in {
      packages.default = pkgs.buildGoModule {
        pname = "msg";
        version = "0.1.0";

        src = self;

        vendorHash = "sha256-Om/C7JxRBtvdRDI2NUQocwrff1Qu7Kq0pfhqZgBqMK8=";

        subPackages = ["."];

        ldflags = [
          "-s"
          "-w"
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
