#!/bin/bash
find . -name "*.go" -type f -exec sed -i '' 's|github.com/ieasydevops/go-stock/internal|go-stock/internal|g' {} \;
