#!/bin/bash

# Fix import syntax errors where imports are on the same line

cd /home/user/Flutter-Database/backend

echo "🔧 Fixing import syntax errors..."

for file in internal/*/validator.go; do
    # Fix pattern: "regexp"	"time" -> "regexp"\n\t"time"
    sed -i 's/"regexp"\t"time"/"regexp"\n\t"time"/g' "$file"
    # Also fix any other double imports on same line
    sed -i 's/"\([^"]*\)"\t"\([^"]*\)"/"\1"\n\t"\2"/g' "$file"
done

echo "✅ Import syntax fixed!"
