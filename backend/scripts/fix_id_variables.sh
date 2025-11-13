#!/bin/bash

# Fix the id variable issues in handlers
cd /home/user/Flutter-Database/backend

echo "🔧 Fixing id variables in handler files..."

for file in internal/*/handler.go; do
    # Remove the broken "_ = id // TODO: implement" lines
    sed -i '/_ = id \/\/ TODO: implement/d' "$file"

    # Fix the pattern: _, err := uuid.Parse(chi.URLParam(r, "id"))
    # Change to: id, err := uuid.Parse(chi.URLParam(r, "id"))
    # But only in GetByID, Update, and Delete methods (not in admin methods)

    # Use awk to fix this properly
    awk '
    /func \(h \*Handler\) (GetByID|Update|Delete)\(/ {
        in_crud = 1
    }
    /func \(h \*Handler\) (AdminList|GetStats|Export|Import|ListDeleted|Restore|PermanentDelete)/ {
        in_crud = 0
    }
    in_crud && /_, err := uuid\.Parse\(chi\.URLParam\(r, "id"\)\)/ {
        gsub(/_, err :=/, "id, err :=")
    }
    { print }
    ' "$file" > "${file}.tmp" && mv "${file}.tmp" "$file"

    echo "✅ $file"
done

echo ""
echo "✅ ID variables fixed!"
