CREATE TABLE IF NOT EXISTS market_symbol (
    guid        VARCHAR PRIMARY KEY,
    symbol VARCHAR NOT NULL,
    unified_symbol VARCHAR NOT NULL,
    inst_type   VARCHAR NOT NULL,
    exchange      VARCHAR NOT NULL,
    chain_id      VARCHAR NOT NULL,
    base      VARCHAR NOT NULL,
    quote      VARCHAR NOT NULL,
    timestamp   INTEGER NOT NULL CHECK (timestamp > 0),
);
CREATE INDEX idx_market_symbol ON symbol_mapping(exchange, chain_id,inst_type);


CREATE TABLE IF NOT EXISTS symbol_spot_prices (
    guid        VARCHAR PRIMARY KEY,
    symbol        VARCHAR NOT NULL,
    unified_symbol        VARCHAR NOT NULL,
    price VARCHAR NOT NULL,
    exchange      VARCHAR NOT NULL,
    chain_id      VARCHAR NOT NULL,
    base      VARCHAR NOT NULL,
    quote      VARCHAR NOT NULL,
    timestamp   INTEGER NOT NULL CHECK (timestamp > 0),
);
CREATE INDEX idx_symbol_spot_prices ON symbol_spot_prices(exchange, chain_id);


CREATE TABLE IF NOT EXISTS symbol_futures_prices (
    guid        VARCHAR PRIMARY KEY,
    symbol        VARCHAR NOT NULL,
    unified_symbol        VARCHAR NOT NULL,
    price VARCHAR NOT NULL,
    mark_price VARCHAR NOT NULL,
    funding_rate VARCHAR NOT NULL,
    exchange      VARCHAR NOT NULL,
    chain_id      VARCHAR NOT NULL,
    base      VARCHAR NOT NULL,
    quote      VARCHAR NOT NULL,
    timestamp   INTEGER NOT NULL CHECK (timestamp > 0),
);
CREATE INDEX idx_symbol_futures_prices ON symbol_futures_prices(exchange, chain_id);
