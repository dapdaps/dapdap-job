CREATE TABLE "t_action_dapp" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "record_id" int NOT NULL,
    "count" INT NOT NULL  DEFAULT 0,
    "participants" INT NOT NULL DEFAULT 0,
    "dapp_id" int NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "idx_action_dapp_record_id" ON "t_action_dapp" ("record_id");
CREATE unique INDEX "idx_action_dapp_dapp_id" ON "t_action_dapp" ("dapp_id");

CREATE TABLE "t_action_chain" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "record_id" int NOT NULL,
    "count" INT NOT NULL DEFAULT 0,
    "action_title" varchar(512) NOT NULL,
    "dapp_id" int NOT NULL DEFAULT 0,
    "network_id" int NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL  DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL  DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "idx_action_chain_record_id" ON "t_action_chain" ("record_id");
CREATE INDEX "idx_action_chain_network_count" ON "t_action_chain" ("network_id","count");
CREATE UNIQUE INDEX "idx_action_chain_dapp_network_title" ON "t_action_chain" ("dapp_id","network_id","action_title");

CREATE TABLE "t_action_quest" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "record_id" int NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL  DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL  DEFAULT CURRENT_TIMESTAMP
);