-- Insert sample permission data into the permissions table
INSERT INTO permissions (id, name, description)
VALUES 
    ('3106c5de-a659-4d50-bbd2-bb6547f83301', 'project.create', 'permission to create project'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83302', 'project.edit', 'permission to edit project details'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83303', 'project.get', 'permission to get project details'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83304', 'tables.get', 'permission to get all tables'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83305', 'tableschema.get', 'permission to get table schema'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83306', 'tableschema.edit', 'permission to edit table schema'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83307', 'tableschema.delete', 'permission to delete table schema'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83308', 'table.row.insert', 'permission to add new data into the table'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83309', 'table.row.update', 'permission to update the data in the table'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83310', 'table.row.delete', 'permission to delete row from table'),
    ('3106c5de-a659-4d50-bbd2-bb6547f83311', 'table.row.get', 'permission to get data of a table');
