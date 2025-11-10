#!/bin/bash

# Script to update all handler files to use safe context extraction
# This replaces MustGetOrganizationID and MustGetUserID with safe helpers

set -e

echo "Updating handler files to use safe context extraction..."

# Find all handler files
HANDLER_FILES=$(find backend/internal/http/rest -name "*_handlers.go" | sort)

for file in $HANDLER_FILES; do
    echo "Processing: $file"

    # Create backup
    cp "$file" "$file.bak"

    # Replace MustGetOrganizationID(ctx) patterns - single line
    sed -i 's/orgID := appctx\.MustGetOrganizationID(ctx)/\/\/ Safe context extraction\n\t\torgID, err := getOrganizationID(r)\n\t\tif err != nil {\n\t\t\trespondError(w, logger, err)\n\t\t\treturn\n\t\t}/g' "$file"

    # Replace MustGetUserID(ctx) patterns - single line
    sed -i 's/userID := appctx\.MustGetUserID(ctx)/userID, err := getUserID(r)\n\t\tif err != nil {\n\t\t\trespondError(w, logger, err)\n\t\t\treturn\n\t\t}/g' "$file"

    echo "  ✓ Updated $file"
done

echo ""
echo "✅ All handler files updated!"
echo ""
echo "Backup files created with .bak extension"
echo "To remove backups: find backend/internal/http/rest -name '*.bak' -delete"
