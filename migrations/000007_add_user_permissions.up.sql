-- Insert sample permission data into the permissions table
INSERT INTO permissions (id, name, description)
VALUES 
    ('3106c5de-a659-4d50-bbd2-bb6547f83312', 'project.delete', 'permission to delete project'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83313', 'project.user.delete', 'permission to delete project user'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83314', 'project.user.get', 'permission to view project users'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83315', 'project.user.edit', 'permission to edit project users');

