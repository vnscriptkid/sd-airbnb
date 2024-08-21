-- seed 3 users: host and guest1, guest2
INSERT INTO users (id, name, email, role, created_at, updated_at) VALUES
(1, 'Host User', 'host@gmail.com', 'Host', NOW(), NOW()),
(2, 'Guest User 1', 'guest1@gmail.com', 'Guest', NOW(), NOW()),
(3, 'Guest User 2', 'guest2@gmail.com', 'Guest', NOW(), NOW());

-- seed 1 hotel by host
INSERT INTO hotels (id, host_id, name, description, address, created_at, updated_at) VALUES
(1, 1, 'Hotel California', 'Such a lovely place', '1234 California St, San Francisco, CA', NOW(), NOW());

-- seed 2 rooms for the hotel
INSERT INTO rooms (id, hotel_id, name, description, type, quantity, price_per_night, created_at, updated_at) VALUES
(1, 1, 'Room 101', 'Single room', 'Single', 2, 100.00, NOW(), NOW()),
(2, 1, 'Room 102', 'Double room', 'Double', 1, 150.00, NOW(), NOW());

-- seed room availability for the next 5 days: date should be in the format 'YYYY-MM-DD'
INSERT INTO room_availabilities (room_id, date, rooms_left, created_at, updated_at) VALUES
(1, '2024-09-01', 2, NOW(), NOW()),
(1, '2024-09-02', 2, NOW(), NOW()),
(1, '2024-09-03', 2, NOW(), NOW()),
(1, '2024-09-04', 2, NOW(), NOW()),
(1, '2024-09-05', 2, NOW(), NOW()),
(2, '2024-09-01', 1, NOW(), NOW()),
(2, '2024-09-02', 1, NOW(), NOW()),
(2, '2024-09-03', 1, NOW(), NOW()),
(2, '2024-09-04', 1, NOW(), NOW()),
(2, '2024-09-05', 1, NOW(), NOW());

-- get all indexes in the database
SELECT * FROM pg_indexes WHERE schemaname = 'public';

-- get all constraints in the database
SELECT * FROM information_schema.table_constraints WHERE table_schema = 'public';

-- try update room availability to negative
UPDATE room_availabilities SET rooms_left = -1 WHERE room_id = 1 AND date = '2024-09-01';

