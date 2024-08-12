CREATE TABLE "user" (
  id serial PRIMARY KEY,
  full_name text NOT NULL,
  phone_number text NOT NULL,
  "password" text not null,
  balance decimal not null,
  created_time timestamp NOT NULL default now(),
  updated_time timestamp,
  successful_login_count int not null default 0,
  CONSTRAINT user_phone_number_uniquekey UNIQUE (phone_number),
  CONSTRAINT balance_non_negative CHECK (balance >= 0)
);

INSERT INTO "user" (full_name, phone_number, "password", balance) VALUES ('name1', '+6281122334455', '$2a$12$3gbfndmoRHh9k0qNlHL78e1tXEFceJqxFKWGKz92D2ibtVt91niM6', 0);
INSERT INTO "user" (full_name, phone_number, "password", balance) VALUES ('name2', '+6285544332211', '$2a$12$3gbfndmoRHh9k0qNlHL78e1tXEFceJqxFKWGKz92D2ibtVt91niM6', 10000000);

CREATE TABLE transaction (
    id UUID PRIMARY KEY,
    user_id integer NOT NULL,
    recipient_id integer NOT NULL,
    amount decimal NOT NULL,
    type text NOT NULL,
    created_time timestamp NOT NULL default now(),
    updated_time timestamp,
    status text NOT NULL,
    description text,

    CONSTRAINT fk_transaction_user_id FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);