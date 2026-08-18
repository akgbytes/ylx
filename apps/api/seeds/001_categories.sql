INSERT INTO categories (id, name, slug) VALUES
('a1000000-0000-4000-8000-000000000001', 'Laptops',           'laptops'),
('a1000000-0000-4000-8000-000000000002', 'Monitors',          'monitors'),
('a1000000-0000-4000-8000-000000000003', 'Peripherals',       'peripherals'),
('a1000000-0000-4000-8000-000000000004', 'Audio',             'audio'),
('a1000000-0000-4000-8000-000000000005', 'Furniture',         'furniture'),
('a1000000-0000-4000-8000-000000000006', 'PC Components',     'pc-components'),
('a1000000-0000-4000-8000-000000000007', 'DIY & Electronics', 'diy-electronics'),
('a1000000-0000-4000-8000-000000000008', 'Gaming & VR',       'gaming-vr'),
('a1000000-0000-4000-8000-000000000009', 'Cameras',           'cameras')
ON CONFLICT (id) DO NOTHING;
