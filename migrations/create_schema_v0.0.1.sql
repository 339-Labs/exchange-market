CREATE TABLE IF NOT EXISTS symbol_mapping (
    guid        VARCHAR PRIMARY KEY,
    unified_symbol        VARCHAR NOT NULL,
    symbol VARCHAR NOT NULL,
    exchange      VARCHAR NOT NULL,
    chain      VARCHAR NOT NULL,
    base      VARCHAR NOT NULL,
    quote      VARCHAR NOT NULL,
    timestamp   INTEGER NOT NULL CHECK (timestamp > 0),
);
CREATE INDEX idx_symbol_mapping_exchange ON symbol_mapping(exchange);
CREATE INDEX idx_symbol_mapping ON symbol_mapping(exchange, chain);


CREATE TABLE IF NOT EXISTS symbol_spot_prices (
                                              guid        VARCHAR PRIMARY KEY,
                                              unified_symbol        VARCHAR NOT NULL,
                                              price NUMERIC(18,8) NOT NULL,
                                              exchange      VARCHAR NOT NULL,
                                              chain      VARCHAR NOT NULL,
                                              timestamp   INTEGER NOT NULL CHECK (timestamp > 0),
);
CREATE INDEX idx_symbol_spot_prices_unified_symbol ON symbol_spot_prices(unified_symbol);
CREATE INDEX idx_symbol_spot_prices_exchange ON symbol_spot_prices(exchange);
CREATE INDEX idx_symbol_spot_prices ON symbol_spot_prices(exchange, chain);


CREATE TABLE IF NOT EXISTS symbol_futures_prices (
    guid        VARCHAR PRIMARY KEY,
    unified_symbol        VARCHAR NOT NULL,
    price NUMERIC(38,18) NOT NULL,
    funding_rate NUMERIC(10,8),
    exchange      VARCHAR NOT NULL,
    chain      VARCHAR NOT NULL,
    timestamp   INTEGER NOT NULL CHECK (timestamp > 0),
);
CREATE INDEX idx_symbol_futures_prices_unified_symbol ON symbol_futures_prices(unified_symbol);
CREATE INDEX idx_symbol_futures_prices_exchange ON symbol_futures_prices(exchange);
CREATE INDEX idx_symbol_futures_prices ON symbol_futures_prices(exchange, chain);
