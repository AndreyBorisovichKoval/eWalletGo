-- Вставка данных в таблицу phones
INSERT INTO phones (phone_number, created_at, updated_at) VALUES
('1234567890', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('0987654321', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('5555555555', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('1112223333', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('4445556666', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Вставка данных в таблицу users (идентифицированные)
INSERT INTO users (full_name, username, email, passport_number, password, password_reset_required, role, status, created_at, updated_at) VALUES
('Иван Иванов', 'ivan', 'ivan@example.com', '1234 567890', 'password_hash', true, 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Мария Петрова', 'maria', 'maria@example.com', '2345 678901', 'password_hash', true, 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Петр Сидоров', 'petr', 'petr@example.com', '3456 789012', 'password_hash', false, 'manager', 'blocked', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Анна Козлова', 'anna', 'anna@example.com', '4567 890123', 'password_hash', true, 'user', 'canceled', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Олег Титов', 'oleg', 'oleg@example.com', '5678 901234', 'password_hash', true, 'user', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Вставка данных в таблицу unverified_users (неидентифицированные пользователи)
INSERT INTO unverified_users (phone_id, full_name, email, status, created_at, updated_at) VALUES
(3, 'Алексей Смирнов', 'alex@example.com', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'Ксения Лебедева', 'ksenia@example.com', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(5, 'Николай Воробьев', 'nikolai@example.com', 'canceled', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Вставка данных в таблицу wallets
INSERT INTO wallets (phone_id, user_id, client_type, status, wallet_number, masked_wallet_number, created_at, updated_at) VALUES
(1, 1, 'identified', 'active', '1234567890123456', '1234 **** **** 3456', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 2, 'identified', 'active', '6543210987654321', '6543 **** **** 4321', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, NULL, 'unidentified', 'active', '9876543210987654', '9876 **** **** 7654', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, NULL, 'unidentified', 'blocked', '1234098765432109', '1234 **** **** 2109', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(5, 3, 'identified', 'blocked', '8765432109876543', '8765 **** **** 6543', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Вставка данных в таблицу accounts
INSERT INTO accounts (wallet_id, user_id, balance, created_at, updated_at) VALUES
(1, 1, 1000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 2, 500.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, NULL, 750.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, NULL, 200.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(5, 3, 0.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Вставка данных в таблицу transactions
INSERT INTO transactions (account_id, amount, type, created_at) VALUES
(1, 200.00, 'recharge', CURRENT_TIMESTAMP),
(1, -50.00, 'withdrawal', CURRENT_TIMESTAMP),
(1, 150.00, 'recharge', CURRENT_TIMESTAMP),
(2, 300.00, 'recharge', CURRENT_TIMESTAMP),
(2, -100.00, 'withdrawal', CURRENT_TIMESTAMP),
(3, 150.00, 'recharge', CURRENT_TIMESTAMP),
(4, 50.00, 'recharge', CURRENT_TIMESTAMP),
(5, -200.00, 'withdrawal', CURRENT_TIMESTAMP),
(5, 400.00, 'recharge', CURRENT_TIMESTAMP);

-- Вставка данных в таблицу user_settings
INSERT INTO user_settings (user_id, add_confirmation, update_confirmation, delete_confirmation, display_language, desktop_theme, dark_mode_theme, font, font_size, accessibility_options, notification_sound, email_notifications, notification_frequency, created_at, updated_at) VALUES
(1, true, true, true, 'Russian', 'Green animation', false, 'Helvetica', 11, 'High contrast', true, false, 'daily', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, true, true, true, 'English', 'Blue calm', true, 'Arial', 12, 'None', true, true, 'weekly', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, false, false, false, 'Russian', 'Dark mode', false, 'Times New Roman', 13, 'None', false, false, 'monthly', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, true, false, true, 'Russian', 'Default', true, 'Courier', 10, 'High contrast', true, true, 'never', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(5, true, true, false, 'English', 'Colorful', false, 'Verdana', 14, 'None', false, false, 'daily', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Вставка данных в таблицу limit_settings
INSERT INTO limit_settings (client_type, default_limit, custom_limit, agreement_date, created_at, updated_at) VALUES
('identified', 100000.00, NULL, NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('unidentified', 10000.00, NULL, NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
