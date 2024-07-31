{
  inputs = {
    templ.url = "github:a-h/templ/v0.2.663";
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
  };

  outputs = inputs@{ templ, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      templPkg = templ.packages.${system}.templ;
      pkgs = inputs.nixpkgs.legacyPackages.${system};
    in
    {
      devShell.${system} = pkgs.mkShell {
        name = "redmage-shell";
        buildInputs = with pkgs; [
          templPkg
          go
          modd
          nodePackages_latest.nodejs
          goose
          air
          upx
          buf
          buf-language-server
          protoc-gen-go
          protoc-gen-go-grpc
          protoc-gen-connect-go
        ];
      };
    };
}
