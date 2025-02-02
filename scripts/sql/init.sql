-- uuid-ossp é uma extensão do PostgreSQL que permite a geração de UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    opt_out BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cptec_id INTEGER UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    state VARCHAR(2) NOT NULL,
    is_coastal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notification_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    location_id UUID NOT NULL REFERENCES locations(id),
    frequency VARCHAR(50) NOT NULL,
    time_of_day TIME NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    location_id UUID NOT NULL REFERENCES locations(id),
    content JSONB NOT NULL,
    status VARCHAR(50) NOT NULL,
    scheduled_for TIMESTAMP WITH TIME ZONE NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_scheduled_for ON notifications(scheduled_for);
CREATE INDEX idx_notification_schedules_user_location ON notification_schedules(user_id, location_id);

INSERT INTO users (name, email) VALUES
    ('Matehus', 'matheus@teste.com'),
    ('Joao', 'joao@litoral.com');

INSERT INTO locations (cptec_id, name, state, is_coastal) VALUES
    (456, 'Rio de Janeiro', 'RJ', FALSE);
    (123, 'São Paulo', 'SP', TRUE),