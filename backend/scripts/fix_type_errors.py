#!/usr/bin/env python3
"""Fix all pointer type assignment errors"""

import re
import glob

def fix_pointer_assignments(filepath):
    """Fix incorrect pointer assignments in service files"""
    with open(filepath, 'r') as f:
        content = f.read()

    # Pattern 1: entity.Field = *req.Field where Field is a pointer
    # Should be: entity.Field = req.Field
    # Match lines like: entity.SomeField = *req.SomeField
    content = re.sub(
        r'entity\.(\w+)\s*=\s*\*req\.(\w+)',
        r'entity.\1 = req.\2',
        content
    )

    # Pattern 2: cannot use req.Field (pointer) as uuid.UUID
    # Should dereference: *req.Field
    # This is trickier - look for validator patterns
    content = re.sub(
        r'v\.validate(\w+)\(req\.(\w+)\)',
        r'v.validate\1(*req.\2)',
        content
    )

    with open(filepath, 'w') as f:
        f.write(content)

def fix_bank_statement_line_dto():
    """Fix duplicate Status field in bank_statement_line"""
    filepath = 'internal/bank_statement_line/dto.go'
    with open(filepath, 'r') as f:
        lines = f.readlines()

    new_lines = []
    status_seen = False

    for line in lines:
        if 'Status' in line and '`json:"status"`' in line:
            if not status_seen:
                new_lines.append(line)
                status_seen = True
            # Skip duplicates
        else:
            new_lines.append(line)

    with open(filepath, 'w') as f:
        f.writelines(new_lines)

def fix_background_job_orgid():
    """Fix OrganizationId type issue in background_job"""
    filepath = 'internal/background_job/service.go'
    with open(filepath, 'r') as f:
        content = f.read()

    # Fix: OrganizationId: orgID (where orgID is uuid.UUID but field expects *uuid.UUID)
    # Change to: OrganizationId: &orgID
    content = content.replace(
        'OrganizationId: orgID,',
        'OrganizationId: &orgID,'
    )

    # Fix comparison: entity.OrganizationId != orgID
    # Change to: *entity.OrganizationId != orgID
    content = content.replace(
        'entity.OrganizationId != orgID',
        '*entity.OrganizationId != orgID'
    )

    with open(filepath, 'w') as f:
        f.write(content)

def main():
    print("🔧 Fixing Type Assignment Errors")
    print("=" * 50)

    print("\n📦 Fixing pointer assignments in service files...")
    count = 0
    for filepath in glob.glob('internal/*/service.go'):
        try:
            fix_pointer_assignments(filepath)
            count += 1
        except Exception as e:
            print(f"  ❌ {filepath}: {e}")

    print(f"  ✅ Fixed {count} service files")

    print("\n📦 Fixing bank_statement_line Status duplicates...")
    try:
        fix_bank_statement_line_dto()
        print("  ✅ Fixed")
    except Exception as e:
        print(f"  ❌ Error: {e}")

    print("\n📦 Fixing background_job OrganizationId...")
    try:
        fix_background_job_orgid()
        print("  ✅ Fixed")
    except Exception as e:
        print(f"  ❌ Error: {e}")

    print("\n✅ Type fixes completed!")

if __name__ == '__main__':
    main()
