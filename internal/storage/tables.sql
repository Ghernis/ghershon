CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    id_ticket TEXT,
    description TEXT,
    problem_statement TEXT,
    architecture TEXT,
    evidence TEXT,
    expected_finish_date TEXT,
    completed_at TEXT,
    time_before_automation INTEGER,
    time_after_automation INTEGER,
    tags TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS snippets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    description TEXT,
    language TEXT,
    tags TEXT,
    content TEXT,
    source_file TEXT,
    start_line INTEGER,
    end_line INTEGER,
    documentation_url TEXT,
    project_id INTEGER, 
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS project_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    is_done BOOLEAN DEFAULT FALSE,
    created_at TEXT DEFAULT (datetime('now')),
    due_date TEXT,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS secrets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER,
    environment TEXT DEFAULT 'default',
    name TEXT NOT NULL,
    description TEXT,
    secret_type TEXT,
    encoded_value TEXT NOT NULL,
    is_encrypted BOOLEAN DEFAULT 1,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    UNIQUE (project_id, environment, name)
);
