-- =====================================================
-- Accounting Seed Data: Fixed Assets & Depreciation
-- Description: Fixed assets with realistic depreciation schedules
-- Includes: 15 assets, 12 months of depreciation entries
-- =====================================================

\echo 'Loading fixed assets and depreciation schedules...';

-- =====================================================
-- SECTION 1: Fixed Asset Categories and Initial Assets
-- =====================================================

DO $$
DECLARE
    v_org_id UUID;
    v_fiscal_year_id UUID;
    v_asset_id UUID;
    v_furniture_account UUID;
    v_computer_account UUID;
    v_equipment_account UUID;
    v_vehicle_account UUID;
    v_building_account UUID;
    v_land_account UUID;
    v_accum_dep_furniture UUID;
    v_accum_dep_computer UUID;
    v_accum_dep_equipment UUID;
    v_accum_dep_vehicle UUID;
    v_accum_dep_building UUID;
    v_admin_user_id UUID;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_fiscal_year_id FROM fiscal_years WHERE organization_id = v_org_id AND fiscal_year = '2024' LIMIT 1;
    SELECT id INTO v_admin_user_id FROM users WHERE email = 'admin@demoretail.com' LIMIT 1;

    -- Get account IDs
    SELECT id INTO v_furniture_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1510' LIMIT 1;
    SELECT id INTO v_computer_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1520' LIMIT 1;
    SELECT id INTO v_equipment_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1530' LIMIT 1;
    SELECT id INTO v_vehicle_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1540' LIMIT 1;
    SELECT id INTO v_building_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1550' LIMIT 1;
    SELECT id INTO v_land_account FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1560' LIMIT 1;

    SELECT id INTO v_accum_dep_furniture FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1511' LIMIT 1;
    SELECT id INTO v_accum_dep_computer FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1521' LIMIT 1;
    SELECT id INTO v_accum_dep_equipment FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1531' LIMIT 1;
    SELECT id INTO v_accum_dep_vehicle FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1541' LIMIT 1;
    SELECT id INTO v_accum_dep_building FROM chart_of_accounts WHERE organization_id = v_org_id AND account_code = '1551' LIMIT 1;

    -- ========================================
    -- FURNITURE & FIXTURES (7-year life)
    -- ========================================

    -- Asset 1: Store Display Shelving
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-FUR-001', 'Store Display Shelving System', 'furniture',
        v_furniture_account, v_accum_dep_furniture,
        '2022-06-15', 18000.00, 1000.00,
        'straight_line', 7, 84,
        18000.00 - (17000.00 / 84 * 30), (17000.00 / 84 * 30), -- 30 months depreciation
        '2022-07-01', 'active', 'Main Store Floor', 'Custom retail shelving with LED lighting', v_admin_user_id
    );

    -- Asset 2: Office Furniture Set
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-FUR-002', 'Office Furniture - Desks & Chairs', 'furniture',
        v_furniture_account, v_accum_dep_furniture,
        '2021-03-10', 12500.00, 500.00,
        'straight_line', 7, 84,
        12500.00 - (12000.00 / 84 * 45), (12000.00 / 84 * 45), -- 45 months depreciation
        '2021-04-01', 'active', 'Back Office', 'Executive desks and ergonomic chairs', v_admin_user_id
    );

    -- Asset 3: Customer Waiting Area Furniture
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-FUR-003', 'Customer Lounge Furniture', 'furniture',
        v_furniture_account, v_accum_dep_furniture,
        '2023-01-20', 8500.00, 500.00,
        'straight_line', 7, 84,
        8500.00 - (8000.00 / 84 * 18), (8000.00 / 84 * 18), -- 18 months depreciation
        '2023-02-01', 'active', 'Store Entrance', 'Lounge chairs and coffee tables', v_admin_user_id
    );

    -- ========================================
    -- COMPUTER EQUIPMENT (5-year life)
    -- ========================================

    -- Asset 4: Point of Sale System
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-COMP-001', 'POS System - 5 Terminals', 'computer',
        v_computer_account, v_accum_dep_computer,
        '2022-09-01', 15000.00, 1000.00,
        'straight_line', 5, 60,
        15000.00 - (14000.00 / 60 * 27), (14000.00 / 60 * 27), -- 27 months depreciation
        '2022-09-01', 'active', 'Checkout Counters', 'Modern touchscreen POS terminals', v_admin_user_id
    );

    -- Asset 5: Office Computers & Laptops
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-COMP-002', 'Office Computers - 8 Units', 'computer',
        v_computer_account, v_accum_dep_computer,
        '2022-01-15', 12000.00, 800.00,
        'straight_line', 5, 60,
        12000.00 - (11200.00 / 60 * 35), (11200.00 / 60 * 35), -- 35 months depreciation
        '2022-02-01', 'active', 'Office', 'Dell workstations and laptops', v_admin_user_id
    );

    -- Asset 6: Network Equipment
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-COMP-003', 'Network Infrastructure', 'computer',
        v_computer_account, v_accum_dep_computer,
        '2021-11-10', 8000.00, 500.00,
        'straight_line', 5, 60,
        8000.00 - (7500.00 / 60 * 37), (7500.00 / 60 * 37), -- 37 months depreciation
        '2021-12-01', 'active', 'IT Room', 'Routers, switches, and servers', v_admin_user_id
    );

    -- ========================================
    -- STORE EQUIPMENT (10-year life)
    -- ========================================

    -- Asset 7: HVAC System
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-EQP-001', 'Commercial HVAC System', 'equipment',
        v_equipment_account, v_accum_dep_equipment,
        '2020-05-20', 35000.00, 3000.00,
        'straight_line', 10, 120,
        35000.00 - (32000.00 / 120 * 55), (32000.00 / 120 * 55), -- 55 months depreciation
        '2020-06-01', 'active', 'Building Roof', 'Trane commercial HVAC unit', v_admin_user_id
    );

    -- Asset 8: Security System
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-EQP-002', 'Security & Surveillance System', 'equipment',
        v_equipment_account, v_accum_dep_equipment,
        '2021-08-15', 18000.00, 2000.00,
        'straight_line', 10, 120,
        18000.00 - (16000.00 / 120 * 40), (16000.00 / 120 * 40), -- 40 months depreciation
        '2021-09-01', 'active', 'Store-wide', '24 HD cameras and monitoring system', v_admin_user_id
    );

    -- Asset 9: Lighting System
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-EQP-003', 'LED Lighting System Upgrade', 'equipment',
        v_equipment_account, v_accum_dep_equipment,
        '2022-03-01', 22000.00, 2000.00,
        'straight_line', 10, 120,
        22000.00 - (20000.00 / 120 * 33), (20000.00 / 120 * 33), -- 33 months depreciation
        '2022-03-01', 'active', 'Store-wide', 'Energy-efficient LED system', v_admin_user_id
    );

    -- ========================================
    -- VEHICLES (5-year life)
    -- ========================================

    -- Asset 10: Delivery Van
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-VEH-001', 'Delivery Van - Ford Transit', 'vehicle',
        v_vehicle_account, v_accum_dep_vehicle,
        '2022-07-10', 45000.00, 8000.00,
        'straight_line', 5, 60,
        45000.00 - (37000.00 / 60 * 29), (37000.00 / 60 * 29), -- 29 months depreciation
        '2022-08-01', 'active', 'Parking Lot', '2022 Ford Transit 250 - VIN: 1FTBR2XM7NKA12345', v_admin_user_id
    );

    -- Asset 11: Company Truck
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-VEH-002', 'Pickup Truck - Chevrolet Silverado', 'vehicle',
        v_vehicle_account, v_accum_dep_vehicle,
        '2021-04-20', 38000.00, 10000.00,
        'straight_line', 5, 60,
        38000.00 - (28000.00 / 60 * 44), (28000.00 / 60 * 44), -- 44 months depreciation
        '2021-05-01', 'active', 'Parking Lot', '2021 Chevrolet Silverado 1500 - VIN: 1GCUYGEL3MZ654321', v_admin_user_id
    );

    -- ========================================
    -- BUILDING IMPROVEMENTS (15-year life)
    -- ========================================

    -- Asset 12: Store Renovation
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-BLD-001', 'Store Front Renovation', 'building_improvement',
        v_building_account, v_accum_dep_building,
        '2020-01-15', 125000.00, 10000.00,
        'straight_line', 15, 180,
        125000.00 - (115000.00 / 180 * 59), (115000.00 / 180 * 59), -- 59 months depreciation
        '2020-02-01', 'active', 'Store Building', 'Complete storefront modernization', v_admin_user_id
    );

    -- Asset 13: Warehouse Expansion
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-BLD-002', 'Warehouse Storage Expansion', 'building_improvement',
        v_building_account, v_accum_dep_building,
        '2021-06-01', 85000.00, 5000.00,
        'straight_line', 15, 180,
        85000.00 - (80000.00 / 180 * 42), (80000.00 / 180 * 42), -- 42 months depreciation
        '2021-07-01', 'active', 'Back Warehouse', 'Additional 2000 sq ft storage space', v_admin_user_id
    );

    -- ========================================
    -- LAND (No depreciation)
    -- ========================================

    -- Asset 14: Store Property Land
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-LAND-001', 'Commercial Property - 1.5 Acres', 'land',
        v_land_account, NULL, -- Land does not depreciate
        '2019-03-15', 250000.00, 250000.00, -- Salvage = Cost (no depreciation)
        'none', 0, 0,
        250000.00, 0,
        NULL, 'active', '123 Main Street', 'Commercial zoned property, 1.5 acres', v_admin_user_id
    );

    -- Asset 15: Parking Lot Paving
    v_asset_id := gen_random_uuid();
    INSERT INTO fixed_assets (
        id, organization_id, asset_number, asset_name, asset_category,
        asset_account_id, accumulated_depreciation_account_id,
        acquisition_date, acquisition_cost, salvage_value,
        depreciation_method, useful_life_years, useful_life_months,
        current_book_value, accumulated_depreciation,
        depreciation_start_date, status, location, notes, created_by
    )
    VALUES (
        v_asset_id, v_org_id, 'FA-BLD-003', 'Parking Lot - Asphalt & Striping', 'building_improvement',
        v_building_account, v_accum_dep_building,
        '2022-04-10', 45000.00, 2000.00,
        'straight_line', 15, 180,
        45000.00 - (43000.00 / 180 * 32), (43000.00 / 180 * 32), -- 32 months depreciation
        '2022-05-01', 'active', 'Parking Area', '50-space customer parking lot', v_admin_user_id
    );

    RAISE NOTICE 'Created 15 fixed assets';

END $$;

-- =====================================================
-- SECTION 2: Depreciation Schedule for 2024
-- =====================================================

DO $$
DECLARE
    v_org_id UUID;
    v_fiscal_year_id UUID;
    v_asset_record RECORD;
    v_period_id UUID;
    v_schedule_date DATE;
    v_monthly_depreciation NUMERIC(20,4);
    v_schedule_count INTEGER := 0;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1;
    SELECT id INTO v_fiscal_year_id FROM fiscal_years WHERE organization_id = v_org_id AND fiscal_year = '2024' LIMIT 1;

    -- Generate depreciation schedule for each depreciable asset for all 12 months of 2024
    FOR v_asset_record IN (
        SELECT id, asset_number, asset_name, acquisition_cost, salvage_value,
               useful_life_months, accumulated_depreciation_account_id
        FROM fixed_assets
        WHERE organization_id = v_org_id
            AND depreciation_method != 'none'
            AND status = 'active'
    ) LOOP

        -- Calculate monthly depreciation
        v_monthly_depreciation := (v_asset_record.acquisition_cost - v_asset_record.salvage_value) /
                                  NULLIF(v_asset_record.useful_life_months, 0);

        -- Create depreciation schedule entries for each month of 2024
        FOR month_num IN 1..12 LOOP
            v_schedule_date := ('2024-' || LPAD(month_num::TEXT, 2, '0') || '-28')::DATE;

            SELECT id INTO v_period_id
            FROM accounting_periods
            WHERE fiscal_year_id = v_fiscal_year_id
                AND period_number = month_num
            LIMIT 1;

            INSERT INTO asset_depreciation_schedule (
                id, fixed_asset_id, organization_id,
                fiscal_year_id, accounting_period_id,
                depreciation_date, depreciation_amount,
                accumulated_depreciation_account_id,
                status, notes, created_by
            )
            SELECT
                gen_random_uuid(),
                v_asset_record.id,
                v_org_id,
                v_fiscal_year_id,
                v_period_id,
                v_schedule_date,
                v_monthly_depreciation,
                v_asset_record.accumulated_depreciation_account_id,
                'posted',
                'Monthly depreciation - ' || v_asset_record.asset_name,
                u.id
            FROM users u
            WHERE u.email = 'admin@demoretail.com'
            LIMIT 1;

            v_schedule_count := v_schedule_count + 1;

        END LOOP;
    END LOOP;

    RAISE NOTICE 'Created % depreciation schedule entries (12 months × 14 assets)', v_schedule_count;

END $$;

\echo '';
\echo '==========================================';
\echo 'Fixed Assets Summary:';
\echo '==========================================';
\echo 'Furniture & Fixtures:        3 assets';
\echo 'Computer Equipment:          3 assets';
\echo 'Store Equipment:             3 assets';
\echo 'Vehicles:                    2 assets';
\echo 'Building Improvements:       3 assets';
\echo 'Land:                        1 asset (no depreciation)';
\echo '==========================================';
\echo 'TOTAL FIXED ASSETS:         15 assets';
\echo '==========================================';
\echo '';
\echo 'Depreciation Schedule:';
\echo '  ✓ 12 months of depreciation for 2024';
\echo '  ✓ 168 depreciation schedule entries';
\echo '  ✓ Straight-line depreciation method';
\echo '  ✓ Realistic useful life periods';
\echo '  ✓ Prior depreciation calculated';
\echo '';
\echo 'Fixed assets data loaded successfully!';

-- =====================================================
-- End of seed file 005
-- =====================================================
