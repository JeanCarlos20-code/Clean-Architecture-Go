CREATE TABLE orders (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    price FLOAT NOT NULL,
    tax FLOAT NOT NULL,
    final_price FLOAT NOT NULL,
    issue_date DATETIME NOT NULL,
    type_requisition VARCHAR(10) NOT NULL,
    delete_at DATETIME
);