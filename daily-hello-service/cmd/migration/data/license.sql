-- delete from roles;
-- alter sequence roles_id_seq restart with 1;
INSERT INTO roles (
    name, code, description 
) VALUES 
  ( 'HDB Creation', 'HDB_CREATION', 'Creation Role for HDB'),
  ( 'HDB Approval', 'HDB_APPROVAL', 'Approval Role for HDB'),
  ( 'Agent Creation', 'AG_CREATION', 'Creation Role for Agent'),
  ( 'Agent Approval', 'AG_APPROVAL', 'Approval Role for Agent');

INSERT INTO permissions (
    name, code, description 
) VALUES 
  ( 'Admin View', 'ADMIN_VIEW', 'Allows viewing in Admin Portal'),
  ( 'Admin Create', 'ADMIN_CREATE', 'Allows creating in Admin Portal'),
  ( 'Admin Update', 'ADMIN_UPDATE', 'Allows updating in Admin Portal'),
  ( 'Admin Approve', 'ADMIN_APPROVE', 'Allows approving in Admin Portal'),
  ( 'Agent View', 'AGENT_VIEW', 'Allows viewing in Agent Portal'),
  ( 'Agent Create', 'AGENT_CREATE', 'Allows creating in Agent Portal'),
  ( 'Agent Update', 'AGENT_UPDATE', 'Allows updating in Agent Portal'),
  ( 'Agent Approve', 'AGENT_APPROVE', 'Allows approving in Agent Portal');

-- FOR ADMIN --
INSERT INTO action_details (
    type, code, path, http_method 
) VALUES 
  ('API', 'ADMIN_COMPANY_CREATE', '/api/v1/admin/company/create', 'POST'),
  ('API', 'ADMIN_COMPANY_VIEW', '/api/v1/admin/company/:id', 'GET'),
  ('API', 'ADMIN_COMPANY_UPDATE', '/api/v1/admin/company/:id', 'PATCH'),
  ('API', 'ADMIN_COMPANY_LIST', '/api/v1/admin/company/list', 'GET'),
  ('API', 'ADMIN_COMPANY_SYNC', '/api/v1/admin/company/sync', 'POST'),
  ('API', 'ADMIN_COMPANY_FEES', '/api/v1/admin/company/:id/fees', 'GET'),
  ('API', 'ADMIN_FEE_LIST_TYPES', '/api/v1/admin/fee/list-types', 'GET'),
  ('API', 'ADMIN_FEE_CREATE', '/api/v1/admin/fee/create', 'POST'),
  ('API', 'ADMIN_FEE_VIEW', '/api/v1/admin/fee/:id', 'GET'),
  ('API', 'ADMIN_FEE_UPDATE', '/api/v1/admin/fee/:id', 'PATCH'),
  ('API', 'ADMIN_FEE_LIST', '/api/v1/admin/fee/list', 'GET'),
  ('API', 'ADMIN_FEE_UPSERT_BATCH', '/api/v1/admin/fee/upsert-batch', 'POST'),
  ('API', 'ADMIN_LIMIT_CREATE', '/api/v1/admin/limit/create', 'POST'),
  ('API', 'ADMIN_LIMIT_CREATE_BATCH', '/api/v1/admin/limit/create-batch', 'POST'),
  ('API', 'ADMIN_LIMIT_VIEW', '/api/v1/admin/limit/:id', 'GET'),
  ('API', 'ADMIN_LIMIT_UPDATE', '/api/v1/admin/limit/:id', 'PATCH'),
  ('API', 'ADMIN_LIMIT_LIST', '/api/v1/admin/limit/list', 'GET'),
  ('API', 'ADMIN_LOCATION_GET_LIST', '/api/v1/admin/location/get-list', 'POST'),
  ('API', 'ADMIN_PROFILE_ME', '/api/v1/admin/profile/me', 'GET'),
  ('API', 'ADMIN_STAFF_CREATE', '/api/v1/admin/staff/create', 'POST'),
  ('API', 'ADMIN_STAFF_VIEW', '/api/v1/admin/staff/:id', 'GET'),
  ('API', 'ADMIN_STAFF_UPDATE', '/api/v1/admin/staff/:id', 'PATCH'),
  ('API', 'ADMIN_STAFF_UPDATE_STATUS', '/api/v1/admin/staff/:id/status', 'PATCH'),
  ('API', 'ADMIN_STAFF_LIST', '/api/v1/admin/staff/list', 'GET'),
  ('API', 'ADMIN_STORE_CREATE', '/api/v1/admin/store/create', 'POST'),
  ('API', 'ADMIN_STORE_VIEW', '/api/v1/admin/store/:id', 'GET'),
  ('API', 'ADMIN_STORE_UPDATE', '/api/v1/admin/store/:id', 'PATCH'),
  ('API', 'ADMIN_STORE_LIST', '/api/v1/admin/store/list', 'GET'),
  ('API', 'ADMIN_STORE_APPROVE_STORES', '/api/v1/admin/store/approve-stores', 'POST'),
  ('API', 'ADMIN_STORE_EXPORT_DATA', '/api/v1/admin/store/export-data', 'POST'),
  ('API', 'ADMIN_TRANSACTION_CREATE', '/api/v1/admin/transaction/create', 'POST'),
  ('API', 'ADMIN_TRANSACTION_VIEW', '/api/v1/admin/transaction/:id', 'GET'),
  ('API', 'ADMIN_TRANSACTION_UPDATE', '/api/v1/admin/transaction/:id', 'PATCH'),
  ('API', 'ADMIN_TRANSACTION_LIST', '/api/v1/admin/transaction/list', 'POST'),
  ('API', 'ADMIN_TRANSACTION_APPROVE_TRANSACTIONS', '/api/v1/admin/transaction/approve-transactions', 'POST'),
  ('API', 'ADMIN_TRANSACTION_LIST_TYPES', '/api/v1/admin/transaction/list-types', 'GET');

INSERT INTO permission_has_actions (
    permission_code, action_detail_code, description, is_active
) VALUES 
  ('ADMIN_VIEW', 'ADMIN_COMPANY_VIEW', 'Allows viewing in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_COMPANY_LIST', 'Allows viewing in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_COMPANY_FEES', 'Allows viewing in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_COMPANY_CREATE', 'Allows creating in Admin Portal', true),
  ('ADMIN_UPDATE', 'ADMIN_COMPANY_UPDATE', 'Allows updating in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_COMPANY_SYNC', 'Allows syncing in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_FEE_LIST_TYPES', 'Allows viewing in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_FEE_CREATE', 'Allows creating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_FEE_VIEW', 'Allows viewing in Admin Portal', true),
  ('ADMIN_UPDATE', 'ADMIN_FEE_UPDATE', 'Allows updating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_FEE_LIST', 'Allows viewing in Admin Portal', true),
  ('ADMIN_UPDATE', 'ADMIN_FEE_UPSERT_BATCH', 'Allows upserting in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_LIMIT_CREATE', 'Allows creating in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_LIMIT_CREATE_BATCH', 'Allows creating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_LIMIT_VIEW', 'Allows viewing in Admin Portal', true),
  ('ADMIN_UPDATE', 'ADMIN_LIMIT_UPDATE', 'Allows updating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_LIMIT_LIST', 'Allows viewing in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_LOCATION_GET_LIST', 'Allows viewing in Admin Portal', true),          
  ('ADMIN_VIEW', 'ADMIN_PROFILE_ME', 'Allows viewing in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_STAFF_CREATE', 'Allows creating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_STAFF_VIEW', 'Allows viewing in Admin Portal', true),
  ('ADMIN_UPDATE', 'ADMIN_STAFF_UPDATE', 'Allows updating in Admin Portal', true),
  ('ADMIN_UPDATE', 'ADMIN_STAFF_UPDATE_STATUS', 'Allows updating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_STAFF_LIST', 'Allows viewing in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_STORE_CREATE', 'Allows creating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_STORE_VIEW', 'Allows viewing in Admin Portal', true),
  ('ADMIN_UPDATE', 'ADMIN_STORE_UPDATE', 'Allows updating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_STORE_LIST', 'Allows viewing in Admin Portal', true),
  ('ADMIN_APPROVE', 'ADMIN_STORE_APPROVE_STORES', 'Allows approving in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_STORE_EXPORT_DATA', 'Allows exporting in Admin Portal', true),
  ('ADMIN_CREATE', 'ADMIN_TRANSACTION_CREATE', 'Allows creating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_TRANSACTION_VIEW', 'Allows viewing in Admin Portal', true),
  ('ADMIN_UPDATE', 'ADMIN_TRANSACTION_UPDATE', 'Allows updating in Admin Portal', true),
  ('ADMIN_VIEW', 'ADMIN_TRANSACTION_LIST', 'Allows viewing in Admin Portal', true),     
  ('ADMIN_APPROVE', 'ADMIN_TRANSACTION_APPROVE_TRANSACTIONS', 'Allows approving in Admin Portal', true),  
  ('ADMIN_VIEW', 'ADMIN_TRANSACTION_LIST_TYPES', 'Allows viewing in Admin Portal', true);

INSERT INTO role_has_permissions (
permission_code, role_code, is_active
) VALUES 
  ('ADMIN_VIEW', 'HDB_CREATION', true),
  ('ADMIN_CREATE', 'HDB_CREATION', true),
  ('ADMIN_UPDATE', 'HDB_CREATION', true),
  ('ADMIN_VIEW', 'HDB_APPROVAL', true),
  ('ADMIN_APPROVE', 'HDB_APPROVAL', true);


-- FOR AGENT --
INSERT INTO action_details (
    type, code, path, http_method 
) VALUES 
  ('API', 'AGENT_COMPANY_VIEW','/api/v1/agent/company', 'GET'),
  ('API', 'AGENT_COMPANY_LIST_FEES','/api/v1/agent/company/fees', 'GET'),
  ('API', 'AGENT_LIMIT_CREATE','/api/v1/agent/limit/create', 'POST'),
  ('API', 'AGENT_LIMIT_CREATE_BATCH','/api/v1/agent/limit/create-batch', 'POST'),
  ('API', 'AGENT_LIMIT_VIEW','/api/v1/agent/limit/:id', 'GET'),
  ('API', 'AGENT_LIMIT_UPDATE','/api/v1/agent/limit/:id', 'PATCH'),
  ('API', 'AGENT_LIMIT_LIST','/api/v1/agent/limit/list', 'GET'),
  ('API', 'AGENT_LOCATION_LIST','/api/v1/agent/location/get-list', 'POST'),
  ('API', 'AGENT_PROFILE_VIEW','/api/v1/agent/profile/me', 'GET'),
  ('API', 'AGENT_STAFF_CREATE','/api/v1/agent/staff/create', 'POST'),
  ('API', 'AGENT_STAFF_VIEW','/api/v1/agent/staff/:id', 'GET'),
  ('API', 'AGENT_STAFF_UPDATE','/api/v1/agent/staff/:id', 'PATCH'),
  ('API', 'AGENT_STAFF_UPDATE_STATUS','/api/v1/agent/staff/:id/status', 'PATCH'),
  ('API', 'AGENT_STAFF_LIST','/api/v1/agent/staff/list', 'GET'),
  ('API', 'AGENT_STORE_CREATE','/api/v1/agent/store/create', 'POST'),
  ('API', 'AGENT_STORE_VIEW','/api/v1/agent/store/:id', 'GET'),
  ('API', 'AGENT_STORE_UPDATE','/api/v1/agent/store/:id', 'PATCH'),
  ('API', 'AGENT_STORE_LIST','/api/v1/agent/store/list', 'GET'),
  ('API', 'AGENT_STORE_APPROVE','/api/v1/agent/store/approve-stores', 'POST'),
  ('API', 'AGENT_STORE_EXPORT_DATA','/api/v1/agent/store/export-data', 'POST'),
  ('API', 'AGENT_TRANSACTION_CREATE','/api/v1/agent/transaction/create', 'POST'),
  ('API', 'AGENT_TRANSACTION_VIEW','/api/v1/agent/transaction/:id', 'GET'),
  ('API', 'AGENT_TRANSACTION_UPDATE','/api/v1/agent/transaction/:id', 'PATCH'),
  ('API', 'AGENT_TRANSACTION_LIST','/api/v1/agent/transaction/list', 'POST'),
  ('API', 'AGENT_TRANSACTION_APPROVE','/api/v1/admin/transaction/approve-transactions', 'POST'),
  ('API', 'AGENT_TRANSACTION_LIST_TYPES','/api/v1/agent/transaction/list-types', 'GET');

INSERT INTO permission_has_actions (
    permission_code, action_detail_code, description, is_active
) VALUES 
  ('AGENT_VIEW', 'AGENT_COMPANY_VIEW', 'Allows creating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_COMPANY_LIST_FEES', 'Allows creating in Agent Portal', true),
  ('AGENT_CREATE', 'AGENT_LIMIT_CREATE', 'Allows creating in Agent Portal', true),
  ('AGENT_CREATE', 'AGENT_LIMIT_CREATE_BATCH', 'Allows creating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_LIMIT_VIEW', 'Allows viewing in Agent Portal', true),
  ('AGENT_UPDATE', 'AGENT_LIMIT_UPDATE', 'Allows updating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_LIMIT_LIST', 'Allows viewing in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_LOCATION_LIST', 'Allows viewing in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_PROFILE_VIEW', 'Allows viewing in Agent Portal', true),
  ('AGENT_CREATE', 'AGENT_STAFF_CREATE', 'Allows creating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_STAFF_VIEW', 'Allows viewing in Agent Portal', true),
  ('AGENT_UPDATE', 'AGENT_STAFF_UPDATE', 'Allows updating in Agent Portal', true),
  ('AGENT_UPDATE', 'AGENT_STAFF_UPDATE_STATUS', 'Allows updating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_STAFF_LIST', 'Allows viewing in Agent Portal', true),
  ('AGENT_CREATE', 'AGENT_STORE_CREATE', 'Allows creating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_STORE_VIEW', 'Allows viewing in Agent Portal', true),
  ('AGENT_UPDATE', 'AGENT_STORE_UPDATE', 'Allows updating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_STORE_LIST', 'Allows viewing in Agent Portal', true),
  ('AGENT_APPROVE', 'AGENT_STORE_APPROVE', 'Allows approving in Agent Portal', true),
  ('AGENT_CREATE', 'AGENT_STORE_EXPORT_DATA', 'Allows exporting in Agent Portal', true),
  ('AGENT_CREATE', 'AGENT_TRANSACTION_CREATE', 'Allows creating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_TRANSACTION_VIEW', 'Allows viewing in Agent Portal', true),
  ('AGENT_UPDATE', 'AGENT_TRANSACTION_UPDATE', 'Allows updating in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_TRANSACTION_LIST', 'Allows viewing in Agent Portal', true),
  ('AGENT_APPROVE', 'AGENT_TRANSACTION_APPROVE', 'Allows approving in Agent Portal', true),
  ('AGENT_VIEW', 'AGENT_TRANSACTION_LIST_TYPES', 'Allows viewing in Agent Portal', true);

INSERT INTO role_has_permissions (
permission_code, role_code, is_active
) VALUES 
  ('AGENT_VIEW', 'AG_CREATION', true),
  ('AGENT_CREATE', 'AG_CREATION', true),   
  ('AGENT_UPDATE', 'AG_CREATION', true),
  ('AGENT_VIEW', 'AG_APPROVAL', true),
  ('AGENT_APPROVE', 'AG_APPROVAL', true);  