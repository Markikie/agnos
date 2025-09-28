# Agnos Hospital Middleware API Specification

## Base URL
- **Development**: `https://localhost:443`
- **Production**: `https://hospital-a.api.co.th`

## Authentication
The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <access_token>
```

## Content Type
All requests and responses use JSON format:
```
Content-Type: application/json
```

---

## Staff Management APIs

### 1. Create Staff Member
Creates a new hospital staff member with login credentials.

**Endpoint**: `POST /staff/create`

**Request Body**:
```json
{
    "username": "string (required)",
    "password": "string (required)",
    "hospital": "string (required)"
}
```

**Response**:
- **201 Created**:
```json
{
    "message": "Staff created successfully",
    "staff_id": "uuid"
}
```

- **400 Bad Request**:
```json
{
    "error": "username, password, and hospital are required"
}
```

- **409 Conflict**:
```json
{
    "error": "staff with this username already exists in this hospital"
}
```

**Example**:
```bash
curl -X POST https://localhost:443/staff/create \
  -H "Content-Type: application/json" \
  -d '{
    "username": "doctor001",
    "password": "securePassword123",
    "hospital": "hospital-a"
  }'
```

---

### 2. Staff Login
Authenticates a staff member and returns an access token.

**Endpoint**: `POST /staff/login`

**Query Parameters**:
- `hospital` (required): Hospital identifier

**Request Body**:
```json
{
    "username": "string (required)",
    "password": "string (required)"
}
```

**Response**:
- **200 OK**:
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

- **400 Bad Request**:
```json
{
    "error": "hospital parameter is required"
}
```

- **401 Unauthorized**:
```json
{
    "error": "invalid credentials"
}
```

**Example**:
```bash
curl -X POST "https://localhost:443/staff/login?hospital=hospital-a" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "doctor001",
    "password": "securePassword123"
  }'
```

---

## Patient Search API

### 3. Search Patients
Searches for patients based on provided criteria. Staff can only search for patients in their assigned hospital.

**Endpoint**: `POST /patient/search`

**Headers**:
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body** (all fields are optional):
```json
{
    "national_id": "string",
    "passport_id": "string", 
    "first_name": "string",
    "middle_name": "string",
    "last_name": "string",
    "date_of_birth": "string (YYYY-MM-DD)",
    "phone_number": "string",
    "email": "string"
}
```

**Response**:
- **200 OK**:
```json
{
    "patients": [
        {
            "id": "uuid",
            "first_name_th": "สมชาย",
            "middle_name_th": "",
            "last_name_th": "ใจดี",
            "first_name_en": "Somchai",
            "middle_name_en": "",
            "last_name_en": "Jaidee",
            "date_of_birth": "1990-01-01T00:00:00Z",
            "patient_hn": "HN001234",
            "national_id": "1234567890123",
            "passport_id": "",
            "phone_number": "0812345678",
            "email": "somchai@example.com",
            "gender": "M"
        }
    ],
    "count": 1
}
```

- **400 Bad Request**:
```json
{
    "error": "Invalid request format"
}
```

- **401 Unauthorized**:
```json
{
    "error": "Authorization header required"
}
```

- **500 Internal Server Error**:
```json
{
    "error": "Internal server error message"
}
```

**Example**:
```bash
curl -X POST https://localhost:443/patient/search \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "national_id": "1234567890123"
  }'
```

---

## Health Check

### 4. Health Check
Returns the API health status.

**Endpoint**: `GET /`

**Response**:
- **200 OK**:
```json
{
    "message": "Agnos Hospital Middleware API",
    "status": "healthy"
}
```

**Example**:
```bash
curl https://localhost:443/
```

---

## Error Handling

### HTTP Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `500 Internal Server Error`: Server error

### Error Response Format
All error responses follow this format:
```json
{
    "error": "Human-readable error message"
}
```

---

## Data Integration

### Hospital A API Integration
The system automatically integrates with Hospital A's external API when searching for patients:

- **External API**: `GET https://hospital-a.api.co.th/patient/search/{id}`
- **Trigger**: When searching by `national_id` or `passport_id` with no local results
- **Caching**: Retrieved patient data is stored locally for future searches
- **Fallback**: Local database search if external API is unavailable

### Patient Data Flow
1. Search request received from staff
2. Query local database first
3. If no results and searching by ID, query Hospital A API
4. Cache external API results locally
5. Return combined results to staff

---

## Security Considerations

### Authentication
- JWT tokens expire after 24 hours
- Tokens include staff ID, username, and hospital information
- All patient search endpoints require valid authentication

### Authorization
- Staff can only search for patients in their assigned hospital
- Hospital isolation is enforced at the service layer
- Cross-hospital data access is prevented

### Data Protection
- Passwords are hashed using bcrypt
- HTTPS/TLS encryption for all communications
- Input validation and sanitization
- SQL injection prevention through parameterized queries

---

## Rate Limiting
Currently no rate limiting is implemented. Consider adding rate limiting for production deployment.

## API Versioning
Current API version: v1 (implicit)
Future versions should include version in the URL path: `/v2/staff/create`

---

## Testing

### Sample Test Data
```json
{
    "staff": {
        "username": "testdoctor",
        "password": "testpass123",
        "hospital": "hospital-a"
    },
    "patient": {
        "national_id": "1234567890123",
        "first_name": "John",
        "last_name": "Doe",
        "phone_number": "0812345678"
    }
}
```

### Test Scenarios
1. **Staff Creation**: Valid and invalid staff data
2. **Authentication**: Valid/invalid credentials, missing hospital
3. **Patient Search**: Various search criteria combinations
4. **Authorization**: Cross-hospital access attempts
5. **Error Handling**: Invalid JSON, missing headers, server errors
