CREATE TABLE IF NOT EXISTS TotalAmount (
    id BIGSERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    spent_amount INT NOT NULL,
    remaining_amount INT NOT NULL,
    statuss VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS TotalAmountCategories (
    id BIGSERIAL PRIMARY KEY,
    included_categories VARCHAR(255) NOT NULL,
    labels VARCHAR(255) NOT NULL,
    UNIQUE(included_categories)
);