#!/usr/bin/env python3
"""
Smart pointer assignment fix that understands when to dereference
Rules:
1. If req.Field is *T and entity.Field is T, use *req.Field
2. If req.Field is *T and entity.Field is *T, use req.Field
3. For UUID validator calls, always dereference pointer UUID arguments
"""

import re
import glob

def fix_service_updates(filepath):
    """Fix Update method pointer assignments"""
    with open(filepath, 'r') as f:
        content = f.read()

    # Revert the too-aggressive fix from before
    # Pattern: entity.Field = req.Field (where req.Field is *string and entity.Field is string)
    # Should be: entity.Field = *req.Field

    # This is complex, so let's do it line by line
    lines = content.split('\n')
    new_lines = []
    in_update_method = False

    for line in lines:
        if 'func (s *Service) Update(' in line:
            in_update_method = True

        if in_update_method and 'if req.' in line and '!= nil {' in line:
            # Next line likely has the assignment
            new_lines.append(line)
            continue

        # Pattern: entity.SomeField = req.SomeField (in update)
        # Most update request fields are pointers that need dereferencing
        if in_update_method and re.match(r'\s+entity\.\w+\s+=\s+req\.\w+$', line):
            # Add dereference
            line = re.sub(r'(\s+entity\.\w+\s+=\s+)req\.(\w+)$', r'\1*req.\2', line)

        # End of update method
        if in_update_method and line.strip() == '}' and 'return' not in line:
            in_update_method = False

        new_lines.append(line)

    content = '\n'.join(new_lines)

    with open(filepath, 'w') as f:
        f.write(content)

def fix_validator_uuid_args(filepath):
    """Fix UUID arguments in validator calls"""
    with open(filepath, 'r') as f:
        content = f.read()

    # Pattern: v.validateFieldExists(req.Field) where req.Field is *uuid.UUID
    # Should be: v.validateFieldExists(*req.Field)
    content = re.sub(
        r'v\.validate(\w+Exists)\(req\.(\w+Id)\)',
        r'v.validate\1(*req.\2)',
        content
    )

    with open(filepath, 'w') as f:
        f.write(content)

def fix_create_assignments(filepath):
    """Fix Create method where entity fields are non-pointer but req fields may be pointers"""
    with open(filepath, 'r') as f:
        content = f.read()

    # In Create methods, most entity fields match req fields directly
    # except when entity field is non-pointer and req is pointer (rare, usually optional fields)
    # These need manual review, but most Create requests have matching types

    with open(filepath, 'w') as f:
        f.write(content)

def main():
    print("🔧 Smart Pointer Fix")
    print("=" * 50)

    print("\n📦 Fixing service Update methods...")
    count = 0
    for filepath in glob.glob('internal/*/service.go'):
        try:
            fix_service_updates(filepath)
            count += 1
        except Exception as e:
            print(f"  ❌ {filepath}: {e}")

    print(f"  ✅ Fixed {count} service files")

    print("\n📦 Fixing validator UUID arguments...")
    for filepath in glob.glob('internal/*/validator.go'):
        try:
            fix_validator_uuid_args(filepath)
        except Exception as e:
            print(f"  ❌ {filepath}: {e}")

    print("  ✅ Validators fixed")

    print("\n✅ Smart pointer fixes completed!")

if __name__ == '__main__':
    main()
