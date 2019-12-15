#!/usr/bin/env bash

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [[ -h "$SOURCE" ]] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd "$DIR"

# Change project name
grep -rli --exclude-dir=node_modules 'golang-starter' * | xargs -I@ sed -i "" 's/golang-starter/'"${PWD##*/}"'/g' @

# Install and run Commitizen locally
npm install --save-dev commitizen
# initialize the conventional changelog adapter
npx commitizen init cz-conventional-changelog --save-dev --save-exact

# Install commitlint cli and conventional config
npm install --save-dev @commitlint/{config-conventional,cli}
echo "module.exports = {extends: ['@commitlint/config-conventional']};" > commitlint.config.js

# Install husky as devDependency, a handy git hook helper available on npm
# This allows us to add git hooks directly into our package.json via the husky.hooks field
npm install --save-dev husky
npm install --save-dev semantic-release

