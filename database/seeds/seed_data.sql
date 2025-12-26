-- Seed data for development and testing

-- Insert demo organization
INSERT INTO organizations (name, description) VALUES 
('Demo Company', 'A demo organization for testing the task management system');

-- Insert demo users (password: password123 - hashed with bcrypt)
INSERT INTO users (organization_id, email, password_hash, full_name, timezone) VALUES 
(1, 'admin@demo.com', '$2a$10$9ut3U/drv7cPiRAmZXQjPe0HXbv7iLoCzDy2gi5BwvoV6LWBzkEgi', 'Admin User', 'Asia/Jakarta'),
(1, 'manager@demo.com', '$2a$10$9ut3U/drv7cPiRAmZXQjPe0HXbv7iLoCzDy2gi5BwvoV6LWBzkEgi', 'Manager User', 'Asia/Jakarta'),
(1, 'developer1@demo.com', '$2a$10$9ut3U/drv7cPiRAmZXQjPe0HXbv7iLoCzDy2gi5BwvoV6LWBzkEgi', 'Developer One', 'Asia/Jakarta'),
(1, 'developer2@demo.com', '$2a$10$9ut3U/drv7cPiRAmZXQjPe0HXbv7iLoCzDy2gi5BwvoV6LWBzkEgi', 'Developer Two', 'Asia/Jakarta');

-- Insert demo teams
INSERT INTO teams (organization_id, name, description) VALUES 
(1, 'Engineering', 'Software development team'),
(1, 'Frontend Team', 'Frontend development sub-team');

-- Make Frontend Team a sub-team of Engineering
UPDATE teams SET parent_team_id = 1 WHERE id = 2;

-- Insert team members with roles
INSERT INTO team_members (team_id, user_id, role) VALUES 
(1, 1, 'manager'),      -- Admin is manager of Engineering
(1, 2, 'assistant'),    -- Manager is assistant
(1, 3, 'member'),       -- Developer 1 is member
(1, 4, 'member'),       -- Developer 2 is member
(2, 2, 'manager'),      -- Manager leads Frontend Team
(2, 3, 'member'),       -- Developer 1 in Frontend
(2, 4, 'member');       -- Developer 2 in Frontend

-- Insert workflow statuses for Organization (shared by all teams)
INSERT INTO issue_statuses (organization_id, name, position, is_final, color) VALUES 
(1, 'WAITING', 1, false, '#9CA3AF'),
(1, 'IN_PROGRESS', 2, false, '#3B82F6'),
(1, 'QA', 3, false, '#F59E0B'),
(1, 'READY_TO_DEPLOY', 4, false, '#8B5CF6'),
(1, 'DONE', 5, true, '#10B981'),
(1, 'HOLD', 6, false, '#EF4444');

-- Insert demo issues
INSERT INTO issues (team_id, status_id, title, description, priority, created_by) VALUES 
(1, 1, 'Setup CI/CD Pipeline', 'Configure automated deployment pipeline', 'HIGH', 1),
(1, 2, 'Implement Authentication', 'Add JWT-based authentication system', 'URGENT', 1),
(1, 3, 'Fix Database Performance', 'Optimize slow queries in reports', 'NORMAL', 2),
(2, 1, 'Design Landing Page', 'Create modern landing page design', 'HIGH', 2),
(2, 2, 'Build Dashboard UI', 'Implement main dashboard interface', 'NORMAL', 2);

-- Insert assignments
INSERT INTO issue_assignments (issue_id, user_id, start_date, end_date, assigned_by, is_active) VALUES 
(1, 3, CURRENT_DATE, CURRENT_DATE + INTERVAL '5 days', 1, true),
(2, 3, CURRENT_DATE, CURRENT_DATE + INTERVAL '7 days', 1, true),
(3, 4, CURRENT_DATE - INTERVAL '2 days', CURRENT_DATE + INTERVAL '3 days', 2, true),
(4, 3, CURRENT_DATE + INTERVAL '1 day', CURRENT_DATE + INTERVAL '4 days', 2, true),
(5, 4, CURRENT_DATE + INTERVAL '2 days', CURRENT_DATE + INTERVAL '8 days', 2, true);

-- Insert some work logs
INSERT INTO issue_work_logs (issue_id, user_id, work_date, minutes_spent, notes) VALUES 
(2, 3, CURRENT_DATE, 240, 'Implemented JWT token generation and validation'),
(3, 4, CURRENT_DATE - INTERVAL '1 day', 180, 'Analyzed slow queries and added indexes');

-- Insert activity logs
INSERT INTO issue_activities (issue_id, user_id, activity_type, description) VALUES 
(1, 1, 'created', 'Issue created'),
(1, 1, 'assigned', 'Assigned to Developer One'),
(2, 1, 'created', 'Issue created'),
(2, 1, 'assigned', 'Assigned to Developer One'),
(2, 3, 'status_changed', 'Moved to IN_PROGRESS'),
(3, 2, 'created', 'Issue created'),
(3, 2, 'assigned', 'Assigned to Developer Two'),
(3, 4, 'status_changed', 'Moved to QA');
