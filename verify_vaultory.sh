#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
EMAIL_SUFFIX=$(date +%s)
MANAGER_EMAIL="manager_${EMAIL_SUFFIX}@vaultory.com"
SUPERVISOR_EMAIL="supervisor_${EMAIL_SUFFIX}@vaultory.com"
STAFF_EMAIL="staff_${EMAIL_SUFFIX}@vaultory.com"
PASSWORD="SecurePass123!"

echo "🚀 Starting Vaultory Verification..."

# 1. Register Manager
echo -e "\n1️⃣  Registering Manager..."
REGISTER_RES=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"full_name\": \"Vaultory Manager\",
    \"email\": \"${MANAGER_EMAIL}\",
    \"password\": \"${PASSWORD}\",
    \"company_name\": \"Vaultory Inc\"
  }")
echo "Response: $REGISTER_RES"

# 2. Login Manager
echo -e "\n2️⃣  Logging in Manager..."
LOGIN_RES=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${MANAGER_EMAIL}\",
    \"password\": \"${PASSWORD}\"
  }")
MANAGER_TOKEN=$(echo $LOGIN_RES | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
echo "Manager Token: ${MANAGER_TOKEN:0:20}..."

# 3. Create Warehouse (Manager)
echo -e "\n3️⃣  Creating Warehouse (Manager)..."
WAREHOUSE_RES=$(curl -s -X POST "${BASE_URL}/manager/warehouse/create" \
  -H "Authorization: Bearer $MANAGER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alpha Warehouse",
    "location": "Sector 7"
  }')
echo "Response: $WAREHOUSE_RES"
WAREHOUSE_ID=$(echo $WAREHOUSE_RES | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Warehouse ID: $WAREHOUSE_ID"

# 4. Create Supervisor (Manager)
echo -e "\n4️⃣  Creating Supervisor (Manager)..."
SUPERVISOR_RES=$(curl -s -X POST "${BASE_URL}/manager/supervisor/create" \
  -H "Authorization: Bearer $MANAGER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"full_name\": \"Alpha Supervisor\",
    \"email\": \"${SUPERVISOR_EMAIL}\",
    \"password\": \"${PASSWORD}\",
    \"warehouse_id\": \"${WAREHOUSE_ID}\"
  }")
echo "Response: $SUPERVISOR_RES"

# 5. Login Supervisor
echo -e "\n5️⃣  Logging in Supervisor..."
SUP_LOGIN_RES=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${SUPERVISOR_EMAIL}\",
    \"password\": \"${PASSWORD}\"
  }")
SUPERVISOR_TOKEN=$(echo $SUP_LOGIN_RES | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
echo "Supervisor Token: ${SUPERVISOR_TOKEN:0:20}..."

# 6. Create Staff (Supervisor)
echo -e "\n6️⃣  Creating Staff (Supervisor)..."
STAFF_RES=$(curl -s -X POST "${BASE_URL}/supervisor/staff/create" \
  -H "Authorization: Bearer $SUPERVISOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"full_name\": \"Alpha Staff\",
    \"email\": \"${STAFF_EMAIL}\",
    \"password\": \"${PASSWORD}\",
    \"warehouse_id\": \"${WAREHOUSE_ID}\"
  }")
echo "Response: $STAFF_RES"

# 7. Login Staff
echo -e "\n7️⃣  Logging in Staff..."
STAFF_LOGIN_RES=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${STAFF_EMAIL}\",
    \"password\": \"${PASSWORD}\"
  }")
STAFF_TOKEN=$(echo $STAFF_LOGIN_RES | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
echo "Staff Token: ${STAFF_TOKEN:0:20}..."

# 8. Create Item (Staff)
echo -e "\n8️⃣  Creating Item (Staff)..."
ITEM_RES=$(curl -s -X POST "${BASE_URL}/staff/item/add" \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"sku\": \"ITEM-${EMAIL_SUFFIX}\",
    \"name\": \"Secure Box\",
    \"quality\": \"New\",
    \"price\": 99.99,
    \"warehouse_id\": \"${WAREHOUSE_ID}\",
    \"quantity\": 50
  }")
echo "Response: $ITEM_RES"

# 9. List Audit Logs (Manager)
echo -e "\n9️⃣  Listing Audit Logs (Manager)..."
AUDIT_RES=$(curl -s -X GET "${BASE_URL}/manager/audit-logs" \
  -H "Authorization: Bearer $MANAGER_TOKEN")
# Check if response contains "CREATE" action (from previous steps)
if [[ $AUDIT_RES == *"CREATE"* ]]; then
  echo "✅ Audit Logs retrieved successfully"
else
  echo "❌ Failed to retrieve audit logs or no logs found"
  echo "Response: ${AUDIT_RES:0:100}..."
fi

echo -e "\n✅ Verification Complete!"
