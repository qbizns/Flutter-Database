# POS Backend API - Postman Collection

This directory contains a comprehensive Postman collection and environment for testing the POS Backend API.

## Contents

- **POS_Backend_API.postman_collection.json** - Complete API collection with 40+ requests
- **POS_Backend.postman_environment.json** - Environment variables for local/staging/production
- **README.md** - This file

## Features

### Organized Folder Structure

- 🔐 **Authentication** (3 requests)
  - Login
  - Register
  - Refresh Token

- 💊 **Health & Monitoring** (2 requests)
  - Health Check
  - Prometheus Metrics

- 📦 **Products** (5 requests)
  - List, Create, Get, Update, Delete

- 👥 **Customers** (5 requests)
  - List, Create, Get, Update, Delete

- 🏢 **Suppliers** (5 requests)
  - List, Create, Get, Update, Delete

- 📂 **Categories** (5 requests)
  - List, Create, Get, Update, Delete

- 📍 **Locations** (5 requests)
  - List, Create, Get, Update, Delete

- 💰 **Sales** (3 requests)
  - List, Create, Get

### Automated Testing

Every request includes comprehensive test scripts that:
- Validate HTTP status codes
- Check response structure
- Verify data integrity
- Measure response times
- Store IDs in environment variables

### Smart Automation

- **Auto-authentication**: Tokens automatically stored after login
- **Variable chaining**: IDs from LIST requests auto-populate subsequent requests
- **Dynamic data**: Uses Postman's dynamic variables ({{$randomEmail}}, etc.)
- **Error handling**: Global test scripts provide helpful error messages

## Quick Start

### 1. Import Collection and Environment

```bash
# In Postman:
1. Click "Import" button
2. Select both JSON files from this directory
3. Collection and Environment will be imported
```

### 2. Configure Environment

1. Select "POS Backend Environment" from the environment dropdown
2. Click the eye icon to view/edit variables
3. Update required variables:
   - `organization_id` - Your organization UUID
   - `test_user_email` - Your test user email (default: admin@example.com)
   - `test_user_password` - Your test user password

For different environments:
- **Local**: Use `base_url` = `http://localhost:8080`
- **Staging**: Enable `staging_url` variable
- **Production**: Enable `production_url` variable

### 3. Run Your First Request

1. Open "🔐 Authentication" folder
2. Click "Login" request
3. Click "Send"
4. Check "Tests" tab - all tests should pass ✅
5. Access token is now stored automatically

### 4. Try Other Endpoints

After logging in, try these workflows:

**Create a Product:**
```
1. Create Category → Stores category_id
2. Create Product → Uses category_id, stores product_id
3. Get Product → Retrieves the created product
```

**Create a Sale:**
```
1. List Customers → Stores customer_id
2. List Locations → Stores location_id
3. List Products → Stores product_id
4. Create Sale → Uses all stored IDs
```

## Running Collection Tests

### Via Postman Runner

1. Click "Runner" button in Postman
2. Select "POS Backend API - Complete Collection"
3. Select "POS Backend Environment"
4. Click "Run POS Backend API"
5. Watch all tests execute sequentially

### Via Newman (CLI)

```bash
# Install Newman
npm install -g newman

# Run collection
newman run POS_Backend_API.postman_collection.json \
  -e POS_Backend.postman_environment.json \
  --reporters cli,json

# Run with detailed HTML report
newman run POS_Backend_API.postman_collection.json \
  -e POS_Backend.postman_environment.json \
  --reporters cli,htmlextra \
  --reporter-htmlextra-export report.html
```

## Environment Variables

### Base Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| `base_url` | API base URL | `http://localhost:8080` |
| `organization_id` | Organization UUID | Set manually |
| `test_user_email` | Test user login | `admin@example.com` |
| `test_user_password` | Test user password | Set manually (secret) |

### Authentication (Auto-populated)

| Variable | Description | Set By |
|----------|-------------|--------|
| `access_token` | JWT access token (1hr) | Login request |
| `refresh_token` | JWT refresh token (7d) | Login request |
| `user_id` | Current user ID | Login request |
| `csrf_token` | CSRF protection token | Auto-extracted |

### Entity IDs (Auto-populated)

| Variable | Description | Set By |
|----------|-------------|--------|
| `product_id` | Sample product ID | List Products |
| `created_product_id` | New product ID | Create Product |
| `customer_id` | Sample customer ID | List Customers |
| `created_customer_id` | New customer ID | Create Customer |
| `supplier_id` | Sample supplier ID | List Suppliers |
| `created_supplier_id` | New supplier ID | Create Supplier |
| `category_id` | Sample category ID | List Categories |
| `created_category_id` | New category ID | Create Category |
| `location_id` | Sample location ID | List Locations |
| `created_location_id` | New location ID | Create Location |
| `sale_id` | Sample sale ID | List Sales |
| `created_sale_id` | New sale ID | Create Sale |

## Testing Best Practices

### 1. Run Requests in Order

For best results, run requests in this order:
1. **Login** - Get authentication token
2. **List endpoints** - Populate entity IDs
3. **Create endpoints** - Create test data
4. **Get/Update endpoints** - Test with created data
5. **Delete endpoints** - Clean up test data

### 2. Check Test Results

Every request includes tests. Check the "Test Results" tab to see:
- ✅ Passed tests (green)
- ❌ Failed tests (red)
- Test execution time

### 3. Monitor Console Output

The Console tab (View → Show Postman Console) displays:
- Request/response logs
- Test script console.log() output
- Variable changes
- Error messages with helpful hints

### 4. Rate Limiting

Some endpoints have rate limits:
- **Auth endpoints**: 5 requests/minute
- **Regular endpoints**: 100 requests/minute

If you hit rate limits:
- Wait 60 seconds
- Check X-RateLimit headers in responses
- Use delays in Collection Runner

## Troubleshooting

### "Unauthorized" Errors (401)

**Problem**: Access token expired or missing

**Solution**:
1. Run the "Login" request again
2. Check that `access_token` is set in environment
3. Try "Refresh Token" request

### "Forbidden" Errors (403)

**Problem**: User lacks required permissions

**Solution**:
1. Check user roles in database
2. Verify required permissions in AUTHORIZATION_GUIDE.md
3. Use admin/superadmin account for testing

### "Rate Limited" Errors (429)

**Problem**: Too many requests

**Solution**:
1. Wait 60 seconds before retrying
2. Add delays in Collection Runner settings
3. Check rate limit headers in response

### Variables Not Auto-Populating

**Problem**: IDs not saved from requests

**Solution**:
1. Check "Tests" tab for errors
2. Ensure response returns expected data
3. Verify test scripts are executing

### Connection Refused

**Problem**: Backend server not running

**Solution**:
1. Start backend server: `go run cmd/server/main.go`
2. Check `base_url` matches server address
3. Verify no firewall blocking localhost:8080

## Advanced Usage

### Pre-request Scripts

Global pre-request script automatically:
- Checks for access token
- Logs request details
- Warns if token missing

You can add custom pre-request scripts to individual requests.

### Response Examples

To add response examples:
1. Send request successfully
2. Click "Save Response" → "Save as Example"
3. Examples help document expected responses

### Chaining Requests

Requests already chain automatically via test scripts. To manually chain:

```javascript
// In test script of Request A:
const responseData = pm.response.json();
pm.environment.set("custom_variable", responseData.some_value);

// In Request B:
// Use {{custom_variable}} in URL/body
```

### Custom Test Scripts

Add custom tests to any request:

```javascript
pm.test("Custom validation", function () {
    const response = pm.response.json();
    pm.expect(response.custom_field).to.equal("expected_value");
});
```

## API Documentation

For complete API documentation, see:
- **openapi.yaml** - OpenAPI 3.0.3 specification
- **AUTHORIZATION_GUIDE.md** - RBAC permission system
- **DEPLOYMENT_GUIDE.md** - Deployment instructions

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review test script output in Console
3. Check backend logs for error details
4. Refer to OpenAPI spec for endpoint documentation

## Collection Maintenance

To keep the collection updated:

1. **New endpoints**: Add to appropriate folder with tests
2. **Deprecated endpoints**: Disable or remove
3. **Updated schemas**: Update request bodies and tests
4. **New variables**: Add to environment file with description

## Version History

- **v1.0.0** (2025-01-11)
  - Initial comprehensive collection
  - 40+ requests across 8 folders
  - Complete test coverage
  - Auto-variable management
  - Multi-environment support
