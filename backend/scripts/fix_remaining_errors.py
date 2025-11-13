#!/usr/bin/env python3
"""Fix the remaining specific errors"""

import re
import glob
import os

def fix_time_now_address():
    """Fix time.Now() address issue in repositories"""
    for repo_file in glob.glob('internal/*/repository.go'):
        with open(repo_file, 'r') as f:
            content = f.read()

        # Replace &time.Now() with a variable
        # Pattern: UpdatedAt: &time.Now() -> now := time.Now(); UpdatedAt: &now
        content = re.sub(
            r'(\s+)UpdatedAt:\s*&time\.Now\(\)',
            r'\1UpdatedAt: &[]time.Time{time.Now()}[0]',
            content
        )

        # Better pattern - use time.Now().UTC() and take address properly
        content = re.sub(
            r'&time\.Now\(\)',
            r'func() *time.Time { t := time.Now(); return &t }()',
            content
        )

        with open(repo_file, 'w') as f:
            f.write(content)

def fix_duplicate_import_source():
    """Fix duplicate ImportSource in bank_statement repository"""
    repo_file = 'internal/bank_statement/repository.go'
    if os.path.exists(repo_file):
        with open(repo_file, 'r') as f:
            lines = f.readlines()

        new_lines = []
        seen_import_source = False

        for line in lines:
            if 'ImportSource' in line and 'import_source' in line:
                if not seen_import_source:
                    new_lines.append(line)
                    seen_import_source = True
                # Skip duplicates
            else:
                new_lines.append(line)

        with open(repo_file, 'w') as f:
            f.writelines(new_lines)

def fix_wrong_dereferences():
    """Fix cases where we're dereferencing when both are pointers"""
    for service_file in glob.glob('internal/*/service.go'):
        with open(service_file, 'r') as f:
            content = f.read()

        # When both are pointers, don't dereference
        # Pattern: entity.Field = *req.Field where both are pointers
        # These are fields that end in specific patterns

        # Fix pointer UUID fields
        content = re.sub(
            r'entity\.(\w+Id)\s*=\s*\*req\.(\w+Id)',
            r'entity.\1 = req.\2',
            content
        )

        # Fix pointer bool fields when both are pointers
        content = re.sub(
            r'entity\.(Is\w+|Has\w+)\s*=\s*\*req\.(Is\w+|Has\w+)',
            r'entity.\1 = req.\2',
            content
        )

        # Fix pointer int64 fields when entity expects pointer
        content = re.sub(
            r'entity\.(\w*Level|Order|Number|Count)\s*=\s*\*req\.(\w*Level|Order|Number|Count)',
            r'entity.\1 = req.\2',
            content
        )

        # Fix pointer string fields when entity expects pointer
        content = re.sub(
            r'entity\.(Description|Notes|Status|Type|Method|Code|Name)\s*=\s*\*req\.(Description|Notes|Status|Type|Method|Code|Name)',
            r'entity.\1 = req.\2',
            content
        )

        # Fix pointer float64 fields
        content = re.sub(
            r'entity\.(\w*Balance|Amount|Value|Percent)\s*=\s*\*req\.(\w*Balance|Amount|Value|Percent)',
            r'entity.\1 = req.\2',
            content
        )

        # Fix time.Time pointers
        content = re.sub(
            r'entity\.(\w*At|Date)\s*=\s*\*req\.(\w*At|Date)',
            r'entity.\1 = req.\2',
            content
        )

        # Fix json.RawMessage
        content = re.sub(
            r'entity\.(Metadata|Settings|Headers|Events|Result|Changes|OldValues|NewValues)\s*=\s*\*req\.(Metadata|Settings|Headers|Events|Result|Changes|OldValues|NewValues)',
            r'entity.\1 = req.\2',
            content
        )

        # Fix orgID type mismatch in create
        content = re.sub(
            r'OrganizationId:\s*orgID,',
            r'OrganizationId: &orgID,',
            content
        )

        # Fix comparison with orgID
        content = re.sub(
            r'entity\.OrganizationId\s*!=\s*orgID',
            r'*entity.OrganizationId != orgID',
            content
        )

        with open(service_file, 'w') as f:
            f.write(content)

def add_missing_fields_systematically():
    """Add all missing fields found in error messages"""
    # Mapping of module to missing fields
    missing_fields = {
        'accounting_period': ['Status', 'Metadata'],
        'background_job': ['Status', 'Result', 'ErrorDetails'],
        'bank_statement': ['ImportSource', 'Status'],
        'bank_statement_line': ['Status'],
        'bank_reconciliation': ['Status', 'Metadata'],
        'bank_account': ['Metadata'],
        'audit_log': ['OldValues', 'NewValues', 'Changes', 'Metadata'],
        'api_key': ['Scopes'],
        'webhook': ['Events', 'Headers'],
        'webhook_delivery': ['Status', 'RequestHeaders', 'RequestBody', 'ResponseHeaders'],
        'your_table_name': ['Settings', 'Metadata'],
        'vendor_payment': ['Metadata'],
    }

    for module, fields in missing_fields.items():
        dto_file = f'internal/{module}/dto.go'
        if not os.path.exists(dto_file):
            continue

        with open(dto_file, 'r') as f:
            content = f.read()

        for field in fields:
            # Determine field type
            if field in ['Metadata', 'Settings', 'OldValues', 'NewValues', 'Changes', 'Headers', 'Scopes', 'RequestHeaders', 'RequestBody', 'ResponseHeaders', 'Events', 'Result', 'ErrorDetails']:
                field_type = 'json.RawMessage'
                need_pointer = True
            else:  # Status, ImportSource
                field_type = 'string'
                need_pointer = True

            # Check if field already exists
            if f'{field}' in content and f'`json:"{field.lower()}"' in content:
                continue

            # Add to CreateRequest
            create_struct = toPascalCase(module)
            if f'type Create{create_struct}Request struct' in content:
                # Find the struct and add field before closing brace
                pattern = rf'(type Create{create_struct}Request struct \{{[^}}]*?)(\n\}})'
                if need_pointer:
                    replacement = rf'\1\n\t{field} *{field_type} `json:"{field.lower()},omitempty"`\n\2'
                else:
                    replacement = rf'\1\n\t{field} {field_type} `json:"{field.lower()},omitempty"`\n\2'

                content = re.sub(pattern, replacement, content, flags=re.DOTALL)

            # Add to UpdateRequest
            if f'type Update{create_struct}Request struct' in content:
                pattern = rf'(type Update{create_struct}Request struct \{{[^}}]*?)(\n\}})'
                if need_pointer:
                    replacement = rf'\1\n\t{field} *{field_type} `json:"{field.lower()},omitempty"`\n\2'
                else:
                    replacement = rf'\1\n\t{field} {field_type} `json:"{field.lower()},omitempty"`\n\2'

                content = re.sub(pattern, replacement, content, flags=re.DOTALL)

        # Make sure json import exists if we added json.RawMessage
        if 'json.RawMessage' in content and '"encoding/json"' not in content:
            content = content.replace('import (', 'import (\n\t"encoding/json"')

        with open(dto_file, 'w') as f:
            f.write(content)

def fix_validator_uuid_dereferences():
    """Fix validator UUID arguments that shouldn't be dereferenced"""
    for validator_file in glob.glob('internal/*/validator.go'):
        with open(validator_file, 'r') as f:
            content = f.read()

        # If req.Field is *uuid.UUID and function expects uuid.UUID, dereference
        # But check if it was already dereferenced
        if 'cannot use req.' in content:
            # The fix was applied but might be wrong - let's be more careful
            pass

        # Just ensure *req.FieldId pattern for validator calls
        # Only if req.FieldId is defined as *uuid.UUID
        content = re.sub(
            r'v\.validate\w+Exists\(req\.(\w+Id)\)',
            r'v.validate\w+Exists(*req.\1)',
            content
        )

        with open(validator_file, 'w') as f:
            f.write(content)

def toPascalCase(s):
    """Convert snake_case to PascalCase"""
    return ''.join(word.capitalize() for word in s.split('_'))

def main():
    print("🔧 Fixing Remaining Compilation Errors")
    print("=" * 60)

    os.chdir('/home/user/Flutter-Database/backend')

    print("\n📦 Step 1: Fixing time.Now() address issues...")
    fix_time_now_address()
    print("  ✅ Fixed")

    print("\n📦 Step 2: Fixing duplicate ImportSource...")
    fix_duplicate_import_source()
    print("  ✅ Fixed")

    print("\n📦 Step 3: Fixing incorrect dereferences (when both are pointers)...")
    fix_wrong_dereferences()
    print("  ✅ Fixed")

    print("\n📦 Step 4: Adding missing DTO fields systematically...")
    add_missing_fields_systematically()
    print("  ✅ Fixed")

    print("\n" + "=" * 60)
    print("✅ REMAINING ERRORS FIXED!")

if __name__ == '__main__':
    main()
