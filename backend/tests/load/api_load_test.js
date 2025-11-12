import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const postingDuration = new Trend('posting_duration');

// Test configuration
export const options = {
  stages: [
    { duration: '2m', target: 50 },   // Ramp up to 50 users
    { duration: '5m', target: 50 },   // Stay at 50 users
    { duration: '2m', target: 100 },  // Ramp up to 100 users
    { duration: '5m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 200 },  // Ramp up to 200 users
    { duration: '5m', target: 200 },  // Stay at 200 users
    { duration: '2m', target: 500 },  // Spike to 500 users
    { duration: '1m', target: 500 },  // Hold spike
    { duration: '3m', target: 0 },    // Ramp down
  ],
  thresholds: {
    'http_req_duration': ['p(95)<200'],   // 95% of requests under 200ms
    'http_req_failed': ['rate<0.01'],     // Error rate under 1%
    'errors': ['rate<0.05'],              // Custom error rate under 5%
    'posting_duration': ['p(95)<500'],    // Posting under 500ms
  },
};

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const API_VERSION = '/api/v1';

// Test data
const testUser = {
  email: `loadtest_${__VU}@example.com`,
  password: 'LoadTest123!@#',
};

let authToken = '';
let orgId = '';

// Setup: Register and login once per VU
export function setup() {
  console.log('Setting up load test...');
  console.log(`Target: ${BASE_URL}`);
  console.log(`VUs: ${__VU}`);

  return { baseUrl: BASE_URL };
}

export default function (data) {
  // If not authenticated, login first
  if (!authToken) {
    const loginRes = http.post(`${BASE_URL}${API_VERSION}/auth/login`, JSON.stringify(testUser), {
      headers: { 'Content-Type': 'application/json' },
    });

    if (loginRes.status === 200) {
      const body = JSON.parse(loginRes.body);
      authToken = body.token;
      orgId = body.organization_id;
    } else {
      // If login fails, try to register
      const registerRes = http.post(`${BASE_URL}${API_VERSION}/auth/register`, JSON.stringify({
        ...testUser,
        name: `Load Test User ${__VU}`,
      }), {
        headers: { 'Content-Type': 'application/json' },
      });

      if (registerRes.status === 201) {
        const body = JSON.parse(registerRes.body);
        authToken = body.token;
        orgId = body.organization_id;
      } else {
        errorRate.add(1);
        console.error('Failed to register/login');
        return;
      }
    }
  }

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${authToken}`,
  };

  // Test 1: Health check (lightweight)
  {
    const res = http.get(`${BASE_URL}/health`);
    check(res, {
      'health check status is 200': (r) => r.status === 200,
      'health check response time < 50ms': (r) => r.timings.duration < 50,
    }) || errorRate.add(1);
  }

  sleep(0.5);

  // Test 2: List products (cached after first request)
  {
    const res = http.get(`${BASE_URL}${API_VERSION}/organizations/${orgId}/products`, {
      headers: headers,
    });

    check(res, {
      'list products status is 200': (r) => r.status === 200,
      'list products response time < 200ms': (r) => r.timings.duration < 200,
    }) || errorRate.add(1);
  }

  sleep(0.5);

  // Test 3: Get product (should hit cache)
  {
    const res = http.get(`${BASE_URL}${API_VERSION}/organizations/${orgId}/products/1`, {
      headers: headers,
    });

    check(res, {
      'get product response time < 100ms': (r) => r.timings.duration < 100,
    });
  }

  sleep(0.5);

  // Test 4: List customers
  {
    const res = http.get(`${BASE_URL}${API_VERSION}/organizations/${orgId}/customers`, {
      headers: headers,
    });

    check(res, {
      'list customers status is 200': (r) => r.status === 200,
      'list customers response time < 200ms': (r) => r.timings.duration < 200,
    }) || errorRate.add(1);
  }

  sleep(0.5);

  // Test 5: Create sale (write operation)
  {
    const sale = {
      customer_id: '00000000-0000-0000-0000-000000000001',
      total_amount: 99.99,
      payment_method: 'CASH',
      items: [
        { product_id: '1', quantity: 2, price: 49.995 },
      ],
    };

    const res = http.post(
      `${BASE_URL}${API_VERSION}/organizations/${orgId}/sales`,
      JSON.stringify(sale),
      { headers: headers }
    );

    check(res, {
      'create sale status is 201': (r) => r.status === 201,
      'create sale response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);
  }

  sleep(1);

  // Test 6: Posting engine (most critical operation)
  {
    const postingReq = {
      document_type: 'POS_SALE',
      document_id: '00000000-0000-0000-0000-000000000001',
      event: 'on_post',
    };

    const startTime = Date.now();
    const res = http.post(
      `${BASE_URL}${API_VERSION}/organizations/${orgId}/posting/post`,
      JSON.stringify(postingReq),
      { headers: headers }
    );
    const duration = Date.now() - startTime;

    postingDuration.add(duration);

    check(res, {
      'posting status is 200 or 422': (r) => r.status === 200 || r.status === 422,
      'posting response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);
  }

  sleep(1);

  // Test 7: Dashboard/Reports (complex queries)
  {
    const res = http.get(
      `${BASE_URL}${API_VERSION}/organizations/${orgId}/reports/trial-balance`,
      { headers: headers }
    );

    check(res, {
      'report response time < 1000ms': (r) => r.timings.duration < 1000,
    });
  }

  sleep(2);
}

export function teardown(data) {
  console.log('Load test completed');
}

// Handle termination
export function handleSummary(data) {
  return {
    'summary.json': JSON.stringify(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}

function textSummary(data, options) {
  const indent = options.indent || '';
  const lines = [];

  lines.push(`${indent}Test Results:`);
  lines.push(`${indent}  Checks................: ${(data.metrics.checks.values.passes / data.metrics.checks.values.count * 100).toFixed(2)}% passed`);
  lines.push(`${indent}  HTTP Requests.........: ${data.metrics.http_reqs.values.count}`);
  lines.push(`${indent}  Request Duration (avg): ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms`);
  lines.push(`${indent}  Request Duration (p95): ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms`);
  lines.push(`${indent}  Error Rate...........: ${(data.metrics.http_req_failed.values.rate * 100).toFixed(2)}%`);
  lines.push(`${indent}  Posting Duration (p95): ${data.metrics.posting_duration.values['p(95)'].toFixed(2)}ms`);

  return lines.join('\n');
}
