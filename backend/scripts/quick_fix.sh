#!/bin/bash
cd /home/user/Flutter-Database/backend

echo "🔧 Quick Fix: Adding missing json imports..."

# Fix missing json imports in dto files
for file in internal/audit_log/dto.go \
            internal/api_request_log/dto.go \
            internal/bank_account/dto.go \
            internal/bank_reconciliation/dto.go; do
    if [ -f "$file" ] && ! grep -q '"encoding/json"' "$file"; then
        # Check if there's already an import block
        if grep -q '^import ($' "$file"; then
            sed -i '/^import ($/a\	"encoding/json"' "$file"
        else
            # Add import after package line
            sed -i '/^package /a\\nimport (\n\t"encoding/json"\n)' "$file"
        fi
        echo "✅ Fixed $file"
    fi
done

# Fix Background Job duplicate Status fields
echo ""
echo "🔧 Fixing background_job Status duplicates..."
file="internal/background_job/dto.go"
if [ -f "$file" ]; then
    # Remove duplicate Status lines (keep first, comment others)
    awk '
    /Status.*string.*`json:"status"`/ {
        if (!seen[$0]++) {
            print
        } else {
            print "\t// Duplicate removed: " $0
        }
        next
    }
    {print}
    ' "$file" > "${file}.tmp" && mv "${file}.tmp" "$file"

    # Also add json import
    if ! grep -q '"encoding/json"' "$file"; then
        sed -i '/^import ($/a\	"encoding/json"' "$file"
    fi
    echo "✅ Fixed $file"
fi

# Fix bank_statement_line duplicate Status
echo ""
echo "🔧 Fixing bank_statement_line Status duplicates..."
file="internal/bank_statement_line/dto.go"
if [ -f "$file" ]; then
    awk '
    /Status.*string.*`json:"status"`/ {
        if (!seen[$0]++) {
            print
        } else {
            next
        }
    }
    {print}
    ' "$file" > "${file}.tmp" && mv "${file}.tmp" "$file"
    echo "✅ Fixed $file"
fi

# Fix bank_statement duplicate ImportSource
file="internal/bank_statement/dto.go"
if [ -f "$file" ]; then
    awk '
    /ImportSource.*string/ {
        if (!seen_import++) {
            print
        } else {
            next
        }
    }
    {print}
    ' "$file" > "${file}.tmp" && mv "${file}.tmp" "$file"
    echo "✅ Fixed bank_statement ImportSource"
fi

# Fix unused id variables in handlers
echo ""
echo "🔧 Fixing unused id variables..."
for file in internal/*/handler.go; do
    sed -i 's/^\(\s*\)id, err := uuid\.Parse(chi\.URLParam(r, "id"))$/\1_, err := uuid.Parse(chi.URLParam(r, "id"))/' "$file"
done
echo "✅ Fixed unused id variables"

echo ""
echo "✅ Quick fixes applied!"
