INSERT INTO listings (id, title, description, price, city, seller_id, category_id, status) VALUES
-- Laptops
('c1000000-0000-4000-8000-000000000001', 'MacBook Pro M3 14-inch', '16GB RAM, 512GB SSD. Purchased this year with original box and charger.', 14800000, 'Bangalore', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000001', 'active'),
('c1000000-0000-4000-8000-000000000002', 'MacBook Air M2', '8GB RAM, 256GB SSD. Excellent battery health, barely used.', 7600000, 'Pune', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000001', 'active'),
('c1000000-0000-4000-8000-000000000003', 'ThinkPad X1 Carbon Gen 11', 'Intel Core Ultra 7, 32GB RAM, ideal for development work.', 12800000, 'Hyderabad', 'b1000000-0000-4000-8000-000000000003', 'a1000000-0000-4000-8000-000000000001', 'active'),

-- Monitors
('c1000000-0000-4000-8000-000000000004', 'Dell UltraSharp 27-inch Monitor', '4K IPS monitor in excellent condition with original stand.', 3400000, 'Delhi', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000002', 'active'),
('c1000000-0000-4000-8000-000000000005', 'LG UltraWide 34-inch Monitor', '3440x1440 ultrawide monitor for productivity.', 5200000, 'Bangalore', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000002', 'active'),

-- Peripherals
('c1000000-0000-4000-8000-000000000006', 'Keychron K2 Mechanical Keyboard', 'Hot-swappable, Gateron Brown switches.', 650000, 'Gurgaon', 'b1000000-0000-4000-8000-000000000003', 'a1000000-0000-4000-8000-000000000003', 'active'),
('c1000000-0000-4000-8000-000000000007', 'Keychron Q1 Keyboard', 'Aluminum body with tactile switches.', 1250000, 'Mumbai', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000003', 'active'),
('c1000000-0000-4000-8000-000000000008', 'Logitech MX Keys S', 'Wireless keyboard with USB receiver.', 780000, 'Noida', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000003', 'active'),
('c1000000-0000-4000-8000-000000000009', 'Logitech MX Master 3S', 'Excellent condition with USB receiver included.', 650000, 'Delhi', 'b1000000-0000-4000-8000-000000000003', 'a1000000-0000-4000-8000-000000000003', 'active'),
('c1000000-0000-4000-8000-000000000010', 'Apple Magic Trackpad 2', 'Original Apple accessory in mint condition.', 850000, 'Chennai', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000003', 'active'),
('c1000000-0000-4000-8000-000000000011', 'BenQ ScreenBar Monitor Light', 'USB powered monitor light for desk setup.', 650000, 'Pune', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000003', 'active'),

-- Audio
('c1000000-0000-4000-8000-000000000012', 'Sony WH-1000XM5 Headphones', 'Noise cancelling headphones with carrying case.', 2100000, 'Mumbai', 'b1000000-0000-4000-8000-000000000003', 'a1000000-0000-4000-8000-000000000004', 'active'),
('c1000000-0000-4000-8000-000000000013', 'AirPods Pro (2nd Gen)', 'USB-C charging case, original accessories included.', 1650000, 'Delhi', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000004', 'active'),

-- Furniture
('c1000000-0000-4000-8000-000000000014', 'Herman Miller Aeron Chair', 'Size B ergonomic chair with adjustable lumbar support.', 7200000, 'Bangalore', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000005', 'active'),
('c1000000-0000-4000-8000-000000000015', 'IKEA Bekant Standing Desk', 'Height adjustable desk, works perfectly.', 2800000, 'Hyderabad', 'b1000000-0000-4000-8000-000000000003', 'a1000000-0000-4000-8000-000000000005', 'active'),
('c1000000-0000-4000-8000-000000000016', 'Ergonomic Office Chair', 'Mesh back with adjustable armrests.', 550000, 'Gurgaon', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000005', 'active'),
('c1000000-0000-4000-8000-000000000017', 'Solid Wood Study Desk', '120cm desk with cable management tray.', 850000, 'Jaipur', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000005', 'active'),

-- PC Components
('c1000000-0000-4000-8000-000000000018', 'NVIDIA RTX 4070 Graphics Card', 'Used for light gaming and AI experiments.', 5200000, 'Bangalore', 'b1000000-0000-4000-8000-000000000003', 'a1000000-0000-4000-8000-000000000006', 'active'),
('c1000000-0000-4000-8000-000000000019', 'Ryzen 7 7800X3D Processor', 'Excellent condition with original packaging.', 2800000, 'Delhi', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000006', 'active'),
('c1000000-0000-4000-8000-000000000020', 'Synology DS224+ NAS', '2-bay NAS with 8TB storage installed.', 3400000, 'Chandigarh', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000006', 'active'),

-- DIY & Electronics
('c1000000-0000-4000-8000-000000000021', 'Raspberry Pi 5 8GB Kit', 'Includes official power adapter and case.', 950000, 'Ahmedabad', 'b1000000-0000-4000-8000-000000000003', 'a1000000-0000-4000-8000-000000000007', 'active'),
('c1000000-0000-4000-8000-000000000022', 'Arduino Uno Starter Kit', 'Complete beginner kit with sensors and components.', 250000, 'Lucknow', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000007', 'active'),
('c1000000-0000-4000-8000-000000000023', 'ESP32 Development Board Pack', 'Set of five ESP32 boards for IoT projects.', 320000, 'Indore', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000007', 'active'),

-- Gaming & VR
('c1000000-0000-4000-8000-000000000024', 'Steam Deck OLED 512GB', 'Barely used, includes carrying case and charger.', 4700000, 'Pune', 'b1000000-0000-4000-8000-000000000003', 'a1000000-0000-4000-8000-000000000008', 'active'),
('c1000000-0000-4000-8000-000000000025', 'Meta Quest 3', '128GB VR headset with controllers.', 4200000, 'Hyderabad', 'b1000000-0000-4000-8000-000000000001', 'a1000000-0000-4000-8000-000000000008', 'active'),

-- Cameras
('c1000000-0000-4000-8000-000000000026', 'GoPro Hero 13 Black', 'Action camera with spare battery and carrying case.', 3600000, 'Goa', 'b1000000-0000-4000-8000-000000000002', 'a1000000-0000-4000-8000-000000000009', 'active')
ON CONFLICT (id) DO NOTHING;
