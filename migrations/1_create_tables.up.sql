
CREATE TABLE IF NOT EXISTS transactions (
    tx_time CHAR(19) NOT NULL,
    tx_type CHAR(1) NOT NULL,
    tx_hash CHAR(128) NOT NULL,
    tx_data TEXT NOT NULL,
    tx_prev CHAR(128) NOT NULL,
    tx_epoc TEXT NOT NULL,
    tx_subg CHAR(128) NOT NULL,
    tx_prnt CHAR(128),
    tx_mile BOOLEAN NOT NULL,
    tx_lead BOOLEAN NOT NULL
);