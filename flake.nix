{
  inputs = {
    templ.url = "github:a-h/templ/v0.2.648";
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
        ];
      };
    };
}
