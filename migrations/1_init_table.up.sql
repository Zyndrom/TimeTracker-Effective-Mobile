CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  passport_number varchar(50) NOT NULL,
  name varchar(250) not null, 
  surname varchar(250)not null,
  patronymic varchar(250),
  address varchar(255) NOT NULL
);



CREATE TABLE IF NOT EXISTS tasks (
	id serial PRIMARY KEY,
	owner int REFERENCES users(id) ON DELETE CASCADE, 
	name varchar(250), 
	created_at timestamp DEFAULT CURRENT_TIMESTAMP, 
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP, 
	active boolean DEFAULT TRUE,
	duration int DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_task_owner ON tasks(owner);